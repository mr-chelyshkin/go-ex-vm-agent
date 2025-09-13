package worker

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
	})
	return validate
}
