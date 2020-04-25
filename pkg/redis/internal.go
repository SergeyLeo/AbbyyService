package redis

import (
	"fmt"
	slRedis "github.com/mediocregopher/radix/v3"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"os"
	"time"
)

var pool *slRedis.Pool
var err error

type poolConfig struct {
	auth         string
	authRequired bool
	count        int
	port         int
	base         int
	host         string
	network      string
	timeout      int
}

func (pc *poolConfig) dsn() string {
	return fmt.Sprintf("%s:%d", pc.host, pc.port)
}

/// Создаем конфигурацию для пула
func makePoolConfig() *poolConfig {
	isValid := ValidateKeys(slConfig.GetCurrentConfig(), GetConfigKeys())
	if !isValid {
		fmt.Printf(appError.FlushErrors())
		os.Exit(1)
	}

	config := new(poolConfig)
	config.authRequired, _ = slConfig.GetValueAsBool(keyRedisPoolReqAuth)
	config.auth, _ = slConfig.GetValue(keyRedisPoolAuth)
	config.base, _ = slConfig.GetValueAsInt(keyRedisPoolBase)
	config.host, _ = slConfig.GetValue(keyRedisPoolHost)
	config.port, _ = slConfig.GetValueAsInt(keyRedisPoolPort)
	config.count, _ = slConfig.GetValueAsInt(keyRedisPoolCount)
	config.network, _ = slConfig.GetValue(keyRedisPoolNetwork)
	// контроль максимального номера бызы
	if config.base > maxDB {
		config.base = 0
	}

	return config
}

func makePoolWithAuth(config *poolConfig) (*slRedis.Pool, error) {
	var optTimeout, optAuthPass, optConnectBase slRedis.DialOpt

	if config.timeout > 0 {
		if config.timeout < 4 {
			optTimeout = slRedis.DialTimeout(time.Minute * time.Duration(config.timeout))
		} else {
			optTimeout = slRedis.DialTimeout(time.Minute * 3)
		}
	} else {
		optTimeout = slRedis.DialTimeout(time.Minute * 3)
	}

	optAuthPass = slRedis.DialAuthPass(config.auth)
	optConnectBase = slRedis.DialSelectDB(config.base)

	funcConnect := func(network, addr string) (slRedis.Conn, error) {
		return slRedis.Dial(network, addr, optTimeout, optAuthPass, optConnectBase)
	}
	return slRedis.NewPool(
		config.network,
		config.dsn(),
		config.count,
		slRedis.PoolConnFunc(funcConnect))
}

func makePool(config *poolConfig) (*slRedis.Pool, error) {
	var optConnectBase slRedis.DialOpt

	optConnectBase = slRedis.DialSelectDB(config.base)
	funcConnect := func(network, addr string) (slRedis.Conn, error) {
		return slRedis.Dial(network, addr, optConnectBase)
	}

	return slRedis.NewPool(
		config.network,
		config.dsn(),
		config.count,
		slRedis.PoolConnFunc(funcConnect))
}

func makePoolClassic(config *poolConfig) (*slRedis.Pool, error) {
	return slRedis.NewPool(
		config.network,
		config.dsn(),
		config.count)
}
