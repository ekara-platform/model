package model

import (
	"encoding/json"
	"reflect"
	"strings"
)

type (
	//ErrorType type used to represent the type of a validation error
	ErrorType int

	//ValidationErrors represents a list of all error resulting of the construction
	//or the validation of an environment
	ValidationErrors struct {
		Errors []ValidationError
	}

	//ValidationError represents an error created during the construction of the
	// validation or an environment
	ValidationError struct {
		// ErrorType represents the type of the error
		ErrorType ErrorType
		// Location represents the place, within the descriptor, where the error occurred
		Location DescriptorLocation
		// Message represents a human readable message telling what need to be
		// fixed into the descriptor to get rid of this error
		Message string
	}

	// validatableContent represents any structs which can be validated and then
	// produce ValidationErrors
	//
	// If a structure implements validatableContent then it will be validated invoking:
	//  ErrorOnInvalid(myStruct)
	validatableContent interface {
		validate() ValidationErrors
	}
)

// WarningOnEmptyOrInvalid allows to validate interfaces matching the following content:
//
// The created validation errors will be warnings.
//
// A string
//
// The string must not be empty.
//
// The string content will be trimmed before the validation.
//
// Any Map
//
// The map cannot be empty.
//
// If the map content implements validatableContent or validatableReference then it will be validated.
//
// Any Slice
//
// The slice cannot be empty.
//
// If the slice content implements validatableContent or validatableReference then it will be validated.
//
//
var WarningOnEmptyOrInvalid = validNotEmpty(Warning)

// ErrorOnEmptyOrInvalid allows to validate interfaces matching the following content:
//
// The created validation errors will be errors.
//
// A String
//
// The string must not be empty.
//
// The string content will be trimmed before the validation.
//
// Any Map
//
// The map cannot be empty.
//
// If the map content implements validatableContent or validatableReference then it will be validated.
//
// Any Slice
//
// The slice cannot be empty.
//
// If the slice content implements validatableContent or validatableReference then it will be validated.
var ErrorOnEmptyOrInvalid = validNotEmpty(Error)

// ErrorOnInvalid allows to validate interfaces matching the following content:
//
// The created validation errors will be errors.
//
// Maps of structs implementing validatableContent or validatableReference
//
// map[interface{}]validatableContent.
//
// map[interface{}]validatableReference.
//
// Slices of structs implementing validatableContent or validatableReference
//
// []validatableContent.
//
// []validatableReference.
//
// Any struct
//
// The structure must implement validatableContent or validatableReference
//
//
var ErrorOnInvalid = valid(Error)

// Error returns the message resulting of the concatenation of all included ValidationError(s)
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
	//Warning allows to mark validation error as Warning
	Warning ErrorType = 0
	//Error allows to mark validation error as Error
	Error ErrorType = 1
)

// String return the name of the given ErrorType
func (r ErrorType) String() string {
	names := [...]string{
		"Warning",
		"Error"}
	if r < Warning || r > Error {
		return "Unknown"
	}
	return names[r]
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

// HasErrors returns true if the ValidationErrors contains at least one error
func (ve ValidationErrors) HasErrors() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Error {
			return true
		}
	}
	return false
}

// HasWarnings returns true if the ValidationErrors contains at least one warning
func (ve ValidationErrors) HasWarnings() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Warning {
			return true
		}
	}
	return false
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

		validatorType := reflect.TypeOf((*validatableContent)(nil)).Elem()
		validatorRefType := reflect.TypeOf((*validatableReferencer)(nil)).Elem()

		for _, in := range ins {

			vOf := reflect.ValueOf(in)
			switch vOf.Kind() {
			case reflect.Map:
				validMap(in, vOf, validatorType, validatorRefType, &vErrs, t)
			case reflect.Slice:
				validSlice(in, vOf, validatorType, validatorRefType, &vErrs, t)
			default:
				validDefault(in, validatorType, validatorRefType, &vErrs, t)
			}
		}
		return vErrs
	}
}

func validMap(in interface{}, vOf reflect.Value, validatorType reflect.Type, validatorRefType reflect.Type, vErrs *ValidationErrors, t ErrorType) {
	for _, key := range vOf.MapKeys() {
		val := vOf.MapIndex(key)
		okImpl := reflect.TypeOf(val.Interface()).Implements(validatorType)
		if okImpl {
			concreteVal, ok := val.Interface().(validatableContent)
			if ok {
				vErrs.merge(checkValidContent(t, concreteVal))
			}
		}

		okImpl = reflect.TypeOf(val.Interface()).Implements(validatorRefType)
		if okImpl {
			concreteVal, ok := val.Interface().(validatableReferencer)
			if ok {
				vErrs.merge(checkValidReference(t, concreteVal))
			}
		}
	}
}

func validSlice(in interface{}, vOf reflect.Value, validatorType reflect.Type, validatorRefType reflect.Type, vErrs *ValidationErrors, t ErrorType) {
	for i := 0; i < vOf.Len(); i++ {
		val := vOf.Index(i)
		okImpl := reflect.TypeOf(val.Interface()).Implements(validatorType)
		if okImpl {
			concreteVal, ok := val.Interface().(validatableContent)
			if ok {
				vErrs.merge(checkValidContent(t, concreteVal))
			}
		}
		okImpl = reflect.TypeOf(val.Interface()).Implements(validatorRefType)
		if okImpl {
			concreteVal, ok := val.Interface().(validatableReferencer)
			if ok {
				vErrs.merge(checkValidReference(t, concreteVal))
			}
		}
	}

}

func validDefault(in interface{}, validatorType reflect.Type, validatorRefType reflect.Type, vErrs *ValidationErrors, t ErrorType) {
	okImpl := reflect.TypeOf(in).Implements(validatorType)
	if okImpl {
		concreteVal, ok := in.(validatableContent)
		if ok {
			vErrs.merge(checkValidContent(t, concreteVal))
		}
	}
	okImpl = reflect.TypeOf(in).Implements(validatorRefType)
	if okImpl {
		concreteVal, ok := in.(validatableReferencer)
		if ok {
			vErrs.merge(checkValidReference(t, concreteVal))
		}
	}
}

//checkValidContent validates a validatableContent
func checkValidContent(t ErrorType, c validatableContent) ValidationErrors {
	return c.validate()
}

//checkValidReference validates a validatableReference
func checkValidReference(t ErrorType, c validatableReferencer) ValidationErrors {
	vErrs := ValidationErrors{}
	ref := c.reference()
	if ref.Id == "" {
		if ref.Mandatory {
			vErrs.append(t, "empty "+ref.Type+" reference", ref.Location)
		}
	} else {
		if _, ok := ref.Repo[ref.Id]; !ok {
			vErrs.append(t, "reference to unknown "+ref.Type+": "+ref.Id, ref.Location)
		}
	}
	return vErrs
}
