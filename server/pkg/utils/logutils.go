package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

type contextKey string

const subjectContextKey contextKey = "subject"
const counterContextKey contextKey = "counter"

var (
	subjectCounters = sync.Map{} // map[string]*atomic.Uint64
)

func getOrCreateCounter(subject string) *atomic.Uint64 {
	if counter, ok := subjectCounters.Load(subject); ok {
		return counter.(*atomic.Uint64)
	}

	counter := &atomic.Uint64{}
	actual, loaded := subjectCounters.LoadOrStore(subject, counter)
	if loaded {
		return actual.(*atomic.Uint64)
	}
	return counter
}

func MustGetLogger(module string) *ContextLogger {
	return &ContextLogger{
		logger: logging.MustGetLogger(module),
	}
}

type ContextLogger struct {
	logger *logging.Logger
}

func appendContext(ctx context.Context, format string) string {
	subject := GetSubject(ctx)
	counter := getCounter(ctx)

	// Calculate padding: 25 chars total, minus subject length, minus counter digits, minus '#'
	counterStr := fmt.Sprintf("#%d", counter)
	padding := 25 - len(subject) - len(counterStr)
	if padding < 0 {
		padding = 0
	}

	subjectWithCounter := fmt.Sprintf("%s%s%s", subject, strings.Repeat(" ", padding), counterStr)
	return fmt.Sprintf("[%s] %s", subjectWithCounter, format)
}

func (l *ContextLogger) Warningf(ctx context.Context, format string, args ...any) {
	l.logger.Warningf(appendContext(ctx, format), args...)
}

func (l *ContextLogger) Errorf(ctx context.Context, format string, args ...any) {
	l.logger.Errorf(appendContext(ctx, format), args...)
}

func (l *ContextLogger) Debugf(ctx context.Context, format string, args ...any) {
	l.logger.Debugf(appendContext(ctx, format), args...)
}

func (l *ContextLogger) Infof(ctx context.Context, format string, args ...any) {
	l.logger.Infof(appendContext(ctx, format), args...)
}

func ConfigureLogger(ctx context.Context) error {
	logFormat := `[%{time:2006-01-02 15:04:05.000}] %{color}%{level:-7s}%{color:reset} %{message} [%{module} - %{shortfile}]`
	formatter, err := logging.NewStringFormatter(logFormat)
	if err != nil {
		return err
	}

	fileLogFormat := `[%{time:2006-01-02 15:04:05.000}] %{level:-7s} %{message} [%{module} - %{shortfile}]`
	fileFormatter, err := logging.NewStringFormatter(fileLogFormat)
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile("mycollection.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	consoleBackend := logging.NewLogBackend(os.Stdout, "", 0)
	consoleFormatted := logging.NewBackendFormatter(consoleBackend, formatter)

	fileBackend := logging.NewLogBackend(logFile, "", 0)
	fileFormatted := logging.NewBackendFormatter(fileBackend, fileFormatter)

	logging.SetBackend(consoleFormatted, fileFormatted)
	logger.Debugf(ctx, "Logger initialized with format %v", logFormat)
	return nil
}

func LogError(ctx context.Context, logger *ContextLogger, message string, err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		logger.Errorf(ctx, "Error: %s %v", message, e.ErrorStack())
	} else {
		logger.Errorf(ctx, "Error: %s %s", message, squashString(fmt.Sprintf("%v", err)))
	}
}

func LogWarning(ctx context.Context, logger *ContextLogger, message string, err error) {
	if err == nil {
		return
	}

	logger.Warningf(ctx, "Warning: %s %s", message, squashString(fmt.Sprintf("%v", err)))
}

func squashString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "\r", "")
}

func ContextWithSubject(parent context.Context, subject string) context.Context {
	counter := getOrCreateCounter(subject)
	count := counter.Add(1)

	ctx := context.WithValue(parent, subjectContextKey, subject)
	ctx = context.WithValue(ctx, counterContextKey, count)
	return ctx
}

func HasSubject(ctx context.Context) bool {
	return ctx.Value(subjectContextKey) != nil
}

func GetSubject(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value(subjectContextKey))
}

func getCounter(ctx context.Context) uint64 {
	if count, ok := ctx.Value(counterContextKey).(uint64); ok {
		return count
	}
	return 0
}
