package model

import (
	"github.com/jinzhu/gorm"
	"html/template"
	"log"
)

type Page struct {
	gorm.Model
	Title      string
	RawContent string
	Content    template.HTML
	Date       string
	Comments   []Comment
	Session    Session
}

type Comment struct {
	gorm.Model
	Content string
}

func (p *Page) GetByPageID(pageGUID string, page *Page) error {
	//sql := "select page_title,page_content,page_date from pages where page_guid=?"
	//err := r.DB.QueryRow(sql, pageGUID).Scan(&p.Title, &p.Content, &p.Date)
	//if err != nil {
	//	return err
	//}
	if err := DB.Where("page_guid = ?", pageGUID).First(page).Error; err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
