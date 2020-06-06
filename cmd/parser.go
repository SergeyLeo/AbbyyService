package main

import (
	"github.com/integrii/flaggy"
	slParser "kallaur.ru/libs/abbyyservice/pkg/abbyyJsonParser"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	slRedis "kallaur.ru/libs/abbyyservice/pkg/redis"
	"os"
)

func init() {
	flaggy.SetName("Abbyy json file parser")
	flaggy.SetDescription(`
Разбираем ранее полученные файлы от сервиса Abbyy.
По-умолчанию рабочая среда prod. Если мы указываем dev, по-умолчанию файл .env будем искать рядом с бинарником.
`)

	flaggy.String(&Env, "e", "env", "A variable just for environment type!")

	flaggy.SetVersion(Version)
	flaggy.Parse()

}

// Требует при компиляции CommonVars
func main() {
	var err error
	var json string
	//var all []string
	//var lexem = "мела"
	var lexem = "прибыль"
	var lang = 1049
	if Env == "dev" {
		ConfigPath = "./bin"
	} else {
		err = keeper.GetDefaultConfigPath(&ConfigPath)
		if err != nil {
			os.Exit(1)
		}
	}
	slConfig.InitConfigEnv(ConfigPath)
	appErr := slRedis.InitRedisPool()
	if appErr != nil {
		os.Exit(1)
	}
	keyRedis, _ := slConfig.GetValue(keeper.KeyAbbyyJsonList)
	err = slRedis.HGet(keyRedis, lexem, &json)
	if err != nil || len(json) < 1 {
		os.Exit(1)
	}
	ajd, err := slParser.MarshalAbbyyJsonData(json, lexem, lang)
	if err != nil {
		os.Exit(1)
	}
	lErr := slParser.FetchWords(ajd)
	if lErr != nil {
		os.Exit(1)
	}
}
