package main

import (
	"fmt"
	"kallaur.ru/libs/abbyyservice/pkg/abbyyJsonParser"
	"strings"
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

func TestRegExpForWord(t *testing.T) {
	err := abbyyJsonParser.RegExpForWords()
	if err != nil {
		t.Fail()
	}
}

func TestBitesOperations(t *testing.T) {
	left := 47
	right := 127
	fmt.Printf("Left  = %08b\nRight = %08b\n", left, right)
	fmt.Printf("L & R = %08b\n", left&right)
	fmt.Printf("R & L = %08b\n", right&left)
	fmt.Printf("%s\n", "================= ==============")
	left = 47
	right = 17
	fmt.Printf("Left  = %08b\nRight = %08b\n", left, right)
	fmt.Printf("L & R = %08b\n", left&right)
	fmt.Printf("R & L = %08b\n", right&left)
	fmt.Printf("%s\n", "================= ==============")
}

func TestMapToString(t *testing.T) {
	testmap := make(map[string]string, 10)
	var line string

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		testmap[key] = value
	}

	for key, value := range testmap {
		line = fmt.Sprintf("%s %s %s", line, key, value)
	}
	line = strings.Trim(line, " ")
	fmt.Printf(line)
}
