package model

import (
	"fmt"

	"blog.ka1em.site/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var database *gorm.DB

const (
	DBHost = "47.93.11.105"
	DBPort = "3306"
	DBUser = "lgwd"
	DBPass = "queeheeChiegiusheeD4"
	DBBase = "lgwd"
	DBParm = "charset=utf8mb4&parseTime=True&loc=Local"
)

func GetDB() *gorm.DB {
	var err error
	if database == nil {
		database, err = connDatabase()
		if err != nil {
			common.Suggar.Error(err.Error())
			return nil
		}
	}
	return database
}

func connDatabase() (*gorm.DB, error) {
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", DBUser, DBPass, DBHost, DBPort, DBBase, DBParm)

	if database, err = gorm.Open("mysql", dbConn); err != nil {
		common.Suggar.Error(err.Error())
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
		common.Suggar.Error(err.Error())
		return nil, err
	}

	return database, nil
}
