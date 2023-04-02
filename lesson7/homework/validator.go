package homework

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	// TODO: implement this
	res := ""
	for i, ve := range v {
		res = res + fmt.Sprint(ve.Err)
		if i != len(v)-1 {
			res = res + ", "
		}
	}
	return res
}

//func validateSlice(v any, v ValidationErrors) {
//	switch t.Field(i).Type.Elem().Kind(){
//	case reflect.String:
//		validateString()
//	case reflect.Int:
//		validateInt()
//	}
//}

func validateString(val reflect.Value, st string, ve *ValidationErrors) error {
	str := val.String()
	if strings.Contains(st, "min") {
		var s string
		_, err := fmt.Sscanf(st, "min:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		min, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}

		if len(str) < min {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	if strings.Contains(st, "max") {
		var s string
		_, err := fmt.Sscanf(st, "max:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		max, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if len(str) > max {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	if strings.Contains(st, "len") {
		var s string
		_, err := fmt.Sscanf(st, "len:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		l, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if len(str) != l {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	if strings.Contains(st, "in") && !strings.Contains(st, "min") {
		var s string
		_, err := fmt.Sscanf(st, "in:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		set := strings.Split(s, ",")
		for i, v := range set {
			set[i] = v
		}
		if !contains(set, str) {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	return nil
}

func validateInt(val reflect.Value, st string, ve *ValidationErrors) error {
	num := val.Int()
	if strings.Contains(st, "min") {
		var s string
		_, err := fmt.Sscanf(st, "min:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		min, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if num < int64(min) {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	if strings.Contains(st, "max") {
		var s string
		_, err := fmt.Sscanf(st, "max:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		max, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if num > int64(max) {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	if strings.Contains(st, "len") {
		l := 0
		_, err := fmt.Sscanf(st, "len:%d", &l)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		*ve = append(*ve, ValidationError{errors.New("cannot use len-validation for int value")})
	}
	if strings.Contains(st, "in") && !strings.Contains(st, "min") {
		var s string
		_, err := fmt.Sscanf(st, "in:%s", &s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		set := strings.Split(s, ",")
		iSet := make([]int, len(set))
		for i, v := range set {
			n, err := strconv.Atoi(v)
			if err != nil {
				return ErrInvalidValidatorSyntax
			}
			iSet[i] = n
		}
		if !contains(iSet, int(num)) {
			*ve = append(*ve, ValidationError{errors.New("field not valid")})
		}
	}
	return nil
}

func contains[T comparable](set []T, val T) bool {
	for _, v := range set {
		if v == val {
			return true
		}
	}
	return false
}

func Validate(v any) error {
	// TODO: implement this
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	ve := ValidationErrors{}
	if t.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("validate") == "" {
			continue
		}
		if !t.Field(i).IsExported() {
			return append(ValidationErrors{}, ValidationError{ErrValidateForUnexportedFields})
		}
		switch t.Field(i).Type.Kind() {
		case reflect.String:
			err := validateString(val.Field(i), t.Field(i).Tag.Get("validate"), &ve)
			if err != nil {
				ve = append(ve, ValidationError{err})
			}
		case reflect.Int:
			err := validateInt(val.Field(i), t.Field(i).Tag.Get("validate"), &ve)
			if err != nil {
				ve = append(ve, ValidationError{err})
			}
			//case reflect.Slice:
			//	validateSlice();
		}
	}
	if len(ve) != 0 {
		return ve
	}
	return nil
}
