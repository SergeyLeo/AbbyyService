// Пакет не должен иметь ссылок на внутренние пакеты приложения за исключением пакета с общими константами
// keeper
// Исходя из вышесказанного считывает значение ENV только из памяти а не из файла
package appError

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type AppError struct {
	Class     uint32
	Code      string // код через нижние подчеркивания
	Message   string
	Throwable error
}

func (ae *AppError) Error() string {
	return fmt.Sprintf(
		"Error class: %d | Error code: %s | Error message: %s | Original error: %+v.",
		ae.Class,
		ae.Code,
		ae.Message,
		ae.Throwable)
}

func (ae *AppError) ErrorFmt() string {
	return fmt.Sprintf(
		"Error class: %d\nError code: %s\nError message: %s\nOriginal error: %+v\n",
		ae.Class,
		ae.Code,
		ae.Message,
		ae.Throwable)
}

const (
	ErrorClassBootstrap   uint32 = 1 << iota // перед стартом
	ErrorClassConfiguring                    // конфигурирование компонентов
	ErrorClassConnecting                     // подключения к бд и другим сервисам
	ErrorClassHttp                           // эксклюзивно для http
	ErrorClassDb                             // работа с БД
	ErrorClassRedis                          // работа с редис
	ErrorClassMongo                          // работа с mongo
)

func CreateAppError(
	class uint32,
	code string,
	message string,
	throwable error) *AppError {

	if isDev() {
		return CreateAppErrorWithStack(class, code, message, throwable)
	}

	return &AppError{
		Class:     class,
		Code:      code,
		Message:   message,
		Throwable: throwable,
	}

}

func CreateAppErrorWithStack(
	class uint32,
	code string,
	message string,
	throwable error) *AppError {

	return &AppError{
		Class:     class,
		Code:      code,
		Message:   message,
		Throwable: errors.WithStack(throwable),
	}
}

/*
******** ИНТЕРФЕЙС ДЛЯ ОБСЛУЖИВАНИЯ МАССИВА ОШИБОК ********
 */

func InitErrorsArray(err *AppError) {
	if errorsElements == nil {
		errorsElements = new(appErrorArray)
	}
	errorsElements.count = 0
	errorsElements.errElements = make([]*AppError, 0, 8)

	if err != nil {
		errorsElements.count += 1
		errorsElements.errElements = append(errorsElements.errElements, err)
	}
}

func AddError(err *AppError) {
	if err == nil {
		return
	}
	if errorsElements == nil {
		InitErrorsArray(nil)
	}
	errorsElements.count += 1
	errorsElements.errElements = append(errorsElements.errElements, err)
}

func FlushErrors() string {
	if errorsElements.count == 0 {
		return ""
	}

	tmpString := make([]string, errorsElements.count, errorsElements.count)
	for idx, err := range errorsElements.errElements {
		tmpString[idx] = err.Error()
	}

	errorsElements.count = 0
	errorsElements.errElements = make([]*AppError, 0, 8)

	return strings.Join(tmpString, "\n")
}

func HasErrors() bool {
	if errorsElements == nil {
		return false
	}
	return errorsElements.count > 0
}
func GetAllElements() []*AppError {
	if errorsElements.count == 0 || errorsElements == nil {
		return []*AppError{}
	}
	return errorsElements.errElements
}
