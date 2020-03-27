package config

import (
	// . "fmt"

	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-yaml/yaml"

	. "github.com/outerdev/algoc/errors"
)

var configFilename string
var isConfigFilenameSet bool
var configNameMux sync.Mutex

func SetConfigFileName(filename string) {
	configNameMux.Lock()
	defer configNameMux.Unlock()

	if isConfigFilenameSet {
		panic(ErrOnlySetConfigFilenameOnce)
	} else if len(filename) <= 1 {
		panic(ErrFileNameInvalid)
	} else {
		configFilename = filename
		isConfigFilenameSet = true
	}
}

func GetConfigFilename() string {
	configNameMux.Lock()
	defer configNameMux.Unlock()

	return configFilename
}

func IsConfigNotPresent(err error) bool {
	return err == ErrFileNotFound
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

func fileExists(filePath string) bool {
	return checkFileSystem(filePath, false)
}

func DirExists(dirPath string) bool {
	return checkFileSystem(dirPath, true)
}

func getDirPaths(basename, extension string) ([]string, error) {

	home, _ := os.UserHomeDir()
	dirPath1 := home + "/.config/"
	dirPath2 := home + "/"

	var configFolder string
	if basename[0] != '.' {
		configFolder = "." + basename
	} else {
		configFolder = basename
	}

	dirPaths := []string{
		dirPath1 + configFolder[1:],
		dirPath2 + configFolder,
	}

	return dirPaths, nil
}

func getBaseAndExtension(filename string) (string, string) {
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

func defaultConfigDir(filename string) (string, error) {

	basename, extension := getBaseAndExtension(filename)

	dirPaths, err := getDirPaths(basename, extension)
	if err != nil {
		return "", nil
	}

	return dirPaths[1], nil
}

func LocateConfigDir() (string, error) {

	filename := GetConfigFilename()
	if len(filename) == 0 {
		return "", ErrConfigFilenameNotSet
	}

	basename, extension := getBaseAndExtension(filename)

	dirPaths, err := getDirPaths(basename, extension)
	if err != nil {
		return "", err
	}

	var existingDirPath string
	for _, dirPath := range dirPaths {
		if DirExists(dirPath) {
			existingDirPath = dirPath
			break
		}
	}

	if len(existingDirPath) == 0 {
		return "", ErrFileNotFound
	}

	return existingDirPath, nil
}

func locateConfigPath(filename string) (string, error) {

	basename, extension := getBaseAndExtension(filename)

	existingDirPath, err := LocateConfigDir()
	if err != nil {
		return "", err
	}

	fullPath := existingDirPath + "/" + basename + extension
	if !fileExists(fullPath) {
		return "", ErrFileNotFound
	}

	return fullPath, nil
}

func ReadConfigData() ([]byte, error) {

	filename := GetConfigFilename()
	if len(filename) == 0 {
		return nil, ErrConfigFilenameNotSet
	}

	filePath, err := locateConfigPath(filename)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(filePath)
}

func WriteConfig(config interface{}) error {

	filename := GetConfigFilename()
	if len(filename) == 0 {
		return ErrConfigFilenameNotSet
	}

	filePath, err := locateConfigPath(filename)
	if IsConfigNotPresent(err) {
		fileDir, err := defaultConfigDir(filename)
		if err != nil {
			return err
		}

		if !DirExists(fileDir) {
			err = os.Mkdir(fileDir, 0755)
			if err != nil {
				return err
			}
		}

		basename, extension := getBaseAndExtension(filename)
		filePath = fileDir + "/" + basename + extension
	}

	fileData, err := yaml.Marshal(config)
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(filePath, fileData, 0644)
}

func LoadConfig(config interface{}) error {

	fileData, err := ReadConfigData()
	if err != nil {
		return err
	}

	return yaml.Unmarshal(fileData, config)
}
