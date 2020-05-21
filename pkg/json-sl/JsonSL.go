package json_sl

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

func IdentifyTypeJson(json easyjson.RawMessage) int {
	const TypeJsonUndefined = 0x00
	const TypeJsonObject = 0x01
	const TypeJsonArrayObjects = 0x02
	const TypeJsonArrayArray = 0x04

	l := jlexer.Lexer{
		Data: json,
	}

	start := l.IsStart()
	if !start {
		return TypeJsonUndefined
	}

	if l.IsDelim(123) {
		return TypeJsonObject
	}

	if l.IsDelim(91) {
		l.Skip()
		if l.IsDelim(91) {
			return TypeJsonArrayArray
		} else if l.IsDelim(123) {
			return TypeJsonArrayObjects
		}
	}
	return TypeJsonUndefined
}

func Marshal(json easyjson.RawMessage, value *interface{}) []*jlexer.LexerError {
	l := jlexer.Lexer{
		Data: json,
	}

	start := l.IsStart()
	if !start {
		return l.GetNonFatalErrors()
	}

	return l.GetNonFatalErrors()
}

func MarshalArray(json easyjson.RawMessage) ([][]byte, []*jlexer.LexerError) {
	var objects = make([][]byte, 0, 2)
	l := jlexer.Lexer{
		Data: json,
	}

	start := l.IsStart()
	if !start {
		return nil, l.GetNonFatalErrors()
	}

	if l.IsDelim(91) {
		l.Skip()
	} else {
		l.AddNonFatalError(fmt.Errorf("not correct type json, migtht be = %s", "Array of objects"))
		return nil, l.GetNonFatalErrors()
	}

	for l.Ok() {
		if l.IsDelim(123) {
			element := l.Raw()
			objects = append(objects, element)
		}
		l.WantComma()
		if l.IsDelim(93) {
			break
		}
	}

	return objects, nil
}

func MarshalArrayArray(json easyjson.RawMessage) ([][][]byte, []*jlexer.LexerError) {
	var objects = make([][][]byte, 0, 2)
	l := jlexer.Lexer{
		Data: json,
	}

	start := l.IsStart()
	if !start {
		return nil, l.GetNonFatalErrors()
	}

	if l.IsDelim(91) {
		l.Skip()
	} else {
		l.AddNonFatalError(fmt.Errorf("not correct type json, migtht be = %s", "Array array of objects"))
		return nil, l.GetNonFatalErrors()
	}
	// Обнаружено зацикливание на таком варианте.
	// Контролируем через состояние. Пока >= 0 все хорошо
	// Иначе позитивной нагрузки нет.
	var state = 0
	for l.Ok() {
		if state < 0 {
			break
		}
		if l.IsDelim(91) {
			element := l.Raw()
			arr, lErr := MarshalArray(element)
			objects = append(objects, arr)
			if lErr != nil {
				return nil, lErr
			}
			state++
		}
		l.WantComma()
		state--
	}

	return objects, nil
}

// lexer должен стоять перед позицией открывающей фигурной скобки
func MarshalJsonObject(l *jlexer.Lexer) []byte {
	if l.IsDelim(123) {
		return l.Raw()
	}
	return nil
}

// вызываем осознанно, что бы забрать имя = значение. Сразу проверяем валидность полученного имени поля
func MarshalElementString(l *jlexer.Lexer, validKeyName string, withClose bool) string {
	if !readAndValidateKey(l, validKeyName) {
		if withClose {
			l.WantComma()
		}
		return ""
	}
	l.WantColon()
	value := l.String()
	if withClose {
		l.WantComma()
	}
	return value
}

func MarshalElementInt(l *jlexer.Lexer, validKeyName string, withClose bool) int {
	if !readAndValidateKey(l, validKeyName) {
		if withClose {
			l.WantComma()
		}
		return 0
	}
	l.WantColon()
	value := l.Int()
	if withClose {
		l.WantComma()
	}
	return value
}

// первое это значение, второе говорит что это реально было найдено а не по умолчанию
func MarshalElementBool(l *jlexer.Lexer, validKeyName string, withClose bool) (bool, bool) {
	if !readAndValidateKey(l, validKeyName) {
		if withClose {
			l.WantComma()
		}
		return false, false
	}
	l.WantColon()
	value := l.Bool()
	if withClose {
		l.WantComma()
	}
	return value, true
}

// Пропустить ключ элемента json-а
func SkipKey(l *jlexer.Lexer) {
	l.FetchToken()
	l.WantColon()
	l.Skip()
}
