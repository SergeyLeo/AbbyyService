package main

import (
	slParser "kallaur.ru/libs/abbyyservice/pkg/abbyyJsonParser"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	slRedis "kallaur.ru/libs/abbyyservice/pkg/redis"
	"os"
)

// Требует при компиляции CommonVars
func main() {
	var json string
	//var all []string
	var lexem = "прибыль"
	var lang = 1049
	//var lexem = "мела"

	err := keeper.GetDefaultConfigPath(&ConfigPath)
	if err != nil {
		os.Exit(1)
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
