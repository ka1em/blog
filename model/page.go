package model

import (
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-xorm/xorm"
)

// Page 文章
type Page struct {
	ID      uint64 `json:"id,string" xorm:"pk notnull"`
	Guid    string `json:"guid" xorm:"varchar(64)"`
	UserID  uint64 `json:"user_id,string"`
	Title   string `json:"title" xorm:"varchar(256)"`
	Content string `json:"content" xorm:"text"`

	Created     time.Time `json:"-" xorm:"-" redis:"-"`
	CreatedUnix int64     `json:"created_unix"`
	Updated     time.Time `json:"-" xorm:"-" redis:"-"`
	UpdatedUnix int64     `json:"updated_unix"`

	XDB       *xorm.Engine `json:"-" xorm:"-" redis:"-"`
	RedisPool *redis.Pool  `json:"-" redis:"-" json:"-" xorm:"-" `

	Comments []*Comment `json:"comments,omitempty" xorm:"-" redis:"-"`
}

const contentLen = 20

func (p *Page) BeforeInsert() {
	p.CreatedUnix = time.Now().Unix()
	p.UpdatedUnix = p.CreatedUnix
}

func (p *Page) BeforeUpdate() {
	p.UpdatedUnix = time.Now().Unix()
}

func (p *Page) AfterSet(colName string, _ xorm.Cell) {
	switch colName {
	case "created_unix":
		p.Created = time.Unix(p.CreatedUnix, 0).Local()
	case "updated_unix":
		p.Updated = time.Unix(p.UpdatedUnix, 0).Local()
	}
}

// Add 添加page
func (p *Page) Add() error {
	if p.XDB == nil {
		p.XDB = xdb
	}
	_, err := p.XDB.Insert(p)
	return err
}

// GetPageByID 获取
func (p *Page) GetPageByID() (Page, error) {
	if p.XDB == nil {
		p.XDB = xdb
	}
	page := Page{}
	exist, err := p.XDB.Where("id = ?", p.ID).Get(&page)
	if !exist {
		return page, errors.New("not exist")
	}
	return page, err
}

// GetPages 获取page
func (p *Page) GetPages(pIndex, pSize int) (pages []*Page, err error) {
	if p.XDB == nil {
		p.XDB = xdb
	}
	err = p.XDB.Limit(pSize, (pIndex-1)*pSize).Desc("created").Find(&pages)
	return
}

// TruncatedText 截短字符串
func TruncatedText(s string) string {
	chars := 0
	for i := range s {
		chars++
		if chars > contentLen {
			return s[:i] + `...`
		}
	}
	return s
}
