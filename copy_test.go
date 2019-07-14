package main

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "otus_lesson9_test-")

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tempFile.Name())

	if err = fillTempFileWithRandomString(tempFile); err != nil {
		t.Fatal(err)
	}

	t.Logf("Created temp file: %s (length: %d)", tempFile.Name(), getFileSize(t, tempFile))

	dupFile, err := ioutil.TempFile("", "otus_lesson9_test-")

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(dupFile.Name())

	err = copyContent(tempFile, dupFile, 1000, 0)

	t.Logf("Created temp file: %s (length: %d)", dupFile.Name(), getFileSize(t, dupFile))

	tempFileBytes, err := ioutil.ReadAll(tempFile)

	if err != nil {
		t.Fatal(err)
	}

	dupFileBytes, err := ioutil.ReadAll(dupFile)

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(tempFileBytes, dupFileBytes) != 0 {
		t.Fatalf("File content doesn't equals")
	}
}

func getFileSize(t *testing.T, file *os.File) int64 {
	stat, err := file.Stat()

	if err != nil {
		t.Fatal(err)
	}

	return stat.Size()
}

func fillTempFileWithRandomString(file *os.File) error {
	for i := 0; i < 10; i++ {
		_, err := file.Write(randomBytes(1024))

		if err != nil {
			return err
		}
	}
	return nil
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
