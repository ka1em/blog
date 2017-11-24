package model

import (
	"fmt"

	"blog/common/setting"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	SchemaDecoder *schema.Decoder // schema decoder
	db            *gorm.DB
	redisConn     redis.Conn
)

func init() {
	SchemaDecoder = schema.NewDecoder()
}
func DBInit() {
	connDB()
	connRedis()
}

func connDB() {
	var err error

	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)

	if db, err = gorm.Open("mysql", address); err != nil {
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
}

func connRedis() {
	var err error
	address := fmt.Sprintf("%s:%s", setting.RedisHost, setting.RedisPort)

	redisConn, err = redis.Dial("tcp", address)
	if err != nil {
		panic(err.Error())
	}
}
