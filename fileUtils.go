package bingo

import (
	"os"
)

func isFileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}
