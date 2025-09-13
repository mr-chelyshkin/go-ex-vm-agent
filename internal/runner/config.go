package runner

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

func defaultConfig() Config {
	return Config{
		ShutdownTimeout: 30 * time.Second,
		RestartDelay:    5 * time.Second,
		MaxRestarts:     3,
		EnableRestart:   true,
	}
}

type Config struct {
	// ShutdownTimeout максимальное время ожидания graceful shutdown
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"min=1s,max=5m"`
	// RestartDelay задержка перед рестартом worker'а
	RestartDelay time.Duration `mapstructure:"restart_delay" validate:"min=0s,max=1m"`
	// MaxRestarts максимальное количество попыток рестарта (0 = без лимита)
	MaxRestarts int `mapstructure:"max_restarts" validate:"min=0,max=100"`
	// EnableRestart включает автоматический рестарт при падении worker'а
	EnableRestart bool `mapstructure:"enable_restart"`
}

func (c *Config) Validate() error {
	c.setDefaults()

	if err := getValidator().Struct(c); err != nil {
		return c.formatValidationErr(err)
	}
	return c.validateRules()
}

func (c *Config) setDefaults() {
	defaults := defaultConfig()

	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = defaults.ShutdownTimeout
	}
	if c.RestartDelay == 0 {
		c.RestartDelay = defaults.RestartDelay
	}
	if c.MaxRestarts == 0 {
		c.MaxRestarts = defaults.MaxRestarts
	}
}

func (c *Config) formatValidationErr(err error) error {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "ShutdownTimeout":
				return initError("shutdown timeout must be between 1s and 5m, got: %v", c.ShutdownTimeout)
			case "RestartDelay":
				return initError("restart delay must be between 0s and 1m, got: %v", c.RestartDelay)
			case "MaxRestarts":
				return initError("max restarts must be between 0 and 100, got: %d", c.MaxRestarts)
			default:
				return initError("validation failed for field '%s': %s", fieldError.Field(), fieldError.Tag())
			}
		}
	}
	return initError("validation failed: %v", err)
}

func (c *Config) validateRules() error {
	return nil
}
