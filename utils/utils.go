package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
)

// Sha1Reader get file sha1
func Sha1Reader(reader io.Reader) string {
	_sha1 := sha1.New()
	data := make([]byte, 10240)
	for {
		n, err := reader.Read(data)

		if io.EOF == err {
			break
		}

		if err != nil {
			log.Printf("read file fiald: %v", err)
			return ""
		}

		_sha1.Write(data[:n])
	}

	return hex.EncodeToString(_sha1.Sum(nil))
}

// Sha1File get file sha1
func Sha1File(path string) string {
	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		log.Printf("open file fiald: %v", err)
		return ""
	}

	defer f.Close()

	return Sha1Reader(f)
}
