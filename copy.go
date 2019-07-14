package main

import (
	"bufio"
	"errors"
	"io"
	"os"
)

func copyContent(from *os.File, to *os.File, bufferSize int, skipBytes int64) error {
	fromFileStat, err := from.Stat()

	if err != nil {
		return err
	}

	fromFileSize := fromFileStat.Size()

	if skipBytes >= fromFileSize {
		return errors.New("skipped bytes >= in file size")
	}

	_, err = from.Seek(skipBytes, io.SeekStart)

	if err != nil {
		panic(err)
	}

	readBuff := make([]byte, bufferSize)
	var stop bool
	bufferedWriter := bufio.NewWriterSize(to, bufferSize)
	var processedBytes int64

	for !stop {
		bytesRead, err := from.Read(readBuff)

		if err != nil {
			if err == io.EOF {
				stop = true
			} else {
				return err
			}
		}

		if bytesRead > 0 {
			if _, err = bufferedWriter.Write(readBuff[:bytesRead]); err != nil {
				return err
			}
			processedBytes += int64(bytesRead)
			printProgress(processedBytes, fromFileSize, false)
		}
	}

	err = bufferedWriter.Flush()
	printProgress(fromFileSize, fromFileSize, true)

	if err != nil {
		return err
	}

	return to.Sync()
}
