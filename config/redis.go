package config

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		redisHost := GetConfig().GetString("REDIS_HOST")
		redisPort := GetConfig().GetString("REDIS_PORT")
		return redis.Dial("tcp", redisHost+":"+redisPort)
	},
}

func NewRedis() *redis.Pool {
	return redisPool
}

func NewFakeRedis() *redis.Pool {
	return &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			redisHost := GetConfig().GetString("REDIS_HOST") + "dsada"
			redisPort := GetConfig().GetString("REDIS_PORT") + "12"
			return redis.Dial("tcp", redisHost+":"+redisPort)
		},
	}
}

var workEnqueuer *work.Enqueuer

func GetEnqueuer() *work.Enqueuer {
	if workEnqueuer == nil {
		workEnqueuer = work.NewEnqueuer(GetConfig().GetString("REDIS_QUEUE"), redisPool)
	}
	return workEnqueuer
}
