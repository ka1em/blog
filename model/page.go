package model

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Page struct {
	gorm.Model
	Title      string `gorm:varchar(256)`
	RawContent string `gorm:"text"`
	Content    string `gorm:"text"`
	Comments   []Comment
	Session    Session `json:"-" gorm:"-"`

	//Content    template.HTML `gorm:"text"`
	//Date       string
}

type Comment struct {
	gorm.Model
	Content string `gorm:"text"`
}

func (p *Page) GetByPageID(pageGUID string, page *Page) error {
	if err := DB.Where("page_guid = ?", pageGUID).First(page).Error; err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
