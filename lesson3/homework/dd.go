package main

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Executer interface {
	Execute() error
}

type Converter interface {
	Convert([]byte) []byte
}

type DD struct {
	Options
	infile        *DDFile
	outfile       *DDFile
	currentOffset int64
	start         bool
	bufSpaces     []byte
	bufSize       int
	trim          bool
	upper         bool
	lower         bool
}

func (dd *DD) setConvStatus(conversions *StringSlice) error {
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
	dd.trim, dd.upper, dd.lower = trim, upper, lower
	return err
}

func New(options *Options) (*DD, error) {
	if options.Offset < 0 {
		return nil, errors.New("invalid offset")
	}
	var dd = DD{Options: *options, start: true}
	err := dd.setConvStatus(&options.Conv)
	if err != nil {
		return nil, err
	}
	return &dd, nil
}

func (dd *DD) Execute() error {
	var err error

	dd.infile, err = Open(dd.From, dd.BlockSize)
	if err != nil {
		return err
	}
	defer func() {
		if err := dd.infile.Close(); err != nil {
			panic(err)
		}
	}()

	dd.outfile, err = Create(dd.To, dd.BlockSize)
	if err != nil {
		return err
	}
	defer func() {
		if err := dd.outfile.Close(); err != nil {
			panic(err)
		}
	}()

	err = dd.readConvertWrite()
	if err != nil {
		return err
	}
	return nil
}

func (dd *DD) readConvertWrite() error {
	var (
		buf = make([]byte, dd.BlockSize+4)
		err error
	)
	for dd.Limit == -1 || dd.currentOffset < dd.Offset+dd.Limit {
		dd.bufSize, err = dd.read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		if dd.currentOffset < dd.Offset {
			continue
		}
		bufLimited := dd.adjustToLimits(buf)
		bufValid := dd.trimInvalidBytes(bufLimited)
		bufTrimmed, err := dd.trimSpaces(bufValid)
		if err != nil {
			return err
		}
		bufConverted := dd.Convert(bufTrimmed)
		if _, err = dd.outfile.Write(bufConverted); err != nil {
			return err
		}
	}
	if _, err = dd.outfile.Write(dd.infile.carryOver[:dd.infile.carryLen]); err != nil {
		return err
	}
	return nil
}

func (dd *DD) read(buf []byte) (int, error) {
	n, err := dd.infile.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	}
	if n == 0 || err == io.EOF {
		if dd.currentOffset < dd.Offset {
			return 0, errors.New("offset index out of range")
		}
		return 0, io.EOF
	}
	dd.currentOffset += int64(n)
	return n, err
}

func (dd *DD) adjustToLimits(b []byte) []byte {
	bufOffset := b
	if delta := int(dd.currentOffset - dd.Offset); delta < dd.bufSize && dd.currentOffset >= dd.Offset {
		bufOffset = b[dd.bufSize-delta : dd.bufSize]
		dd.bufSize = delta
	}
	if dd.Limit != -1 && dd.currentOffset > dd.Offset+dd.Limit {
		dd.bufSize -= int(dd.currentOffset - (dd.Offset + dd.Limit))
	}
	return bufOffset[:dd.bufSize]
}

func (dd *DD) trimInvalidBytes(b []byte) []byte {
	i := dd.bufSize - 1
	for ; i >= 0; i-- {
		if utf8.RuneStart(b[i]) {
			if !utf8.Valid(b[i:]) {
				dd.infile.carryOver = b[i:]
				b = b[:i]
				dd.infile.carryLen = dd.bufSize - i
				dd.currentOffset -= int64(dd.infile.carryLen)
				dd.bufSize = i
			}
			break
		}
	}
	return b
}

func (dd *DD) trimSpaces(b []byte) ([]byte, error) {
	strTrimmed := string(b[:dd.bufSize])
	if dd.trim && dd.start {
		strTrimmed = strings.TrimLeftFunc(strTrimmed, unicode.IsSpace)
		dd.bufSize = len(strTrimmed)
		if dd.bufSize != 0 {
			dd.start = false
		}
	}
	if dd.trim && !dd.start {
		strTrimmed = strings.TrimRightFunc(strTrimmed, unicode.IsSpace)
		dd.bufSize = len(strTrimmed)
		if dd.bufSize != 0 {
			if _, err := dd.outfile.Write(dd.bufSpaces); err != nil {
				return nil, err
			}
			dd.bufSpaces = nil
		} else {
			if len(dd.bufSpaces) == 0 {
				dd.bufSpaces = make([]byte, 0)
			}
			dd.bufSpaces = append(dd.bufSpaces, b[dd.bufSize:]...)
		}
	}
	return []byte(strTrimmed[:dd.bufSize]), nil
}

func (dd *DD) Convert(b []byte) []byte {
	converted := b
	switch {
	case dd.upper:
		converted = bytes.ToUpper(converted[:dd.bufSize])
	case dd.lower:
		converted = bytes.ToLower(converted[:dd.bufSize])
	}
	return converted
}
