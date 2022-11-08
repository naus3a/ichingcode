package main

import (
	"errors"
	"io/ioutil"
	"os"
)

// DoesFileExist returns true if file does exist
func DoesFileExist(pth string) bool {
	if _, err := os.Stat(pth); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

// LoadFile loads a file or returns an error
func LoadFile(pth string) (data []byte, err error) {
	if !DoesFileExist(pth) {
		err = errors.New("file does not exist")
		return
	}
	data, err = ioutil.ReadFile(pth)
	return
}

// SaveFile saves a file or returns an error
func SaveFile(data []byte, pth string) error {
	return ioutil.WriteFile(pth, data, 0644)
}
