package model

import (
	"time"

	"blog.ka1em.site/common"
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
	//Id   int
	ID         uint64 `json:"id,string"     gorm:"primary_key"`
	UserName   string `json:"user_name"     gorm:"not null; varchar(256)" form:"user_name"`
	UserGuid   string `json:"user_guid"     gorm:"not null; varchar(256)" form:"user_guid"`
	UserEmail  string `json:"user_email"    gorm:"not null; varchar(256)" form:"user_email"`
	UserPasswd string `json:"user_passwd"   gorm:"not null; varchar(256)" form:"user_passwd"`
	UserSalt   string `json:"-"             gorm:"varchar(256)`

	Role string `json:"role" gorm:"not nulll; varchar(64)"` //角色 admin:管理员 users:用户

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (u *User) TabelName() string {
	if u.Role == "admin" {
		return "admin_users"
	} else {
		return "users"
	}
}

func (u *User) CreateUser() error {
	err := DB.Create(u).Error
	if err != nil {
		common.Suggar.Error(err.Error)
		return err
	}
	return nil
}
