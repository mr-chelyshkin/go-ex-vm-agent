package runner

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	maxRestartDelay = 10 * time.Minute
)

// defaultConfig returns a Config object with predefined default values for all configurable parameters.
func defaultConfig() Config {
	return Config{
		ShutdownTimeout:    60 * time.Second,
		RestartDelay:       10 * time.Second,
		MaxRestarts:        0,
		EnableRestart:      false,
		ExponentialBackoff: false,
	}
}

// Config defines configuration options for managing shutdown behavior, restart attempts, and delay strategies.
type Config struct {
	// ShutdownTimeout specifies the maximum time duration to wait for graceful shutdown. Valid range: 1s to 5m.
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"min=1s,max=5m"`

	// RestartDelay defines the duration to wait before attempting a restart, validated between 1s and 1m.
	RestartDelay time.Duration `mapstructure:"restart_delay" validate:"min=1s,max=1m"`

	// MaxRestarts specifies the maximum number of restart attempts. A value of 0 allows infinite restarts.
	MaxRestarts int `mapstructure:"max_restarts" validate:"min=0,max=100"`

	// EnableRestart determines whether the restart mechanism for stopped or failed processes is enabled or disabled.
	EnableRestart bool `mapstructure:"enable_restart"`

	// ExponentialBackoff determines whether to apply exponential delay strategy for restarts.
	ExponentialBackoff bool `mapstructure:"exponential_backoff"`
}

// Validate validates the Config object to ensure all fields comply with defined constraints and sets default values.
func (c *Config) Validate() error {
	c.setDefaults()

	if err := getValidator().Struct(c); err != nil {
		return c.formatValidationErr(err)
	}
	return c.validateRules()
}

// GetRestartDelay calculates and returns the restart delay based on the attempt number and exponential backoff settings.
func (c *Config) GetRestartDelay(attempt int) time.Duration {
	if !c.ExponentialBackoff {
		return c.RestartDelay
	}
	delay := c.RestartDelay * time.Duration(1<<attempt)

	if delay > maxRestartDelay {
		delay = maxRestartDelay
	}
	return delay
}

// setDefaults sets default values for Config fields if they are not already specified.
func (c *Config) setDefaults() {
	defaults := defaultConfig()

	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = defaults.ShutdownTimeout
	}
	if c.RestartDelay == 0 {
		c.RestartDelay = defaults.RestartDelay
	}
}

// formatValidationErr processes validation errors for Config fields and returns detailed error messages.
func (c *Config) formatValidationErr(err error) error {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "ShutdownTimeout":
				return initError("shutdown timeout must be between 1s and 5m, got: %v", c.ShutdownTimeout)
			case "RestartDelay":
				return initError("restart delay must be between 1s and 1m, got: %v", c.RestartDelay)
			case "MaxRestarts":
				return initError("max restarts must be between 0 and 100, got: %d", c.MaxRestarts)
			default:
				return initError("validation failed for field '%s': %s", fieldError.Field(), fieldError.Tag())
			}
		}
	}
	return initError("validation failed: %v", err)
}

// validateRules performs additional custom validation logic for the Config struct to ensure its integrity.
func (c *Config) validateRules() error {
	return nil
}
