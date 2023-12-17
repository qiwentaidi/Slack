package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// reference: https://raw.githubusercontent.com/artdarek/go-unzip/master/pkg/unzip/unzip.go

type Unzip struct{}

func NewUnzip() *Unzip {
	return &Unzip{}
}

func (uz Unzip) Extract(source, destination string) ([]string, error) {
	r, err := zip.OpenReader(source)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			logger.NewDefaultLogger().Fatal(err.Error())
		}
	}()

	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return nil, err
	}

	var extractedFiles []string
	for _, f := range r.File {
		err := uz.extractAndWriteFile(destination, f)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, f.Name)
	}

	return extractedFiles, nil
}

func (Unzip) extractAndWriteFile(destination string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			logger.NewDefaultLogger().Fatal(err.Error())
		}
	}()

	path := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(path, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", path)
	}

	if f.FileInfo().IsDir() {
		err = os.MkdirAll(path, f.Mode())
		if err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(path), f.Mode())
		if err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				logger.NewDefaultLogger().Fatal(err.Error())
			}
		}()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
