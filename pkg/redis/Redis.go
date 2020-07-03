package redis

import (
	"fmt"
	slRedis "github.com/mediocregopher/radix/v3"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	"strings"
)

func InitRedisPool() *appError.AppError {

	if pool != nil {
		// пул уже авторизован
		return nil
	}

	config := makePoolConfig()
	var err error

	if config.authRequired {
		pool, err = makePoolWithAuth(config)
		if err != nil {
			return appError.CreateAppError(
				appError.ErrorClassConnecting,
				"redis_pool_not_init",
				"init redis pool with auth has problem",
				err)
		}
	} else {
		pool, err = makePool(config)
		//pool, err = makePoolClassic(config)
		if err != nil {
			return appError.CreateAppError(
				appError.ErrorClassConnecting,
				"redis_pool_not_init",
				"init redis pool without auth has problem",
				err)
		}
	}

	return nil
}

// команды redis
func Set(key string, value string) error {
	err := pool.Do(slRedis.Cmd(nil, "SET", key, value))
	return err
}

func Get(key string, value *string) error {
	err := pool.Do(slRedis.Cmd(value, "GET", key))
	return err
}

// ********* Работаем со списками. Хорошо добавлять в начало или конец, и доставать с удалением ***********

func LPush(key string, value string) error {
	err := pool.Do(slRedis.Cmd(nil, "LPUSH", key, value))
	return err
}

func RPush(key string, value string) error {
	err := pool.Do(slRedis.Cmd(nil, "RPUSH", key, value))
	return err
}

func LPop(key string, value *string) error {
	err := pool.Do(slRedis.Cmd(value, "LPOP", key))
	return err
}

func RPop(key string, value *string) error {
	err := pool.Do(slRedis.Cmd(value, "RPOP", key))
	return err
}

func LLen(key string, value *int) error {
	err = pool.Do(slRedis.Cmd(value, "LLEN", key))
	return err
}

func HSet(hash string, field string, value string) error {
	err := pool.Do(slRedis.Cmd(nil, "HSET", hash, field, value))
	return err
}

func HMSetMap(hash string, values map[string]string) error {
	var redisOptionsLine string

	for key, value := range values {
		redisOptionsLine = fmt.Sprintf("%s %s %s", redisOptionsLine, key, value)
	}
	redisOptionsLine = strings.Trim(redisOptionsLine, " ")
	return pool.Do(
		slRedis.Cmd(nil, "HMSET", redisOptionsLine))
}

func HGet(hash string, field string, value *string) error {
	err := pool.Do(slRedis.Cmd(value, "HGET", hash, field))
	return err
}

func HGetAll(hash string, values *[]string) error {
	err := pool.Do(slRedis.Cmd(values, "HGETALL", hash))
	return err
}