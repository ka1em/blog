package model

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	zlog "blog/common/zlog"

	"blog/common/setting"

	"strconv"

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

func GetUserID(sessionId string, active int) (uint64, error) {
	s := Session{}
	if err := DataBase().Select("user_id").Where("session_id = ? and session_active = ?", sessionId, active).First(&s).Error; err != nil {
		return 0, err
	}
	return s.UserId, nil
}

func UpdateSession(userId uint64, sessionId string) error {
	sess := DataBase().Begin()

	if err := sess.Exec("update sessions set session_active = 0 where user_id = ?", userId).Error; err != nil {
		sess.Rollback()
		return err
	}

	sql := "INSERT INTO sessions (session_id,user_id,session_update,session_active) VALUES (?,?,?,?)" +
		"ON DUPLICATE KEY UPDATE user_id=?, session_update=?,session_active=?"

	if err := sess.Exec(sql, sessionId, userId, time.Now().Format(time.RFC3339), 0,
		userId, time.Now().Format(time.RFC3339), 1).Error; err != nil {
		sess.Rollback()
		return err
	}

	return sess.Commit().Error
}

func generateSessionId() (string, error) {
	sid := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, sid); err != nil {
		return "", err
	}
	zlog.ZapLog.Debugf("generate session id  %s", base64.URLEncoding.EncodeToString(sid))
	return base64.StdEncoding.EncodeToString(sid), nil
}

// CreateSession 创建session
func CreateSession(w http.ResponseWriter, r *http.Request) (string, error) {
	sessionStore := GetSessionStore()
	session, err := sessionStore.Get(r, "app-session")
	if err != nil {
		return "", err
	}

	if sid, valid := session.Values["sid"]; valid {
		userId, err := GetUserID(sid.(string), 0)
		if err != nil {
			return "", err
		}
		UpdateSession(userId, sid.(string))
		zlog.ZapLog.Debugf("sid = %s", sid)
		return sid.(string), nil
	} else {
		sessionId, err := generateSessionId()
		if err != nil {
			return "", err
		}
		session.Values["sid"] = sessionId
		session.Save(r, w)
		UpdateSession(0, sessionId)
		return sessionId, nil
	}
}

// PreCreateSession 验证用户名之前，先行创建session
func PreCreateSession(w http.ResponseWriter, r *http.Request) (string, error) {
	sessionStore := GetSessionStore()
	session, err := sessionStore.Get(r, "app-session")
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		return "", err
	}
	sessionId, err := generateSessionId()
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		return "", err
	}
	session.Values["sid"] = sessionId
	session.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   setting.SSL_ON,
	}

	session.Save(r, w)
	UpdateSession(0, sessionId)
	return sessionId, nil
}

func (s *Session) Close() error {
	s.SessionActive = 0
	return DataBase().Model(s).Where("user_id = ?", s.UserId).Update("session_active").Error
}

func SessionGetUserID(r *http.Request) (uint64, error) {
	var err error
	var uid uint64
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		return 0, err
	}

	uid, err = strconv.ParseUint(userIds.(string), 10, 64)
	if err != nil {
		return 0, err
	}
	return uid, nil
}
