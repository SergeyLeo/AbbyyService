package redis

import slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"

// Расширение функционала под действующее приложение

func GetWord(word *string) error {
	key := getAbbyyListKey()
	return LPop(key, word)
}

func AddWord(word string) error {
	key := getAbbyyListKey()
	return RPush(key, word)
}

func GetBadWord(word *string) error {
	key := getAbbyyBadListKey()
	return RPop(key, word)
}

func AddBadWord(word string) error {
	key := getAbbyyBadListKey()
	return LPush(key, word)
}

func GetLlen(keyParam string) int {
	var ll int
	keyValue, _ := slConfig.GetValue(keyParam)
	err := LLen(keyValue, &ll)
	if err != nil {
		return 0
	}
	return ll
}

func AddJsonData(json string, lexem string) error {
	key := getAbbyyJsonListKey()
	return HSet(key, lexem, json)
}

func GetAjdUuid() (string, error) {
	key := getAbbyyJsonDataListKey()
	var value string

	err := LPop(key, &value)
	return value, err
}
