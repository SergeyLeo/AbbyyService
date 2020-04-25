package envconfig

import (
	"fmt"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	"os"
)

// Собираем несколько конфигураций в одну
func mergeAllConfig(defaults []EnvConfig) EnvConfig {
	mergedConfig := map[string]string{}
	for _, mapItem := range defaults {
		for key, value := range mapItem {
			mergedConfig[key] = value
		}
	}

	return mergedConfig
}

// валидируем существование указанного каталога
// валидируем значение переменной среды
func validateInstallerParams(path string, env string) error {
	if env != keeper.EnvProd && env != keeper.EnvDev && env != keeper.EnvExample {
		return fmt.Errorf("env value not in list:\n- %s\n- %s", keeper.EnvProd, keeper.EnvDev)
	}
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	return nil
}
