package model

import (
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

	// GlobalVolume contains the specifications of a shared volume to create
	GlobalVolume struct {
		Content []VolumeContent
	}

	VolumeContent struct {
		// The component holding the content to copy into the volume
		Component componentRef
		// The path, whithin the component, of the content to copy
		Path string
	}

	// GlobalVolume represents all the volumes shared across the whole environment
	GlobalVolumes map[string]*GlobalVolume
)

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

func createGlobalVolumes(env *Environment, location DescriptorLocation, yamlEnv *yamlEnvironment) GlobalVolumes {
	res := make(map[string]*GlobalVolume)
	for name, yamlVol := range yamlEnv.Volumes {
		gv := GlobalVolume{}
		gv.Content = make([]VolumeContent, len(yamlVol.Content))
		for _, v := range yamlVol.Content {
			gv.Content = append(gv.Content, VolumeContent{
				Component: createComponentRef(env, location.appendPath("component"), v.Component, false),
				Path:      v.Path,
			})
		}
		res[name] = &gv
	}
	return res
}
