package log

import (
	"context"
)

type Behavior interface {
	Resolve(ctx context.Context) Logger
	NewContext(ctx context.Context) context.Context
}

type CustomLogger struct {
	logger   Logger
	behavior Behavior
}

func NewCustomLogger(
	logger Logger,
	behavior Behavior,
) *CustomLogger {
	return &CustomLogger{
		logger:   logger,
		behavior: behavior,
	}
}

func (that *CustomLogger) NewContext(ctx context.Context) context.Context {
	return that.behavior.NewContext(ctx)
}

func (that *CustomLogger) WithField(ctx context.Context, key string, value interface{}) Logger {
	return that.behavior.Resolve(ctx).WithField(ctx, key, value)
}

func (that *CustomLogger) WithFields(ctx context.Context, fields Fields) Logger {
	return that.behavior.Resolve(ctx).WithFields(ctx, fields)
}

func (that *CustomLogger) WithError(ctx context.Context, err error) Logger {
	return that.behavior.Resolve(ctx).WithError(ctx, err)
}

func (that *CustomLogger) Tracef(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Tracef(ctx, format, args...)
}

func (that *CustomLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Debugf(ctx, format, args...)
}

func (that *CustomLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Infof(ctx, format, args...)
}

func (that *CustomLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Warningf(ctx, format, args...)
}

func (that *CustomLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Errorf(ctx, format, args...)
}

func (that *CustomLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Fatalf(ctx, format, args...)
}

func (that *CustomLogger) Panicf(ctx context.Context, format string, args ...interface{}) {
	that.behavior.Resolve(ctx).Panicf(ctx, format, args...)
}

func (that *CustomLogger) Trace(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Trace(ctx, args...)
}

func (that *CustomLogger) Debug(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Debug(ctx, args...)
}

func (that *CustomLogger) Info(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Info(ctx, args...)
}

func (that *CustomLogger) Warning(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Warning(ctx, args...)
}

func (that *CustomLogger) Error(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Error(ctx, args...)
}

func (that *CustomLogger) Fatal(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Fatal(ctx, args...)
}

func (that *CustomLogger) Panic(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Panic(ctx, args...)
}

func (that *CustomLogger) Traceln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Traceln(ctx, args...)
}

func (that *CustomLogger) Debugln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Debugln(ctx, args...)
}

func (that *CustomLogger) Infoln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Infoln(ctx, args...)
}

func (that *CustomLogger) Warningln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Warningln(ctx, args...)
}

func (that *CustomLogger) Errorln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Errorln(ctx, args...)
}

func (that *CustomLogger) Fatalln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Fatalln(ctx, args...)
}

func (that *CustomLogger) Panicln(ctx context.Context, args ...interface{}) {
	that.behavior.Resolve(ctx).Panicln(ctx, args...)
}

type FieldsBehavior struct {
	Logger
	fields Fields
}

func NewFieldsBehavior(logger Logger, fields Fields) *FieldsBehavior {
	return &FieldsBehavior{
		Logger: logger,
		fields: fields,
	}
}

func (that *FieldsBehavior) Resolve(ctx context.Context) Logger {
	return that.Logger.WithFields(ctx, that.fields)
}

//
//type EntityActionLogger struct {
//	Logger
//}
//
//func NewEntityActionLogger(logger Logger, actions map[EntityAction]string) *EntityActionLogger {
//	if actions == nil {
//		actions = EntityActions
//	}
//
//	return &EntityActionLogger{
//		Logger:  logger,
//		actions: actions,
//	}
//}
//
//func (that *EntityActionLogger) getAction(action EntityAction) string {
//	val, _ := that.actions[action]
//	return val
//}
//
//func (that *EntityActionLogger) WithRequestReceived(ctx context.Context) Logger {
//	return that.WithEntityAction(ctx, that.getAction(ActionRequestReceived))
//}
//
//func (that *EntityActionLogger) WithResponseSent(ctx context.Context) Logger {
//	return that.WithEntityAction(ctx, that.getAction(ActionResponseSent))
//}
//
//func (that *EntityActionLogger) WithRequestSent(ctx context.Context) Logger {
//	return that.WithEntityAction(ctx, that.getAction(ActionRequestSent))
//}
//
//func (that *EntityActionLogger) WithResponseReceived(ctx context.Context) Logger {
//	return that.WithEntityAction(ctx, that.getAction(ActionResponseReceived))
//}
