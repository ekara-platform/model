package model

import (
	"os"
)

//FileExist returns true if a file corresponding to the given path exixts
func FileExist(path string) (bool, os.FileInfo) {
	i, e := os.Stat(path)
	if e != nil {
		return false, i
	}
	return true, i
}

//DirExist returns true if a directory corresponding to the given path exixts
func DirExist(path string) bool {
	if b, i := FileExist(path); b && i.IsDir() {
		return true
	}
	return false
}
