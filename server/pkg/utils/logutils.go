package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-errors/errors"
	"github.com/op/go-logging"
)

type contextKey string

const subjectContextKey contextKey = "subject"

func ConfigureLogger() error {
	logFormat := `[%{time:2006-01-02 15:04:05.000}] %{color}%{level:-7s}%{color:reset} %{message} [%{module} - %{shortfile}]`
	formatter, err := logging.NewStringFormatter(logFormat)
	if err != nil {
		return err
	}

	logging.SetBackend(logging.NewLogBackend(os.Stdout, "", 0))
	logging.SetFormatter(formatter)

	logger.Debugf("Logger initialized with format %v", logFormat)
	return nil
}

func LogError(message string, err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		logger.Errorf("Error: %s %v", message, e.ErrorStack())
	} else {
		logger.Errorf("Error: %s %s", message, squashString(fmt.Sprintf("%v", err)))
	}
}

func LogWarning(message string, err error) {
	if err == nil {
		return
	}

	logger.Warningf("Warning: %s %s", message, squashString(fmt.Sprintf("%v", err)))
}

func squashString(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), "\r", "")
}

func ContextWithSubject(parent context.Context, subject string) context.Context {
	return context.WithValue(parent, subjectContextKey, subject)
}

func GetSubject(ctx context.Context) string {
	return fmt.Sprintf("%s", ctx.Value(subjectContextKey))
}
