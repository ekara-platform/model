package model

import (
	"errors"
)

// Volume contains the specifications of a volume to create
type Volume struct {
	// The mounting path of the created volume
	Path string
	// The parameters required to create the volume.
	Parameters Parameters `yaml:"params"`
}

func createVolumes(vErrs *ValidationErrors, location string, yamlRef []yamlVolume) []Volume {
	volumes := make([]Volume, 0, 00)
	for _, v := range yamlRef {
		if len(v.Path) == 0 {
			vErrs.AddError(errors.New("empty volume path"), location+".path")
		} else {
			volumes = append(volumes, Volume{Parameters: createParameters(v.Params), Path: v.Path})
		}
	}
	return volumes
}
