package main

import (
	"errors"
	"io"
	"os"
	"unicode/utf8"
)

const MaxCarry = 3

type BlockSizeGetter interface {
	GetBlockSize() int
}

type CarryOverGetter interface {
	GetCarryOver() []byte
}

type BlockReadSeekWriter interface {
	io.ReadSeekCloser
	io.Writer
	BlockSizeGetter
	CarryOverGetter
}

type DDFile struct {
	file          *os.File
	currentOffset int64
	mode          byte
	// Options
	blockSize int
	limit     int64
	// Attribute for limit function
	bytesRead int64
	// Attributes for invalid bytes transferring
	carryOver []byte
	carryLen  int
}

func Open(path string, blockSize int, limit int64) (*DDFile, error) {
	var file = os.Stdin
	if path != "stdin" {
		var err error
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}
	return &DDFile{file: file, mode: 'r', blockSize: blockSize, carryOver: make([]byte, MaxCarry), limit: limit}, nil
}

func Create(path string, blockSize int) (*DDFile, error) {
	var file = os.Stdout
	if path != "stdout" {
		_, err := os.Stat(path)
		if !errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("outfile already exists")
		}
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	return &DDFile{file: file, blockSize: blockSize, mode: 'w'}, nil
}

func (ddf *DDFile) Close() error {
	if err := ddf.file.Close(); err != nil {
		return err
	}
	ddf.carryOver = nil
	return nil
}

func (ddf *DDFile) Read(b []byte) (int, error) {
	if ddf.mode == 'w' {
		return 0, errors.New("can not read: file opened for writing only")
	}
	buf := make([]byte, ddf.blockSize+ddf.carryLen)
	readLen := ddf.blockSize
	if bufLen := len(b); bufLen-ddf.carryLen < ddf.blockSize {
		readLen = bufLen - ddf.carryLen
	}
	if ddf.limit != -1 && ddf.limit-ddf.bytesRead < int64(readLen) {
		readLen = int(ddf.limit - ddf.bytesRead)
	}
	n, err := ddf.file.Read(buf[ddf.carryLen : readLen+ddf.carryLen])
	if err != nil && err != io.EOF {
		return 0, err
	}
	if n == 0 || err == io.EOF {
		return 0, io.EOF
	}

	ddf.bytesRead += int64(n)
	copy(buf, ddf.carryOver[:ddf.carryLen])
	n += ddf.carryLen
	ddf.carryLen = 0
	ddf.currentOffset += int64(n)

	bufValid := buf[:n]
	bufValid = ddf.trimInvalidBytes(bufValid)
	copy(b, bufValid)
	n -= ddf.carryLen
	return n, err
}

func (ddf *DDFile) Write(b []byte) (int, error) {
	if ddf.mode == 'r' {
		return 0, errors.New("can not write: file opened for reading")
	}
	var (
		bufLen   = len(b)
		writeLen int
	)
	for i := 0; i < bufLen; i += ddf.blockSize {
		writeLen := i + ddf.blockSize
		if bufLen < i+ddf.blockSize {
			writeLen = bufLen
		}
		if _, err := ddf.file.Write(b[i:writeLen]); err != nil {
			return writeLen, err
		}
	}
	return writeLen, nil
}

func (ddf *DDFile) Seek(offset int64, whence int) (int64, error) {
	if ddf.mode == 'w' {
		return 0, errors.New("can not seek: file opened for writing only")
	}
	if offset < 0 {
		return ddf.currentOffset, errors.New("invalid offset")
	}
	if ddf.file == os.Stdin {
		switch {
		case whence == 2:
			return ddf.currentOffset, errors.New("can not seek from the end in stdin")
		case whence == 0 && ddf.currentOffset != 0:
			return ddf.currentOffset, errors.New("can not seek from the start in stdin (if current offset != 0)")
		}
		b := make([]byte, ddf.blockSize)
		for blockSize := int64(ddf.blockSize); ddf.currentOffset < offset; {
			seekLen := blockSize
			if offset < ddf.currentOffset+blockSize {
				seekLen = offset - ddf.currentOffset
			}
			n, err := ddf.file.Read(b[:seekLen])
			if err != nil && err != io.EOF {
				return 0, err
			}
			if n == 0 || err == io.EOF {
				return 0, errors.New("offset index out of range")
			}
			ddf.currentOffset += int64(n)
		}
		return offset, nil
	}
	return ddf.file.Seek(offset, whence)
}

func (ddf *DDFile) trimInvalidBytes(b []byte) []byte {
	i := len(b) - 1
	buf := b
	for ; i >= 0; i-- {
		if utf8.RuneStart(b[i]) {
			if !utf8.Valid(b[i:]) {
				copy(ddf.carryOver, b[i:])
				buf = b[:i]
				ddf.carryLen = len(b) - i
				ddf.currentOffset -= int64(ddf.carryLen)
			}
			break
		}
	}
	return buf
}

func (ddf *DDFile) GetBlockSize() int {
	return ddf.blockSize
}

func (ddf *DDFile) GetCarryOver() []byte {
	if ddf.mode == 'w' {
		return nil
	}
	return ddf.carryOver[:ddf.carryLen]
}
