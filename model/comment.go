package model

import (
	"time"
)

type Comment struct {
	Id        uint64     `json:"id,string"       gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"               sql:"index"`

	PageId uint64 `json:"page_id,string"    gorm:"type:bigint(20)"`
	UserId uint64 `json:"user_id,string"    gorm:"type:bigint(20)"`
	//CommentGuid  string `json:"comment_guid"      gorm:"type:varchar(256)"`
	CommentName  string `json:"comment_name"      gorm:"type:varchar(64)"`
	CommentEmail string `json:"comment_email"     gorm:"type:varchar(256)"`
	CommentText  string `json:"comment_text"      gorm:"type:mediumtext"`
}

func (c *Comment) AddComment() error {
	return DataBase().Create(c).Error
}

func (c *Comment) GetComment(pIndex, pSize int) (comments []*Comment, err error) {
	if err := DataBase().Order("created_at desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *Comment) UpdateComment() error {
	return DataBase().Exec("update comments set comment_text = ? where id = ? and user_id = ?", c.CommentText, c.Id, ).Error
}
