package abbyyJsonParser

import (
	"strings"
)

// файл токенов русского языка для добычи слов
const (
	tableOneRow       uint32 = 0x01 // таблица имеет только одну строку
	paosIsVerb        uint32 = 0x02 // часть речи глагол
	currentIdxRowZero uint32 = 0x04 // текущий индекс ряда
	currentIdxColZero uint32 = 0x08 // текущий индекс колонки
	isFirstTable      uint32 = 0x10 // текущая таблица является первой в группе таблиц
	haveFirstValue    uint32 = 0x20 // есть значение в нулевой строке в нулевой колонке
	tableHasName      uint32 = 0x40 // у таблицы есть имя
	participleTable   uint32 = 0x80 // глагол таблица причастий. В каждом ряду токен
)

type addressTypeFunc func(uint32) uint32

var mapProperties map[uint32]addressTypeFunc

// Координата состоит из трех составляющих:
// - Z - индекс таблицы
// - X - индекс строки
// - Y - нидекс колонки
// Адресация хранения токенов
// - 0x0x0 - токены по этому адресу распространяются на все слова
// - Zx0x0 - токен распространяется на все слова таблицы с индексом Z
// - ZxXx0 - токен распространяется на все слова в таблице с индексом Z в строке с индексом X
// - Zx0xY - токен распространяется на все слова в таблице с индексом Z в колонке с индексом Y
func makeAddress(table uint32, row uint32, col uint32) uint32 {
	// увеличим все индексы на 1
	table++
	row++
	col++
	// в случае не штатной ситуации создаем максимальный индекс, вероятность существования его очень низкая
	// а значит минимизирована возможность лишних токенов
	if table > 0x0f || row > 0x0f || col > 0x0f {
		return 0xfff
	}
	// каждую часть адреса смещаем на 4 бита
	return (table << 8) | (row << 4) | col
}

// makeAddressProperties - формирует идентификатор свойств адреса. По ним мы определяем тип адреса для токена
func makeAddressProperties(
	idxTable int,
	countRows int,
	idxRow int,
	idxCol int,
	isVerb bool,
	hasFirstValue bool,
	nameTable bool,
	twoColumns bool) uint32 {
	var out uint32 = 0

	if countRows == 1 {
		out |= tableOneRow
	}
	if isVerb {
		out |= paosIsVerb
	}
	if idxRow == 0 {
		out |= currentIdxRowZero
	}
	if idxCol == 0 {
		out |= currentIdxColZero
	}
	if idxTable == 0 {
		out |= isFirstTable
	}
	if hasFirstValue {
		out |= haveFirstValue
	}
	if nameTable {
		out |= tableHasName
	}
	if twoColumns {
		out |= participleTable
	}

	return out
}

// Возвращаем в следующей последовательности
// индекс таблицы, индекс строки, индекс колонки
func parseAddress(address uint32) (uint32, uint32, uint32) {
	var col, row, table uint32
	col = (address & 0x0f) - 1
	row = ((address >> 4) & 0x0f) - 1
	table = ((address >> 8) & 0x0f) - 1
	return table, row, col

}

// Параметры
// - isNamed - у таблицы заполнено название
// - isVerb - мы разбираем лексему глагола
// - zeroElementHasValue - нулевой элемент имеет значение
func isTokenAddress(address uint32, addressProperties uint32) bool {
	table, row, col := parseAddress(address)

	isVerb := (addressProperties & paosIsVerb) > 0
	isNamed := (addressProperties & tableHasName) > 0
	zeroElementHasValue := (addressProperties & haveFirstValue) > 0

	if isVerb {
		if addressProperties == makePropertyInfinitiv() {
			// это инфинитив глагола
			return true
		}
		if !isNamed {
			if zeroElementHasValue {
				if col == 0 && table > 0 {
					return true
				}
			} else {
				if row == 0 && table > 0 {
					return true
				}
			}
		}
	}
	// обрабатываем другие части речи
	if table == 0 && row == 0 {
		return true
	}
	if row > 0 && col == 0 {
		return true
	}

	return false
}

// Есть следующие типы адресов:
// - адрес таблицы, токен распространяется на все слова таблицы
// - адрес ряда, токен распространяется на все слова ряда конкретной таблицы
// - адрес колонки, токен распространяется на все слова всех рядов конкретной таблицы
// Возвращает адрес по которому мы записываем токен
func getTokenTypeAddress(addressProperties uint32, address uint32) uint32 {
	if mapProperties == nil {
		makeTokenTypesMap()
	}
	value, ok := mapProperties[addressProperties]
	if !ok {
		return 0 // указанный набор свойств не существует
	}
	return value(address)
}

func makeTokenTypesMap() {
	mapProperties = map[uint32]addressTypeFunc{}

	// инфинитив глагола
	properties := makePropertyInfinitiv()
	mapProperties[properties] = getAddressRow

	// существительное формы ед и множ числа
	// прилагательное полные и краткие формы
	properties = isFirstTable | currentIdxRowZero
	mapProperties[properties] = getAddressCol

	// существительное формы падежей
	// прилагательное изменения по родам
	properties = isFirstTable | currentIdxColZero
	mapProperties[properties] = getAddressRow

	// прилагательное степени сравнения
	properties = tableOneRow | currentIdxRowZero | currentIdxColZero
	mapProperties[properties] = getAddressRow

	// глагол сопуствующие части речи
	properties = paosIsVerb | currentIdxColZero | haveFirstValue
	mapProperties[properties] = getAddressRow

	// глагол токены на колонки как у существительного
	properties = paosIsVerb | currentIdxRowZero
	mapProperties[properties] = getAddressCol

	// глагол токены на строки как у существительного
	properties = paosIsVerb | currentIdxColZero
	mapProperties[properties] = getAddressRow

	// глагол эксклюзивная ситуация когда на причастия идет по 2 колонки
	properties = paosIsVerb | currentIdxColZero | participleTable
	mapProperties[properties] = getAddressRow
}

func makePropertyInfinitiv() uint32 {
	return tableOneRow | paosIsVerb | makePropertyZeroElementWithValue()
}

func makePropertyZeroElementWithValue() uint32 {
	return currentIdxRowZero | currentIdxColZero | haveFirstValue
}

// Получить адрес токена уровня таблица
func getAddressTable(address uint32) uint32 {
	return address & 0xf00
}

// Получить адрес токена уровня строка
func getAddressRow(address uint32) uint32 {
	table, row, _ := parseAddress(address)
	// должны совпадать таблица и строка
	addressRow := makeAddress(table, row, 0)

	return addressRow
}

// Получить адрес токена уровня колонка
func getAddressCol(address uint32) uint32 {
	table, _, col := parseAddress(address)
	// должны совпадать таблица и строка
	addressCol := makeAddress(table, 0, col)

	return addressCol
}

// Линкуем слова и токены
func linkWordsAndTokens(tokens map[uint32]string, wordsDraft map[uint32]string, wordsMap *map[string][]string) {
	// сначала нужно выбрать базовые формы, затем осуществить привязку слов к нему
	// word - может быть несколько слов через пробел или пробел и запятую
	// слово может начинаться со звездочки.
	// в токенах такая же история может наблюдаться

	for addr, wordLine := range wordsDraft {
		addrRow := getAddressRow(addr)
		addrCol := getAddressCol(addr)
		tokensList := make([]string, 0)
		wordsList := make([]string, 0)
		tokensListRow, ok := tokens[addrRow]
		if ok {
			tokenElements := trimmingWords(tokensListRow)
			appendElements(&tokensList, tokenElements)
		}
		tokensListCol, ok := tokens[addrCol]
		if ok {
			tokenElements := trimmingWords(tokensListCol)
			appendElements(&tokensList, tokenElements)
		}
		wordsList = trimmingWords(wordLine)
		// добавляем слова в общую map
		for _, word := range wordsList {
			// проверяем наличие в карте текущего слова
			realTokens, ok := (*wordsMap)[word]
			if ok {
				appendElements(&realTokens, tokensList)
				(*wordsMap)[word] = realTokens
			} else {
				newTokenList := make([]string, len(tokensList))
				copy(newTokenList, tokensList)
				(*wordsMap)[word] = newTokenList
			}
		}
	}
}

// В json может быть вставлено по два слова вместо одного. Требуется выбросить все лишние пробелы и *
func trimmingWords(word string) []string {
	words := strings.Split(word, ",")
	for idx, word := range words {
		word = strings.Trim(word, "* ")
		words[idx] = word
	}

	return words
}

func appendElements(dst *[]string, src []string) {
	for _, value := range src {
		*dst = append(*dst, value)
	}
}
