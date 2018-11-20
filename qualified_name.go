package model

import (
	"regexp"
)

// The Qualified name of an environment.
// This name can be used to identify, using for example Tags or Labels, all the
// content created relatively to the environment on the infrastructure of the desired
// cloud provider.
type QualifiedName string

func (qn QualifiedName) String() string {
	return string(qn)
}

var IsAValidQualifier = regexp.MustCompile(`^[a-zA-Z_0-9]+$`).MatchString

func (qn QualifiedName) ValidQualifiedName() bool {
	return IsAValidQualifier(qn.String())
}

// QualifiedName returns the concatenation of the environment name and qualifier
// separated by a "_".
// If the environment qualifier is not defined it will return just the name
func (r Environment) QualifiedName() QualifiedName {
	return qualify(r.Name, r.Qualifier)
}

// QualifiedName returns the concatenation of the environment name and qualifier
// separated by a "_".
// If the environment qualifier is not defined it will return just the name
func (r yamlEnvironment) QualifiedName() QualifiedName {
	return qualify(r.Name, r.Qualifier)
}

func qualify(n, q string) QualifiedName {
	if len(q) == 0 {
		return QualifiedName(n)
	} else {
		return QualifiedName(n + "_" + q)
	}
}
