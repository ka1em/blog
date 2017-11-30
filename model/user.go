package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
	uuid "github.com/satori/go.uuid"
)

// User 用户表
type User struct {
	ID     uint64 `json:"id,string" xorm:"pk notnull"`
	Name   string `json:"name" xorm:"varchar(64) notnull unique"`
	Email  string `json:"email" xorm:"varchar(64) notnull"`
	Passwd string `json:"-" xorm:"varchar(256) notnull" redis:"-"`
	Salt   string `json:"-" xorm:"varchar(256)" redis:"-"`
	Role   string `json:"role" xorm:"varchar(64)"` //角色 admin:管理员 users:用户

	Created     time.Time `xorm:"-" redis:"-"`
	CreatedUnix int64     `json:"created_unix"`
	Updated     time.Time `xorm:"-" redis:"-"`
	UpdatedUnix int64     `json:"updated_unix"`

	XDB       *xorm.Engine `json:"-" xorm:"-" redis:"-"`
	RedisPool *redis.Pool  `json:"-" xorm:"-" redis:"-"`
}

func (u *User) BeforeInsert() {
	u.CreatedUnix = time.Now().Unix()
	u.UpdatedUnix = u.CreatedUnix
}

func (u *User) BeforeUpdate() {
	u.UpdatedUnix = time.Now().Unix()
}

func (u *User) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		u.Created = time.Unix(u.CreatedUnix, 0).Local()
	case "updated_unix":
		u.Updated = time.Unix(u.UpdatedUnix, 0).Local()
	}
}

// Create 创建用户
func (u *User) Create() error {
	if u.XDB == nil {
		u.XDB = xdb
	}
	exist, err := u.NameWasExist()
	if err != nil {
		return err
	}
	if exist {
		return errors.New(ErrMap[UserNameExist])
	}
	u.Salt = uuid.NewV4().String()
	u.Passwd, err = passwordHash(u.Passwd, u.Salt)
	if err != nil {
		return err
	}
	_, err = u.XDB.Insert(u)
	return err
}

// SetCache 缓存用户
func (u *User) SetCache() (string, error) {
	r := RedisDao{
		Pool: u.RedisPool,
	}
	key := REDIS_KEY_USER + fmt.Sprintf("%d", u.ID)
	return r.HMSet(key, u)
}

// GetCache 获取缓存用户信息
func (u *User) GetCache(id uint64) (User, error) {
	user := User{}
	r := RedisDao{
		Pool: u.RedisPool,
	}
	key := REDIS_KEY_USER + fmt.Sprintf("%d", id)
	v, err := r.HGetAll(key)
	if err != nil {
		return user, err
	}
	if err := redis.ScanStruct(v, &user); err != nil {
		return user, err
	}
	return user, nil
}

// NameWasExist 用户名存在
func (u *User) NameWasExist() (bool, error) {
	return u.XDB.Where("name = ?", u.Name).Get(u)
}

func passwordHash(p, salt string) (string, error) {
	hash := sha256.New()
	s := p + salt
	if _, err := hash.Write([]byte(s)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// getValidInfo 获取需要确认用户的信息
func (u *User) getValidInfo(userName string) (*User, error) {
	if u.XDB == nil {
		u.XDB = xdb
	}
	user := &User{}
	exist, err := u.XDB.Where("name = ?", userName).Get(user)
	if err != nil {
		return user, err
	}
	if !exist {
		return user, errors.New(ErrMap[NoUserName])
	}
	return user, err
}

// CheckPassWord 检查密码
func (u *User) CheckPassWord() (*User, error) {
	if u.XDB == nil {
		u.XDB = xdb
	}
	realUser, err := u.getValidInfo(u.Name)
	if err != nil {
		return realUser, errors.New(ErrMap[NoUserName])
	}
	tmp, err := passwordHash(u.Passwd, realUser.Salt)
	if err != nil {
		return realUser, errors.New(ErrMap[PasswordHashErr])
	}
	if realUser.Passwd != tmp {
		return realUser, errors.New(ErrMap[PasswordErr])
	}
	return realUser, nil
}
