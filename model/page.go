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

	PageGuid   string    `json:"page_guid"           gorm:"type:varchar(64);unique_index"`
	Title      string    `json:"title"               gorm:"type:varchar(256)"`
	RawContent string    `json:"-"                   gorm:"type:text"`
	Content    string    `json:"content"             gorm:"type:text"`
	Comments   []Comment `json:"comments"            gorm:"-"`
	Session    Session   `json:"-"                   gorm:"-"`
}

type Comment struct {
	Id        uint64     `json:"id,string"           gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"                   sql:"index"`

	PageId  uint64 `json:"page_id,string"            gorm:"type:bigint(20)"`
	Content string `json:"content"                   gorm:"type:text"`
}

const TRUNCNUM = 20

// GET page by page_guid
func (p *Page) GetByPageGUID(pageGUID string) error {
	if err := DB.Where("page_guid = ?", pageGUID).First(p).Error; err != nil {
		common.Suggar.Error(err.Error())
		return err
	}
	return nil
}

func (p *Page) GetAllPage(pIndex, pSize int) (pages []*Page, err error) {
	if err := DB.Order("created_at  desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&pages).Error; err != nil {
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
