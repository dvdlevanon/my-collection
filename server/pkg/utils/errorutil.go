package utils

import "github.com/go-errors/errors"

func LogError(err error) {
	if err == nil {
		return
	}

	var e *errors.Error
	if errors.As(err, &e) {
		logger.Errorf("Error: %v", e.ErrorStack())
	} else {
		logger.Errorf("Error: %v", err)
	}
}
