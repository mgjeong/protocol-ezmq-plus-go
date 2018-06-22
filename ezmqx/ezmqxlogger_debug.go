// +build debug

package ezmqx

import (
	"go.uber.org/zap"

	"fmt"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	Logger, err = zap.NewDevelopment()
	if nil != err {
		_ = fmt.Errorf("\nlogger creation failed")
	}
}
