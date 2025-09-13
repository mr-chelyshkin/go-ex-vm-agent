package logger

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func getValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()

		_ = validate.RegisterValidation("log_level", func(fl validator.FieldLevel) bool {
			level := LogLevel(fl.Field().String())
			return level.IsValid()
		})

		_ = validate.RegisterValidation("log_format", func(fl validator.FieldLevel) bool {
			format := LogFormat(fl.Field().String())
			return format.IsValid()
		})

		_ = validate.RegisterValidation("log_output", func(fl validator.FieldLevel) bool {
			output := LogOutput(fl.Field().String())
			return output.IsValid()
		})
	})
	return validate
}
