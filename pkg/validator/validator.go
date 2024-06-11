package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func Validate(model interface{}) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("json")
	})

	return validate.Struct(model)
}
