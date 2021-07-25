package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func NewRedisConn() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     viper.GetInt("redis.max-idle"),
		IdleTimeout: viper.GetDuration("redis.timeout"),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", viper.GetString("redis.host"))
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", viper.GetString("redis.password")); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
	}
}
