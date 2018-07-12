package model

import (
	"errors"
)

// Volume contains the specifications of a volume to create
type Volume struct {
	// The mounting path of the created volume
	Name string
	// The parameters required to create the volume.
	Parameters attributes `yaml:"params"`
}

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
