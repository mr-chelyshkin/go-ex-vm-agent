package runner

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

func defaultConfig() Config {
	return Config{
		ShutdownTimeout:    30 * time.Second,
		RestartDelay:       5 * time.Second,
		MaxRestarts:        0, // 0 = бесконечно
		EnableRestart:      true,
		ExponentialBackoff: false,
	}
}

type Config struct {
	// ShutdownTimeout максимальное время ожидания graceful shutdown всего приложения
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"min=1s,max=5m"`
	// RestartDelay базовая задержка перед рестартом worker'а
	RestartDelay time.Duration `mapstructure:"restart_delay" validate:"min=1s,max=1m"`
	// MaxRestarts максимальное количество попыток рестарта (0 = без лимита, бесконечно)
	MaxRestarts int `mapstructure:"max_restarts" validate:"min=0,max=100"`
	// EnableRestart включает автоматический рестарт при падении worker'а
	EnableRestart bool `mapstructure:"enable_restart"`
	// ExponentialBackoff увеличивает RestartDelay экспоненциально при каждом рестарте (2x, 4x, 8x...)
	ExponentialBackoff bool `mapstructure:"exponential_backoff"`
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
	// MaxRestarts может быть 0 (бесконечно), поэтому не устанавливаем дефолт
}

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

func (c *Config) validateRules() error {
	return nil
}

// GetRestartDelay возвращает задержку для указанной попытки рестарта
func (c *Config) GetRestartDelay(attempt int) time.Duration {
	if !c.ExponentialBackoff {
		return c.RestartDelay
	}

	// Экспоненциальное увеличение: 2^attempt
	multiplier := 1 << attempt // 2^attempt
	delay := c.RestartDelay * time.Duration(multiplier)

	// Ограничиваем максимальную задержку 5 минутами
	maxDelay := 5 * time.Minute
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}
