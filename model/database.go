package model

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"database/sql"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open("mysql", "lgwd:queeheeChiegiusheeD4@tcp(47.93.11.105:3306)/lgwd?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err.Error())
		return
	}

	//更新最新数据
	if err := DB.AutoMigrate(&User{}, &Session{}, &Cookie{}, &Page{}).Error; err != nil {
		log.Println(err.Error())
		return
	}
}
