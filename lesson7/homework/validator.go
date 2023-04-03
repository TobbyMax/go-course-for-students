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
			res = res + ", "
		}
	}
	return res
}

type Validator struct {
	Errors ValidationErrors
}

func contains[T comparable](set []T, val T) bool {
	for _, v := range set {
		if v == val {
			return true
		}
	}
	return false
}

func (v *Validator) ValidateField(sf reflect.StructField, val reflect.Value) ValidationErrors {
	if sf.Tag.Get("validate") == "" {
		return nil
	}
	if !sf.IsExported() {
		return append(ValidationErrors{}, ValidationError{ErrValidateForUnexportedFields})
	}
	switch sf.Type.Kind() {
	case reflect.String:
		opts, err := ParseOptions[string](sf.Tag.Get("validate"))
		if err != nil {
			return append(ValidationErrors{}, ValidationError{err})
		}
		v.validateValue(val, opts)
	case reflect.Int:
		opts, err := ParseOptions[int](sf.Tag.Get("validate"))
		if err != nil {
			return append(ValidationErrors{}, ValidationError{err})
		}
		v.validateValue(val, opts)
		//case reflect.Slice:
		//	err := v.validateSlice(sf, val)
	}
	return nil
}

//func (v* Validator) validateSlice(sf reflect.StructField, val reflect.Value) error {
//	switch sf.Type.Kind() {}
//}

func (v *Validator) validateValue(val reflect.Value, opts Options) {
	var n int
	if val.Kind() == reflect.Int {
		n = int(val.Int())
		if opts.InInt != nil && !contains(opts.InInt, n) {
			v.Errors = append(v.Errors, ValidationError{errors.New("field not valid")})
		}
	} else if val.Kind() == reflect.String {
		n = len(val.String())
		if opts.InStr != nil && !contains(opts.InStr, val.String()) {
			v.Errors = append(v.Errors, ValidationError{errors.New("field not valid")})
		}
	}
	for k, l := range opts.Numeric {
		switch {
		case k == "min" && n < l:
			v.Errors = append(v.Errors, ValidationError{errors.New("field not valid")})
		case k == "max" && n > l:
			v.Errors = append(v.Errors, ValidationError{errors.New("field not valid")})
		case k == "len" && n != l:
			v.Errors = append(v.Errors, ValidationError{errors.New("field not valid")})
		}
	}
}

func Validate(v any) error {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	vld := Validator{}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < t.NumField(); i++ {
		err := vld.ValidateField(t.Field(i), val.Field(i))
		if err != nil {
			return err
		}
	}
	if len(vld.Errors) != 0 {
		return vld.Errors
	}
	return nil
}
