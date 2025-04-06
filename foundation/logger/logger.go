package logger

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, keyValues ...interface{})
	Error(ctx context.Context, msg string, keyValues ...interface{})
}
