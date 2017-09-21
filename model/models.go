package model

import (
	"fmt"

	"blog/common/setting"
	"blog/common/zlog"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var decoder *schema.Decoder
var database *gorm.DB

//SchemaDecoder 获取decoder
func SchemaDecoder() *schema.Decoder {
	if decoder == nil {
		decoder = schema.NewDecoder()
	}
	return decoder
}

//DataBase 数据库
func DataBase() *gorm.DB {
	var err error
	if database == nil {
		database, err = connDatabase()
		if err != nil {
			zlog.ZapLog.Error(err.Error())
			return nil
		}
	}
	return database
}

func connDatabase() (*gorm.DB, error) {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DB_USER, setting.DB_PASS, setting.DB_HOST, setting.DB_PORT, setting.DB_BASE, setting.DB_PARM)

	if database, err = gorm.Open("mysql", dbConn); err != nil {
		zlog.ZapLog.Error(err.Error())
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
		zlog.ZapLog.Error(err.Error())
		return nil, err
	}

	return database, nil
}
