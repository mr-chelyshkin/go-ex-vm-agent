package logger

import "fmt"

const (
	ErrConfigValidation = "config validation error: %s"
	ErrInitializeLogger = "failed to initialize logger: %s"
)

func validateError(format string, args ...any) error {
	return fmt.Errorf(ErrConfigValidation, fmt.Sprintf(format, args...))
}

func initError(format string, args ...any) error {
	return fmt.Errorf(ErrInitializeLogger, fmt.Sprintf(format, args...))
}
