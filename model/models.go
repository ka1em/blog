package model

import (
	"fmt"

	"blog/common/setting"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	SchemaDecoder *schema.Decoder
	db            *gorm.DB
)

func init() {
	SchemaDecoder = schema.NewDecoder()

	err := initDB()
	if err != nil {
		panic(err)
	}
}

func initDB() error {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DB_USER, setting.DB_PASS, setting.DB_HOST, setting.DB_PORT, setting.DB_BASE, setting.DB_PARM)

	if db, err = gorm.Open("mysql", dbConn); err != nil {
		return err
	}

	//update  tabel
	if err := db.AutoMigrate(
		&User{},
		&Session{},
		&Page{},
		&Comment{},
	).Error; err != nil {
		return err
	}
	return nil
}
