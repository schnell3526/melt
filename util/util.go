package util

import (
	"fmt"
	"os"
)

// if file does not exists, return false.
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// create dirPath directory with mode dirMode
func Mkdir(dirPath string, dirMode os.FileMode) error {
	err := os.MkdirAll(dirPath, dirMode)
	if err != nil {
		return fmt.Errorf("%s: making directory: %v", dirPath, err)
	}
	return nil
}
