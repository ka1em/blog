package model

import (
	"blog/common/setting"
	"blog/common/zlog"

	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/gorilla/schema"
	"github.com/jinzhu/gorm"
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

const (
	REDIS_KEY_PREFIX  = "BLOG:"
	REDIS_KEY_SESSION = REDIS_KEY_PREFIX + "SESSION:"
	REDIS_KEY_USER    = REDIS_KEY_PREFIX + "USER:"
	REDIS_KEY_PAGE    = REDIS_KEY_PREFIX + "PAGE:"
)

func init() {
	SchemaDecoder = schema.NewDecoder()
}

func DBInit() {
	connXDB()
	connRedisPool()
}

func connDB() {
	//var err error
	//address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
	//	setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)
	//
	//db, err = openMysql(address)
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//if err := db.AutoMigrate(
	//	&User{},
	//	&Session{},
	//	&Page{},
	//	&Comment{},
	//).Error; err != nil {
	//	panic(err.Error())
	//}
	//zlog.ZapLog.Debug("gorm connect mysql ok")
}

func connXDB() {
	var err error
	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DBUser, setting.DBPass, setting.DBHost, setting.DBPort, setting.DBBase, setting.DBParm)

	xdb, err = xormEngine(address)
	if err != nil {
		panic(err.Error())
	}

	xdb.ShowSQL(true)
	xdb.Logger().SetLevel(core.LOG_DEBUG)
	xdb.SetMapper(core.GonicMapper{})
	xdb.Sync2(
		new(User),
		new(Page),
	)

	//f, err := os.Create("sql.log")
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
	//xdb.SetLogger(xorm.NewSimpleLogger(f))

	zlog.ZapLog.Debug("xorm connect mysql ok")
}

func openMysql(address string) (*gorm.DB, error) {
	return gorm.Open("mysql", address)
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
