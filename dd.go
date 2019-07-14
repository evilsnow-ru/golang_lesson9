package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	//ErrIncorrectBufferSize error throws if incorrect value of buffer size is set in incoming params
	ErrIncorrectBufferSize = errors.New("incorrect buffer size: value must be greater then zero")

	//ErrIncorrectSkipSize error throws if incorrect value of skip size is set in incoming params
	ErrIncorrectSkipSize   = errors.New("incorrect skip size: value must be greater or equals to zero")
	errFileNotSetTemplate  = "%s file doesn't set"
)

var (
	bufferSize   int
	inFilePath   string
	outFilePath  string
	skipBytes    int64
	forceRewrite bool
)

func init() {
	flag.IntVar(&bufferSize, "bs", 1024, "Read buffer size")
	flag.StringVar(&inFilePath, "if", "", "Input file")
	flag.StringVar(&outFilePath, "of", "", "Out file")
	flag.Int64Var(&skipBytes, "skip", 0, "Skip bytes from input file")
	flag.BoolVar(&forceRewrite, "fw", false, "Force rewrite existing out file. Default: false.")
}

func main() {
	flag.Parse()

	if err := validateParams(); err != nil {
		panic(err)
	}

	inFile, err := os.Open(inFilePath)

	if err != nil {
		panic(err)
	}

	var outFile *os.File
	stat, err := os.Stat(outFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			outFile, err = os.Create(outFilePath)

			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		if forceRewrite {
			outFile, err = os.OpenFile(outFilePath, os.O_WRONLY, stat.Mode())
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "File \"%s\" already exist. If you need to rewrite - use \"-fw\" option.", outFilePath)
			os.Exit(-1)
		}
	}

	defer closeFile(inFile, inFilePath)
	defer closeFile(outFile, outFilePath)

	if err = copyContent(inFile, outFile, bufferSize, skipBytes); err != nil {
		panic(err)
	}
}

func closeFile(file *os.File, path string) {
	err := file.Close()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error closing file \"%s\": %s.", path, err)
	}
}

func validateParams() error {
	if bufferSize <= 0 {
		return ErrIncorrectBufferSize
	}

	if skipBytes < 0 {
		return ErrIncorrectSkipSize
	}

	if inFilePath == "" {
		return fmt.Errorf(errFileNotSetTemplate, "input")
	}

	if outFilePath == "" {
		return fmt.Errorf(errFileNotSetTemplate, "output")
	}

	return nil
}
