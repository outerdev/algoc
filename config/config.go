package config

import (
	"io/ioutil"
	"os"
	"sync"

	"github.com/go-yaml/yaml"

	. "github.com/outerdev/algoc/errors"
	"github.com/outerdev/algoc/path"
)

var configFilename string
var isConfigFilenameSet bool
var configNameMux sync.Mutex

func SetConfigFileName(filename string) {
	configNameMux.Lock()
	defer configNameMux.Unlock()

	if isConfigFilenameSet {
		panic(ErrConfigOnlySetFilenameOnce)
	} else if len(filename) <= 1 {
		panic(ErrConfigFileNameInvalid)
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
	return err == ErrConfigFileNotFound
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

func defaultConfigDir(filename string) (string, error) {

	basename, extension := path.BaseAndExtension(filename)

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

	basename, extension := path.BaseAndExtension(filename)

	dirPaths, err := getDirPaths(basename, extension)
	if err != nil {
		return "", err
	}

	var existingDirPath string
	for _, dirPath := range dirPaths {
		if path.DirExists(dirPath) {
			existingDirPath = dirPath
			break
		}
	}

	if len(existingDirPath) == 0 {
		return "", ErrConfigFileNotFound
	}

	return existingDirPath, nil
}

func locateConfigPath(filename string) (string, error) {

	basename, extension := path.BaseAndExtension(filename)

	existingDirPath, err := LocateConfigDir()
	if err != nil {
		return "", err
	}

	fullPath := existingDirPath + "/" + basename + extension
	if !path.FileExists(fullPath) {
		return "", ErrConfigFileNotFound
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

		if !path.DirExists(fileDir) {
			err = os.Mkdir(fileDir, 0755)
			if err != nil {
				return err
			}
		}

		basename, extension := path.BaseAndExtension(filename)
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
