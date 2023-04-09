package homework

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var Min = "min"
var Max = "max"
var Len = "len"
var In = "in"

type Options struct {
	// Numeric map to store 'min', 'max' and 'len' options
	Numeric map[string]int
	// InStr slice of string values in 'in' option
	// also usable for printing values from 'in' values in case of error
	InStr []string
	// InInt slice of integers, if 'in' option is applied to an integer
	InInt []int
}

// getOption parses tag string to get value after option
func (*Options) getOption(str string, opt string) (string, error) {
	var s string
	_, err := fmt.Sscanf(str, opt+":%s", &s)
	if err != nil {
		return "", ErrInvalidValidatorSyntax
	}
	return s, nil
}

// parseNumericOption parses 'len', 'max', 'min' options
func (o *Options) parseNumericOption(str string, opt string) error {
	if strings.Contains(str, opt) {
		s, err := o.getOption(str, opt)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		o.Numeric[opt] = n
	}
	return nil
}

// parseInOption creates only slice InStr if 't' is a string
// or slices InStr and InInt, if 't' is an integer
func (o *Options) parseInOption(t any, st string) error {
	if strings.Contains(st, "in") && !strings.Contains(st, "min") {
		s, err := o.getOption(st, "in")
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		o.InStr = strings.Split(s, ",")
		if reflect.TypeOf(t).Kind() == reflect.Int {
			o.InInt = make([]int, len(o.InStr))
			for i, v := range o.InStr {
				n, err := strconv.Atoi(v)
				if err != nil {
					return ErrInvalidValidatorSyntax
				}
				o.InInt[i] = n
			}
		}
	}
	return nil
}

// ParseOptions parses tag string to retrieve constraint options
func ParseOptions[T int | string](st string) (Options, error) {
	numerical := []string{Min, Max, Len}
	opts := Options{Numeric: make(map[string]int)}
	for _, opt := range numerical {
		err := opts.parseNumericOption(st, opt)
		if err != nil {
			return Options{}, err
		}
	}
	var t T
	err := opts.parseInOption(t, st)
	if err != nil {
		return Options{}, err
	}
	return opts, nil
}
