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
