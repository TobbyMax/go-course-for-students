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
	BlockSize int
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
	flag.IntVar(&opts.BlockSize, "block-size", 1024, "file to write. by default - stdout")
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

type DDFile struct {
	file      *os.File
	blockSize int
	carryOver []byte
	carryLen  int
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

func getConvStatus(conversions *StringSlice) (trim, upper, lower bool, err error) {
	err = nil
	if len(*conversions) > 2 {
		err = errors.New("too many arguments")
		return false, false, false, err
	}

	trim = false
	upper = false
	lower = false
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
			return false, false, false, err
		}
	}
	if lower == true && upper == true {
		err = errors.New("can not apply 'upper_case' and 'lower_case' simultaneously")
		return false, false, false, err
	}
	return trim, upper, lower, err
}

func Open(path string, blockSize int) (*DDFile, error) {
	var (
		file       = os.Stdin
		err  error = nil
	)
	if path != "stdin" {
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}
	return &DDFile{file, blockSize, make([]byte, 4), 0}, err
}

func (ddf *DDFile) Close() error {
	if err := ddf.file.Close(); err != nil {
		return err
	}
	return nil
}

func (ddf *DDFile) Read(b []byte) (int, error) {
	bufSize, err := ddf.file.Read(b[ddf.carryLen : ddf.blockSize+ddf.carryLen])
	if err != nil && err != io.EOF {
		return 0, err
	}
	copy(b, ddf.carryOver[:ddf.carryLen])
	bufSize += ddf.carryLen
	ddf.carryLen = 0
	return bufSize, err
}

func Create(path string, blockSize int) (*DDFile, error) {
	var (
		file       = os.Stdin
		err  error = nil
	)
	file = os.Stdout
	if path != "stdout" {
		_, err = os.Stat(path)
		if !errors.Is(err, os.ErrNotExist) {
			panic("Outfile already exists")
		}
		file, err = os.Create(path)
		if err != nil {
			panic(err)
		}
	}
	return &DDFile{file, blockSize, make([]byte, 4), 0}, err
}

func (ddf *DDFile) Write(b []byte) error {
	var bufLen = len(b)
	for i := 0; i < bufLen; i += ddf.blockSize {
		writeLen := bufLen
		if bufLen > i+ddf.blockSize {
			writeLen = i + ddf.blockSize
		}
		if _, err := ddf.file.Write(b[i:writeLen]); err != nil {
			return err
		}
	}
	return nil
}
func Execute(options *Options) error {
	if options.Offset < 0 {
		panic("Invalid offset")
	}
	var err error = nil
	trim, upper, lower, err := getConvStatus(&options.Conv)
	if err != nil {
		return err
	}

	infile, err := Open(options.From, options.BlockSize)
	if err != nil {
		return err
	}

	defer func() {
		if err := infile.Close(); err != nil {
			panic(err)
		}
	}()

	outfile, err := Create(options.To, options.BlockSize)
	if err != nil {
		return err
	}

	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()

	var (
		offset    int64  = 0
		start            = true
		bufSpaces []byte = nil
		buf              = make([]byte, options.BlockSize+4)
	)

	for options.Limit == -1 || offset < options.Offset+options.Limit {
		var (
			bufSize int
			err     error
		)
		bufSize, err = infile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if bufSize == 0 || err == io.EOF {
			if offset < options.Offset {
				return errors.New("offset index out of range")
			}
			break
		}
		offset += int64(bufSize)
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
					infile.carryOver = bufLimited[i:]
					bufLimited = bufLimited[:i]
					infile.carryLen = bufSize - i
					offset -= int64(infile.carryLen)
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
				if err = outfile.Write(bufSpaces); err != nil {
					return err
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
		if err = outfile.Write(bufConverted); err != nil {
			return err
		}
	}
	if err = outfile.Write(infile.carryOver[:infile.carryLen]); err != nil {
		return err
	}
	return nil
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	//fmt.Println(opts)
	// todo: implement the functional requirements described in read.me
	if err := Execute(opts); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not do all:", err)
		os.Exit(2)
	}
}
