package model

import (
	//"github.com/jinzhu/gorm"
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
	//Id   int
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	UserName   string
	UserGuid   string
	UserEmail  string
	UserPasswd string
	UserSalt   string
}

type Session struct {
	Id              string
	Authenticated   bool
	Unauthenticated bool
	User            User
}

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Raw        string
	Unparsed   []string
}
