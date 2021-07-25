package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

type SetDataNoExpireRedisFn func(key string, value interface{}) error

func NewSetDataNoExpireRedisFn(pool *redis.Pool) SetDataNoExpireRedisFn {
	return func(key string, value interface{}) error {
		conn := pool.Get()
		defer conn.Close()

		_, err := conn.Do("SET", key, value)
		if err != nil {
			return err
		}
		return nil
	}
}

type SetDataWExpireRedisFn func(key string, ttl int, value interface{}) error

func NewSetDataWExpireRedisFn(pool *redis.Pool) SetDataWExpireRedisFn {
	return func(key string, ttl int, value interface{}) error {
		conn := pool.Get()
		defer conn.Close()

		_, err := conn.Do("SETEX", key, ttl, value)
		if err != nil {
			return err
		}
		return nil
	}
}

type SetStructWExpireRedisFn func(key string, ttl int, value interface{}) error

func NewSetStructWExpireRedisFn(pool *redis.Pool) SetStructWExpireRedisFn {
	return func(key string, ttl int, value interface{}) error {
		conn := pool.Get()
		defer conn.Close()

		b, err := json.Marshal(&value)
		if err != nil {
			return err
		}

		if _, err := conn.Do("SETEX", key, ttl, string(b)); err != nil {
			return err
		}
		return nil
	}
}

type GetStringDataRedisFn func(key string) (string, error)

func NewGetStringDataRedisFn(pool *redis.Pool) GetStringDataRedisFn {
	return func(key string) (string, error) {
		conn := pool.Get()
		defer conn.Close()

		data, err := redis.String(conn.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				return "", nil
			} else {
				return "", err
			}
		}
		return data, nil
	}
}

type GetDeleteStringDataRedisFn func(key string) (string, error)

func NewGetDeleteStringDataRedisFn(pool *redis.Pool) GetDeleteStringDataRedisFn {
	return func(key string) (string, error) {
		conn := pool.Get()
		defer conn.Close()

		data, err := redis.String(conn.Do("GETDEL", key))
		if err != nil {
			if err == redis.ErrNil {
				return "", nil
			} else {
				return "", err
			}
		}
		return data, nil
	}
}

type GetDeleteIntDataRedisFn func(key string) (int, error)

func NewGetDeleteIntDataRedisFn(pool *redis.Pool) GetDeleteIntDataRedisFn {
	return func(key string) (int, error) {
		conn := pool.Get()
		defer conn.Close()

		data, err := redis.Int(conn.Do("GETDEL", key))
		if err != nil {
			if err == redis.ErrNil {
				return 0, nil
			} else {
				return 0, err
			}
		}
		return data, nil
	}
}

type GetFloatDataRedisFn func(key string) (float64, error)

func NewGetFloatDataRedisFn(pool *redis.Pool) GetFloatDataRedisFn {
	return func(key string) (float64, error) {
		conn := pool.Get()
		defer conn.Close()

		data, err := redis.Float64(conn.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				return 0, nil
			} else {
				return 0, err
			}
		}
		return data, nil
	}
}

type GetStructDataRedisFn func(key string, dest interface{}) error

func NewGetStructDataRedisFn(pool *redis.Pool) GetStructDataRedisFn {
	return func(key string, dest interface{}) error {
		conn := pool.Get()
		defer conn.Close()

		data, err := redis.String(conn.Do("GET", key))
		if err != nil {
			if err == redis.ErrNil {
				return nil
			} else {
				return err
			}
		}
		return json.Unmarshal([]byte(data), dest)
	}
}
