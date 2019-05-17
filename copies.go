package model

type (
	// Copies represents a list of content  to be copied
	// The key of the map is the path where the content should be copied
	// The map content is an array of path patterns to locate the content to be copied
	Copies struct {
		//Content lists all the content to be copies
		Content map[string]Patterns
	}
)

func (r Copies) inherits(parent Copies) Copies {
	dst := Copies{}
	dst.Content = make(map[string]Patterns)
	for k, v := range r.Content {
		// We copy all the original content
		dst.Content[k] = v
	}
	for k, v := range parent.Content {
		// if the parent content is new then we add it
		if _, ok := dst.Content[k]; !ok {
			dst.Content[k] = v
		} else {
			// if it's not new will merge the patterns from the original content and the parent
			dst.Content[k] = dst.Content[k].inherits(v)
		}
	}
	return dst
}

func createCopies(env *Environment, location DescriptorLocation, copies []yamlCopy) Copies {
	//TODO do somthing with the location if one day we decide to validate the copies content
	res := Copies{}
	res.Content = make(map[string]Patterns)
	for _, vPath := range copies {
		patterns := Patterns{}
		for _, vPattern := range vPath.Patterns {
			patterns.Content = append(patterns.Content, vPattern)
		}
		res.Content[vPath.Path] = patterns
	}
	return res
}
