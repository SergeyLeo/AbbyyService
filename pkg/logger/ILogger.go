package logger

import (
	"go.uber.org/zap"
)

// основная точка входа для вызова логгера
func InstanceLogger() *zap.Logger {
	if logger != nil {
		return logger
	}

	InitLogger()

	return logger
}
