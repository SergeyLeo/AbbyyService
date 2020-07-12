package main

import (
	slParser "kallaur.ru/libs/abbyyservice/pkg/abbyyJsonParser"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"testing"
)

func main() {
	t := &testing.T{}
	slConfig.InitConfigEnv("./bin")
	slParser.SaveAjdRedisTest(t)
}
