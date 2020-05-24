package abbyyJsonParser

import (
	"fmt"
	slCC "kallaur.ru/libs/abbyyservice/pkg/colorConsole"
)

// интерфейс для тестирования функций файла tokinizerRu

func ShowTokenAddressType() {
	makeTokenTypesMap()

	for key := range mapProperties {
		fmt.Printf("Комбинация = %02x %08b\n", key, key)
	}
	fmt.Printf("===========\n")
	console := slCC.InstanceColor(slCC.ColorTypeInfo)
	console.Printf("Всего значений = %d\n", len(mapProperties))
}

func ShowConvertAddresses() {
	var baseRow, line = 3, 0
	var formatString = ""

	addresses := dpAddressElements()
	console := slCC.InstanceColor(slCC.ColorTypeInfo)
	for _, elements := range addresses {
		line++
		table, row, col := elements[0], elements[1], elements[2]
		addressCell := makeAddress(table, row, col)
		addressTable := getAddressTable(addressCell)
		addressRow := getAddressRow(addressCell)
		addressCol := getAddressCol(addressCell)
		if line%baseRow == 0 {
			formatString = "AddressCell: = %02x\t Table = %02x Row = %02x, Col = %02x\n\n"
		} else {
			formatString = "AddressCell: = %02x\t Table = %02x Row = %02x, Col = %02x\n"

		}
		console.Printf(
			formatString,
			addressCell,
			addressTable,
			addressRow,
			addressCol)
	}
}

func dpAddressElements() [][]uint32 {
	return [][]uint32{
		{0, 0, 0}, {0, 0, 1}, {0, 0, 2},
		{0, 1, 0}, {0, 1, 1}, {0, 1, 2},
		{0, 2, 0}, {0, 2, 1}, {0, 2, 2},
		{1, 0, 0}, {1, 0, 1}, {1, 0, 2},
		{1, 1, 0}, {1, 1, 1}, {1, 1, 2},
		{1, 2, 0}, {1, 2, 1}, {1, 2, 2},
	}
}
