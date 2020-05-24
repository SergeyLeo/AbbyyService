package main

import (
	"kallaur.ru/libs/abbyyservice/pkg/abbyyJsonParser"
	"testing"
)

func TestShowTableTypeTokenAddress(t *testing.T) {
	abbyyJsonParser.ShowTokenAddressType()
}

func TestShowAddressConvert(t *testing.T) {
	abbyyJsonParser.ShowConvertAddresses()
}
