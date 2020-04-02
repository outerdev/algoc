package path

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/outerdev/algoc/errors"
)

func WriteStringToFile(filePath, str string) error {
	var file, err = os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	if _, err = file.WriteString(str); err != nil {
		return err
	}

	return nil
}

func ReadStringFromFile(filePath string) (string, error) {
	if !FileExists(filePath) {
		return "", errors.New("file '" + filePath + "' does not exist")
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(fileBytes), nil
}

func DeleteFile(filePath string) error {
	// TODO: Setup logging for this
	// Println("DeleteFile")
	if !FileExists(filePath) {
		return ErrFileNotFound
	}

	return os.Remove(filePath)
}

func BaseAndExtension(filename string) (string, string) {
	var basename string
	extension := filepath.Ext(filename)
	if len(extension) == 0 {
		basename = filename
		extension = ".yaml"
	} else {
		basename = filename[:len(filename)-len(extension)]
	}

	return basename, extension
}

func checkFileSystem(path string, shouldBeDir bool) bool {
	var info os.FileInfo
	var err error
	if info, err = os.Stat(path); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}

	if shouldBeDir {
		return info.IsDir()
	} else {
		return !info.IsDir()
	}
}

func FileExists(filePath string) bool {
	return checkFileSystem(filePath, false)
}

func DirExists(dirPath string) bool {
	return checkFileSystem(dirPath, true)
}
