package hapm

import (
	"errors"
	"fmt"
)

var errHandled = errors.New("[handled]")

// IsHandledError checks whether error is already logged and represented to user.
func IsHandledError(err error) bool {
	return errors.Is(err, errHandled)
}

// HandledError wraps error with marker for already-logged errors.
func HandledError(err error) error {
	if err == nil {
		return errHandled
	}
	return fmt.Errorf("%w %s", errHandled, err)
}

func (a *App) handledError(action string, err error) error {
	if err == nil {
		return HandledError(errors.New(action))
	}
	a.reporter.Exception(action, err)
	return HandledError(err)
}

func (a *App) handledMessage(message string) error {
	a.reporter.Error(message)
	return HandledError(errors.New(message))
}
