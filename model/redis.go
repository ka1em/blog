package model

import "github.com/garyburd/redigo/redis"

// RedisDao redis存储
type RedisDao struct {
	Pool *redis.Pool
}

func (r *RedisDao) Set(key, value interface{}) (string, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.String(c.Do("SET", key, value))
}

func (r *RedisDao) GetString(key interface{}) (string, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.String(c.Do("GET", key))
}

func (r *RedisDao) GetInt(key interface{}) (int, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int(c.Do("GET", key))
}

func (r *RedisDao) GetInt64(key interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("GET", key))
}

// SetNx
// 1 set 0 not set
func (r *RedisDao) SetNx(key, value interface{}) (int, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int(c.Do("SETNX", key, value))
}

func (r *RedisDao) SetEx(key, second, value interface{}) (string, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.String(c.Do("SETEX", key, second, value))
}

func (r *RedisDao) GetSetInt(key, value interface{}) (int, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int(c.Do("GETSET", key, value))
}

func (r *RedisDao) GetSetInt64(key, value interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("GETSET", key, value))
}

func (r *RedisDao) GetSetString(key, value interface{}) (string, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.String(c.Do("GETSET", key, value))
}

// Append 追加字符串
func (r *RedisDao) Append(key, value interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("APPEND", key, value))
}

// Decr 递减
func (r *RedisDao) Decr(key interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("DECR", key))
}

// Incr 递增
func (r *RedisDao) Incr(key interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("INCR", key))
}

func (r *RedisDao) SetRange(key, offset, value interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("SETRANGE", key, offset, value))
}

func (r *RedisDao) GetRange(key, start, end interface{}) (string, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.String(c.Do("GETRANGE", key, start, end))
}

func (r *RedisDao) HSet(key, field, value interface{}) (int64, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int64(c.Do("HSET", key, field, value))
}

func (r *RedisDao) HMSet(key, value interface{}) (string, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.String(c.Do("HMSET", redis.Args{}.Add(key).AddFlat(value)...))
}

func (r *RedisDao) HGet(key, field interface{}) ([]byte, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Bytes(c.Do("HGET", key, field))
}

func (r *RedisDao) HGetAll(key interface{}) ([]interface{}, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Values(c.Do("HGETALL", key))
}

func (r *RedisDao) Del(key interface{}) (int, error) {
	if r.Pool == nil {
		r.Pool = redisPool
	}
	c := r.Pool.Get()
	defer c.Close()
	return redis.Int(c.Do("DEL", key))
}
