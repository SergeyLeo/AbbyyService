package json_sl

import (
	"fmt"
	"github.com/mailru/easyjson/jlexer"
)

// Вспомогательная функция для получения и контроля значений
func readAndValidateKey(l *jlexer.Lexer, validKeyName string) bool {
	key := l.String()
	if key != validKeyName {
		// добавляем фатальную ошибку, которую можно проверить потом как l.Ok()
		l.AddError(fmt.Errorf("key %s != validKeyName %s", key, validKeyName))
		l.WantColon()
		l.Skip()
		return false
	}

	return true
}
