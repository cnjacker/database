package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisURL = "redis://127.0.0.1:6379/0"
	RedisExp = 5 * time.Minute
)

type Redis struct {
	Conn *redis.Client
	Ctx  context.Context
}

func NewRedis() *Redis {
	return NewRedisWithURL(RedisURL)
}

func NewRedisWithURL(url string) *Redis {
	if v, err := redis.ParseURL(url); err == nil {
		return &Redis{
			Conn: redis.NewClient(v),
			Ctx:  context.Background(),
		}
	}

	return nil
}

func (r Redis) Cache(key string, callback func(...any) any, args ...any) any {
	data := r.Get(key)

	if data != nil {
		return data
	}

	data = callback(args...)
	r.Set(key, data)

	return data
}

func (r Redis) Get(key string) any {
	if v, err := r.Conn.Get(r.Ctx, key).Bytes(); err == nil {
		var data any

		if json.Unmarshal(v, &data) == nil {
			return data
		}
	}

	return nil
}

func (r Redis) Set(key string, data any) bool {
	return r.SetWithExp(key, data, RedisExp)
}

func (r Redis) SetWithExp(key string, data any, exp time.Duration) bool {
	if v, err := json.Marshal(data); err == nil {
		if err := r.Conn.Set(r.Ctx, key, v, exp).Err(); err == nil {
			return true
		}
	}

	return false
}

func (r Redis) Delete(keys ...string) int64 {
	if v, err := r.Conn.Del(r.Ctx, keys...).Result(); err == nil {
		return v
	}

	return 0
}

func (r Redis) Scan(key string) []string {
	var keys []string
	var cursor uint64

	for {
		v, cursor, err := r.Conn.Scan(r.Ctx, cursor, key, 100).Result()

		if err != nil {
			break
		}

		keys = append(keys, v...)

		if cursor == 0 {
			break
		}
	}

	return keys
}

func (r Redis) Close() {
	r.Conn.Close()
}
