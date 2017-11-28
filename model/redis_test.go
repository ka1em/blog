package model

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestRedisDao_Set(t *testing.T) {
	var r RedisDao
	r.Pool = getRedisPool("127.0.0.1", "6379")

	_, err := r.Set("daotest", 1)
	if err != nil {
		t.Fatal(err.Error())
	}

	v, err := r.GetInt("daotest")
	if err != nil {
		t.Fatal(err.Error())
	}

	if v != 1 {
		t.Fatal(err.Error())
	}
	t.Log(v)

	v, err = r.GetInt("daotest22")
	if err == nil {
		t.Fatal(err.Error())
	}
	t.Log(v, err.Error())
}

func TestRedisDao_Incr(t *testing.T) {
	var r RedisDao
	r.Pool = getRedisPool("127.0.0.1", "6379")

	v, err := r.SetEx("dsss", 2, 999)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(v)

	v2, err := r.Incr("dsss")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(v2)
}

func TestRedisDao_HSet(t *testing.T) {
	var r RedisDao
	r.Pool = getRedisPool("127.0.0.1", "6379")

	user := User{
		ID:    90909090,
		Name:  "wang",
		Email: "wangkaimin@gma.com",
		Role:  "user",
	}

	//j, err := jsoniter.Marshal(user)
	//if err != nil {
	//	t.Fatal(err.Error())
	//}

	v, err := r.HSet("USER", user.ID, user)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(v)

	b, err := r.HGet("USER", user.ID)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(b))
}

func TestRedisDao_HGetAll(t *testing.T) {
	var r RedisDao
	r.Pool = getRedisPool("127.0.0.1", "6379")
	type ReUser struct {
		ID   int    `redis:"id"`
		Name string `redis:"name"`
		Age  int    `redis:"age"`
	}
	u := ReUser{
		ID:   1,
		Name: "wang kaimin",
		Age:  20,
	}
	v1, err := r.HMSet("id1", &u)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("v1 %+v", v1)

	v, err := r.HGetAll("id1")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(v)
	var u3 ReUser
	err = redis.ScanStruct(v, &u3)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(u3)
	t.Logf("u3.name %s\n", u3.Name)
}

func TestRedisDao_HGetAll2(t *testing.T) {
	var r RedisDao
	r.Pool = getRedisPool("127.0.0.1", "6379")

	var p1, p2 struct {
		Title  string `redis:"title"`
		Author string `redis:"author"`
		Body   string `redis:"body"`
	}

	p1.Title = "Example"
	p1.Author = "Gary"
	p1.Body = "Hello"

	if _, err := r.Pool.Get().Do("HMSET", redis.Args{}.Add("id1").AddFlat(&p1)...); err != nil {
		panic(err)
	}

	v, err := redis.Values(r.Pool.Get().Do("HGETALL", "id1"))
	if err != nil {
		panic(err)
	}

	if err := redis.ScanStruct(v, &p2); err != nil {
		panic(err)
	}
	t.Log(p2)
}
