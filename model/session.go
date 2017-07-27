package model

import "time"

/*
CREATE TABLE `sessions` (
 `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
 `session_id` varchar(256) NOT NULL DEFAULT '',
  `user_id` int(11) DEFAULT NULL,
   `session_start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `session_update` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    `session_active` tinyint(1) NOT NULL,
     PRIMARY KEY (`id`),
     UNIQUE KEY `session_id` (`session_id`)) ENGINE=InnoDB DEFAULT CHARSET=latin1;
*/

type Session struct {
	Id            uint64    `json:"id,string"              gorm:"not null; AUTO_INCREMENT"`
	SessionId     string    `json:"session_id"             gorm:"not null default ''"`
	UserId        uint64    `json:"user_id,string"`
	SessionStart  time.Time `json:"session_start,string"`
	SessionUpdate time.Time `json:"session_update,string"`
	SessionActive uint      `json:"session_active,string"  gorm:"not null"`

	Authenticated   bool `json:"authenticated,string"      gorm:"-"`
	Unauthenticated bool `json:"unauthenticated,string"    gorm:"-"`
	User            User `json:"user"                      gorm:"-"`
}
