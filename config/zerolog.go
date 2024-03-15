package config

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

// Glog is a global logger instance. Prefer using injected one
var Glog Logger

// initialize global logger
func init() {
	logger := zerolog.New(os.Stderr)

	Glog.Logger = logger
}

// NewZeroLog returns a new structure logger
func NewZeroLog(logLevel, environment string) (*Logger, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	if logLevel == "" {
		logLevel = "info"
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(os.Stderr).With().Timestamp().Str("hostname", hostname).Logger()

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger = logger.Level(level)
	Glog.Logger = Glog.Logger.Level(level)

	if environment != "development" {
		logger = logger.Hook(SeverityHook{})
		Glog.Logger = Glog.Logger.Hook(SeverityHook{})
	} else {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		Glog.Logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return &Logger{Logger: logger}, nil
}

// Logger wraps zerolog.Logger to add extra features to make use of
// GCP logging infrastructure
type Logger struct {
	zerolog.Logger

	MethodField  string // Method where the log is coming from
	PackageField string // Package is the package where the log is coming from
	PayloadField string // Payload is the payload received from the request
}

// Helper function to add the logger fields to the Logger
func addLoggerFields(l *Logger, event *zerolog.Event) *zerolog.Event {
	if l.MethodField != "" {
		event = event.Str("method", l.MethodField)
	}

	if l.PackageField != "" {
		event = event.Str("package", l.PackageField)
	}

	if l.PayloadField != "" {
		event = event.Str("payload", l.PayloadField)
	}

	return event
}

// InfoGCP extends zerolog.Info() to put the error.String() inside message field,
// that could include the trace
func (l *Logger) InfoGCP(err error) {
	event := Glog.Info().Str("message", fmt.Sprintf("%s\n", err.Error()))
	event = addLoggerFields(l, event)
	event.Msg("")
}

// WarnGCP extends zerolog.Warn() to put the error.String() inside message field,
// that could include the trace
func (l *Logger) WarnGCP(err error) {
	stack := make([]byte, 1<<16)
	stackSize := runtime.Stack(stack, false)
	event := Glog.Warn().Str("message", fmt.Sprintf("%s\n%s", err.Error(), (stack[:stackSize])))
	event = addLoggerFields(l, event)
	event.Msg("")
}

// ErrorGCP extends zerolog.Error() to put the error.String() inside message field,
// that could include the tracep
func (l *Logger) ErrorGCP(err error) {
	stack := make([]byte, 1<<16)
	stackSize := runtime.Stack(stack, false)
	event := Glog.Error().Str("message", fmt.Sprintf("%s\n%s", err.Error(), (stack[:stackSize])))
	event = addLoggerFields(l, event)
	event.Msg("")
}

// FatalGCP extends zerolog.Fatal() to put the error.String() inside message field,
// that could include the trace
func (l *Logger) FatalGCP(err error) {
	stack := make([]byte, 1<<16)
	stackSize := runtime.Stack(stack, false)
	event := Glog.Fatal().Str("message", fmt.Sprintf("%s\n%s", err.Error(), (stack[:stackSize])))
	event = addLoggerFields(l, event)
	event.Msg("")
}

// PanicGCP extends zerolog.Panic() to put the error.String() inside message field,
// that could include the trace
func (l *Logger) PanicGCP(err error) {
	stack := make([]byte, 1<<16)
	stackSize := runtime.Stack(stack, false)
	event := Glog.Panic().Str("message", fmt.Sprintf("%s\n%s", err.Error(), (stack[:stackSize])))
	event = addLoggerFields(l, event)
	event.Msg("")
}

// SeverityHook is a hook that will set the severity field in the logs.
// This field is used by gcloud logging to determine the log level
type SeverityHook struct{}

// Run sets the severity hook
func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", toGcloudLevel(level))
	}
}

func toGcloudLevel(l zerolog.Level) string {
	switch l {
	case zerolog.DebugLevel:
		return "DEBUG"
	case zerolog.InfoLevel:
		return "INFO"
	case zerolog.WarnLevel:
		return "WARNING"
	case zerolog.ErrorLevel:
		return "ERROR"
	case zerolog.FatalLevel:
		return "CRITICAL"
	case zerolog.PanicLevel:
		return "ALERT"
	case zerolog.NoLevel:
		return "DEFAULT"
	}
	return "DEFAULT"
}

// Method returns a new Event with the info level
func (l *Logger) Method(field string) *Logger {
	l.MethodField = field
	return l
}

// Package returns a new Event with the info level
func (l *Logger) Package(field string) *Logger {
	l.PackageField = field
	return l
}

// Payload returns a new Event with the info level
func (l *Logger) Payload(field string) *Logger {
	l.PayloadField = field
	return l
}
