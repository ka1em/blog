package model

import (
	"time"

	"blog/common/zlog"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

// Page 文章
type Page struct {
	ID          int64     `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
	Guid        string     `json:"guid" gorm:"type:varchar(64);unique_index"`
	UserId      int64     `json:"user_id,string" sql:"type:bigint(20)"'`
	Title       string     `json:"title" gorm:"type:varchar(256)"`
	Content     string     `json:"content" gorm:"type:text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
	CreatedUnix int64      `json:"created_unix,string" sql:"type:bigint(20)"`
	UpdatedUnix int64      `json:"updated_unix,string" sql:"type:bigint(20)"`

	Comments []*Comment `json:"comments,omitempty" gorm:"-"`
}

const TRUNCNUM = 20

func (p *Page) BeforeCreate(scope *gorm.Scope) error {
	id, err := sf.NextID()
	if err != nil {
		return err
	}
	p.ID = int64(id)
	p.Guid = uuid.NewV4().String()
	p.CreatedUnix = time.Now().Unix()
	p.UpdatedUnix = p.CreatedUnix
	return nil
}

func (p *Page) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("updated_unix", time.Now().Unix())
}

// Add 添加page
func (p *Page) Add() error {
	return db.Create(p).Error
}

// GetByID 获取
func GetByID(pageId uint64) (Page, error) {
	p := Page{}
	if err := db.Where("id = ?", pageId).First(&p).Error; err != nil {
		return Page{}, err
	}

	return p, nil
}

// GetAllPage 获取page
func GetAllPage(pIndex, pSize int) (pages []*Page, err error) {
	if err := db.Order("created_at  desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&pages).Error; err != nil {
		zlog.ZapLog.Error(err.Error())
		return nil, err
	}
	return pages, nil
}

// TruncatedText 截短字符串
func TruncatedText(s string) string {
	chars := 0
	for i := range s {
		chars++
		if chars > TRUNCNUM {
			return s[:i] + `...`
		}
	}
	return s
}
