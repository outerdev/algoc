package errors

import (
	. "fmt"

	"errors"
	"os"
)

var (
	// Config errors
	ErrFileNotFound    = errors.New("algoc config does not exist")
	ErrFileNameInvalid = errors.New("config file name must be at least one letter which is not '.'")

	// Prompt errors
	ErrKeyNotFound = func(key string) error {
		return errors.New(Sprintf("Could not find the key '%s'", key))
	}
	ErrUnrecognizedValidation = func(validationTag string) error {
		return errors.New(Sprintf("Unrecognized validation named '%s'", validationTag))
	}
	ErrOneValidationsMap  = errors.New("One validations map is the maximum")
	ErrNilFieldNotAllowed = func(s string) error {
		return errors.New(Sprintf("Field named '%s' can not be nil. Add a default struct value for this field"))
	}
)

func Fatal(err error) {
	Println(err)
	os.Exit(1)
}
