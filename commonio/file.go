package commonio

import (
	"os"
	"io/ioutil"
	"errors"
)

func IsFileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
	    return false
    }
	return true
}

func ReadFile(filepath string) ([]byte, error) {
	if !IsFileExists(filepath) {
		return nil, errors.New("ReadFile: " + filepath + " is not found!")
	}
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func WriteToFile(filepath string, context []byte) error {
    return ioutil.WriteFile(filepath, context, 0644)
}

func GetFileModifyTime(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return -1, errors.New("open file error")
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return -1, errors.New("stat fileinfo error")
	}
	return fi.ModTime().Unix(), nil
}