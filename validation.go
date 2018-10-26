package model

import (
	"encoding/json"
	"log"
)

type ErrorType int

const (
	Warning ErrorType = 0
	Error   ErrorType = 1
)

type ValidationErrors struct {
	Errors []ValidationError
}

type DescriptorLocation struct {
	Descriptor string
	Path       string
}

type ValidationError struct {
	ErrorType ErrorType
	Location  DescriptorLocation
	Message   string
}

func (ve ValidationErrors) Error() string {
	s := "Validation errors or warnings have occurred:\n"
	for _, err := range ve.Errors {
		s = s + "\t" + err.ErrorType.String() + ": " + err.Message + " @" + err.Location.Path + "\n\tin: " + err.Location.Descriptor
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

func (r DescriptorLocation) appendPath(suffix string) DescriptorLocation {
	newLoc := DescriptorLocation{Path: r.Path, Descriptor: r.Descriptor}
	if newLoc.Path == "" {
		newLoc.Path = suffix
	} else {
		newLoc.Path = newLoc.Path + "." + suffix
	}
	return newLoc
}

func (ve *ValidationErrors) merge(other ValidationErrors) {
	ve.Errors = append(ve.Errors, other.Errors...)
}

func (ve *ValidationErrors) addError(err error, location DescriptorLocation) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  location,
		Message:   err.Error(),
		ErrorType: Error})
}

func (ve *ValidationErrors) addWarning(message string, location DescriptorLocation) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  location,
		Message:   message,
		ErrorType: Warning})
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
