package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lamberjack "gopkg.in/natefinch/lumberjack.v2"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	"os"
)

var logger *zap.Logger

func makeCore() zapcore.Core {
	var encoder zapcore.Encoder
	var core zapcore.Core
	isDev, err := slConfig.GetValue(keeper.GetCommonKeyNameEnv())
	if err != nil {
		appErr := appError.CreateAppError(
			appError.ErrorClassBootstrap,
			"error_make_core_logger",
			"Not found env variable in config file",
			err)
		fmt.Printf(appErr.Error())
		os.Exit(1)
	}

	if isDev == keeper.EnvDev {
		encoder = zapcore.NewJSONEncoder(
			zap.NewDevelopmentEncoderConfig())
	} else {
		encoder = zapcore.NewJSONEncoder(
			zap.NewProductionEncoderConfig())
	}

	writeSyncer := zapcore.AddSync(makeLogWriter())
	if isDev == keeper.EnvDev {
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)
	}

	return core
}

func makeLogWriter() *lamberjack.Logger {
	return &lamberjack.Logger{
		Filename:   configLJC.filename,
		MaxSize:    configLJC.maxSize,
		MaxAge:     configLJC.maxAge,
		MaxBackups: configLJC.maxBackups,
		LocalTime:  configLJC.localTime,
		Compress:   configLJC.compress,
	}
}
