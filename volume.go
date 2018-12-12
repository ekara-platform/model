package model

import (
	"encoding/json"
	"errors"
	"reflect"
)

type (
	// Volume contains the specifications of a volume to create
	Volume struct {
		location DescriptorLocation
		// The mounting path of the created volume
		Path string
		// The parameters required to create the volume.
		Parameters Parameters `yaml:"params"`
	}

	// Volume represents all the volumes to create for a Node set
	Volumes map[string]*Volume
)

// MarshalJSON returns the serialized content of a volume as JSON
func (r Volume) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Path       string      `json:",omitempty"`
		Parameters *Parameters `json:",omitempty"`
	}{
		Path:       r.Path,
		Parameters: &r.Parameters,
	})
}

func (r *Volume) merge(other Volume) error {
	if !reflect.DeepEqual(r, &other) {
		if r.Path != other.Path {
			return errors.New("cannot merge unrelated volumes (" + r.Path + " != " + other.Path + ")")
		}
		r.Parameters = r.Parameters.inherits(other.Parameters)
	}
	return nil
}

func (r Volume) validate() ValidationErrors {
	vErrs := ValidationErrors{}
	if r.Path == "" {
		vErrs.addError(errors.New("empty volume path"), r.location.appendPath("path"))
	}
	return vErrs
}

func createVolumes(location DescriptorLocation, yamlRef []yamlVolume) Volumes {
	volumes := Volumes{}
	for i, v := range yamlRef {
		volumeLocation := location.appendIndex(i)
		volumes[v.Path] = &Volume{location: volumeLocation, Parameters: createParameters(v.Params), Path: v.Path}
	}
	return volumes
}

func (r Volumes) merge(other Volumes) error {
	for id, v := range other {
		if volume, ok := r[id]; ok {
			if err := volume.merge(*v); err != nil {
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
		res = append(res, *v)
	}
	return res
}
