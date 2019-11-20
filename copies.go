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
			dst.Content[k] = work
		}
	}
	return dst
}

func createCopies(env *Environment, location DescriptorLocation, copies []yamlCopy) Copies {
	res := Copies{}
	res.Content = make(map[string]Copy)
	for _, yCop := range copies {
		copy := Copy{
			Once: yCop.Target.Once,
		}
		copy.Labels = yCop.Target.Labels
		patterns := Patterns{}
		for _, vPattern := range yCop.Patterns {
			patterns.Content = append(patterns.Content, vPattern)
		}
		copy.Sources = patterns
		res.Content[yCop.Target.Path] = copy
	}
	return res
}
