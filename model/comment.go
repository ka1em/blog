package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Comment 评论
type Comment struct {
	ID          int64      `json:"id,string" gorm:"primary_key" sql:"type:bigint(20)"`
	PageID      int64      `json:"page_id,string" gorm:"type:bigint(20)"`
	UserID      uint64     `json:"user_id,string" gorm:"type:bigint(20)"`
	Text        string     `json:"comment_text" gorm:"type:mediumtext"`
	CreatedUnix int64      `json:"created_unix" gorm:"type:bigint(20)"`
	UpdatedUnix int64      `json:"updated_unix" gorm:"type:bigint(20)"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
}

// BeforeCreate 评论创建前的操作
func (c *Comment) BeforeCreate(scope *gorm.Scope) error {
	id, err := SF.NextID()
	if err != nil {
		return err
	}
	c.ID = int64(id)
	c.CreatedUnix = time.Now().Unix()
	c.UpdatedUnix = c.CreatedUnix
	return nil
}

// BeforeUpdate 评论更新前操作
func (c *Comment) BeforeUpdate(scope *gorm.Scope) error {
	return scope.SetColumn("updated_unix", time.Now().Unix())
}

// Add 添加评论
func (c *Comment) Add() error {
	return db.Create(c).Error
}

// Get 获取评论
func (c *Comment) Get(pIndex, pSize int) ([]*Comment, error) {
	list := []*Comment{}
	if err := db.Order("created_at desc").Limit(pSize).Offset((pIndex - 1) * pSize).Find(&list).Error; err != nil {
		return list, err
	}
	return list, nil
}

// Update 更新评论
func (c *Comment) Update() error {
	return db.Model(c).Where("id = ? and user_id = ?", c.ID, c.UserID).Update("text").Error
}
