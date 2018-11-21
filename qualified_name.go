package model

import (
	"errors"
	"regexp"
)

// The Qualified name of an environment.
// This name can be used to identify, using for example Tags or Labels, all the
// content created relatively to the environment on the infrastructure of the desired
// cloud provider.
type QualifiedName struct {
	name string
	r    Environment
}

func (qn QualifiedName) String() string {
	return qn.name
}

var IsAValidQualifier = regexp.MustCompile(`^[a-zA-Z0-9_a-zA-Z0-90-9]+$`).MatchString

func (n QualifiedName) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	if !n.ValidQualifiedName() {
		vErrs.addError(errors.New("the environment name or the qualifier contains a non alphanumeric character"), n.r.location.appendPath("name|qualifier"))
	}
	return vErrs
}

func (qn QualifiedName) ValidQualifiedName() bool {
	return IsAValidQualifier(qn.String())
}

// QualifiedName returns the concatenation of the environment name and qualifier
// separated by a "_".
// If the environment qualifier is not defined it will return just the name
func (r Environment) QualifiedName() QualifiedName {
	return qualify(r)
}

// QualifiedName returns the concatenation of the environment name and qualifier
// separated by a "_".
// If the environment qualifier is not defined it will return just the name
func (r yamlEnvironment) QualifiedName() QualifiedName {
	if len(r.Qualifier) == 0 {
		return QualifiedName{name: r.Name}
	} else {
		return QualifiedName{name: r.Name + "_" + r.Qualifier}
	}
}

func qualify(r Environment) QualifiedName {

	if len(r.Qualifier) == 0 {
		return QualifiedName{r: r, name: r.Name}
	} else {
		return QualifiedName{r: r, name: r.Name + "_" + r.Qualifier}
	}
}
