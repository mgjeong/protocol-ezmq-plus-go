// +build !debug

package ezmqx

import (
	"go.uber.org/zap"

	"fmt"
)

var Logger *zap.Logger

func InitLogger() {
	var err error
	Logger = zap.NewNop()
	if nil != err {
		_ = fmt.Errorf("\nlogger creation failed")
	}
}
