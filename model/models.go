package model

import (
	"fmt"

	"blog/common/setting"

	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	SchemaDecoder *schema.Decoder // schema decoder
	db            *gorm.DB
)

func init() {
	SchemaDecoder = schema.NewDecoder()
}
func DBInit() {
	err := connDB()
	if err != nil {
		panic(err)
	}
}

func connDB() error {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)

	if db, err = gorm.Open("mysql", dbConn); err != nil {
		panic(err.Error())
	}

	if err := db.AutoMigrate(
		&User{},
		&Session{},
		&Page{},
		&Comment{},
	).Error; err != nil {
		panic(err.Error())
	}
	return nil
}
