package appError

import (
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	"os"
)

type appErrorArray struct {
	count       int
	errElements []*AppError
}

var errorsElements *appErrorArray

// Девелоперская среда должна быть четко обозначена ключом
func isDev() bool {
	value := os.Getenv(keeper.GetCommonKeyNameEnv())
	return value == keeper.EnvDev
}
