package log

import (
	"context"
	"github.com/adverax/core/log/logrus"
)

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	WarnLevel  = logrus.WarnLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
	TraceLevel = logrus.TraceLevel
)

const (
	FieldKeyTraceID  = logrus.FieldKeyTraceID
	FieldKeyEntity   = logrus.FieldKeyEntity
	FieldKeyAction   = logrus.FieldKeyAction
	FieldKeyMethod   = logrus.FieldKeyMethod
	FieldKeySubject  = logrus.FieldKeySubject
	FieldKeyData     = logrus.FieldKeyData
	FieldKeyLogType  = "type"
	FieldKeyDuration = "duration"
	FieldKeyCause    = "cause"
)

const (
	EntityUnknown   = "UNKNOWN"
	EntityHttp      = "HTTP"
	EntityMaster    = "MASTER"
	EntitySlave     = "SLAVE"
	EntityFront     = "FRONT"
	EntityFrontHttp = "FRONT-HTTP"
	EntitySQL       = "SQL"
	EntityBus       = "BUS"
)

const (
	TypeSql    = "sql"
	TypeHttp   = "http"
	TypeSocket = "socket"
)

const (
	EntityActionIncomeRequest       = ">>"
	EntityActionOutcomeRequest      = "<<"
	EntityActionIncomeResponse      = ">"
	EntityActionOutcomeResponse     = "<"
	EntityActionTransitRequest      = ">>>"
	EntityActionTransitResponse     = "<--"
	EntityActionTransitNotification = "<<<"
)

const (
	ActionRequestReceived  = ">>"
	ActionResponseSent     = "<"
	ActionRequestSent      = "<<"
	ActionResponseReceived = ">"
)

type Purifier = logrus.Purifier

type Level = logrus.Level

type Fields map[string]interface{}

type Logger interface {
	NewContext(ctx context.Context) context.Context

	WithField(ctx context.Context, key string, value interface{}) Logger
	WithFields(ctx context.Context, fields Fields) Logger
	WithError(ctx context.Context, err error) Logger

	Tracef(ctx context.Context, format string, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})

	Trace(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warning(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Panic(ctx context.Context, args ...interface{})

	Traceln(ctx context.Context, args ...interface{})
	Debugln(ctx context.Context, args ...interface{})
	Infoln(ctx context.Context, args ...interface{})
	Warningln(ctx context.Context, args ...interface{})
	Errorln(ctx context.Context, args ...interface{})
	Fatalln(ctx context.Context, args ...interface{})
	Panicln(ctx context.Context, args ...interface{})
}

type ErrorTuner interface {
	TuneError(error) error
}

type ErrorTunerFunc func(error) error

func (fn ErrorTunerFunc) TuneError(err error) error {
	return fn(err)
}

type logger struct {
	*logrus.Logger
	errors ErrorTuner
}

func (that *logger) NewContext(ctx context.Context) context.Context {
	return ctx
}

func (that *logger) WithField(ctx context.Context, key string, value interface{}) Logger {
	fs := make(logrus.Fields, 1)
	fs[key] = value
	return &entry{logger: that, Entry: that.Logger.WithFields(fs)}
}

func (that *logger) WithFields(ctx context.Context, fields Fields) Logger {
	fs := make(logrus.Fields, len(fields))
	for k, v := range fields {
		fs[k] = v
	}
	return &entry{logger: that, Entry: that.Logger.WithFields(fs)}
}

func (that *logger) WithError(ctx context.Context, err error) Logger {
	return that.WithField(ctx, logrus.ErrorKey, that.errors.TuneError(err))
}

func (that *logger) WithEntityAction(ctx context.Context, action string) Logger {
	if action == "" {
		return that
	}

	return that.WithField(ctx, logrus.FieldKeyAction, action)
}

func (that *logger) WithRequestReceived(ctx context.Context) Logger {
	return that.WithEntityAction(ctx, ActionRequestReceived)
}

func (that *logger) WithResponseSent(ctx context.Context) Logger {
	return that.WithEntityAction(ctx, ActionResponseSent)
}

func (that *logger) WithRequestSent(ctx context.Context) Logger {
	return that.WithEntityAction(ctx, ActionRequestSent)
}

func (that *logger) WithResponseReceived(ctx context.Context) Logger {
	return that.WithEntityAction(ctx, ActionResponseReceived)
}

func (that *logger) Tracef(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Tracef(format, args...)
}

func (that *logger) Debugf(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Debugf(format, args...)
}

func (that *logger) Infof(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Infof(format, args...)
}

func (that *logger) Warningf(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Warningf(format, args...)
}

func (that *logger) Errorf(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Errorf(format, args...)
}

func (that *logger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Fatalf(format, args...)
}

func (that *logger) Panicf(ctx context.Context, format string, args ...interface{}) {
	that.Logger.Panicf(format, args...)
}

func (that *logger) Trace(ctx context.Context, args ...interface{}) {
	that.Logger.Trace(args...)
}

func (that *logger) Debug(ctx context.Context, args ...interface{}) {
	that.Logger.Debug(args...)
}

func (that *logger) Info(ctx context.Context, args ...interface{}) {
	that.Logger.Info(args...)
}

func (that *logger) Warning(ctx context.Context, args ...interface{}) {
	that.Logger.Warning(args...)
}

func (that *logger) Error(ctx context.Context, args ...interface{}) {
	that.Logger.Error(args...)
}

func (that *logger) Fatal(ctx context.Context, args ...interface{}) {
	that.Logger.Fatal(args...)
}

func (that *logger) Panic(ctx context.Context, args ...interface{}) {
	that.Logger.Panic(args...)
}

func (that *logger) Traceln(ctx context.Context, args ...interface{}) {
	that.Logger.Traceln(args...)
}

func (that *logger) Debugln(ctx context.Context, args ...interface{}) {
	that.Logger.Debugln(args...)
}

func (that *logger) Infoln(ctx context.Context, args ...interface{}) {
	that.Logger.Infoln(args...)
}

func (that *logger) Warningln(ctx context.Context, args ...interface{}) {
	that.Logger.Warningln(args...)
}

func (that *logger) Errorln(ctx context.Context, args ...interface{}) {
	that.Logger.Errorln(args...)
}

func (that *logger) Fatalln(ctx context.Context, args ...interface{}) {
	that.Logger.Fatalln(args...)
}

func (that *logger) Panicln(ctx context.Context, args ...interface{}) {
	that.Logger.Panicln(args...)
}

type entry struct {
	*logrus.Entry
	logger *logger
}

func (that *entry) NewContext(ctx context.Context) context.Context {
	return ctx
}

func (that *entry) WithField(ctx context.Context, key string, value interface{}) Logger {
	fs := make(logrus.Fields, len(that.Data)+1)
	for k, v := range that.Data {
		fs[k] = v
	}
	fs[key] = value
	return &entry{logger: that.logger, Entry: that.Logger.WithFields(fs)}
}

func (that *entry) WithFields(ctx context.Context, fields Fields) Logger {
	fs := make(logrus.Fields, len(fields)+len(that.Data))
	for k, v := range that.Data {
		fs[k] = v
	}
	for k, v := range fields {
		fs[k] = v
	}
	return &entry{logger: that.logger, Entry: that.Logger.WithFields(fs)}
}

func (that *entry) WithError(ctx context.Context, err error) Logger {
	return that.WithField(ctx, logrus.ErrorKey, that.logger.errors.TuneError(err))
}

func (that *entry) Tracef(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Tracef(format, args...)
}

func (that *entry) Debugf(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Debugf(format, args...)
}

func (that *entry) Infof(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Infof(format, args...)
}

func (that *entry) Warningf(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Warningf(format, args...)
}

func (that *entry) Errorf(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Errorf(format, args...)
}

func (that *entry) Fatalf(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Fatalf(format, args...)
}

func (that *entry) Panicf(ctx context.Context, format string, args ...interface{}) {
	that.Entry.Panicf(format, args...)
}

func (that *entry) Trace(ctx context.Context, args ...interface{}) {
	that.Entry.Trace(args...)
}

func (that *entry) Debug(ctx context.Context, args ...interface{}) {
	that.Entry.Debug(args...)
}

func (that *entry) Info(ctx context.Context, args ...interface{}) {
	that.Entry.Info(args...)
}

func (that *entry) Warning(ctx context.Context, args ...interface{}) {
	that.Entry.Warning(args...)
}

func (that *entry) Error(ctx context.Context, args ...interface{}) {
	that.Entry.Error(args...)
}

func (that *entry) Fatal(ctx context.Context, args ...interface{}) {
	that.Entry.Fatal(args...)
}

func (that *entry) Panic(ctx context.Context, args ...interface{}) {
	that.Entry.Panic(args...)
}

func (that *entry) Traceln(ctx context.Context, args ...interface{}) {
	that.Entry.Traceln(args...)
}

func (that *entry) Debugln(ctx context.Context, args ...interface{}) {
	that.Entry.Debugln(args...)
}

func (that *entry) Infoln(ctx context.Context, args ...interface{}) {
	that.Entry.Infoln(args...)
}

func (that *entry) Warningln(ctx context.Context, args ...interface{}) {
	that.Entry.Warningln(args...)
}

func (that *entry) Errorln(ctx context.Context, args ...interface{}) {
	that.Entry.Errorln(args...)
}

func (that *entry) Fatalln(ctx context.Context, args ...interface{}) {
	that.Entry.Fatalln(args...)
}

func (that *entry) Panicln(ctx context.Context, args ...interface{}) {
	that.Entry.Panicln(args...)
}

func NewLogger(
	l *logrus.Logger,
	errors ErrorTuner,
) Logger {
	if errors == nil {
		errors = ErrorTunerFunc(func(err error) error {
			return err
		})
	}

	return &logger{Logger: l, errors: errors}
}
