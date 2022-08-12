package tools

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
)

func ReadFile(file multipart.File) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, file)

	return buf.Bytes(), err
}

func WriteFile(path string, file []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(file)

	return err
}
