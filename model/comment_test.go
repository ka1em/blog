package model

import (
	"fmt"
	"os"
	"testing"

	"blog/common/setting"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func TestMain(m *testing.M) {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		setting.DB_USER, setting.DB_PASS, setting.DB_HOST, setting.DB_PORT, setting.DB_BASE, setting.DB_PARM)

	if db, err = gorm.Open("mysql", dbConn); err != nil {
		panic(err.Error())
	}

	os.Exit(m.Run())
}
func TestCommentAddCommentTrascation(t *testing.T) {
	tx := db.Begin()

	c1 := Comment{Id: 4}
	if err := tx.Create(&c1).Error; err != nil {
		tx.Rollback()
		t.Fatalf("%s", err.Error())
		return
	}

	c2 := Comment{Id: 4}
	if err := tx.Create(&c2).Error; err != nil {
		tx.Rollback()
		t.Fatalf(err.Error())
		return
	}

	err := tx.Commit().Error
	if err != nil {
		t.Fatalf("%s", err.Error())
		return
	}
}
