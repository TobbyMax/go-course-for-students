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
	Numeric map[string]int
	InStr   []string
	InInt   []int
}

func (*Options) getOption(str string, opt string) (string, error) {
	var s string
	_, err := fmt.Sscanf(str, opt+":%s", &s)
	if err != nil {
		return "", ErrInvalidValidatorSyntax
	}
	return s, nil
}

func (o *Options) getOptionNum(str string, opt string) error {
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

func ParseOptions[T int | string](st string) (Options, error) {
	numerical := []string{Min, Max, Len}
	opts := Options{Numeric: make(map[string]int)}
	for _, opt := range numerical {
		err := opts.getOptionNum(st, opt)
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
