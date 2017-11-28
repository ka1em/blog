package model

import (
	"testing"
)

func TestUser_Create(t *testing.T) {
	d, err := openMysql("root:123456@tcp(localhost:3306)/lgwd")
	if err != nil {
		t.Fatal(err.Error())
	}
	u := User{
		ID:   123,
		Name: "test 1",
		DB:   d,
	}
	u.Create()
	u.SetCache()
	u2, err := u.GetCache(123)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(u2)
}
