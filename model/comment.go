package model

import (
	"time"

	"blog.ka1em.site/common"
)

type Comment struct {
	Id        uint64     `json:"id,string"       gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"               sql:"index"`

	PageId uint64 `json:"page_id,string"    gorm:"type:bigint(20)"`
	//CommentGuid  string `json:"comment_guid"      gorm:"type:varchar(256)"`
	CommentName  string `json:"comment_name"      gorm:"type:varchar(64)"`
	CommentEmail string `json:"comment_email"     gorm:"type:varchar(256)"`
	CommentText  string `json:"comment_text"      gorm:"type:mediumtext"`
}

func (c *Comment) AddComment() error {
	return DB.Create(c).Error
}

func (c *Comment) GetComment(pIndex, pSize int) (comments []*Comment, err error) {
	if err := DB.Order("created_at desc").Limit(pSize).Offset(pIndex).Find(&comments).Error; err != nil {
		common.Suggar.Error("%s", err.Error())
		return nil, err
	}
	return comments, nil
}

func (c *Comment) UpdateComment() error {
	return DB.Exec("update comments set comment_text = ? where id = ?", c.CommentText, c.Id).Error
}