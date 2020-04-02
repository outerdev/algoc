package kmd

import (
	. "strconv"
	. "strings"

	"errors"

	"github.com/outerdev/algoc/path"
)

const (
	algocKmdVersion = "1.0"

	algocKmdVersionFlag = "kmdversion"
	algocKmdVersionFile = "algoc.kmd.version"
)

type Version string

func NewVersion() Version {
	return algocKmdVersion
}

func NewVersionFromString(versionStr string) (Version, error) {
	versionArgs := Split(versionStr, "=")

	if versionArgs[0] != "-"+algocKmdVersionFlag {
		return "", errors.New("no '" + algocKmdVersionFlag + "' flag found")
	}
	if len(versionArgs) == 1 {
		return "", errors.New("no version for '" + algocKmdVersionFlag + "' flag")
	}

	version := Split(versionArgs[1], ".")
	if len(version) != 2 {
		return "", errors.New("kmd version not in proper format")
	}

	if _, err := Atoi(version[0]); err != nil {
		return "", errors.New("major version not an integer")
	}

	if _, err := Atoi(version[1]); err != nil {
		return "", errors.New("minor version not an integer")
	}

	return Version(versionStr), nil
}

func (v Version) Major() int {
	major, _ := Atoi(Split(string(v), ".")[0])
	return major
}

func (v Version) Minor() int {
	minor, _ := Atoi(Split(string(v), ".")[1])
	return minor
}

func (v Version) IsLater(v2 Version) bool {
	if v.Major() > v2.Major() {
		return true
	}

	if v.Major() == v2.Major() && v.Minor() > v.Minor() {
		return true
	}

	return false
}

func (v Version) IsEqual(v2 Version) bool {
	return v.Major() == v2.Major() && v.Minor() == v2.Minor()
}

func kmdVersionValues(versionStr string) (int, int, error) {
	kmdVersionArgs := Split(versionStr, "=")

	if kmdVersionArgs[0] != "-"+algocKmdVersionFlag {
		return 0, 0, errors.New("no '" + algocKmdVersionFlag + "' flag found")
	}
	if len(kmdVersionArgs) == 1 {
		return 0, 0, errors.New("no version for '" + algocKmdVersionFlag + "' flag")
	}

	kmdVersion := Split(kmdVersionArgs[1], ".")
	if len(kmdVersion) != 2 {
		return 0, 0, errors.New("kmd version not in proper format")
	}

	major, err := Atoi(kmdVersion[0])
	if err != nil {
		return 0, 0, errors.New("major version not an integer")
	}

	minor, err := Atoi(kmdVersion[1])
	if err != nil {
		return 0, 0, errors.New("minor version not an integer")
	}

	return major, minor, nil
}

func kmdVersion(dataDir string) (int, int, error) {

	versionFile := dataDir + "/" + algocKmdVersionFile

	versionStr, err := path.ReadStringFromFile(versionFile)
	if err != nil {
		return 0, 0, err
	}

	return kmdVersionValues(versionStr)
}
