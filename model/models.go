package model

import (
	"fmt"

	"blog.ka1em.site/common"
	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB
var SchemaDecoder = schema.NewDecoder()

const (
	DBHost = "47.93.11.105"
	DBPort = "3306"
	DBUser = "lgwd"
	DBPass = "queeheeChiegiusheeD4"
	DBBase = "lgwd"
	DBParm = "charset=utf8mb4&parseTime=True&loc=Local"
)

func init() {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", DBUser, DBPass, DBHost, DBPort, DBBase, DBParm)

	if DB, err = gorm.Open("mysql", dbConn); err != nil {
		common.Suggar.Error(err.Error())
		panic(err.Error())
		return
	}

	//update  tabel
	if err := DB.AutoMigrate(
		&User{},
		&Session{},
		&Page{},
		&Comment{},
	).Error; err != nil {
		common.Suggar.Error(err.Error())
		return
	}
}
