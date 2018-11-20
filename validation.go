package model

import (
	"encoding/json"

	"log"
	"reflect"
)

type (
	ErrorType int

	ValidationErrors struct {
		Errors []ValidationError
	}

	ValidationError struct {
		ErrorType ErrorType
		Location  DescriptorLocation
		Message   string
	}

	Valid interface {
		validate() ValidationErrors
	}
)

var WarningOnEmpty = validEmpty(Warning)

var ErrorOnEmpty = validEmpty(Error)

const (
	Warning ErrorType = 0
	Error   ErrorType = 1
)

func (ve ValidationErrors) Error() string {
	s := "Validation errors or warnings have occurred:\n"
	for _, err := range ve.Errors {
		s = s + "\t" + err.ErrorType.String() + ": " + err.Message + " @" + err.Location.Path + "\n\tin: " + err.Location.Descriptor + "\n\t"
	}
	return s
}

// JSonContent returns the serialized content of all validations
// errors as JSON
func (ve ValidationErrors) JSonContent() (b []byte, e error) {
	b, e = json.MarshalIndent(ve.Errors, "", "    ")
	return
}

func (r ErrorType) String() string {
	names := [...]string{
		"Warning",
		"Error"}
	if r < Warning || r > Error {
		return "Unknown"
	} else {
		return names[r]
	}
}

func (ve *ValidationErrors) merge(other ValidationErrors) {
	ve.Errors = append(ve.Errors, other.Errors...)
}

func (ve *ValidationErrors) append(t ErrorType, e string, l DescriptorLocation) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  l,
		Message:   e,
		ErrorType: t,
	})
}

func (ve *ValidationErrors) contains(m string, location DescriptorLocation) bool {
	for _, v := range ve.Errors {
		if v.Message == m && v.Location.equals(location) {
			return true
		}
	}
	return false
}

func (ve *ValidationErrors) addError(err error, location DescriptorLocation) {
	ve.append(Error, err.Error(), location)
}

func (ve *ValidationErrors) addWarning(message string, location DescriptorLocation) {
	ve.append(Warning, message, location)
}

// HasErrors returns true if the ValidationErrors contains at least on error
func (ve ValidationErrors) HasErrors() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Error {
			return true
		}
	}
	return false
}

// HasWarnings returns true if the ValidationErrors contains at least on warning
func (ve ValidationErrors) HasWarnings() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Warning {
			return true
		}
	}
	return false
}

// Log logs all the validation errors to the specified logger
func (ve ValidationErrors) Log(logger *log.Logger) {
	for _, err := range ve.Errors {
		logger.Println(err.ErrorType.String() + " @" + err.Location.Path + ": " + err.Message)
	}
}

func validEmpty(t ErrorType) func(in interface{}, location DescriptorLocation, message string) (ValidationErrors, bool, bool) {
	return func(in interface{}, location DescriptorLocation, message string) (ValidationErrors, bool, bool) {
		vErrs := ValidationErrors{}
		switch v := in.(type) {
		case string:
			if len(v) == 0 {
				vErrs.append(t, message, location)
			}
		default:
			vOf := reflect.ValueOf(in)
			if vOf.Kind() == reflect.Map {
				if len(vOf.MapKeys()) == 0 {
					vErrs.append(t, message, location)
				} else {
					for _, key := range vOf.MapKeys() {
						val := vOf.MapIndex(key)

						validatorType := reflect.TypeOf((*Valid)(nil)).Elem()
						okImpl := reflect.TypeOf(val.Interface()).Implements(validatorType)
						if okImpl {
							concreteval, ok := val.Interface().(Valid)
							if ok {
								vErrs.merge(concreteval.validate())
							}
						}
					}
				}
			}
		}
		return vErrs, vErrs.HasErrors(), vErrs.HasWarnings()
	}
}
