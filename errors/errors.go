package errors

import (
	. "fmt"

	"errors"
	"os"
)

var (
	// Config errors
	ErrConfigFileNotFound        = errors.New("algoc config does not exist")
	ErrConfigFileNameInvalid     = errors.New("config file name must be at least one letter which is not '.'")
	ErrConfigOnlySetFilenameOnce = errors.New("config file name may be set only once")
	ErrConfigFilenameNotSet      = errors.New("config file name must be set")

	// Path errors
	ErrFileNotFound = errors.New("file does not exist")

	// Prompt errors
	ErrKeyNotFound = func(key string) error {
		return errors.New(Sprintf("Could not find the key '%s'", key))
	}
	ErrUnrecognizedValidation = func(validationTag string) error {
		return errors.New(Sprintf("Unrecognized validation named '%s'", validationTag))
	}
	ErrOneValidationsMap  = errors.New("One validations map is the maximum")
	ErrNilFieldNotAllowed = func(s string) error {
		return errors.New(Sprintf("Field named '%s' can not be nil. Add a default struct value for this field", s))
	}
)

func Fatal(err error) {
	panic(err)
	os.Exit(1)
}
