package model

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-errors/errors"
	"github.com/go-xorm/xorm"
)

// Comment 评论
type Comment struct {
	ID       uint64 `json:"id,string" xorm:"pk not null"`
	PageID   uint64 `json:"page_id,string" xorm:"not null"`
	UserID   uint64 `json:"user_id,string" xorm:"not null"`
	Text     string `json:"comment_text" xorm:"text"`
	ToUserID uint64 `json:"to_user_id,string"`

	Created     time.Time `xorm:"-" redis:"-"`
	CreatedUnix int64     `json:"created_unix"`
	Updated     time.Time `xorm:"-" redis:"-"`
	UpdatedUnix int64     `json:"updated_unix"`

	XDB       *xorm.Engine `json:"-" xorm:"-" redis:"-"`
	RedisPool *redis.Pool  `json:"-" xorm:"-" redis:"-"`
}

func (c *Comment) BeforeInsert() {
	c.CreatedUnix = time.Now().Unix()
	c.UpdatedUnix = c.CreatedUnix
}

func (c *Comment) BeforeUpdate() {
	c.UpdatedUnix = time.Now().Unix()
}

func (c *Comment) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		c.Created = time.Unix(c.CreatedUnix, 0).Local()
	case "updated_unix":
		c.Updated = time.Unix(c.UpdatedUnix, 0).Local()
	}
}

// Add 添加评论
func (c *Comment) Add() error {
	if c.XDB == nil {
		c.XDB = xdb
	}
	_, err := c.XDB.Insert(c)
	return err
}

// Get 获取评论
func (c *Comment) Get(pIndex, pSize int) ([]Comment, error) {
	if c.XDB == nil {
		c.XDB = xdb
	}
	var list []Comment
	err := c.XDB.Where("page_id = ?", c.PageID).Limit(pSize, (pIndex-1)*pSize).Find(&list)
	return list, err
}

// Update 更新评论
func (c *Comment) Update() error {
	if c.XDB == nil {
		c.XDB = xdb
	}
	n, err := c.XDB.Where("id = ? and user_id = ?", c.ID, c.UserID).Update(c)
	if err != nil {
		return err
	}
	if n == 0 {
		// TODO
		return errors.New("no recored")
	}
	return nil
}
