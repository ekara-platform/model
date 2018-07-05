package model

import (
	"errors"
)

// Lagoon Platform used to manipulate an environment
type LagoonPlatform struct {
	Component
}

// createLagoonPlatform create the Lagoon Platform based on the given
// repository and version
//
// The yamlRepoVersion must contains a repository and a version! If the reposiroty
// or the version is missing then a  error will be generated
func createLagoonPlatform(vErrs *ValidationErrors, env *Environment, location string, def yamlRepoVersion) LagoonPlatform {

	if def.Repository == "" {
		vErrs.AddError(errors.New("no lagoon plaform repository"), location)
	}
	if def.Version == "" {
		vErrs.AddError(errors.New("no lagoon plaform version"), location)
	}

	res := LagoonPlatform{}
	res.Component = createComponent(vErrs, env, "location", def.Repository, def.Version)
	return res
}
