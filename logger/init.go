package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var _logger *zap.Logger

func init() {
	_logger, _ = zap.NewProduction()
	defer _logger.Sync() // flushes buffer, if any
	sugar := _logger.Sugar()
	fmt.Println(sugar)
}
