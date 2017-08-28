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
     UNIQUE KEY `session_id` (`session_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
*/

type Session struct {
	Id            uint64    `json:"id,string"              gorm:"not null;primary_key; AUTO_INCREMENT"`
	SessionId     string    `json:"session_id"             gorm:"type:varchar(191); unique_index;not null;default '' "`
	UserId        uint64    `json:"user_id,string"         gorm:"not null"`
	SessionStart  time.Time `json:"session_start,string"   sql:"DEFAULT:current_timestamp"`
	SessionUpdate time.Time `json:"session_update,string"`
	SessionActive uint      `json:"session_active,string"  gorm:"type:tinyint(1); not null; default 0"`

	Authenticated   bool `json:"authenticated,string"      gorm:"-"`
	Unauthenticated bool `json:"unauthenticated,string"    gorm:"-"`
	User            User `json:"user"                      gorm:"-"`
}

var sessionStore *sessions.CookieStore

func GetSessionStore() *sessions.CookieStore {
	if sessionStore == nil {
		sessionStore = sessions.NewCookieStore([]byte("our-social-network-application"))
	}
	return sessionStore
}

func (s *Session) GetSessionUID() error {
	return DataBase().Where("session_id = ? and session_active = 1", s.SessionId).Order("session_update desc").First(s).Error
}

func (s *Session) UpdateSession() error {
	return DataBase().Exec("INSERT INTO sessions (session_id,user_id,session_update,session_active) VALUES (?,?,?,?)"+
		"ON DUPLICATE KEY UPDATE user_id=?, session_update=?,session_active=?", s.SessionId, s.UserId, time.Now().Format(time.RFC3339), 1,
		s.UserId, time.Now().Format(time.RFC3339), 1).Error
}

func (s *Session) GenerateSessionId() (string, error) {
	sid := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, sid); err != nil {
		return "", err
	}
	common.Suggar.Debugf("generate session id  %s", base64.URLEncoding.EncodeToString(sid))
	return base64.StdEncoding.EncodeToString(sid), nil
}

func (s *Session) CreateSeesion(w http.ResponseWriter, r *http.Request) error {
	sessionStore := GetSessionStore()
	session, err := sessionStore.Get(r, "app-session")
	if err != nil {
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

func (s *Session) CloseSession() error {
	return DataBase().Exec("update sessions set session_active = 0 where user_id = ?", s.UserId).Error
}
