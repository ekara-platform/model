package model

import (
	"errors"
)

type Volume struct {
	Name       string
	Parameters attributes `yaml:"params"`
}

// TODO comments
func createVolumes(vErrs *ValidationErrors, env *Environment, location string, yamlRef []yamlVolumes) []Volume {
	volumes := make([]Volume, 0, 00)
	for _, v := range yamlRef {
		if len(v.Name) == 0 {
			vErrs.AddError(errors.New("empty volume name"), location)
		} else {
			volumes = append(volumes, Volume{Parameters: createAttributes(v.Params, nil), Name: v.Name})
		}
	}
	return volumes
}
