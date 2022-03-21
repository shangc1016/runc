package utils

import (
	"io/ioutil"
	"os"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// TODO:
func WriteFile(data []byte, filepath string) error {
	return ioutil.WriteFile(filepath, data, 0644)
}

// TODO:
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return content, err
}

// TODO:
func GetDiectoryAll(root string) ([]string, error) {
	var dirList []string
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return []string{}, nil
	}
	for _, v := range files {
		if v.IsDir() {
			dirList = append(dirList, v.Name())
		}
	}
	return dirList, nil
}

func GetFileNameAll(p string) ([]string, error) {
	var fileList []string
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return []string{}, err
	}
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList, nil
}

//
