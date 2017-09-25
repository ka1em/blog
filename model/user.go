package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

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

// User 用户
type User struct {
	ID          uint64     `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
	Name        string     `json:"name" gorm:"not null; type:varchar(256)"`
	Email       string     `json:"email" gorm:"not null; type:varchar(256)"`
	Passwd      string     `json:"-" gorm:"not null; type:varchar(256)"`
	Salt        string     `json:"-" gorm:"type:varchar(256)"`
	Role        string     `json:"role" gorm:"not null; type:varchar(64)"` //角色 admin:管理员 users:用户
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index"`
	CreatedUnix int64      `json:"created_unix" gorm:"type:bigint(20)"`
	UpdatedUnix int64      `json:"updated_unix" gorm:"type:bigint(20)"`
}

// CreateUser 创建用户
func (u *User) CreateUser() error {
	if nameIsExist(u.Name) {
		return errors.New("exists")
	}
	return db.Create(u).Error
}

func nameIsExist(name string) bool {
	return !db.Where("name = ?", name).First(&User{}).RecordNotFound()
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	id, err := sf.NextID()
	if err != nil {
		return err
	}

	u.ID = id
	u.Salt = uuid.NewV4().String()
	u.Passwd = PasswordHash(u.Passwd, u.Salt)
	u.CreatedUnix = time.Now().Unix()
	u.UpdatedUnix = time.Now().Unix()

	return nil
}

func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("updated_unix", time.Now().Unix())
}

func PasswordHash(p, salt string) string {
	hash := sha256.New()
	s := p + salt
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

// GetValidInfo 获取需要确认用户的信息
func GetValidInfo(userName string) (*User, bool) {
	u := &User{}
	info := []string{"id", "name", "salt", "passwd"}
	if db.Select(info).Where("name = ?", userName).First(u).RecordNotFound() {
		return nil, false
	}
	return u, true
}
