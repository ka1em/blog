package model

import (
	"blog/common/setting"
	"blog/common/zlog"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-errors/errors"
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
const cookieName = "app-session"
const cookieSecKey = "our-social-network-application"

var SessionStore *sessions.CookieStore

type Session struct {
	ID        int64      `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
	CreatedAt time.Time  `json:"created_at"  sql:"DEFAULT:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`

	SessionID   string `json:"session_id" gorm:"type:varchar(191); unique_index;not null;default '' "`
	UserID      int64  `json:"user_id,string" gorm:"not null"`
	Active      int    `json:"session_active,string" gorm:"type:tinyint(1); not null; default 0"`
	CreatedUnix int64  `json:"created_unix" gorm:""`
	UpdatedUnix int64  `json:"updated_unix" gorm:""`
}

func init() {
	SessionStore = sessions.NewCookieStore([]byte(cookieSecKey))
}

// GetUserID 通过sessionId从数据库查询userid
func (s *Session) GetUserID(sessionID string, active int) (int64, bool) {
	ok := !db.Select("user_id").Where("session_id = ? and active = ?", sessionID, active).First(&s).RecordNotFound()
	return s.UserID, ok
}

// Close session
func (s *Session) Close() error {
	s.Active = 0
	return db.Model(s).Where("user_id = ?", s.UserID).Update("active").Error
}

// UpdateSession 更新session为活跃状态
func UpdateSession(userId int64, sessionId string) error {
	sess := db.Begin()

	if err := sess.Exec("update sessions set active = 0 where user_id = ?", userId).Error; err != nil {
		sess.Rollback()
		return err
	}

	sql := "INSERT INTO sessions (session_id,user_id,updated_at,active) VALUES (?,?,?,?)" +
		"ON DUPLICATE KEY UPDATE user_id=?, updated_at=?,active=?"

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
	session, err := SessionStore.Get(r, cookieName)
	if err != nil {
		return "", err
	}

	if sid, ok := session.Values["sid"]; ok {
		s := &Session{}
		userId, ok := s.GetUserID(sid.(string), 0)
		if !ok {
			return "", errors.New("no user id in CreateSession")
		}
		if err = UpdateSession(userId, sid.(string)); err != nil {
			return "", err
		}
		return sid.(string), nil
	}

	sessionID, err := generateSessionId()
	if err != nil {
		return "", err
	}
	session.Values["sid"] = sessionID

	if err := session.Save(r, w); err != nil {
		return "", err
	}

	if err := UpdateSession(0, sessionID); err != nil {
		return "", err
	}

	return sessionID, nil

}

// PreCreateSession 验证用户名之前，先行创建session
func PreCreateSession(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := SessionStore.Get(r, cookieName)
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		return "", err
	}

	sessionID, err := generateSessionId()
	if err != nil {
		zlog.ZapLog.Error(err.Error())
		return "", err
	}

	session.Values["sid"] = sessionID

	session.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24,
		HttpOnly: true,
		Secure:   setting.SSLMode,
	}

	if err := session.Save(r, w); err != nil {
		return "", err
	}

	if err := UpdateSession(0, sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

// ValidSessionUID 获取context中的user_id
func ValidSessionUID(r *http.Request) (int64, error) {
	var userIds interface{}

	if userIds = r.Context().Value("user_id"); userIds == nil {
		return -1, errors.New("valid session user id is nil")
	}

	uid, err := strconv.ParseInt(userIds.(string), 10, 64)
	if err != nil {
		return -1, err
	}

	return uid, nil
}
