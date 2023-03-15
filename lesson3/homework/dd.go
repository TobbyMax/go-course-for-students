package main

import "io"

type ReadConvertWriter interface {
	ReadConvertWrite(infile BlockReadSeekWriter, converter Converter, outfile io.Writer) error
}

type Executer interface {
	Execute() error
}

type DDer interface {
	ReadConvertWriter
	Executer
}

type DD struct {
	Options
}

func (dd *DD) ReadConvertWrite(infile BlockReadSeekWriter, converter Converter, outfile io.Writer) error {
	buf := make([]byte, infile.GetBlockSize()+MaxCarry)
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

func (dd *DD) Execute() error {
	infile, err := Open(dd.From, dd.BlockSize, dd.Limit)
	if err != nil {
		return err
	}
	defer func() {
		if err := infile.Close(); err != nil {
			panic(err)
		}
	}()

	outfile, err := Create(dd.To, dd.BlockSize)
	if err != nil {
		return err
	}
	defer func() {
		if err := outfile.Close(); err != nil {
			panic(err)
		}
	}()

	converter, err := NewDDConverter(&dd.Conv)
	if err != nil {
		return err
	}
	_, err = infile.Seek(dd.Offset, 0)
	if err != nil {
		return err
	}
	err = dd.ReadConvertWrite(infile, converter, outfile)
	if err != nil {
		return err
	}
	return nil
}
