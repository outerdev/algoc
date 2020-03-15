package main

import (
	. "fmt"

	"errors"
)

var (
	// Config errors
	errFileNotFound    = errors.New("algoc config does not exist")
	errFileNameInvalid = errors.New("config file name must be at least one letter which is not '.'")

	// Prompt errors
	errKeyNotFound = func(key string) error {
		return errors.New(Sprintf("Could not find the key '%s'", key))
	}
	errUnrecognizedValidation = func(validationTag string) error {
		return errors.New(Sprintf("Unrecognized validation named '%s'", validationTag))
	}
	errOneValidationsMap = errors.New("One validations map is the maximum")
)
