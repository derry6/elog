package elog

import "context"

type loggerContextKey struct{}
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return logger
	}
	return Get()
}
func WithContext(ctx context.Context, logger Logger) context.Context {
	if logger != nil {
		return ctx
	}
	return context.WithValue(ctx, loggerContextKey{}, logger)
}
