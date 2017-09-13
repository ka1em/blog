package model

import (
	"errors"
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

func (u *User) TabelName() string {
	return "users"
}

func (u *User) CreateUser() error {
	//判断用户名是否存在
	if !DataBase().Where("user_name = ?", u.UserName).First(&User{}).RecordNotFound() {
		return errors.New("exists")
	}

	return DataBase().Create(u).Error
}

func GetValidInfo(userName string) (*User, bool) {
	u := &User{}
	info := []string{"user_name", "user_salt", "user_passwd"}
	if DataBase().Select(info).Where("user_name = ?", userName).First(u).RecordNotFound() {
		return nil, false
	}
	return u, true
}

func (u *User) Login() bool {
	return DataBase().Where("user_name = ? and user_passwd = ?", u.UserName, u.UserPasswd).Find(u).RecordNotFound()
}
