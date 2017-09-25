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
}
func DBInit() {
	err := updateDB()
	if err != nil {
		panic(err)
	}
}

func updateDB() error {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)

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
