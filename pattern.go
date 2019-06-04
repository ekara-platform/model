package model

type (
	//Patterns represent a list of path patterns
	Patterns struct {
		//Content lists all the path patterns
		Content []string
	}
)

func (r Patterns) inherit(parent Patterns) Patterns {
	dst := Patterns{}
	dst.Content = make([]string, 0, 0)
	// Set used to avoid duplicated entries
	set := make(map[string]struct{})
	for _, v := range r.Content {
		set[v] = struct{}{}
		dst.Content = append(dst.Content, v)
	}
	for _, v := range parent.Content {
		if _, ok := set[v]; !ok {
			dst.Content = append(dst.Content, v)
		}
	}
	return dst
}

func createPatterns(env *Environment, location DescriptorLocation, paths []string) Patterns {
	//TODO do something with the location if one day we decide to validate the patterns content
	res := Patterns{}
	for _, v := range paths {
		res.Content = append(res.Content, v)
	}
	return res
}
