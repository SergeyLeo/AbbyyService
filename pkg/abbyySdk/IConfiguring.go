package abbyySdk

import (
	"kallaur.ru/libs/abbyyservice/pkg/envconfig"
)

const (
	keyAbbyyWordList     string = "AS_ABBYY_WORD_LIST_KEY"
	keyAbbyyJsonList     string = "AS_ABBYY_JSON_LIST_KEY"
	keyAbbyyBadWordsList string = "AS_ABBYY_BAD_WORDS_LIST_KEY"
	keyAbbyyApiDay       string = "AS_ABBYY_API_DAY_KEY"
)

func GetDefaultConfigValues() envconfig.EnvConfig {
	return envconfig.EnvConfig{
		keyAbbyyWordList:     "abbyy:words:list",
		keyAbbyyJsonList:     "abbyy:json:list",
		keyAbbyyBadWordsList: "abbyy:bad:words:list",
		keyAbbyyApiDay:       "abbyy:api:day:key",
	}
}
