package model

import (
	"errors"
	"log"
	"time"
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
	ID         uint64 `json:"id,string"          gorm:"primary_key"`
	UserName   string `json:"user_name"          gorm:"not null; type:varchar(256)"`
	UserGuid   string `json:"user_guid"          gorm:"type:varchar(256)" `
	UserEmail  string `json:"user_email"         gorm:"not null; type:varchar(256)"`
	UserPasswd string `json:"-"                  gorm:"not null; type:varchar(256)"`
	UserSalt   string `json:"-"                  gorm:"type:varchar(256)"`

	Role string `json:"role" gorm:"not null; type:varchar(64)"` //角色 admin:管理员 users:用户

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// CreateUser 创建用户
func (u *User) CreateUser() error {
	//判断用户名是否存在
	if !db.Where("user_name = ?", u.UserName).First(&User{}).RecordNotFound() {
		return errors.New("exists")
	}

	return db.Create(u).Error
}

// GetValidInfo 获取需要确认用户的信息
func GetValidInfo(userName string) (*User, bool) {
	u := &User{}
	info := []string{"id", "user_name", "user_salt", "user_passwd"}
	if db.Select(info).Where("user_name = ?", userName).First(u).RecordNotFound() {
		log.Printf("%+v", *u)
		return nil, false
	}
	return u, true
}
