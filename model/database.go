package model

import (
	"blog/common/setting"
	"blog/common/zlog"

	"fmt"
	"time"

	"os"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/gorilla/schema"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	SchemaDecoder *schema.Decoder // schema decoder
	//db            *gorm.DB
	xdb       *xorm.Engine
	redisPool *redis.Pool
)

const REDIS_MAX_IDLE = 100
const REDIS_MAX_ACTIVE = 100

func init() {
	SchemaDecoder = schema.NewDecoder()
}

func DBInit() {
	connXDB()
	connRedisPool()
}

func connXDB() {
	var err error
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)

	xdb, err = xormEngine(address)
	if err != nil {
		panic(err.Error())
	}

	if setting.RunMode == setting.DevMode || setting.RunMode == setting.TestMode {
		if setting.SQLLogPath != "stdout" {
			f, err := os.Create(setting.SQLLogPath)
			if err != nil {
				panic(err.Error())
			}
			xdb.SetLogger(xorm.NewSimpleLogger(f))
		}
		xdb.ShowSQL(true)
		xdb.Logger().SetLevel(core.LOG_DEBUG)
	}
	xdb.SetMapper(core.GonicMapper{})
	xdb.Sync2(
		new(User),
		new(Page),
	)

	zlog.ZapLog.Debug("xorm connect mysql ok")
}

func xormEngine(address string) (*xorm.Engine, error) {
	return xorm.NewEngine("mysql", address)
}

func connRedisPool() {
	redisPool = getRedisPool(setting.RedisHost, setting.RedisPort)
	conn := redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		panic(err.Error())
	}
	zlog.ZapLog.Debug("connect redis ok")
}

func getRedisPool(host, port string) *redis.Pool {
	return &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host+":"+port)
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := redis.String(c.Do("PING"))
			return err
		},
		MaxIdle:     REDIS_MAX_IDLE,
		MaxActive:   REDIS_MAX_ACTIVE,
		IdleTimeout: time.Second * 60,
	}
}
