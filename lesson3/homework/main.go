package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Options struct {
	From string
	To   string
	// todo: add required flags
	Offset    int64
	Limit     int64
	BlockSize uint64
	Conv      StringSlice
}

type StringSlice []string

func (ss *StringSlice) String() string {
	return strings.Join(*ss, ", ")
}

func (ss *StringSlice) Set(value string) error {
	*ss = strings.Split(value, ",")
	return nil
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "stdin", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "stdout", "file to write. by default - stdout")

	// todo: parse and validate all flags
	flag.Int64Var(&opts.Offset, "offset", 0, "file to write. by default - stdout")
	flag.Int64Var(&opts.Limit, "limit", -1, "file to write. by default - stdout")
	flag.Uint64Var(&opts.BlockSize, "block-size", 1024, "file to write. by default - stdout")
	flag.Var(&opts.Conv, "conv", "file to write. by default - stdout")

	flag.Parse()

	return &opts, nil
}

//type Converter interface {
//	Convert(slice StringSlice)
//}

type DDReader interface {
	io.Reader
	io.Seeker
	io.Writer
	io.Closer
	//Converter
}

type DD struct {
	base    int64
	options Options
}

func Convert(buf []byte, option func([]byte) []byte) []byte {
	return option(buf)
}

func (ss *StringSlice) ToSet() map[string]struct{} {
	var void struct{}
	set := make(map[string]struct{})
	for _, str := range *ss {
		set[str] = void
	}
	return set
}

func DoAll2(options *Options) {
	if options.Offset < 0 {
		panic("Invalid offset")
	}

	if len(options.Conv) > 2 {
		panic("Too many arguments.")
	}
	trim := false
	upper := false
	lower := false
	for _, key := range options.Conv {
		switch key {
		case "trim_spaces":
			trim = true
		case "upper_case":
			upper = true
		case "lower_case":
			lower = true
		default:
			panic("Invalid conversion")
		}
	}
	if lower == true && upper == true {
		panic("Can not apply upper_case and lower_case simultaneously")
	}

	infile := os.Stdin
	var err error = nil
	if options.From != "stdin" {
		infile, err = os.Open(options.From)
		if err != nil {
			panic(err)
		}
	}

	defer func() {
		if !(infile == os.Stdin) {
			if err := infile.Close(); err != nil {
				panic(err)
			}
		}
	}()

	outfile := os.Stdout
	if options.To != "stdout" {
		_, err = os.Stat(options.To)
		if !errors.Is(err, os.ErrNotExist) {
			panic("Outfile already exists")
		}
		outfile, err = os.Create(options.To)
		if err != nil {
			panic(err)
		}
	}

	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()

	var (
		offset    int64 = 0
		start           = true
		bufSpaces []byte
		//lenSpaces int64 = 0
		carryOver = make([]byte, 4)
		carryLen  int
	)
	buf := make([]byte, options.BlockSize+4)
	for options.Limit == -1 || offset < options.Offset+options.Limit {
		var (
			bufSize int
			err     error
		)
		bufSize, err = infile.Read(buf[carryLen : options.BlockSize+uint64(carryLen)])
		if carryLen > 0 {
			copy(buf, carryOver[:carryLen])
		}

		if err != nil && err != io.EOF {
			panic(err)
		}

		if bufSize == 0 || err == io.EOF {
			if offset < options.Offset {
				panic("Offset out of range")
			}
			break
		}

		bufSize += carryLen
		offset += int64(bufSize)
		carryLen = 0
		if offset < options.Offset {
			continue
		}
		bufOffset := buf
		if delta := int(offset - options.Offset); delta < bufSize && offset >= options.Offset {
			bufOffset = buf[bufSize-delta : bufSize]
			bufSize = delta
		}
		if options.Limit != -1 && offset > options.Offset+options.Limit {
			bufSize -= int(offset - (options.Offset + options.Limit))
		}
		bufLimited := bufOffset[:bufSize]
		i := bufSize - 1
		for ; i >= 0; i-- {
			if utf8.RuneStart(bufLimited[i]) {
				if !utf8.Valid(bufLimited[i:]) {
					carryOver = bufLimited[i:]
					bufLimited = bufLimited[:i]
					carryLen = bufSize - i
					offset -= int64(carryLen)
					bufSize = i
				}
				break
			}
		}

		strTrimmed := string(bufLimited[:bufSize])
		if trim && start {
			strTrimmed = strings.TrimLeftFunc(strTrimmed, unicode.IsSpace)
			bufSize = len(strTrimmed)
			if bufSize != 0 {
				start = false
			}
		}
		if trim && !start {
			strTrimmed = strings.TrimRightFunc(strTrimmed, unicode.IsSpace)
			bufSize = len(strTrimmed)
			if bufSize != 0 {
				if _, err = outfile.Write(bufSpaces); err != nil {
					panic(err)
				}
				bufSpaces = nil
			} else {
				if len(bufSpaces) == 0 {
					bufSpaces = make([]byte, 0)
				}
				bufSpaces = append(bufSpaces, bufLimited[bufSize:]...)
			}
		}
		bufConverted := []byte(strTrimmed[:bufSize])
		switch {
		case upper:
			bufConverted = Convert(bufConverted[:bufSize], bytes.ToUpper)
		case lower:
			bufConverted = Convert(bufConverted[:bufSize], bytes.ToLower)
		}
		bsize := uint64(bufSize)
		for i := uint64(0); i < bsize; i += options.BlockSize {
			writeLen := bsize
			if bsize > i+options.BlockSize {
				writeLen = i + options.BlockSize
			}
			if _, err = outfile.Write(bufConverted[i:writeLen]); err != nil {
				panic(err)
			}
		}
	}
	if _, err = outfile.Write(carryOver[:carryLen]); err != nil {
		panic(err)
	}
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	//fmt.Println(opts)
	// todo: implement the functional requirements described in read.me
	DoAll2(opts)
}
