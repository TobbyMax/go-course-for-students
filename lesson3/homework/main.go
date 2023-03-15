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

func ReadConvertWrite(infile *DDFile, converter *DDConverter, outfile *DDFile) error {
	buf := make([]byte, infile.GetBlockSize()+MAX_CARRY)
	for {
		n, err := infile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		bufTrimmed := converter.Trim(buf[:n])
		bufTailed := converter.AddTail(bufTrimmed)
		bufConverted := converter.Transform(bufTailed)
		if _, err = outfile.Write(bufConverted); err != nil {
			return err
		}
	}
	if _, err := outfile.Write(infile.GetCarryOver()); err != nil {
		return err
	}
	return nil
}

func DD(opts *Options) error {
	infile, err := Open(opts.From, opts.BlockSize, opts.Limit)
	if err != nil {
		return err
	}
	defer func() {
		if err := infile.Close(); err != nil {
			panic(err)
		}
	}()

	outfile, err := Create(opts.To, opts.BlockSize)
	if err != nil {
		return err
	}
	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()

	converter, err := NewDDConverter(&opts.Conv)
	if err != nil {
		return err
	}
	_, err = infile.Seek(opts.Offset, 0)
	if err != nil {
		return err
	}
	err = ReadConvertWrite(infile, converter, outfile)
	if err != nil {
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

	if err := DD(opts); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not perform DD:", err)
		os.Exit(2)
	}
}
