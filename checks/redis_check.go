package checks

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type RedisCheck struct {
	Pool *redis.Pool
}

func (r RedisCheck) Pass() bool {
	s, err := redis.String(r.Pool.Get().Do("PING"))
	fmt.Printf("PING Response = %s\n", s)
	return err == nil
}

func (r RedisCheck) Name() string {
	return "redis"
}
