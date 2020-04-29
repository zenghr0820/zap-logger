package zapLogger

import "go.uber.org/zap"

type logger struct {
	zapSugar *zap.SugaredLogger
}

// NewLoggerOf logger with component field
func NewLoggerOf(zapSugar *zap.SugaredLogger) Logger {
	return &logger{
		zapSugar: zapSugar.With(),
	}
}

// With adds a variadic number of fields to the logging context.
// see https://github.com/uber-go/zap/blob/v1.10.0/sugar.go#L91
func (l *logger) With(args ...interface{}) *logger {
	l.zapSugar = l.zapSugar.With(args...)
	return l
}

// Debug package sugar of zap
func (l *logger) Debug(args ...interface{}) {
	l.zapSugar.Debug(args...)
}

// Debugf package sugar of zap
func (l *logger) Debugf(template string, args ...interface{}) {
	l.zapSugar.Debugf(template, args...)
}

// Info package sugar of zap
func (l *logger) Info(args ...interface{}) {
	l.zapSugar.Info(args...)
}

// Infof package sugar of zap
func (l *logger) Infof(template string, args ...interface{}) {
	l.zapSugar.Infof(template, args...)
}

// Warn package sugar of zap
func (l *logger) Warn(args ...interface{}) {
	l.zapSugar.Warn(args...)
}

// Warnf package sugar of zap
func (l *logger) Warnf(template string, args ...interface{}) {
	l.zapSugar.Warnf(template, args...)
}

// Error package sugar of zap
func (l *logger) Error(args ...interface{}) {
	l.zapSugar.Error(args...)
}

// Errorf package sugar of zap
func (l *logger) Errorf(template string, args ...interface{}) {
	l.zapSugar.Errorf(template, args...)
}

// Fatal package sugar of zap
func (l *logger) Fatal(args ...interface{}) {
	l.zapSugar.Fatal(args...)
}

// Fatalf package sugar of zap
func (l *logger) Fatalf(template string, args ...interface{}) {
	l.zapSugar.Fatalf(template, args...)
}

// Panic package sugar of zap
func (l *logger) Panic(args ...interface{}) {
	l.zapSugar.Panic(args...)
}

// Panicf package sugar of zap
func (l *logger) Panicf(template string, args ...interface{}) {
	l.zapSugar.Panicf(template, args...)
}
