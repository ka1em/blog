package model

import (
	"blog/common/zlog"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/go-errors/errors"
	"github.com/gorilla/sessions"
)

const CookieName = "app-session"
const cookieSecKey = "our-social-network-application"
const UserSession = "USER_SESSION"

var SessionStore *sessions.CookieStore

type Session struct {
	SID         string
	UserID      uint64
	CreatedUnix int64
	RedisPool   *redis.Pool `redis:"-"`
}

func init() {
	SessionStore = sessions.NewCookieStore([]byte(cookieSecKey))
}

func (s *Session) SetCache() error {
	key := REDIS_KEY_SESSION + s.SID
	r := RedisDao{
		Pool: s.RedisPool,
	}
	se := Session{
		SID:         s.SID,
		UserID:      s.UserID,
		CreatedUnix: s.CreatedUnix,
	}
	_, err := r.HMSet(key, &se)
	zlog.ZapLog.Debugf("save session key: %s session: %+v", key, se)
	return err
}

func (s *Session) GetCache() (Session, error) {
	se := Session{}
	key := REDIS_KEY_SESSION + s.SID
	r := RedisDao{
		Pool: s.RedisPool,
	}
	v, err := r.HGetAll(key)
	if err != nil {
		return se, err
	}
	err = redis.ScanStruct(v, &se)
	zlog.ZapLog.Debugf("get session key: %s session: %+v", key, se)
	return se, err
}

func (s *Session) Del() error {
	key := REDIS_KEY_SESSION + s.SID
	r := RedisDao{
		Pool: s.RedisPool,
	}
	n, err := r.Del(key)
	if err != nil {
		return err
	}
	if n == 0 {
		// todo
		return errors.New("No KEY")
	}
	zlog.ZapLog.Debugf("del session key: session: %s", key)
	return nil
}

// GenerateSessionID
func GenerateSessionID() (string, error) {
	sid := make([]byte, 24)
	if _, err := io.ReadFull(rand.Reader, sid); err != nil {
		return "", err
	}
	zlog.ZapLog.Debugf("generate session id  %s", base64.URLEncoding.EncodeToString(sid))
	return base64.StdEncoding.EncodeToString(sid), nil
}

// GetCtxSession 获取context中的session
func GetCtxSession(r *http.Request) (Session, error) {
	us := r.Context().Value(UserSession)
	if us == nil {
		return Session{}, errors.New("No ctx session")
	}
	zlog.ZapLog.Debugf("get ctx session: %+v", us.(Session))
	return us.(Session), nil
}

// GetCtxSessionUID 获取context中的session的user_id
func GetCtxSessionUID(r *http.Request) (uint64, error) {
	s, err := GetCtxSession(r)
	if s.UserID == 0 {
		return 0, errors.New("user_id is 0")
	}
	return s.UserID, err
}
