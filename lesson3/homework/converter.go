package main

import (
	"bytes"
	"errors"
	"strings"
	"unicode"
)

type Transformer interface {
	Transform([]byte) []byte
}

type Trimmer interface {
	Trim(b []byte) []byte
}

type TailAdder interface {
	AddTail(b []byte) []byte
}

type Converter interface {
	Transformer
	Trimmer
	TailAdder
}

type DDConverter struct {
	trim  bool
	upper bool
	lower bool
	// Attributes for trim function
	start     bool
	bufSpaces []byte
	tail      []byte
}

func (ddc *DDConverter) setConvStatus(conversions *StringSlice) error {
	var (
		trim, upper, lower       = false, false, false
		err                error = nil
	)
	if len(*conversions) > 2 {
		err = errors.New("too many arguments")
		return err
	}
	for _, key := range *conversions {
		switch key {
		case "trim_spaces":
			trim = true
		case "upper_case":
			upper = true
		case "lower_case":
			lower = true
		default:
			err = errors.New("invalid conversion")
			return err
		}
	}
	if lower && upper {
		err = errors.New("can not apply 'upper_case' and 'lower_case' simultaneously")
		return err
	}
	ddc.trim, ddc.upper, ddc.lower = trim, upper, lower
	return err
}

func NewDDConverter(conversions *StringSlice) (*DDConverter, error) {
	var ddc = DDConverter{}
	ddc.start = true
	err := ddc.setConvStatus(conversions)
	if err != nil {
		return nil, err
	}
	return &ddc, nil
}

func (ddc *DDConverter) Transform(b []byte) []byte {
	converted := b
	switch {
	case ddc.upper:
		converted = bytes.ToUpper(converted)
	case ddc.lower:
		converted = bytes.ToLower(converted)
	}
	return converted
}

func (ddc *DDConverter) Trim(b []byte) []byte {
	strTrimmed := string(b)
	if ddc.trim && ddc.start {
		strTrimmed = strings.TrimLeftFunc(strTrimmed, unicode.IsSpace)
		b = b[len(b)-len(strTrimmed):]
		if len(strTrimmed) != 0 {
			ddc.start = false
		}
	}
	if ddc.trim && !ddc.start {
		strTrimmed = strings.TrimRightFunc(strTrimmed, unicode.IsSpace)
		if len(strTrimmed) != 0 {
			ddc.tail = ddc.bufSpaces
			ddc.bufSpaces = nil
		}
		if ddc.bufSpaces == nil {
			ddc.bufSpaces = make([]byte, 0)
		}
		ddc.bufSpaces = append(ddc.bufSpaces, b[len(strTrimmed):]...)
		b = b[:len(strTrimmed)]
	}
	return b
}

func (ddc *DDConverter) AddTail(b []byte) []byte {
	if len(b) == 0 {
		return b
	}
	return append(ddc.tail, b...)
}
