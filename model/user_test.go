package model

import (
	"fmt"
	"testing"
)

func TestUser_Create(t *testing.T) {
	//d, err := openMysql("root:123456@tcp(localhost:3306)/lgwd")
	//if err != nil {
	//	t.Fatal(err.Error())
	//}
	//db = d

	d, err := xormEngine("root:123456@tcp(localhost:3306)/lgwd")
	if err != nil {
		t.Fatal(err.Error())
	}
	xdb = d
	id, err := SF.NextID()
	if err != nil {
		t.Fatal(err.Error())
	}

	u := &User{
		ID:        id,
		Name:      fmt.Sprintf("%d", id),
		XDB:       xdb,
		RedisPool: getRedisPool("127.0.0.1", "6379"),
	}

	err = u.Create()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = u.SetCache()
	if err != nil {
		t.Fatal(err.Error())
	}

	u2, err := u.GetCache(id)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(u2)
}
