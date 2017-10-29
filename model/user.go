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

// User 用户表
type User struct {
	ID          int64      `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
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

// Create 创建用户
func (u *User) Create() error {
	if nameIsExist(u.Name) {
		return errors.New(ErrMap[UserNameExist])
	}
	return db.Create(u).Error
}

func nameIsExist(name string) bool {
	return !db.Where("name = ?", name).First(&User{}).RecordNotFound()
}

// BeforeCreate 创建之前
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	id, err := sf.NextID()
	if err != nil {
		return err
	}
	u.ID = int64(id)
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
