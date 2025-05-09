package util

import (
	"github.com/go-playground/validator/v10"
)

var instance structValidator

type structValidator struct {
	validator *validator.Validate
}

func init() {
	instance.validator = validator.New()
}

func Validate[T any](t T) error {
	if err := instance.validator.Struct(t); err != nil {
		return err
	}

	return nil
}
