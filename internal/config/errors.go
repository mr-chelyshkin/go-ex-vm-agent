package config

import "fmt"

const (
	ErrConfigParse      = "failed to parse config: %s"
	ErrInitializeConfig = "failed to initialize config: %s"
)

func initError(format string, args ...any) error {
	return fmt.Errorf(ErrInitializeConfig, fmt.Sprintf(format, args...))
}

func parseError(format string, args ...any) error {
	return fmt.Errorf(ErrConfigParse, fmt.Sprintf(format, args...))
}
