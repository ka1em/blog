package model

import (
	"fmt"

	"blog/common/setting"

	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	SchemaDecoder *schema.Decoder // schema decoder
	db            *gorm.DB
	redisPool     *redis.Pool
)

const REDIS_MAX_IDLE = 100
const REDIS_MAX_ACTIVE = 100

func init() {
	SchemaDecoder = schema.NewDecoder()
}
func DBInit() {
	connDB()
	connRedisPool()
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

func connRedisPool() {
	address := fmt.Sprintf("%s:%s", setting.RedisHost, setting.RedisPort)
	redisPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", address)
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     REDIS_MAX_IDLE,
		MaxActive:   REDIS_MAX_ACTIVE,
		IdleTimeout: time.Second * 60,
	}
}
