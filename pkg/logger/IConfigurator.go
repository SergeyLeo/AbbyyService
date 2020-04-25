// конфигурирование
// 1 - константы для чтения ключей
// 		1.1. - возможно их нужно перенести сразу в интерфейс GetDefaultConfigValue()
// 2 - Интерфейс для формирования дефолтного файла конфигурации GetDefaultConfigValue()
// 3 - Фомирование конфигурации логгера
package logger

import (
	"fmt"
	"go.uber.org/zap"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	"os"
)

const (
	keyLJCPath       string = "AS_LJC_PATH"
	keyLJCFilename   string = "AS_LJC_FILENAME"
	keyLJCMaxSize    string = "AS_LJC_MAXSIZE"
	keyLJCMaxAge     string = "AS_LJC_MAXAGE"
	keyLJCMaxBackups string = "AS_LJC_MAXBACKUPS"
)

var configLJC *lamberJackConfig

func GetDefaultConfigValues() slConfig.EnvConfig {
	defaultPath := keeper.GetDefaultLogPath()
	return slConfig.EnvConfig{
		keyLJCPath:       defaultPath,
		keyLJCFilename:   "aservice.log",
		keyLJCMaxSize:    "8", // размер указан в мегабайтах
		keyLJCMaxAge:     "8", // в днях
		keyLJCMaxBackups: "3",
	}
}

// Уместно вызывать из init() пакета main
func InitLogger() {
	// создаем конфиг логгера
	err := makeConfigLJC()
	if err != nil {
		appErr := appError.CreateAppError(
			appError.ErrorClassBootstrap,
			"config_not_create",
			"lamberjack logger config not create",
			err)
		fmt.Printf(appErr.Error())
		os.Exit(1)
	}

	core := makeCore()
	logger = zap.New(core)
}

func makeConfigLJC() error {
	var err *appError.AppError
	var path, filename string
	config := new(lamberJackConfig)
	path, err = slConfig.GetValue(keyLJCPath)
	filename, err = slConfig.GetValue(keyLJCFilename)

	config.filename = fmt.Sprintf("%s/%s", path, filename)
	config.maxAge, err = slConfig.GetValueAsInt(keyLJCMaxAge)
	config.maxBackups, err = slConfig.GetValueAsInt(keyLJCMaxBackups)
	config.maxSize, err = slConfig.GetValueAsInt(keyLJCMaxSize)
	if err != nil {
		return err
	}
	configLJC = config
	return nil
}
