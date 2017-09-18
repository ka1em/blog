package model

import (
	"time"

	"blog/common"
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

func (p *Page) Add() error {
	return DataBase().Create(p).Error
}

func GetByID(pageId uint64) (Page, error) {
	p := Page{}
	if err := DataBase().Where("id = ?", pageId).First(&p).Error; err != nil {
		return Page{}, err
	}

	return p, nil
}

//// GET page by page_guid
//func (p *Page) GetByPageGUID(pageGUID string) error {
//	return DataBase().Where("page_guid = ?", pageGUID).First(p).Error
//}

func GetAllPage(pIndex, pSize int) (pages []*Page, err error) {
	if err := DataBase().Order("created_at  desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&pages).Error; err != nil {
		common.Suggar.Error(err.Error())
		return nil, err
	}
	return pages, nil
}

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
