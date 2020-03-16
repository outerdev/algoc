package config

import (
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"

	. "github.com/outerdev/algoc/errors"
)

var (
	home, _ = os.UserHomeDir()

	dirPath1 = "./"
	dirPath2 = home + "/.config/"
	dirPath3 = home + "/"
)

func IsConfigNotPresent(err error) bool {
	return err == ErrFileNotFound
}

func checkFileSystem(path string, shouldBeDir bool) bool {
	if info, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		if shouldBeDir {
			return info.IsDir()
		} else {
			return !info.IsDir()
		}
	}
}

func fileExists(filePath string) bool {
	return checkFileSystem(filePath, false)
}

func dirExists(dirPath string) bool {
	return checkFileSystem(dirPath, true)
}

func ReadConfigData(filename string) ([]byte, error) {

	if len(filename) <= 1 {
		return nil, ErrFileNameInvalid
	}

	configFolder := filename + "/"
	if filename[0] == '.' {
		configFolder = configFolder[1:]
	}

	filePath := ""
	paths := []string{dirPath1 + filename, dirPath2 + configFolder + filename, dirPath3 + filename}
	for _, path := range paths {
		if len(filePath) == 0 {
			if fileExists(path) {
				filePath = path
			}
		} else {
			break
		}
	}

	if len(filePath) == 0 {
		return nil, ErrFileNotFound
	}

	return ioutil.ReadFile(filePath)
}

func WriteConfig(filename string, config interface{}) error {

	if len(filename) <= 1 {
		return ErrFileNameInvalid
	}

	configFolder := filename + "/"
	if filename[0] == '.' {
		configFolder = configFolder[1:]
	}

	filePath := ""
	paths := []string{dirPath1 + filename, dirPath2 + configFolder + filename, dirPath3 + filename}
	for _, path := range paths {
		if len(filePath) == 0 {
			if fileExists(path) {
				filePath = path
			}
		} else {
			break
		}
	}

	fileData, err := yaml.Marshal(config)
	if err != nil {
		return nil
	}

	if len(filePath) == 0 {
		filePath = paths[2]
	}

	return ioutil.WriteFile(filePath, fileData, 0644)
}

func LoadConfig(filename string, config interface{}) error {

	fileData, err := ReadConfigData(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(fileData, config)
}
