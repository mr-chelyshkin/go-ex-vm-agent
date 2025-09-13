package config

import "fmt"

const (
	ErrConfigParse      = "failed to parse config: %s"
	ErrInitializeConfig = "failed to initialize config: %s"
)

func initError(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf(ErrInitializeConfig, message)
}

func parseError(format string, args ...any) error {
	message := fmt.Sprintf(format, args...)
	return fmt.Errorf(ErrConfigParse, message)
}
