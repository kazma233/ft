package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// go run test -v
func Test_handler1(t *testing.T) {
	filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		t.Log(path)
		return nil
	})
}
