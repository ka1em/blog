package model

import (
	"blog.ka1em.site/common"
	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var SchemaDecoder = schema.NewDecoder()

func init() {
	var err error
	DB, err = gorm.Open("mysql", "lgwd:queeheeChiegiusheeD4@tcp(47.93.11.105:3306)/lgwd?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		common.Suggar.Error(err.Error())
		panic(err.Error())
		return
	}

	//update  tabel
	err = DB.AutoMigrate(
		&User{},
		&Session{},
		&Page{},
		&Comment{},
	).Error

	if err != nil {
		common.Suggar.Error(err.Error())
		return
	}
}
