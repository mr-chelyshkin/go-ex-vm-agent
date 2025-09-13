package logger

import "fmt"

const (
	// ErrConfigValidation is a constant representing an error message format for configuration validation errors.
	ErrConfigValidation = "config validation error: %s"

	// ErrInitializeLogger indicates an error occurred during the initialization of the logger.
	ErrInitializeLogger = "failed to initialize logger: %s"
)

func validateError(format string, args ...any) error {
	return fmt.Errorf(ErrConfigValidation, fmt.Sprintf(format, args...))
}

func initError(format string, args ...any) error {
	return fmt.Errorf(ErrInitializeLogger, fmt.Sprintf(format, args...))
}
