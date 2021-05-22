package tg

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func createMultipartForm(paramName, path string) (*multipart.Writer, *bytes.Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)

	part, err := w.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, nil, err
	}

	if err = w.Close(); err != nil {
		return nil, nil, err
	}

	return w, body, nil
}
