package model

type (
	// Copies represents a list of content to be copied
	// The key of the map is the path where the content should be copied
	// The map content is an array of path patterns to locate the content to be copied
	Copies struct {
		//Content lists all the content to be copies
		Content map[string]Copy
	}

	// Copy represents a content to be copied
	Copy struct {
		//Once indicates if the copy should be done only on one node matching the targeted labels
		Once bool
		// Labels identifies the nodesets where to copy
		Labels Labels
		// Path identifies the destination path of the copy
		Path string
		//Sources identifies the content to copy
		Sources Patterns
	}
)

func (r Copies) inherit(parent Copies) Copies {
	dst := Copies{}
	dst.Content = make(map[string]Copy)
	for k, v := range r.Content {
		// We copy all the original content
		dst.Content[k] = v
	}
	for k, v := range parent.Content {
		// if the parent content is new then we add it
		if _, ok := dst.Content[k]; !ok {
			dst.Content[k] = v
		} else {
			// if it's not new we will merge the patterns/labels from the original content and the parent
			work := dst.Content[k]
			work.Sources = work.Sources.inherit(v.Sources)
			work.Labels = work.Labels.inherit(v.Labels)
			if work.Path == "" {
				// Only override path if none specified
				work.Path = v.Path
			}
			if v.Once == true {
				// only override once if true (meaning if it's true, it's forever true in children)
				work.Once = true
			}
			dst.Content[k] = work
		}
	}
	return dst
}

func createCopies(env *Environment, location DescriptorLocation, copies map[string]yamlCopy) Copies {
	res := Copies{}
	res.Content = make(map[string]Copy)
	for cpName, yCop := range copies {
		theCopy := Copy{
			Once:   yCop.Once,
			Labels: yCop.Labels,
		}
		sources := Patterns{}
		for _, vSource := range yCop.Sources {
			sources.Content = append(sources.Content, vSource)
		}
		theCopy.Sources = sources
		theCopy.Path = yCop.Path
		res.Content[cpName] = theCopy
	}
	return res
}
