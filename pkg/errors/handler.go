package errors

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
)

type ErrorHandler struct {
	Logger *zap.Logger
	Env    string
}

func (er *ErrorHandler) makeShortReport(e interface{}) {
	switch t := e.(type) {
	case appError.AppError:
		er.reportAppError(&t)
	case error:
		er.Logger.Error(t.Error())
	}
}

func (er *ErrorHandler) makeDevReport(e interface{}) {
	switch t := e.(type) {
	case appError.AppError:
		er.reportAppError(&t)
	case error:
		err := errors.WithStack(t)
		er.Logger.Error(err.Error())
	}
}

func (er *ErrorHandler) reportAppError(e *appError.AppError) {
	err := *e
	er.Logger.Error(err.Message,
		zap.Uint32("Class", err.Class),
		zap.String("Code", err.Code),
		zap.Error(err.Throwable))
}

func (er *ErrorHandler) reportDevAppError(e *appError.AppError) {
	err := *e
	errStack := errors.WithStack(e)
	er.Logger.Error(err.Message,
		zap.Uint32("Class", err.Class),
		zap.String("Code", err.Code),
		zap.Error(errStack))
}
