package model

import (
	"fmt"

	"blog/common/log"
	"blog/common/setting"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var decoder *schema.Decoder
var database *gorm.DB

//获取decoder
func SchemaDecoder() *schema.Decoder {
	if decoder == nil {
		decoder = schema.NewDecoder()
	}
	return decoder
}

//获取 数据库
func DataBase() *gorm.DB {
	var err error
	if database == nil {
		database, err = connDatabase()
		if err != nil {
			log.Suggar.Error(err.Error())
			return nil
		}
	}
	return database
}

func connDatabase() (*gorm.DB, error) {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)

	if database, err = gorm.Open("mysql", dbConn); err != nil {
		log.Suggar.Error(err.Error())
		panic(err.Error())
		return nil, err
	}

	//update  tabel
	if err := database.AutoMigrate(
		&User{},
		&Session{},
		&Page{},
		&Comment{},
	).Error; err != nil {
		log.Suggar.Error(err.Error())
		return nil, err
	}

	return database, nil
}
