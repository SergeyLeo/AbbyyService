package envconfig

import (
	"os"
)

var config *EnvConfig

func getValue(key string) (string, bool) {
	if config == nil {
		return "", false
	}
	tmp := *config
	value, ok := tmp[key]

	// Если пустое значение пробуем достать из памяти
	if !ok {
		value = os.Getenv(key)
		if len(value) > 0 {
			ok = true
		}
	}

	return value, ok
}
