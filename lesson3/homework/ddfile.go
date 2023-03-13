package main

import (
	"errors"
	"io"
	"os"
)

type DDFile struct {
	file      *os.File
	blockSize int
	carryOver []byte
	carryLen  int
}

func Open(path string, blockSize int) (*DDFile, error) {
	var (
		file = os.Stdin
		err  error
	)
	if path != "stdin" {
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}
	return &DDFile{file, blockSize, make([]byte, 4), 0}, err
}

func Create(path string, blockSize int) (*DDFile, error) {
	var (
		file *os.File
		err  error
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
