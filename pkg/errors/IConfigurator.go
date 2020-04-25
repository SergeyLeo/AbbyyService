package errors

import (
	"fmt"
	"kallaur.ru/libs/abbyyservice/pkg/appError"
	slConfig "kallaur.ru/libs/abbyyservice/pkg/envconfig"
	"kallaur.ru/libs/abbyyservice/pkg/keeper"
	appLogger "kallaur.ru/libs/abbyyservice/pkg/logger"
	"os"
)

// вызываем в функции
func InitErrorHandler() {
	var err *appError.AppError
	logger := appLogger.InstanceLogger()

	eh = new(ErrorHandler)
	eh.Logger = logger
	eh.Env, err = slConfig.GetValue(keeper.GetCommonKeyNameEnv())
	if err != nil {
		appErr := appError.CreateAppError(
			appError.ErrorClassBootstrap,
			"not_found_config_value",
			fmt.Sprintf("not found value with key = %s", keeper.GetCommonKeyNameEnv()),
			err)
		fmt.Printf(appErr.Error())
		os.Exit(1)
	}
}
