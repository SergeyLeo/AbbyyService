package redis

import (
	"fmt"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"strconv"
)

const (
	keyRedisPoolAuth    string = "AS_REDIS_POOL_AUTH"
	keyRedisPoolReqAuth string = "AS_REDIS_POOL_REQUIRED_AUTH"
	keyRedisPoolBase    string = "AS_REDIS_POOL_BASE"
	keyRedisPoolHost    string = "AS_REDIS_POOL_HOST"
	keyRedisPoolPort    string = "AS_REDIS_POOL_PORT"
	keyRedisPoolCount   string = "AS_REDIS_POOL_COUNT"
	keyRedisPoolNetwork string = "AS_REDIS_POOL_NETWORK"
	keyRedisPoolTimeout string = "AS_REDIS_POOL_TIMEOUT"

	maxDB = 12
)

func GetDefaultConfigValues() slConfig.EnvConfig {
	return slConfig.EnvConfig{
		keyRedisPoolAuth:    "redis",
		keyRedisPoolBase:    "0",
		keyRedisPoolCount:   "2", // количество подключений в пуле
		keyRedisPoolHost:    "127.0.0.1",
		keyRedisPoolPort:    "6379",
		keyRedisPoolNetwork: "tcp",
		keyRedisPoolReqAuth: strconv.FormatBool(false),
		keyRedisPoolTimeout: "1",
	}
}

func GetConfigKeys() []string {
	config := GetDefaultConfigValues()
	keys := make([]string, 0, len(config))
	for key := range config {
		keys = append(keys, key)
	}
	return keys
}

func ValidateKeys(envStorage *slConfig.EnvConfig, keys []string) bool {
	for _, key := range keys {
		_, ok := (*envStorage)[key]
		if !ok {
			appErr := appError.CreateAppError(
				appError.ErrorClassBootstrap,
				"key_not_found",
				fmt.Sprintf("Key=%s not found", key),
				nil)
			if appError.HasErrors() {
				appError.InitErrorsArray(appErr)
			} else {
				appError.AddError(appErr)
			}
		}
	}
	if appError.HasErrors() {
		return false
	}

	return true
}
