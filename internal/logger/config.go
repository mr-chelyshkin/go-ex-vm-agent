package logger

import (
	"errors"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

func defaultConfig() Config {
	return Config{
		Level:      LevelInfo,
		Format:     FormatJSON,
		Output:     OutputStdout,
		Path:       "",
		MaxBackups: 3,
		MaxAge:     28,
		MaxSize:    100,
		Compress:   false,
	}
}

type Config struct {
	Level      LogLevel  `mapstructure:"level" validate:"required,log_level"`
	Format     LogFormat `mapstructure:"format" validate:"required,log_format"`
	Output     LogOutput `mapstructure:"output" validate:"required,log_output"`
	Path       string    `mapstructure:"path" validate:"required_if=Output file"`
	MaxBackups int       `mapstructure:"max_backups" validate:"omitempty,min=0,max=100"`
	MaxAge     int       `mapstructure:"max_age" validate:"omitempty,min=1,max=365"`
	MaxSize    int       `mapstructure:"max_size" validate:"omitempty,min=1,max=512"`
	Compress   bool      `mapstructure:"compress"`
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

	if c.Level == "" {
		c.Level = defaults.Level
	}
	if c.Format == "" {
		c.Format = defaults.Format
	}
	if c.Output == "" {
		c.Output = defaults.Output
	}
	if c.MaxAge == 0 {
		c.MaxAge = defaults.MaxAge
	}
	if c.MaxSize == 0 {
		c.MaxSize = defaults.MaxSize
	}
	if c.MaxBackups == 0 {
		c.MaxBackups = defaults.MaxBackups
	}
}

func (c *Config) formatValidationErr(err error) error {
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			switch fieldError.Tag() {
			case "log_level":
				return validateError("invalid log level '%s', allowed values: %s", c.Level, logLevelsString())
			case "log_format":
				return validateError("invalid log format '%s', allowed values: %s", c.Format, logFormatsString())
			case "valid_log_output":
				return validateError("invalid log output '%s', allowed values: %s", c.Output, logOutputsString())
			case "required_if":
				return validateError("path is required when output is 'file'")
			default:
				return validateError("validation failed for field '%s': %s", fieldError.Field(), fieldError.Tag())
			}
		}
	}
	return validateError("validation failed: %v", err)
}

func (c *Config) validateRules() error {
	if c.Output == OutputFile {
		if filepath.Ext(c.Path) == "" {
			return validateError("path must include filename, not just directory")
		}
	}
	return nil
}
