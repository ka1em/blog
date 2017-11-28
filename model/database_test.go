package model

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestCreateSessionV2(t *testing.T) {

}

func TestRedis(t *testing.T) {
	p := getRedisPool("127.0.0.1", "6379")
	c := p.Get()
	defer c.Close()

	v, err := c.Do("SET", "key1", "ABC")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(v)

	v2, err := c.Do("GET", "key1")
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(v2)

	v3, err := redis.String(c.Do("GET", "key1"))
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(v3)

	v4, err := redis.String(c.Do("MSET", "key2", "hahh", "key3", "hehe"))
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(v4)

	v5, err := redis.ByteSlices(c.Do("MGET", "key2", "key3"))
	if err != nil {
		t.Error(err.Error())
	}

	for _, n := range v5 {
		t.Log(string(n))
	}

	t.Log(v5)
}
