package main

import (
	"fmt"
	"github.com/integrii/flaggy"
	"go.uber.org/zap"
	slAbbyySdk "kallaur.ru/libs/abbyyservice/pkg/abbyySdk"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	slErrors "kallaur.ru/libs/abbyyservice/pkg/errors"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	slLogger "kallaur.ru/libs/abbyyservice/pkg/logger"
	slRedis "kallaur.ru/libs/abbyyservice/pkg/redis"
	"os"
	"strings"
)

// Так как это один пакет с инсталлером. Переменные объявлены там
func init() {
	flaggy.SetName("Config installer")
	flaggy.SetDescription("You might install config with default values")

	flaggy.String(&ConfigPath, "p", "path", "A variable just for path config!")
	flaggy.String(&Env, "e", "env", "A variable just for environment type!")

	flaggy.SetVersion(Version)
	flaggy.Parse()

	// если передан параметр пути изменяем на него
	if len(ConfigPath) < 1 {
		err := keeper.GetDefaultConfigPath(&ConfigPath)
		if err != nil {
			appErr := appError.CreateAppError(
				appError.ErrorClassBootstrap,
				"not_found_default_dir",
				"default dir not found",
				err)
			fmt.Printf(appErr.ErrorFmt())
			os.Exit(1)
		}
	}
	// Инициализируем bootstrap компоненты.
	slConfig.InitConfigEnv(ConfigPath)
	slLogger.InitLogger()
	slErrors.InitErrorHandler()
}

func main() {
	// подключаемся к редис
	err := slRedis.InitRedisPool()
	if err != nil {
		eh := slErrors.InstanceErrorHandler()
		eh.Report(err)
		os.Exit(1)
	}
	dispatchEvent("Begin process info about words downloading")
	runJob()
	dispatchEvent("End process")
	dispatchOnClose()
}

func runJob() {
	var word string
	var idx, maxWords = 0, 0

	eh := slErrors.InstanceErrorHandler()
	if maxWords = prepareJob(); maxWords == 0 {
		dispatchEvent("words list is empty")
		dispatchOnClose()
		os.Exit(0)
	}
	for {
		err := slRedis.GetWord(&word)
		if err != nil {
			appErr := appError.CreateAppError(
				appError.ErrorClassDb,
				"not_get_word",
				fmt.Sprintf("word = %s not getting", word),
				err)
			eh.Report(appErr)
			err := slRedis.AddBadWord(word)
			if err != nil {
				appErr := appError.CreateAppError(
					appError.ErrorClassDb,
					"not_add_value_in_list",
					fmt.Sprintf("word = %s not added in bad word list", word),
					err)
				eh.Report(appErr)
			}
			break
		}
		response := slAbbyySdk.WordFormsRu(word, false)
		if response.WithError {
			for _, err := range response.Errors {
				eh.Report(err)
			}
			_ = slRedis.AddBadWord(word)
			idx++
			continue
		}
		err = slRedis.AddJsonData(response.Body, strings.ToLower(word))
		strlen := len(response.Body)
		if err != nil {
			appErr := appError.CreateAppError(appError.ErrorClassDb,
				"not_add_value_in_hash",
				fmt.Sprintf("Json data for lexem = %s not added in json list", word),
				err)
			eh.Report(appErr)
			break
		}
		registerOperation(word, strlen)
		idx++
		if idx >= maxWords {
			// контроль на максимальное количество слов из списка
			break
		}
	}
}

// если 0 значит нет работы :)
func prepareJob() int {
	var maxWords int
	maxWords = slRedis.GetLlen(keeper.KeyAbbyyWordList)
	if maxWords == 0 {
		logger := slLogger.InstanceLogger()
		logger.Info("not found job.")
		dispatchOnClose()
		return 0
	}
	return maxWords
}

func registerOperation(word string, strlen int) {
	logger := slLogger.InstanceLogger()
	logger.Info("Lexem data received",
		zap.String("Lexem", word),
		zap.Int("Body len", strlen))
	dispatchOnClose()
}

func dispatchEvent(message string) {
	logger := slLogger.InstanceLogger()
	logger.Info(
		message,
		zap.String("point", "main"))
}

func dispatchOnClose() {
	logger := slLogger.InstanceLogger()
	_ = logger.Sync()
}
