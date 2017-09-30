package model

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
)

var tdb *gorm.DB

func TestMain(m *testing.M) {
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		"lgwd", "queeheeChiegiusheeD4", "47.93.11.105", "3306", "lgwd", "charset=utf8mb4&parseTime=True&loc=Local")

	var err error
	if tdb, err = gorm.Open("mysql", dbConn); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
func TestCommentAddCommentTrascation(t *testing.T) {
	tx := tdb.Begin()

	c1 := Comment{ID: 4}
	if err := tx.Create(&c1).Error; err != nil {
		tx.Rollback()
		t.Fatal(err.Error())
	}

	c2 := Comment{ID: 4}
	if err := tx.Create(&c2).Error; err != nil {
		tx.Rollback()
		t.Fatal(err.Error())
	}

	err := tx.Commit().Error
	if err != nil {
		t.Fatalf("test trascation error")
		return
	}
}
