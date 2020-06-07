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

func TestShowAddressRealConvert(t *testing.T) {
	err := abbyyJsonParser.ShowRealElements()
	if err != nil {
		t.Fail()
	}
}

func TestViewVerbAddressProperties(t *testing.T) {
	abbyyJsonParser.ViewVerbAddressProperties()
}
