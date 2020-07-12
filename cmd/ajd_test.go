package main

import (
	slParser "kallaur.ru/libs/abbyyservice/pkg/abbyyJsonParser"
	"testing"
)

func TestSavingAjdToRedis(t *testing.T) {
	slParser.SaveAjdRedisTest(t)
}
