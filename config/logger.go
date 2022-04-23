package config

import (
	"context"
	"fmt"
	"time"
)

type LoggerInterface interface {
	CtxDebug(ctx context.Context, format string, v ...interface{})
}

var DefaultLogger LoggerInterface = &DefaultLoggerImpl{}

type DefaultLoggerImpl struct{}

func (l *DefaultLoggerImpl) CtxDebug(ctx context.Context, format string, v ...interface{}) {
	timePrefix := time.Now().Format("2006-01-02 15:04:05.999")
	fmt.Printf(timePrefix+" "+format+"\n", v...)
}
