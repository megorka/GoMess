package logger

import (
	"context"
	"go.uber.org/zap"
)

type Logger struct {
	l *zap.Logger
}

const (
	Key       = "logger"
	RequestID = "request_id"
)

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, Key, &Logger{logger})
	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) *Logger {
	return ctx.Value(Key).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Info(msg, fields...)

}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {

	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Fatal(msg, fields...)

}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}
	l.l.Error(msg, fields...)
}
