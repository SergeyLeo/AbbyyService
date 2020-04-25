package errors

import (
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
)

var eh *ErrorHandler

// Основной интерфейс для внешних вызовов
func InstanceErrorHandler() *ErrorHandler {
	if eh == nil {
		InitErrorHandler()
	}
	return eh
}

func (er *ErrorHandler) Report(e interface{}) {
	if er.Env == keeper.EnvProd {
		er.makeShortReport(e)
	} else {
		er.makeDevReport(e)
	}
}

// Если хотим слить в лог весь массив ошибок из appError
func (er *ErrorHandler) ReportCollection() {
	defer er.Logger.Sync()
	errorElements := appError.GetAllElements()
	if len(errorElements) == 0 {
		return
	}
	for _, e := range errorElements {
		er.Report(*e)
	}
}
