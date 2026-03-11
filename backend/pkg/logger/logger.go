package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger is the interface that all loggers must implement.
// This allows for dependency injection and easier testing.
type Logger interface {
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
}

// Event represents a log event that can be enriched with fields.
type Event interface {
	Str(key, val string) Event
	Int(key string, val int) Event
	Float64(key string, val float64) Event
	Dur(key string, val time.Duration) Event
	Err(err error) Event
	Msg(msg string)
}

// ZerologAdapter adapts zerolog to our Logger interface.
type ZerologAdapter struct {
	logger zerolog.Logger
}

// NewZerologAdapter creates a new logger instance with the specified environment configuration.
func NewZerologAdapter(env string) Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var zlog zerolog.Logger

	if env == "development" || env == "dev" || env == "" {
		// Pretty console output for development
		zlog = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		// JSON output for production
		zlog = zerolog.New(os.Stdout).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &ZerologAdapter{logger: zlog}
}

// zerologEvent wraps zerolog.Event to implement our Event interface.
type zerologEvent struct {
	event *zerolog.Event
}

func (z *ZerologAdapter) Debug() Event {
	return &zerologEvent{event: z.logger.Debug()}
}

func (z *ZerologAdapter) Info() Event {
	return &zerologEvent{event: z.logger.Info()}
}

func (z *ZerologAdapter) Warn() Event {
	return &zerologEvent{event: z.logger.Warn()}
}

func (z *ZerologAdapter) Error() Event {
	return &zerologEvent{event: z.logger.Error()}
}

func (z *ZerologAdapter) Fatal() Event {
	return &zerologEvent{event: z.logger.Fatal()}
}

func (e *zerologEvent) Str(key, val string) Event {
	e.event = e.event.Str(key, val)
	return e
}

func (e *zerologEvent) Int(key string, val int) Event {
	e.event = e.event.Int(key, val)
	return e
}

func (e *zerologEvent) Float64(key string, val float64) Event {
	e.event = e.event.Float64(key, val)
	return e
}

func (e *zerologEvent) Dur(key string, val time.Duration) Event {
	e.event = e.event.Dur(key, val)
	return e
}

func (e *zerologEvent) Err(err error) Event {
	e.event = e.event.Err(err)
	return e
}

func (e *zerologEvent) Msg(msg string) {
	e.event.Msg(msg)
}
