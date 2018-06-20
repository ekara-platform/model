package model

import (
	"errors"
)

type LagoonPlateform struct {
	Component
}

func createLagoonPlateform(vErrs *ValidationErrors, env *Environment, location string, def lagoonPlateformDef) LagoonPlateform {

	if def.Repository == "" {
		vErrs.AddError(errors.New("no lagoon plaform repository"), location)
	}
	if def.Version == "" {
		vErrs.AddError(errors.New("no lagoon plaform version"), location)
	}

	res := LagoonPlateform{}
	res.Component = createComponent(vErrs, env, "location.", def.Repository, def.Version)
	return res
}
