package auth

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Validate(i interface{}) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate.Struct(i)
}
