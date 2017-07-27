package common

import (
	"log"

	"blog.ka1em.site/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "lgwd:queeheeChiegiusheeD4@tcp(47.93.11.105:3306)/lgwd?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err.Error())
		return
	}

	//update  tabel
	err = DB.AutoMigrate(
		&model.User{},
		&model.Session{},
		&model.Page{},
		&model.Comment{}).Error

	if err != nil {
		log.Println(err.Error())
		return
	}
}
