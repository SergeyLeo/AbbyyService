package redis

import (
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
)

func getAbbyyListKey() string {
	key, err := slConfig.GetValue(keeper.KeyAbbyyWordList)
	if err != nil {
		return ""
	}
	return key
}

func getAbbyyBadListKey() string {
	key, err := slConfig.GetValue(keeper.KeyAbbyyBadWordsList)
	if err != nil {
		return ""
	}
	return key
}

func getAbbyyJsonListKey() string {
	key, err := slConfig.GetValue(keeper.KeyAbbyyJsonList)
	if err != nil {
		return ""
	}
	return key
}
