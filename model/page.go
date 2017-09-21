package model

import (
	"time"

	"blog/common/zlog"
)

type Page struct {
	Id        uint64     `json:"id,string"           gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"                   sql:"index"`

	PageGuid string     `json:"page_guid"           gorm:"type:varchar(64);unique_index"`
	Title    string     `json:"title"               gorm:"type:varchar(256)"`
	Content  string     `json:"content"             gorm:"type:text"`
	Comments []*Comment `json:"comments,omitempty"  gorm:"-"`
	Session  Session    `json:"-"                   gorm:"-"`
}

const TRUNCNUM = 20

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
			return s[:i] + ` ...`
		}
	}
	return s
}
