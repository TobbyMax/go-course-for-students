package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.From, "from", "stdin", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "stdout", "file to write. by default - stdout")

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

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}
	
	dd, err := New(opts)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not do all:", err)
		os.Exit(2)
	}

	if err := dd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not do all:", err)
		os.Exit(3)
	}
}
