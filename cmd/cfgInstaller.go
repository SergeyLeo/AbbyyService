// По умолчанию на Raspberry Pi установка производится со следующим вложением каталогов
// opt/kallaur.ru.
// - настройки устанавливаем в каталог /etc/opt/kallaur.ru/xxx, где xxx - имя сервиса
// - каталог для логов /var/log/opt/kallaur.ru/xxx, где xxx - имя сервиса
// - каталог для рабочих файлов сервисов /var/lib/opt/kallaur.ru/xxx, где xxx - имя сервиса
// - каталог для бинарников /opt/kallaur.ru/xxx, где xxx - имя сервиса
package main

import (
	"fmt"
	"github.com/integrii/flaggy"
	slAbbyy "kallaur.ru/libs/abbyyservice/pkg/abbyySdk"
	ae "kallaur.ru/libs/abbyyservice/pkg/appError"
	slCC "kallaur.ru/libs/abbyyservice/pkg/colorConsole"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	slLogger "kallaur.ru/libs/abbyyservice/pkg/logger"
	slRedis "kallaur.ru/libs/abbyyservice/pkg/redis"
	"os"
	"path/filepath"
)

var makeInfra = false

func getPaths() []string {
	return []string{
		keeper.SettingsPath,
		keeper.BinaryPath,
		keeper.LogsPath,
		keeper.LibPath,
	}
}
func init() {
	flaggy.SetName("Config installer")
	flaggy.SetDescription(`

Инсталлируем файл конфигурации по-умолчанию. 
С флагом -i может создать инфраструктуру каталогов:

/etc/opt/kallaur.ru/xxx, где xxx - имя сервиса, для хранения настроек
/opt/kallaur.ru/xxx, где xxx - имя сервиса, для хранения бинарников
/etc/log/kallaur.ru/xxx, где xxx - имя сервиса, для хранения логов
/etc/lib/kallaur.ru/xxx, где xxx - имя сервиса, для хранения рабочих файлов сервиса

Далее инсталлирует файл настроек в /etc/opt/kallaur/ru/xxx
`)

	flaggy.Bool(&makeInfra, "i", "infra", "Build default dir infra!")
	flaggy.String(&ConfigPath, "p", "path", "A variable just for path config!")
	flaggy.String(&Env, "e", "env", "A variable just for environment type!")

	flaggy.SetVersion(Version)
	flaggy.Parse()
}

func main() {
	var err error
	// ключ создать инфраструктуру является ведущим
	if makeInfra {
		result := makeDefaultInfra()
		if !result {
			os.Exit(0)
		}
	}

	if len(ConfigPath) < 1 {
		// значит требуется установить в каталог по умолчанию
		err = keeper.GetDefaultConfigPath(&ConfigPath)
		if err != nil {
			ConfigPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				console := slCC.InstanceColor(slCC.ColorTypeError)
				appErr := ae.CreateAppError(ae.ErrorClassBootstrap,
					"error_default_dir",
					"all defalut dir is ends",
					err)
				_, _ = console.Printf(appErr.ErrorFmt())
				os.Exit(1)
			}

		}
	}

	if len(Env) < 1 {
		Env = "example"
	}

	var configs []slConfig.EnvConfig

	configs = append(configs, slLogger.GetDefaultConfigValues())
	configs = append(configs, slRedis.GetDefaultConfigValues())
	configs = append(configs, slConfig.GetDefaultConfigValues())
	configs = append(configs, slAbbyy.GetDefaultConfigValues())

	err = slConfig.Install(configs, ConfigPath, Env)
	if err != nil {
		console := slCC.InstanceColor(slCC.ColorTypeError)
		appError := ae.CreateAppError(
			ae.ErrorClassConfiguring,
			"default_config_not_install",
			"Default .env file not create",
			err)
		_, _ = console.Printf(appError.ErrorFmt())
		os.Exit(1)
	}
}

// создает инфраструктуру каталогов по умолчанию
func makeDefaultInfra() bool {
	if os.Geteuid() != 0 {
		console := slCC.InstanceColor(slCC.ColorTypeError)
		appError := ae.CreateAppError(
			ae.ErrorClassBootstrap,
			"no_sudo_user",
			"Run program with sudo",
			nil)
		_, _ = console.Printf(appError.ErrorFmt())
		return false
	}
	for _, prefix := range getPaths() {
		path := fmt.Sprintf("%s/%s", prefix, keeper.AppName)
		_, err := os.Stat(path)
		if err != nil {
			err := os.MkdirAll(path, 0755)
			if err == nil {
				// местная функция для смены владельца
				err = chownR(prefix, 1000, 1000)
				if err != nil {
					appError := ae.CreateAppError(
						ae.ErrorClassBootstrap,
						"not_change_owner",
						fmt.Sprintf("not change owner to uid 1000 for path %s", path),
						err)
					ae.AddError(appError)
				}
			} else {
				appError := ae.CreateAppError(
					ae.ErrorClassBootstrap,
					"mkdir_error",
					fmt.Sprintf("not create dir - %s", path),
					err)
				ae.AddError(appError)
			}
		}
	}
	if ae.HasErrors() {
		line := ae.FlushErrors()
		console := slCC.InstanceColor(slCC.ColorTypeError)
		_, _ = console.Printf(line)
		return false
	}
	return true
}

func chownR(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
}
