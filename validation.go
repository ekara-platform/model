package model

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
)

type (
	//ErrorType type used to represent the type of a validation error
	ErrorType int

	// ValidationErrors represents a list of all error resulting of the construction
	// or the validation of an environment
	ValidationErrors struct {
		Errors []ValidationError
	}

	//ValidationError represents an error created during the construction of the
	// validation of an environment
	ValidationError struct {
		// ErrorType represents the type of the error
		ErrorType ErrorType
		// Location represents the place, within the descriptor, where the error occurred
		Location DescriptorLocation
		// Message represents a human readable message telling what need to be
		// fixed into the descriptor to get rid of this error
		Message string
	}

	// ValidableContent represents any structs which can be validated and then
	// produce ValidationErrors
	ValidableContent interface {
		validate() ValidationErrors
	}
)

// WarningOnEmptyOrInvalid allows to validate interfaces maching the following content:
//
//  // - "string": The string must not be empty.
//  //   The string content will be trimmed before the validation.
//  //
//  // - Any Map: The map cannot be empty
//  //   if the map content implementing ValidableContent or ValidableReference then it will be validated
//  //
//  // - Any Slice: The slice cannot be empty
//  //   if the slice content implementing ValidableContent or ValidableReference then it will be validated
//  //
//
// The Created validation errors will be warnings
var WarningOnEmptyOrInvalid = validNotEmpty(Warning)

// ErrorOnEmptyOrInvalid allows to validate interfaces maching the following content:
//
//  // - "string": The string must not be empty.
//  //   The string content will be trimmed before the validation.
//  //
//  // - Any Map: The map cannot be empty
//  //   if the map content implementing ValidableContent or ValidableReference then it will be validated
//  //
//  // - Any Slice: The slice cannot be empty
//  //   if the slice content implementing ValidableContent or ValidableReference then it will be validated
//  //
//
// The Created validation errors will be errors
var ErrorOnEmptyOrInvalid = validNotEmpty(Error)

// ErrorOnInvalid allows to validate interfaces maching the following content:
//
//  // - Maps of structs implementing ValidableContent or ValidableReference
//  //
//  // map[interface{}]ValidableContent
//  // map[interface{}]ValidableReference
//  //
//  // - Slices of structs implementing ValidableContent or ValidableReference
//  //
//  // []ValidableContent
//  // []ValidableReference
//  //
//  //
//  // - Any struct implementing ValidableContent or ValidableReference
//  //
//
// The Created validation errors will be errors
var ErrorOnInvalid = valid(Error)

// Error returns the message resulting of the concatenation of all included ValidationError
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

const (
	// Allows to mark validation error as Warning
	Warning ErrorType = 0
	// Allows to mark validation error as Error
	Error ErrorType = 1
)

// String return the name of the given ErrorType
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

func (ve *ValidationErrors) contains(ty ErrorType, m string, path string) bool {
	for _, v := range ve.Errors {
		if v.ErrorType.String() == ty.String() && v.Message == m && v.Location.Path == path {
			return true
		}
	}
	return false
}

func (ve *ValidationErrors) locate(m string) []ValidationError {
	result := make([]ValidationError, 0)
	for _, v := range ve.Errors {
		if v.Message == m {
			result = append(result, v)
		}
	}
	return result
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

func validNotEmpty(t ErrorType) func(in interface{}, location DescriptorLocation, message string) (ValidationErrors, bool, bool) {
	return func(in interface{}, location DescriptorLocation, message string) (ValidationErrors, bool, bool) {
		vErrs := ValidationErrors{}
		switch v := in.(type) {
		case string:
			if len(strings.Trim(v, " ")) == 0 {
				vErrs.append(t, message, location)
			}
		default:
			vOf := reflect.ValueOf(in)
			if vOf.Kind() == reflect.Map {
				if len(vOf.MapKeys()) == 0 {
					vErrs.append(t, message, location)
				} else {
					vErrs.merge(valid(t)(in))
				}
			} else if vOf.Kind() == reflect.Slice {
				if vOf.Len() == 0 {
					vErrs.append(t, message, location)
				} else {
					vErrs.merge(valid(t)(in))
				}

			}
		}
		return vErrs, vErrs.HasErrors(), vErrs.HasWarnings()
	}
}

func valid(t ErrorType) func(ins ...interface{}) ValidationErrors {
	return func(ins ...interface{}) ValidationErrors {
		vErrs := ValidationErrors{}

		validatorType := reflect.TypeOf((*ValidableContent)(nil)).Elem()
		validatorRefType := reflect.TypeOf((*ValidableReference)(nil)).Elem()

		for _, in := range ins {

			vOf := reflect.ValueOf(in)
			switch vOf.Kind() {
			case reflect.Map:
				for _, key := range vOf.MapKeys() {
					val := vOf.MapIndex(key)
					okImpl := reflect.TypeOf(val.Interface()).Implements(validatorType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidableContent)
						if ok {
							vErrs.merge(checkValidContent(t, concreteVal))
						}
					}

					okImpl = reflect.TypeOf(val.Interface()).Implements(validatorRefType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidableReference)
						if ok {
							vErrs.merge(checkValidReference(t, concreteVal))
						}
					}
				}
			case reflect.Slice:
				for i := 0; i < vOf.Len(); i++ {
					val := vOf.Index(i)
					okImpl := reflect.TypeOf(val.Interface()).Implements(validatorType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidableContent)
						if ok {
							vErrs.merge(checkValidContent(t, concreteVal))
						}
					}
					okImpl = reflect.TypeOf(val.Interface()).Implements(validatorRefType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidableReference)
						if ok {
							vErrs.merge(checkValidReference(t, concreteVal))
						}
					}
				}
			default:
				okImpl := reflect.TypeOf(in).Implements(validatorType)
				if okImpl {
					concreteVal, ok := in.(ValidableContent)
					if ok {
						vErrs.merge(checkValidContent(t, concreteVal))
					}
				}
				okImpl = reflect.TypeOf(in).Implements(validatorRefType)
				if okImpl {
					concreteVal, ok := in.(ValidableReference)
					if ok {
						vErrs.merge(checkValidReference(t, concreteVal))
					}
				}
			}
		}
		return vErrs
	}
}

//checkValidContent validates a ValidableContent
func checkValidContent(t ErrorType, c ValidableContent) ValidationErrors {
	return c.validate()
}

//checkValidReference validates a ValidableReference
func checkValidReference(t ErrorType, c ValidableReference) ValidationErrors {
	vErrs := ValidationErrors{}
	if c.Reference().Id == "" {
		if c.Reference().Mandatory {
			vErrs.append(t, "empty "+c.Reference().Type+" reference", c.Reference().Location)
		}
	} else {
		if _, ok := c.Reference().Repo[c.Reference().Id]; !ok {
			vErrs.append(t, "reference to unknown "+c.Reference().Type+": "+c.Reference().Id, c.Reference().Location)
		}
	}
	return vErrs
}
