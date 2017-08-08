package model

import (
	"time"

	"crypto/rand"
	"encoding/base64"
	"io"

	"net/http"

	"blog.ka1em.site/common"
	"github.com/gorilla/sessions"
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
	SessionId     string    `json:"session_id"             gorm:"not null; default ''"`
	UserId        uint64    `json:"user_id,string"         gorm:"not null"`
	SessionStart  time.Time `json:"session_start,string"`
	SessionUpdate time.Time `json:"session_update,string"`
	SessionActive uint      `json:"session_active,string"  gorm:"not null"`

	Authenticated   bool `json:"authenticated,string"      gorm:"-"`
	Unauthenticated bool `json:"unauthenticated,string"    gorm:"-"`
	User            User `json:"user"                      gorm:"-"`
}

//var UserSession = new(Session)

var SessionStore = sessions.NewCookieStore([]byte("our-social-network-application"))

func (s *Session) GetSessionUID() error {
	if err := DB.Where("session_id = ?", s.SessionId).First(s).Error; err != nil {
		common.Suggar.Error(err.Error())
		return err
	}
	common.Suggar.Debugf("session userid = %d", s.UserId)
	return nil
}

func (s *Session) UpdateSession() error {
	const timeFmt = "2006-01-02T15:04:05.999999999"
	tstamp := time.Now().Format(timeFmt)
	if err := DB.Exec("INSERT INTO sessions SET session_id=?, user_id=?, session_update=? "+
		"ON DUPLICATE KEY UPDATE user_id=?, session_update=?", s.SessionId, s.UserId, tstamp, s.UserId, tstamp).Error; err != nil {
		//if err := DB.Exec("update sessions SET session_id=?,  session_update=?  where user_id=?,", s.SessionId, tstamp, s.UserId).Error; err != nil {
		common.Suggar.Error(err.Error())
		return err
	}
	return nil
}

func (s *Session) GenerateSessionId() (string, error) {
	sid := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, sid); err != nil {
		common.Suggar.Error(err.Error())
		return "", err
	}
	common.Suggar.Debugf("base 64 session id %s", base64.URLEncoding.EncodeToString(sid))
	return base64.URLEncoding.EncodeToString(sid), nil
}

func (s *Session) CreateSeesion(w http.ResponseWriter, r *http.Request) error {
	session, err := SessionStore.Get(r, "app-session")
	if err != nil {
		common.Suggar.Error(err.Error())
		return err
	}

	if sid, valid := session.Values["sid"]; valid {
		s.SessionId = sid.(string)
		err := s.GetSessionUID()
		if err != nil {
			return err
		}
		s.UpdateSession()
		common.Suggar.Debugf("sid = %s", sid)
	} else {
		s.SessionId, _ = s.GenerateSessionId()
		session.Values["sid"] = s.SessionId
		session.Save(r, w)
		s.UpdateSession()
		common.Suggar.Debugf("newSid = %s", s.SessionId)
	}
	return nil
}
