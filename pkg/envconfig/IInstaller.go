package envconfig

import (
	"fmt"
	slCC "kallaur.ru/libs/abbyyservice/pkg/colorConsole"
)

// Логика работы:
// 1. Через интерфейс GetDefaultConfigValues() собираем словарь со всеми нужными ключами
// 2. Создаем .env файл со значениями по умолчанию
// 3. Приоритет имеет файл, если в нем оставить пустое значение, будет предпринята попытка считать из памяти
//    для динамичности
func Install(defaults []EnvConfig, path string, env string) error {
	mergedConfig := mergeAllConfig(defaults)
	err := validateInstallerParams(path, env)
	if err != nil {
		return err
	}
	envFilename := fmt.Sprintf("%s/.env.%s", path, env)
	err = mergedConfig.Write(envFilename)
	if err != nil {
		return err
	}
	console := slCC.InstanceColor(slCC.ColorTypeInfoEffect)
	console.Printf("Config file %s is installed\n", envFilename)
	return nil
}
