package model

import (
	"encoding/json"
	"errors"
)

type (
	// Volume contains the specifications of a volume to create
	Volume struct {
		// The mounting path of the created volume
		Path string
		// The parameters required to create the volume.
		Parameters Parameters `yaml:"params"`
	}

	// Volume represents all the volumes to create for a Node set
	Volumes map[string]Volume
)

// MarshalJSON returns the serialized content of a volume as JSON
func (r Volume) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name       string      `json:",omitempty"`
		Parameters *Parameters `json:",omitempty"`
	}{
		Name:       r.Path,
		Parameters: &r.Parameters,
	})
}

func (r *Volume) merge(other Volume) error {
	if r.Path != other.Path {
		return errors.New("cannot merge unrelated volumes (" + r.Path + " != " + other.Path + ")")
	}
	r.Parameters = r.Parameters.inherits(other.Parameters)
	return nil
}

func createVolumes(env *Environment, location DescriptorLocation, yamlRef []yamlVolume) Volumes {
	volumes := Volumes{}
	for _, v := range yamlRef {
		if len(v.Path) == 0 {
			env.errors.addError(errors.New("empty volume path"), location.appendPath("path"))
		} else {
			volumes[v.Path] = Volume{Parameters: createParameters(v.Params), Path: v.Path}
		}
	}
	return volumes
}

func (r Volumes) merge(other Volumes) error {
	for id, v := range other {
		if volume, ok := r[id]; ok {
			if err := volume.merge(v); err != nil {
				return err
			}
		} else {
			r[id] = v
		}
	}
	return nil
}

//AsArray return an array of volumes
func (r Volumes) AsArray() []Volume {
	res := make([]Volume, 0)
	for _, v := range r {
		res = append(res, v)
	}
	return res
}
