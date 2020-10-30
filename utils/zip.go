package utils

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ZipDir zip dir
func ZipDir(dir, zipFile string, hidden bool) error {
	err := os.MkdirAll(filepath.Dir(zipFile), 0755)
	if err != nil {
		return err
	}

	fz, err := os.Create(zipFile)
	if err != nil {
		return err
	}

	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		sub := path[len(dir):]
		if !hidden && strings.HasPrefix(sub, ".") {
			return nil
		}

		dst, err := w.Create(sub)
		if err != nil {
			return err
		}

		src, err := os.Open(path)
		if err != nil {
			return err
		}

		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			log.Printf("Copy failed: %s\n", err.Error())
			return err
		}

		return nil
	})

	return nil
}

// Unzip unzip
func Unzip(zipFile, dir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil
	}

	defer r.Close()

	for _, f := range r.File {
		path := dir + string(filepath.Separator) + f.Name
		os.MkdirAll(filepath.Dir(path), 0755)
		dst, err := os.Create(path)
		if err != nil {
			log.Printf("Create failed: %s\n", err.Error())
			continue
		}

		src, err := f.Open()
		if err != nil {
			log.Printf("Open failed: %s\n", err.Error())
			continue
		}

		defer src.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			log.Printf("Copy failed: %s\n", err.Error())
		}

		defer dst.Close()
	}

	return nil
}
