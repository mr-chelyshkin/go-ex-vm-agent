package worker

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

func defaultConfig() Config {
	return Config{
		TaskStopTimeout: 10 * time.Second,
		TaskTimeout:     5 * time.Minute,
		MaxTasks:        100,
		StopOnError:     false,
	}
}

type Config struct {
	// TaskStopTimeout максимальное время ожидания остановки одной задачи
	TaskStopTimeout time.Duration `mapstructure:"task_stop_timeout" validate:"min=1s,max=2m"`
	// TaskTimeout максимальное время выполнения задачи (0 = без лимита)
	TaskTimeout time.Duration `mapstructure:"task_timeout" validate:"min=0"`
	// MaxTasks максимальное количество одновременно работающих задач
	MaxTasks int `mapstructure:"max_tasks" validate:"min=1,max=1000"`
	// StopOnError останавливать ли воркер при ошибке в задаче
	StopOnError bool `mapstructure:"stop_on_error"`
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

	if c.TaskStopTimeout == 0 {
		c.TaskStopTimeout = defaults.TaskStopTimeout
	}
	if c.TaskTimeout == 0 {
		c.TaskTimeout = defaults.TaskTimeout
	}
	if c.MaxTasks == 0 {
		c.MaxTasks = defaults.MaxTasks
	}
}

func (c *Config) formatValidationErr(err error) error {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			switch fieldError.Field() {
			case "TaskStopTimeout":
				return initError("task stop timeout must be between 1s and 2m, got: %v", c.TaskStopTimeout)
			case "TaskTimeout":
				return initError("task timeout must be non-negative, got: %v", c.TaskTimeout)
			case "MaxTasks":
				return initError("max tasks must be between 1 and 1000, got: %d", c.MaxTasks)
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
