package abbyyJsonParser

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	jsonsl "kallaur.ru/libs/abbyyservice/pkg/json-sl"
	"strings"
)

// Основная точка входа
func getLexems(json string) ([]*AbbyyLexem, []*jlexer.LexerError) {
	var al = make([]*AbbyyLexem, 0, 2)
	jsonByte := easyjson.RawMessage(json)
	als, lErrors := jsonsl.MarshalArray(jsonByte)
	if lErrors != nil {
		return nil, lErrors
	}
	for _, pLexem := range als {
		l := jlexer.Lexer{
			Data: pLexem,
		}
		objectLexem, lErrors := parseLexem(&l)
		if lErrors != nil {
			return nil, lErrors
		}
		al = append(al, objectLexem)
	}
	return al, nil
}

/// разбираем полностью лексему
func parseLexem(l *jlexer.Lexer) (*AbbyyLexem, []*jlexer.LexerError) {
	var al *AbbyyLexem
	if l.IsDelim(123) {
		l.Skip()
	}
	l.FetchToken()
	structLexem := jsonsl.MarshalElementString(l, "Lexem", true)
	structPartOfSpeech := jsonsl.MarshalElementString(l, "PartOfSpeech", true)
	if l.GetNonFatalErrors() != nil {
		return nil, l.GetNonFatalErrors()
	}
	// создаем объект AbbyyLexem и возвращаем на него ссылку
	al = NewAbbyyLexem(strings.ToLower(structLexem), structPartOfSpeech)
	// вырезаем ключ в котором лежит объект ParadigmJson
	jsonsl.SkipKey(l)
	// вырезаем объект ParadigmJson
	paradigmJson := jsonsl.MarshalJsonObject(l)
	if l.GetNonFatalErrors() != nil {
		return nil, l.GetNonFatalErrors()
	}
	// данные для разбора парадигмы
	var name string
	var grammar []string

	groups, lErr := parseParadigmJson(paradigmJson, &name, &grammar)
	if lErr != nil {
		return nil, lErr
	}
	al.ParadigmName = strings.ToLower(name)
	al.Grammar = grammar
	al.Groups = groups
	return al, nil
}

// Возвращаем группы которые составляют парадигму
func parseParadigmJson(
	paradigmJson []byte,
	name *string,
	grammar *[]string) ([]*AbbyyGroup, []*jlexer.LexerError) {
	var groups = make([]*AbbyyGroup, 0, 2)

	l := jlexer.Lexer{
		Data: paradigmJson,
	}
	if l.IsDelim(123) {
		l.Skip()
	}
	tmpName := jsonsl.MarshalElementString(&l, "Name", true)
	grammarString := jsonsl.MarshalElementString(&l, "Grammar", true)
	grammarString = strings.ToLower(grammarString)
	if l.GetNonFatalErrors() != nil {
		return nil, l.GetNonFatalErrors()
	}

	*name = tmpName
	*grammar = strings.Split(grammarString, ",")
	// вырезаем ключ Groups
	jsonsl.SkipKey(&l)
	var groupsRaw []byte
	if l.IsDelim(91) {
		groupsRaw = l.Raw()
	}
	if groupsRaw == nil {
		return nil, l.GetNonFatalErrors()
	}
	// вырезаем массив групп
	groupsJson, lErrors := jsonsl.MarshalArray(groupsRaw)
	if lErrors != nil {
		return nil, l.GetNonFatalErrors()
	}

	// разбираем группы
	for idx, groupJson := range groupsJson {
		group, lErrors := parseGroup(groupJson, idx)
		if lErrors != nil {
			return nil, lErrors
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func parseGroup(group []byte, idx int) (*AbbyyGroup, []*jlexer.LexerError) {
	var tableJson []byte
	l := jlexer.Lexer{
		Data: group,
	}
	if l.IsDelim(123) {
		l.Skip()
	}

	name := jsonsl.MarshalElementString(&l, "Name", true)
	jsonsl.SkipKey(&l)
	var rows [][][]byte
	var lErr []*jlexer.LexerError
	if l.IsDelim(91) {
		tableJson = l.Raw()
		rows, lErr = jsonsl.MarshalArrayArray(tableJson)
		if lErr != nil {
			return nil, lErr
		}
	}
	l.WantComma()
	columnCount := jsonsl.MarshalElementInt(&l, "ColumnCount", true)
	rowCount := jsonsl.MarshalElementInt(&l, "RowCount", true)
	if columnCount == 0 || rowCount == 0 {
		l.AddNonFatalError(
			fmt.Errorf("table have columns = %d and rows = %d"))
	}

	ag := new(AbbyyGroup)
	ag.Name = name
	ag.Idx = idx
	ag.Columns = columnCount
	ag.Rows = rowCount
	ag.Data = rows

	return ag, nil
}

// В группу входит имя группы и таблица со словами
func parseTable(table *AbbyyGroup, isVerb bool, words *map[string][]string) []*jlexer.LexerError {
	var value, prefix string
	// сохраняем токены и их адреса
	var tokens = map[uint32]string{}
	// сохраняем слова по адресам затем будем привязывать токены
	var wordsDraft = map[uint32]string{}
	// переменные уровня таблицы
	var hasTableName, hasTwoColumns, hasFirstValue, hasValue = false, false, false, false

	hasTwoColumns = table.Columns == 2
	hasTableName = len(table.Name) > 0
	if hasTableName {
		addr := makeAddress(uint32(table.Idx), 0, 0)
		addrTable := getAddressTable(addr)
		tokens[addrTable] = strings.ToLower(table.Name)
	}
	// разбираем файл на токены и слова
	// элемент colsArr - это
	for rowIdx, colsArr := range table.Data {
		for colIdx, json := range colsArr {
			l := jlexer.Lexer{
				Data: json,
			}
			value, prefix, _ = marshalAbbyyTableCell(&l)
			hasValue = len(value) > 0
			if !hasValue {
				continue // нет значения в ячейке анализировать нечего
			}
			if hasValue && rowIdx == 0 && colIdx == 0 {
				hasFirstValue = true
			}
			address := makeAddress(uint32(table.Idx), uint32(rowIdx), uint32(colIdx))
			addressProperties := makeAddressProperties(
				table.Idx,
				table.Rows,
				rowIdx,
				colIdx,
				isVerb,
				hasFirstValue,
				hasTableName,
				hasTwoColumns)

			if isTokenAddress(address, addressProperties) {
				tokenAddress := getTokenTypeAddress(addressProperties, address)
				tokens[tokenAddress] = strings.ToLower(value)
				continue
			} else {
				// эксклюзивный вариант таблица времен глагола
				if isVerb && hasTableName {
					tokens[address] = strings.ToLower(prefix)
				}
			}
			// записываем слово
			if hasValue {
				wordsDraft[address] = value
			}
		}
	}

	// соединяем слова и токены. Предварительно обрабатываем токены и слова
	linkWordsAndTokens(tokens, wordsDraft, words)
	return nil
}

func marshalAbbyyTableCell(l *jlexer.Lexer) (string, string, string) {
	var value, prefix, row = "", "", ""

	if !l.IsDelim(123) {
		return value, prefix, row
	}
	l.Skip()
	if value = jsonsl.MarshalElementString(l, "Value", true); value == "null" {
		value = ""
	} else {
		value = trimmingWordsStr(value)
	}
	// конвертируем значения что бы остались только слова

	if prefix = jsonsl.MarshalElementString(l, "Prefix", true); prefix == "null" {
		prefix = ""
	} else {
		prefix = trimmingWordsStr(prefix)
	}
	if row = jsonsl.MarshalElementString(l, "Row", false); row == "null" {
		row = ""
	} else {
		row = trimmingWordsStr(row)
	}

	// переводим все в нижний регистр
	return strings.ToLower(value), strings.ToLower(prefix), strings.ToLower(row)
}
