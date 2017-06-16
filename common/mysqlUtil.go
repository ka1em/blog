package common

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	DBHost = "127.0.0.1"
	DBPort = ":3306"
	DBUser = "cmsuser"
	DBPass = "cmspass"
	DBBase = "cms"
)

var database *sql.DB

func GetDB() *sql.DB {
	var db *sql.DB
	var err error

	dbConn := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", DBUser, DBPass, DBHost, DBPort, DBBase)
	if database == nil {
		db, err = sql.Open("mysql", dbConn)
		if err != nil {
			panic(err)
			return nil
		}
		database = db
	}

	return database
}

func init() {
	database = GetDB()
	//database.SetMaxOpenConns(100)
	//database.SetMaxIdleConns(10)
	//database.SetConnMaxLifetime(1 * time.Second)

	go func() {
		for {
			time.Sleep(1 * time.Hour)
			err := database.Ping()
			if err != nil {
				log.Println(err.Error())
			}
			log.Println("ping database")
		}
	}()
}