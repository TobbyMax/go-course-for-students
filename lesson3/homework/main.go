package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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
	flag.Int64Var(&opts.Limit, "limit", math.MaxInt64, "file to write. by default - stdout")
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

func Convert(buf []byte, option string) []byte {
	switch option {
	case "lower_case":
		return bytes.ToLower(buf)
	case "upper_case":
		return bytes.ToUpper(buf)
	}
	return buf
}
func DoAll(options *Options) {
	infile := os.Stdin
	var err error = nil
	if options.From != "stdin" {
		infile, err = os.Open(options.From)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		if err := infile.Close(); err != nil {
			panic(err)
		}
	}()

	outfile := os.Stdout
	if options.To != "stdout" {
		infile, err = os.Create(options.To)
		if err != nil {
			panic(err)
		}
	}
	//defer func() {
	//	if err := outfile.Close(); err != nil {
	//		panic(err)
	//	}
	//}()
	offset, err := infile.Seek(options.Offset, 0)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, options.BlockSize)
	for offset < options.Limit {
		n, err := infile.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		offset += int64(n)
		if offset > options.Limit {
			n -= int(offset - options.Limit)
		}
		//if !utf8.Valid(buf) {
		//	n -= 1
		//	offset, err = infile.Seek(-1, 1)
		//	if err != nil {
		//		panic(err)
		//	}
		//}
		//bufConverted := Convert(buf[:n], options.Conv[0])
		if _, err := outfile.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
}

func DoAll2(options *Options) {
	if options.Offset < 0 {
		panic("Invalid offset")
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
		if err := infile.Close(); err != nil {
			panic(err)
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
	if err != nil {
		panic(err)
	}
	var offset int64 = 0
	buf := make([]byte, options.BlockSize)
	for offset < options.Offset+options.Limit {
		n, err := infile.Read(buf)
		//print("N: ")
		//print(n)
		//print("\n")
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 || err == io.EOF {
			if offset < options.Offset {
				panic("Offset out of range")
			}
			break
		}
		offset += int64(n)
		//print("Offset: ")
		//print(offset)
		//print("\n")
		if offset < options.Offset {
			continue
		}
		bufOffset := buf
		if delta := int(offset - options.Offset); delta < n && offset >= options.Offset {
			bufOffset = buf[n-delta : n]
			n = delta
		}
		if offset > options.Offset+options.Limit {
			n -= int(offset - (options.Offset + options.Limit))
		}
		//print("N2: ")
		//print(n)
		//print("\n")
		if _, err := outfile.Write(bufOffset[:n]); err != nil {
			panic(err)
		}
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
