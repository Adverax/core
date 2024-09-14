package log

import (
	"adverax/core/log/logrus"
	"io"
	log2 "log"
)

type LogrusBuilder struct {
	logger *logrus.Logger
}

func NewLogrusBuilder() *LogrusBuilder {
	return &LogrusBuilder{
		logger: logrus.New(),
	}
}

func (that *LogrusBuilder) Output(output io.Writer) *LogrusBuilder {
	that.logger.SetOutput(output)
	return that
}

func (that *LogrusBuilder) Level(level Level) *LogrusBuilder {
	that.logger.SetLevel(level)
	return that
}

func (that *LogrusBuilder) Formatter(formatter logrus.Formatter) *LogrusBuilder {
	that.logger.SetFormatter(formatter)
	return that
}

func (that *LogrusBuilder) Build() (*logrus.Logger, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	if err := that.updateDefaultFields(); err != nil {
		return nil, err
	}

	return that.logger, nil
}

func (that *LogrusBuilder) checkRequiredFields() error {
	return nil
}

func (that *LogrusBuilder) updateDefaultFields() error {
	if that.logger.Out == nil {
		that.logger.SetOutput(log2.Writer())
	}

	if that.logger.Formatter == nil {
		that.logger.SetFormatter(defaultFormatter)
	}

	return nil
}

//func NewLogrus(output io.Writer, level Level) *logrus.Logger {
//	if output == nil {
//		output = log2.Writer()
//	}
//
//	lr := logrus.New()
//
//	// Log as JSON instead of the default ASCII formatter.
//	lr.SetFormatter(&logrus.TemplateFormatter{
//		TimestampFormat: "2006/01/02 15:04:05",
//	})
//
//	// Only log the warning severity or above.
//	lr.SetLevel(level)
//
//	// Output to stdout instead of the default stderr
//	// Can be any io.Writer, see below for File example
//	lr.SetOutput(output)
//
//	return lr
//}

var defaultFormatter = &logrus.TemplateFormatter{
	TimestampFormat: "2006/01/02 15:04:05",
}

func NewDefaultFormatter() *logrus.TemplateFormatter {
	return defaultFormatter
}

func NewTemplateFormatter(purifier logrus.Purifier) *logrus.TemplateFormatter {
	return &logrus.TemplateFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		Purifier:        purifier,
	}
}
