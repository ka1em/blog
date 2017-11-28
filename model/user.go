package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

/*
CREATE TABLE `users` (
 `id` int(11)
  unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(32) NOT NULL DEFAULT '',
   `user_guid` varchar(256) NOT NULL DEFAULT '',
    `user_email` varchar(128) NOT NULL DEFAULT '',
    `user_password` varchar(128) NOT NULL DEFAULT '',
      `user_salt` varchar(128) NOT NULL DEFAULT '',
       `user_joined_timestamp` timestamp NULL DEFAULT NULL,
       PRIMARY KEY (`id`)
       ) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/

// User 用户表
type User struct {
	ID          uint64      `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
	Name        string      `json:"name" gorm:"not null; type:varchar(256)"`
	Email       string      `json:"email" gorm:"not null; type:varchar(256)"`
	Passwd      string      `json:"-" gorm:"not null; type:varchar(256)" redis:"-"`
	Salt        string      `json:"-" gorm:"type:varchar(256)" redis:"-"`
	Role        string      `json:"role" gorm:"not null; type:varchar(64)"` //角色 admin:管理员 users:用户
	CreatedAt   time.Time   `json:"created_at" redis:"-"`
	UpdatedAt   time.Time   `json:"updated_at" redis:"-"`
	DeletedAt   *time.Time  `sql:"index" redis:"-"`
	CreatedUnix int64       `json:"created_unix" gorm:"type:bigint(20)"`
	UpdatedUnix int64       `json:"updated_unix" gorm:"type:bigint(20)"`
	DB          *gorm.DB    `json:"-" gorm:"-" redis:"-"`
	RedisPool   *redis.Pool `redis:"-" json:"-" gorm:"-" `
}

// Create 创建用户
func (u *User) Create() error {
	if u.DB == nil {
		u.DB = db
	}
	if u.nameIsExist(u.Name) {
		return errors.New(ErrMap[UserNameExist])
	}
	return u.DB.Create(u).Error
}

// SetCache 缓存用户
func (u *User) SetCache() (string, error) {
	if u.RedisPool == nil {
		u.RedisPool = redisPool
	}
	key := REDIS_KEY_USER + fmt.Sprintf("%d", u.ID)
	r := RedisDao{
		Pool: u.RedisPool,
	}
	return r.HMSet(key, u)
}

// GetCache 获取缓存用户信息
func (u *User) GetCache(id uint64) (User, error) {
	user := User{}
	if u.RedisPool == nil {
		u.RedisPool = redisPool
	}
	r := RedisDao{
		Pool: u.RedisPool,
	}
	key := REDIS_KEY_USER + fmt.Sprintf("%d", id)
	v, err := r.HGetAll(key)
	if err != nil {
		return user, err
	}
	if err := redis.ScanStruct(v, &user); err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) nameIsExist(name string) bool {
	return !u.DB.Where("name = ?", name).First(&User{}).RecordNotFound()
}

// BeforeCreate 创建之前
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	var err error
	u.Salt = uuid.NewV4().String()
	u.Passwd, err = passwordHash(u.Passwd, u.Salt)
	if err != nil {
		return err
	}
	u.CreatedUnix = time.Now().Unix()
	u.UpdatedUnix = time.Now().Unix()
	return nil
}

// BeforeUpdate 更新之前
func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("updated_unix", time.Now().Unix())
}

func passwordHash(p, salt string) (string, error) {
	hash := sha256.New()
	s := p + salt
	if _, err := hash.Write([]byte(s)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// getValidInfo 获取需要确认用户的信息
func getValidInfo(userName string) (*User, bool) {
	u := &User{}
	info := []string{"id", "name", "salt", "passwd"}
	if db.Select(info).Where("name = ?", userName).First(u).RecordNotFound() {
		return nil, false
	}
	return u, true
}

// CheckPassWord 检查密码
func CheckPassWord(name, passwd string) (*User, bool, error) {
	u, ok := getValidInfo(name)
	if !ok {
		return nil, false, errors.New(ErrMap[NoUserName])
	}
	tmp, err := passwordHash(passwd, u.Salt)
	if err != nil {
		return nil, false, err
	}

	if u.Passwd != tmp {
		return nil, false, errors.New(ErrMap[PasswordErr])
	}
	return u, true, nil
}
