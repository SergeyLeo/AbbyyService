// Храним сущности необходимые всем пакетам
// Создается для избежания получения перекрестных ссылок
package keeper

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	AppName               string = "abbyy-service"
	AppEnvPrefix          string = "AS"
	AppEnvPrefixDelimiter string = "_"
	EnvProd               string = "prod"
	EnvDev                string = "dev"
	EnvExample            string = "example"
	// ключи для работы с редис
	KeyAbbyyWordList     string = "AS_ABBYY_WORD_LIST_KEY"
	KeyAbbyyJsonList     string = "AS_ABBYY_JSON_LIST_KEY"
	KeyAbbyyBadWordsList string = "AS_ABBYY_BAD_WORDS_LIST_KEY"
	KeyAbbyyApiDay       string = "AS_ABBYY_API_DAY_KEY"
	KeyAJDList           string = "AS_AJD_LIST_IN_REDIS"

	// каталоги инфраструктуры приложения
	SettingsPath string = "/etc/opt/kallaur.ru"
	BinaryPath   string = "/opt/kallaur.ru"
	LogsPath     string = "/var/log/opt/kallaur.ru"
	LibPath      string = "/var/lib/opt/kallaur.ru"
)

func GetCommonKeyNameEnv() string {
	return fmt.Sprintf("%s%s%s", AppEnvPrefix, AppEnvPrefixDelimiter, "ENV")
}

// Получаем каталог по умолчанию исходя из существующих вариантов
func GetDefaultConfigPath(path *string) error {
	var err error
	pathDef := fmt.Sprintf("%s/%s", SettingsPath, AppName)
	_, err = os.Stat(pathDef)
	if err != nil {
		pathDef, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}

	}
	*path = pathDef
	return nil
}

// Получаем каталог по умолчанию исходя из существующих вариантов
func GetDefaultLogPath() string {
	return fmt.Sprintf("%s/%s", LogsPath, AppName)
}
