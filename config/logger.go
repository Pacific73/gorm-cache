package config

import (
	"context"
	"fmt"
	"time"
)

type LoggerInterface interface {
	CtxInfo(ctx context.Context, format string, v ...interface{})
	CtxError(ctx context.Context, format string, v ...interface{})
}

var DefaultLogger LoggerInterface = &DefaultLoggerImpl{}

type DefaultLoggerImpl struct{}

func (l *DefaultLoggerImpl) CtxInfo(ctx context.Context, format string, v ...interface{}) {
	timePrefix := time.Now().Format("2006-01-02 15:04:05.999")
	fmt.Printf(timePrefix+" [INFO] "+format+"\n", v...)
}

func (l *DefaultLoggerImpl) CtxError(ctx context.Context, format string, v ...interface{}) {
	timePrefix := time.Now().Format("2006-01-02 15:04:05.999")
	fmt.Printf(timePrefix+" [ERROR] "+format+"\n", v...)
}
