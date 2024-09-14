package log

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (that *LoggerMock) NewContext(ctx context.Context) context.Context {
	that.Called(ctx)
	return ctx
}

func (that *LoggerMock) WithField(ctx context.Context, key string, value interface{}) Logger {
	that.Called(ctx, key, value)
	return that
}

func (that *LoggerMock) WithFields(ctx context.Context, fields Fields) Logger {
	that.Called(ctx, fields)
	return that
}

func (that *LoggerMock) WithError(ctx context.Context, err error) Logger {
	that.Called(ctx, err)
	return that
}

func (that *LoggerMock) Tracef(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Debugf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Infof(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Printf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Warnf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Warningf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Errorf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Fatalf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Panicf(ctx context.Context, format string, args ...interface{}) {
	that.Called(ctx, format, args)
}

func (that *LoggerMock) Trace(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Debug(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Info(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Print(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Warn(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Warning(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Error(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Fatal(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Panic(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Traceln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Debugln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Infoln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Println(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Warnln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Warningln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Errorln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Fatalln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func (that *LoggerMock) Panicln(ctx context.Context, args ...interface{}) {
	that.Called(ctx, args)
}

func NewLoggerMock() *LoggerMock {
	return new(LoggerMock)
}

type DummyLogger struct {
}

func (that *DummyLogger) NewContext(ctx context.Context) context.Context {
	return ctx
}

func (that *DummyLogger) WithField(ctx context.Context, key string, value interface{}) Logger {
	return that
}

func (that *DummyLogger) WithFields(ctx context.Context, fields Fields) Logger {
	return that
}

func (that *DummyLogger) WithError(ctx context.Context, err error) Logger {
	return that
}

func (that *DummyLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Infof(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Printf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
}

func (that *DummyLogger) Trace(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Debug(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Info(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Print(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Warn(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Warning(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Error(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Fatal(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Panic(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Traceln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Debugln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Infoln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Println(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Warnln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Warningln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Errorln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Fatalln(ctx context.Context, args ...interface{}) {
}

func (that *DummyLogger) Panicln(ctx context.Context, args ...interface{}) {
}

func NewDummyLogger() *DummyLogger {
	return new(DummyLogger)
}
