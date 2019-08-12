package model

type (
	// GlobalVolume contains the specifications of a shared volume to create
	GlobalVolume struct {
		Content []VolumeContent
	}

	//VolumeContent represents the detail content of a volume
	VolumeContent struct {
		// The component holding the content to copy into the volume
		Component componentRef
		// The path, whithin the component, of the content to copy
		Path string
	}

	//GlobalVolumes represents all the volumes shared across the whole environment
	GlobalVolumes map[string]*GlobalVolume
)

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
