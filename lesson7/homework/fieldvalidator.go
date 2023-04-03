package homework

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type FieldValidator interface {
	ValidateField(sf reflect.StructField, val reflect.Value)
}

type Validator struct {
	Errors ValidationErrors
}

func (v *Validator) ValidateField(sf reflect.StructField, val reflect.Value) {
	if sf.Tag.Get("validate") == "" {
		return
	}
	if !sf.IsExported() {
		v.Errors = append(v.Errors, ValidationError{ErrValidateForUnexportedFields})
		return
	}
	if sf.Type.Kind() == reflect.Slice {
		v.validateSlice(sf, val)
		return
	}
	opts := v.getFieldOptions(sf.Type.Kind(), sf.Tag.Get("validate"))
	v.validateValue(val, sf, opts)
}

func (v *Validator) getFieldOptions(kind reflect.Kind, tag string) Options {
	var (
		opts Options
		err  error
	)
	switch kind {
	case reflect.String:
		opts, err = ParseOptions[string](tag)
	case reflect.Int:
		opts, err = ParseOptions[int](tag)
	}
	if err != nil {
		v.Errors = append(v.Errors, ValidationError{err})
	}
	return opts
}

func (v *Validator) validateValue(val reflect.Value, sf reflect.StructField, opts Options) {
	v.validateIn(val, sf, opts)
	v.validateNumeric(val, sf, opts)
}

func (v *Validator) validateSlice(sf reflect.StructField, sl reflect.Value) {
	opts := v.getFieldOptions(sf.Type.Elem().Kind(), sf.Tag.Get("validate"))
	for i := 0; i < sl.Len(); i++ {
		v.validateValue(sl.Index(i), sf, opts)
	}
}

func (v *Validator) validateIn(val reflect.Value, sf reflect.StructField, opts Options) {
	var errStrIn = "field '%s' is not valid: '%s' constraint expected %s from set {%s}, but got %v"
	switch val.Kind() {
	case reflect.Int:
		if opts.InInt != nil && !contains(opts.InInt, int(val.Int())) {
			v.Errors = append(v.Errors, ValidationError{
				errors.Errorf(errStrIn, sf.Name, In, "integer",
					strings.Join(opts.InStr, ","), int(val.Int()))})
		}
	case reflect.String:
		if opts.InStr != nil && !contains(opts.InStr, val.String()) {
			v.Errors = append(v.Errors, ValidationError{
				errors.Errorf(errStrIn, sf.Name, In, "string",
					strings.Join(opts.InStr, ","), val.String())})
		}
	}
}

func (v *Validator) validateNumeric(val reflect.Value, sf reflect.StructField, opts Options) {
	var errStr = "field '%s' is not valid: '%s' constraint expected %s %s= %d, but got %d"
	n, mes := v.getNumericValueAndMes(val)
	for k, l := range opts.Numeric {
		switch {
		case k == Min && n < l:
			v.Errors = append(v.Errors,
				ValidationError{errors.Errorf(errStr, sf.Name, k, mes, ">", l, n)})
		case k == Max && n > l:
			v.Errors = append(v.Errors,
				ValidationError{errors.Errorf(errStr, sf.Name, k, mes, "<", l, n)})
		case k == Len && n != l:
			v.Errors = append(v.Errors,
				ValidationError{errors.Errorf(errStr, sf.Name, k, mes, "=", l, n)})
		}
	}
}

func (v *Validator) getNumericValueAndMes(val reflect.Value) (int, string) {
	switch val.Kind() {
	case reflect.Int:
		return int(val.Int()), "int"
	case reflect.String:
		return len(val.String()), "len(string)"
	}
	return 0, ""
}

func contains[T comparable](set []T, val T) bool {
	for _, v := range set {
		if v == val {
			return true
		}
	}
	return false
}
