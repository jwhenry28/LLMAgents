package utils

import (
	"os"
)

type FileManager struct {
}

func NewFileManager() FileManager {
	return FileManager{}
}

func (fm *FileManager) Read(filename string) string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}

	return string(bytes)
}
