package log

import (
	"context"
	"log"
)

var SysLog Logger = NewDummyLogger()

func WithField(ctx context.Context, key string, value interface{}) Logger {
	return getLogger(ctx, SysLog).WithField(ctx, key, value)
}

func WithFields(ctx context.Context, fields Fields) Logger {
	return getLogger(ctx, SysLog).WithFields(ctx, fields)
}

func WithError(ctx context.Context, err error) Logger {
	return getLogger(ctx, SysLog).WithError(ctx, err)
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Tracef(ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Debugf(ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Infof(ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Warningf(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Errorf(ctx, format, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Fatalf(ctx, format, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	getLogger(ctx, SysLog).Panicf(ctx, format, args...)
}

func Trace(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Trace(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Debug(ctx, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Info(ctx, args...)
}

func Warning(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Warning(ctx, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Error(ctx, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Fatal(ctx, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Panic(ctx, args...)
}

func Traceln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Traceln(ctx, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Debugln(ctx, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Infoln(ctx, args...)
}

func Warningln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Warningln(ctx, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Errorln(ctx, args...)
}

func Fatalln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Fatalln(ctx, args...)
}

func Panicln(ctx context.Context, args ...interface{}) {
	getLogger(ctx, SysLog).Panicln(ctx, args...)
}

// deprecated functions
func DoPrintln(ctx context.Context, v ...interface{}) {
	if SysLog == nil {
		log.Println(v...)
		return
	}

	SysLog.Info(ctx, v...)
}

func DoPrintf(ctx context.Context, format string, v ...interface{}) {
	if SysLog == nil {
		log.Printf(format, v...)
		return
	}

	SysLog.Infof(ctx, format, v...)
}

func DoDebugln(ctx context.Context, v ...interface{}) {
	if SysLog == nil {
		log.Println(v...)
		return
	}

	SysLog.Debug(ctx, v...)
}

func DoTraceln(ctx context.Context, v ...interface{}) {
	if SysLog == nil {
		log.Println(v...)
		return
	}

	SysLog.Trace(ctx, v...)
}
