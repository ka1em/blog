package model

import (
	"time"

	"blog.ka1em.site/common"
)

type Page struct {
	Id        uint64     `json:"id,string"           gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"                   sql:"index"`

	PageGuid string     `json:"page_guid"           gorm:"type:varchar(64);unique_index"`
	Title    string     `json:"title"               gorm:"type:varchar(256)"`
	Content  string     `json:"content"             gorm:"type:text"`
	Comments []*Comment `json:"comments,omitempty"            gorm:"-"`
	Session  Session    `json:"-"                   gorm:"-"`
}

const TRUNCNUM = 20

func (p *Page) AddPage() error {
	db := GetDB()

	return db.Create(p).Error
}

func (p *Page) GetByID() error {
	var err error
	db := GetDB()

	if err := db.Where("id = ?", p.Id).First(p).Error; err != nil {
		return err
	}

	c := &Comment{
		PageId: p.Id,
	}

	p.Comments, err = c.GetComment(1, 10)
	if err != nil {
		return err
	}

	common.Suggar.Debugf("p.Comments %+v", p.Comments)

	return nil
}

// GET page by page_guid
func (p *Page) GetByPageGUID(pageGUID string) error {
	db := GetDB()

	return db.Where("page_guid = ?", pageGUID).First(p).Error
}

func (p *Page) GetAllPage(pIndex, pSize int) (pages []*Page, err error) {
	db := GetDB()

	if err := db.Order("created_at  desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&pages).Error; err != nil {
		common.Suggar.Error(err.Error())
		return nil, err
	}
	return pages, nil
}

func (p *Page) TruncatedText() string {
	chars := 0
	for i, _ := range p.Content {
		chars++
		if chars > TRUNCNUM {
			return p.Content[:i] + ` ...`
		}
	}
	return p.Content
}
