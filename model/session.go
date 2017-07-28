package model

import (
	"time"

	"crypto/rand"
	"encoding/base64"
	"io"

	"blog.ka1em.site/common"
)

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

func (s *Session) GenerateSessionId() (string, error) {
	sid := make([]byte, 24)
	_, err := io.ReadFull(rand.Reader, sid)
	if err != nil {
		common.Suggar.Error(err.Error())
		return "", err
	}
	return base64.URLEncoding.EncodeToString(sid), nil
}

func (s *Session) GetSessionUID(sid string) (uint64, error) {
	if err := DB.Where("session_id = ?", sid).First(s).Error; err != nil {
		common.Suggar.Error(err.Error())
		return 0, err
	}
	return s.UserId, nil
}

func (s *Session) UpdateSession(sid string, uid uint64) error {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	//err := DB.Model(s).Update("session_update", tstamp).Error
	if err := DB.Exec("INSERT INTO sessions SET session_id=?, user_id=?, session_update=? "+
		"ON DUPLICATE KEY UPDATE user_id=?, session_update=?", sid, uid, tstamp, uid, tstamp).Error; err != nil {
		common.Suggar.Error(err.Error())
		return err
	}
	return nil
}
