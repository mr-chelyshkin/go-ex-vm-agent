package config

import "fmt"

const (
	// ErrConfigParse represents an error message format for failures encountered during configuration parsing.
	ErrConfigParse = "failed to parse config: %s"

	// ErrInitializeConfig represents an error message format for failures during the initialization of a configuration.
	ErrInitializeConfig = "failed to initialize config: %s"
)

func initError(format string, args ...any) error {
	return fmt.Errorf(ErrInitializeConfig, fmt.Sprintf(format, args...))
}

func parseError(format string, args ...any) error {
	return fmt.Errorf(ErrConfigParse, fmt.Sprintf(format, args...))
}
