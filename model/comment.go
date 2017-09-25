package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Comment 评论
type Comment struct {
	ID          uint64     `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
	PageId      uint64     `json:"page_id,string" gorm:"type:bigint(20)"`
	UserId      uint64     `json:"user_id,string" gorm:"type:bigint(20)"`
	Text        string     `json:"comment_text" gorm:"type:mediumtext"`
	CreatedUnix int64      `json:"created_unix" gorm:"type:bigint(20)"`
	UpdatedUnix int64      `json:"updated_unix" gorm:"type:bigint(20)"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
}

func (c *Comment) BeforeCreate(scope *gorm.Scope) error {
	id, err := sf.NextID()
	if err != nil {
		return err
	}
	c.ID = id
	c.CreatedUnix = time.Now().Unix()
	c.UpdatedUnix = c.CreatedUnix
	return nil
}

func (c *Comment) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("updated_unix", time.Now().Unix())
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

// Update 更新评论
func (c *Comment) Update() error {
	return db.Model(c).Where("id = ? and user_id = ?", c.ID, c.UserId).Update("text").Error
}
