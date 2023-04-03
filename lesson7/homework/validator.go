package homework

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for i, ve := range v {
		res = res + fmt.Sprint(ve.Err)
		if i != len(v)-1 {
			res = res + "; "
		}
	}
	return res
}

func Validate(v any) error {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	vld := Validator{}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < t.NumField(); i++ {
		vld.ValidateField(t.Field(i), val.Field(i))
	}
	if len(vld.Errors) != 0 {
		return vld.Errors
	}
	return nil
}
