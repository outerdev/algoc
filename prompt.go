package main

import (
	. "fmt"
	. "strconv"
	. "strings"

	"errors"
	"net/url"
	"reflect"
	"regexp"

	"github.com/manifoldco/promptui"
)

type ValidateFunc func(string) error

func IsAlphaNumeric(input string) error {
	isAlphaNumeric := regexp.MustCompile(`^[0-9a-zA-Z]+$`).MatchString
	if !isAlphaNumeric(input) {
		return errors.New("Key must contain only numbers and letters")
	}
	return nil
}

func IsRealNumber(input string) error {
	_, err := ParseFloat(input, 64)
	if err != nil {
		return errors.New("Invalid number")
	}
	return nil
}

func IsInteger(input string) error {
	_, err := ParseInt(input, 10, 64)
	if err != nil {
		return errors.New("Invalid integer")
	}
	return nil
}

func IsURL(input string) error {
	parsedURL, err := url.Parse(input)
	if err != nil || len(parsedURL.Scheme) == 0 || len(parsedURL.Hostname()) == 0 {
		return errors.New("Invalid URL")
	}
	return nil
}

var validationFuncs = map[string]ValidateFunc{
	"alphanumeric": IsAlphaNumeric,
	"real":         IsRealNumber,
	"integer":      IsInteger,
	"url":          IsURL,
}

func promptValue(valueLabel string, validations ...ValidateFunc) (string, error) {

	validate := func(input string) error {
		for _, validation := range validations {
			if err := validation(input); err != nil {
				return err
			}
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    valueLabel,
		Validate: validate,
		Templates: &promptui.PromptTemplates{
			Success:         "",
			ValidationError: "",
			Valid:           valueLabel + ": ",
			Invalid:         valueLabel + ": ",
		},
	}

	return prompt.Run()
}

func hasPrompt(v reflect.StructField) bool {
	actionTags, ok := v.Tag.Lookup("action")
	if ok {
		for _, actionTag := range Split(actionTags, ",") {
			return TrimSpace(actionTag) == "prompt"
		}
	}
	return false
}

func getValidations(v reflect.StructField) []ValidateFunc {
	anyString := func(string) error {
		return nil
	}

	actionTags, _ := v.Tag.Lookup("action")
	validationTags := Split(actionTags, ",")[1:]
	var validations []ValidateFunc

	if len(validationTags) == 0 {
		validations = append(validations, anyString)
	}

	for _, validationTag := range validationTags {
		if validation, ok := validationFuncs[validationTag]; ok {
			validations = append(validations, validation)
		} else {
			panic(errUnrecognizedValidation(validationTag))
		}
	}

	return validations
}

func valueFromString(kind reflect.Kind, valueStr string) (reflect.Value, error) {
	if kind == reflect.String {
		return reflect.ValueOf(valueStr), nil
	} else if kind == reflect.Int {
		value, err := ParseInt(valueStr, 10, 32)
		if err != nil {
			return reflect.ValueOf(0), err
		}
		return reflect.ValueOf(int(value)), nil
	}

	return reflect.Value{}, nil
}

func printValueDetails(i int, f reflect.StructField, v reflect.Value) {
	Printf("%d: %s %s = %v\n", i, f.Name, f.Type, v.Field(i).Interface())
}

func fillKeyByPrompt(v reflect.Value, prefix string, key string) error {
	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		if f.Type.Kind() == reflect.Struct {
			if len(key) > 0 {
				return fillKeyByPrompt(v.Field(i), prefix+f.Name+".", key)
			}

			if err := fillKeyByPrompt(v.Field(i), prefix+f.Name+".", key); err != nil {
				return err
			}
		}

		if f.Type.Kind() == reflect.Ptr {
			if v.Field(i).IsNil() {
				panic(Sprintf("Field named '%s' can not be nil. Add a default struct value for this field", prefix+f.Name))
				// Fails because the value is unaddressable, is this fixable?
				// r := reflect.New(f.Type.Elem()).Elem()
				// v.Field(i).Set(r)
			}

			if len(key) > 0 {
				return fillKeyByPrompt(v.Field(i).Elem(), prefix+f.Name+".", key)
			}

			if err := fillKeyByPrompt(v.Field(i).Elem(), prefix+f.Name+".", key); err != nil {
				return err
			}
		}

		if hasPrompt(f) {
			// If we are looking for a specific key and this is not it,
			// continue looking through the keys
			if len(key) > 0 && key != prefix+f.Name {
				continue
			}

			validations := getValidations(f)
			valueStr, err := promptValue(prefix+f.Name, validations...)
			if err != nil {
				panic(err)
			}
			value, err := valueFromString(f.Type.Kind(), valueStr)
			if err != nil {
				panic(err)
			}
			v.Field(i).Set(value)

			// We found the key we were looking to set so we are done
			if len(key) > 0 {
				return nil
			}
		}
	}

	// If we are looking for a specific key and didn't see it
	// after looking through all the keys then there was an error
	if len(key) > 0 {
		return errKeyNotFound(key)
	} else { // We have set all the keys with action:"prompt" tags
		return nil
	}
}

func fillAllKeysByPrompt(v reflect.Value) error {
	return fillKeyByPrompt(v, "", "")
}

func evaluateValidationsMap(validationsMap ...map[string]ValidateFunc) {
	if len(validationsMap) > 1 {
		panic(errOneValidationsMap)
	} else if len(validationsMap) == 1 {
		for name, validation := range validationsMap[0] {
			validationFuncs[name] = validation
		}
	}
}
func PromptForValues(config *Config, validationsMap ...map[string]ValidateFunc) error {

	evaluateValidationsMap(validationsMap...)

	// var v reflect.Value
	// if reflect.ValueOf(config).Kind() == reflect.Ptr {
	v := reflect.ValueOf(config).Elem()
	// } else {
	// v = reflect.ValueOf(config)
	// }

	return fillAllKeysByPrompt(v)
}

// PromptForValuesWithKey will prompt and fill a series of values in config
func PromptForValuesWithKeys(config *Config, keys []string, validationsMap ...map[string]ValidateFunc) error {

	evaluateValidationsMap(validationsMap...)

	v := reflect.ValueOf(config).Elem()
	for _, key := range keys {
		if err := fillKeyByPrompt(v, "", key); err != nil {
			return err
		}
	}
	return nil
}

// PromptForValuesWithKey will prompt and fill one value in config
func PromptForValuesWithKey(config *Config, key string, validationsMap ...map[string]ValidateFunc) error {
	return PromptForValuesWithKeys(config, []string{key}, validationsMap...)
}
