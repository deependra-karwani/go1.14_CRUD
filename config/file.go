package config

import (
	"bytes"
	"io"
	"os"
)

func SaveFile(file io.Reader, dir string, filename string) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	fileBytes := buf.Bytes()
	buf.Reset()

	f, err := os.Create(dir + filename)
	if err != nil {
		return err
	}
	f.Write(fileBytes)
	f.Close()

	return nil
}
