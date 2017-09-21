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

// Add 添加评论
func (c *Comment) Add() error {
	return db.Create(c).Error
}

// Get 获取评论
func (c *Comment) Get(pIndex, pSize int) (comments []*Comment, err error) {
	if err := db.Order("created_at desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// UpdateComment 更新评论
func (c *Comment) Update() error {
	return db.Model(c).Where("id = ? and user_id = ?", c.Id, c.UserId).Update("comment_text").Error
}
