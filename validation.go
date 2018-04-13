package model

type ErrorType int

const (
	Warning ErrorType = 0
	Error   ErrorType = 1
)

type ValidationErrors struct {
	Errors []ValidationError
}

type ValidationError struct {
	ErrorType ErrorType
	Message   string
	Location  string
}

func (t ErrorType) String() string {
	names := [...]string{
		"Warning",
		"Error"}
	if t < Warning || t > Error {
		return "Unknown"
	} else {
		return names[t]
	}
}

func (ve ValidationErrors) HasErrors() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Error {
			return true
		}
	}
	return false
}

func (ve ValidationErrors) HasWarnings() bool {
	for _, v := range ve.Errors {
		if v.ErrorType == Warning {
			return true
		}
	}
	return false
}

func (ve *ValidationErrors) AddError(err error, loc string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  loc,
		Message:   err.Error(),
		ErrorType: Error})
}

func (ve *ValidationErrors) AddWarning(message string, loc string) {
	ve.Errors = append(ve.Errors, ValidationError{
		Location:  loc,
		Message:   message,
		ErrorType: Warning})
}
