package model

import (
	"encoding/json"
	"errors"
)

// Volume contains the specifications of a volume to create
type Volume struct {
	// The mounting path of the created volume
	Name string
	// The parameters required to create the volume.
	Parameters Parameters `yaml:"params"`
}

func createVolumes(env *Environment, location DescriptorLocation, yamlRef []yamlVolume) []Volume {
	volumes := make([]Volume, 0, 0)
	for _, v := range yamlRef {
		if len(v.Path) == 0 {
			env.errors.addError(errors.New("empty volume path"), location.appendPath("path"))
		} else {
			volumes = append(volumes, Volume{Parameters: createParameters(v.Params), Name: v.Path})
		}
	}
	return volumes
}

func (r *Volume) merge(other Volume) {
	if r.Name != other.Name {
		panic(errors.New("cannot merge unrelated volumes (" + r.Name + " != " + other.Name + ")"))
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
}

func (r Volume) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string      `json:",omitempty"`
		Parameters *Parameters `json:",omitempty"`
	}{
		Name:       r.Name,
		Parameters: &r.Parameters,
	})
}
