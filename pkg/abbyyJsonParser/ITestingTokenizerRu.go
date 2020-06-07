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

func ShowRealElements() error {
	wordsAddr := dpAddressRealElements()
	mapTokens := dpAddressRealToken()
	console := slCC.InstanceColor(slCC.ColorTypeInfo)
	for _, addr := range wordsAddr {
		addrRow := getAddressRow(addr)
		addrCol := getAddressCol(addr)
		_, ok := mapTokens[addrRow]
		if !ok {
			return fmt.Errorf("for address %02x not found row address", addr)
		}
		console.Printf("For address %d found row addr = %d\n",
			addr,
			addrRow)
		_, ok = mapTokens[addrCol]
		if !ok {
			return fmt.Errorf("for address %02x not found col address", addr)
		}
		console.Printf("For address %d found col addr = %d\n\n",
			addr,
			addrCol)
	}
	return nil
}

func ViewVerbAddressProperties() {
	// позиционирование на нулевом элементе таблицы
	// zeroElement := currentIdxColZero | currentIdxRowZero
	// нулевой элемент имеет значение
	zeroElementWithValue := currentIdxColZero | currentIdxRowZero | haveFirstValue
	// нулевой элемент не имеет значения
	//zeroElementWithoutValue := zeroElement

	properties := isFirstTable | tableOneRow | paosIsVerb | zeroElementWithValue
	fmt.Printf(
		"\nТаблица:    Инфинитив\nЗначение:   %b\n",
		properties)
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

func dpAddressRealElements() []uint32 {
	return []uint32{
		306, 307, 322, 338,
		370, 371, 290, 291,
		323, 339, 354, 355,
	}
}

func dpAddressRealToken() map[uint32]uint32 {
	return map[uint32]uint32{
		274: 274, 275: 275, 289: 289,
		305: 305, 321: 321, 337: 337,
		353: 353, 369: 369,
	}
}
