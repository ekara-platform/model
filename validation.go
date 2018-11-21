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

	ValidContent interface {
		validate() ValidationErrors
	}
)

var WarningOnEmpty = validEmpty(Warning)

var ErrorOnEmpty = validEmpty(Error)

var ErrorOn = valid(Error)

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

		validatorType := reflect.TypeOf((*ValidContent)(nil)).Elem()
		validatorRefType := reflect.TypeOf((*ValidReference)(nil)).Elem()

		for _, in := range ins {

			vOf := reflect.ValueOf(in)
			switch vOf.Kind() {
			case reflect.Map:
				for _, key := range vOf.MapKeys() {
					val := vOf.MapIndex(key)
					okImpl := reflect.TypeOf(val.Interface()).Implements(validatorType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidContent)
						if ok {
							vErrs.merge(checkValidContent(t, concreteVal))
						}
					}

					okImpl = reflect.TypeOf(val.Interface()).Implements(validatorRefType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidReference)
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
						concreteVal, ok := val.Interface().(ValidContent)
						if ok {
							vErrs.merge(checkValidContent(t, concreteVal))
						}
					}
					okImpl = reflect.TypeOf(val.Interface()).Implements(validatorRefType)
					if okImpl {
						concreteVal, ok := val.Interface().(ValidReference)
						if ok {
							vErrs.merge(checkValidReference(t, concreteVal))
						}
					}
				}
			default:
				okImpl := reflect.TypeOf(in).Implements(validatorType)
				if okImpl {
					concreteVal, ok := in.(ValidContent)
					if ok {
						vErrs.merge(checkValidContent(t, concreteVal))
					}
				}
				okImpl = reflect.TypeOf(in).Implements(validatorRefType)
				if okImpl {
					concreteVal, ok := in.(ValidReference)
					if ok {
						vErrs.merge(checkValidReference(t, concreteVal))
					}
				}
			}
		}
		return vErrs
	}
}

func checkValidContent(t ErrorType, c ValidContent) ValidationErrors {
	return c.validate()
}

func checkValidReference(t ErrorType, c ValidReference) ValidationErrors {
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
