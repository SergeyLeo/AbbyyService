package envconfig

import (
	"fmt"
	"github.com/joho/godotenv"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	"os"
	"sort"
	"strconv"
)

type EnvConfig map[string]string

// Записываем конфигурацию в файл
func (config *EnvConfig) Write(filename string) error {
	fileConfig, err := os.Create(filename)
	defer fileConfig.Close()

	if err != nil {
		return appError.CreateAppError(
			appError.ErrorClassBootstrap,
			"config_file_not_created",
			".env file not created",
			err)
	}

	tmpConfig := *config
	keys := make([]string, 0, len(tmpConfig))
	for k := range tmpConfig {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := tmpConfig[key]
		line := fmt.Sprintf("%s=%s\n", key, value)
		_, err := fileConfig.WriteString(line)
		if err != nil {
			return appError.CreateAppError(
				appError.ErrorClassConfiguring,
				"line_not_write",
				"Line not write to file",
				err)
		}
	}
	return nil
}

// 0 - каталог для поиска файла конфигурации
// 1 - тип конфигурации:
//		- dev = .env.dev = имя файла, для среды develop
//		- prod = .env.prod = имя файла, для среды prod
func InitConfigEnv(params ...string) {
	dir := "."
	if len(params) >= 1 {
		dir = params[0]
	}
	envFilename := fmt.Sprintf("%s/.env", dir)
	fromFileMap, err := godotenv.Read(envFilename)
	if err != nil {
		appErr := appError.CreateAppError(
			appError.ErrorClassConfiguring,
			"error_read",
			"config file not read",
			err)
		fmt.Printf(appErr.ErrorFmt())
		os.Exit(1)
	}

	tmp := EnvConfig(fromFileMap)
	config = &tmp
	// закачиваем значение среды из файла. Лучше устанавливать эту переменную перед запуском программы
	value, err := GetValue(keeper.GetCommonKeyNameEnv())
	if err == nil {
		value = keeper.EnvProd
	}

	err = os.Setenv(keeper.GetCommonKeyNameEnv(), value)
	if err != nil {
		appErr := appError.CreateAppError(
			appError.ErrorClassConfiguring,
			"error_set_env",
			"os.Setenv() ended with err",
			err)
		fmt.Printf(appErr.ErrorFmt())
		os.Exit(1)
	}
}

func GetCurrentConfig() *EnvConfig {
	return config
}

func GetValueAsInt(key string) (int, *appError.AppError) {
	var out int
	var err error
	var appErr *appError.AppError

	value, ok := getValue(key)
	if !ok {
		out = 0
		appErr = appError.CreateAppError(
			appError.ErrorClassConfiguring,
			"not_found_key_in_map",
			fmt.Sprintf("value not found. Key = %s", key),
			nil)
		return out, appErr
	}

	out, err = strconv.Atoi(value)
	if err != nil {
		appErr = appError.CreateAppError(
			appError.ErrorClassConfiguring,
			"convert_int_value",
			fmt.Sprintf("error convert to int. Key = %s, value = %s", key, value),
			err)
	}

	return out, appErr
}

func GetValueAsBool(key string) (bool, *appError.AppError) {
	var out = false
	var err error
	var appErr *appError.AppError

	value, ok := getValue(key)
	if ok {
		out, err = strconv.ParseBool(value)
		if err != nil {
			appErr = appError.CreateAppError(
				appError.ErrorClassConfiguring,
				"converted_bool",
				fmt.Sprintf("error convert bool param: key = %s, value = %s", key, value),
				err)
		}
	} else {
		out = false
		appErr = appError.CreateAppError(
			appError.ErrorClassConfiguring,
			"not_found_key_in_map",
			fmt.Sprintf("value not found. Key = %s", key),
			nil)
	}

	return out, appErr
}

// Получить строковое представление ключа.
// По умолчанию
func GetValue(key string) (string, *appError.AppError) {
	var err *appError.AppError

	out, ok := getValue(key)
	if !ok {
		err = appError.CreateAppError(appError.ErrorClassConfiguring,
			"error_read_value",
			fmt.Sprintf("error reading for key = %s", key),
			nil)
	}

	return out, err
}
