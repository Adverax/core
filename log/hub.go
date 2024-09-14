package log

type IHub interface {
	Logger
	Errors() Logger
	Warnings() Logger
}

type Hub struct {
	Logger
	errors   Logger
	warnings Logger
}

func (that *Hub) Errors() Logger {
	return that.errors
}

func (that *Hub) Warnings() Logger {
	return that.warnings
}

func NewHub(
	logger Logger,
	warnings Logger,
	errors Logger,
) *Hub {
	return &Hub{
		Logger:   logger,
		warnings: warnings,
		errors:   errors,
	}
}
