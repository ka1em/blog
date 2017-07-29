package model

import (
	"time"

	"blog.ka1em.site/common"
)

type Page struct {
	Id        uint64     `json:"id,string" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	PageGuid   string    `json:"page_guid" gorm:"varchar(64)"`
	Title      string    `json:"title" gorm:"varchar(256)"`
	RawContent string    `json:"raw_content" gorm:"text"`
	Content    string    `json:"content" gorm:"text"`
	Comments   []Comment `json:"comments"`
	Session    Session   `json:"-" gorm:"-"`

	//Content    template.HTML `gorm:"text"`
	//Date       string
}

type Comment struct {
	Id        uint64     `json:"id,string" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`

	Content string `json:"content" gorm:"text"`
}

// GET page by page_guid
func (p *Page) GetByPageGUID(pageGUID string) error {
	if err := DB.Where("page_guid = ?", pageGUID).First(p).Error; err != nil {
		common.Suggar.Error(err.Error())
		return err
	}

	return nil
}
