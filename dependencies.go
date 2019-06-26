package model

type (
	//Dependencies specifies the stack references on which we depends on
	Dependencies struct {
		//Content lists all the stack references on which we depends on
		Content []StackRef
	}
)

func (r Dependencies) inherit(parent Dependencies) Dependencies {
	dst := Dependencies{}
	dst.Content = make([]StackRef, 0, 0)
	// Set used to avoid duplicated entries
	set := make(map[string]struct{})
	for _, v := range r.Content {
		set[v.ref] = struct{}{}
		dst.Content = append(dst.Content, v)
	}
	for _, v := range parent.Content {
		if _, ok := set[v.ref]; !ok {
			dst.Content = append(dst.Content, v)
		}
	}

	return dst
}

func createDependencies(env *Environment, location DescriptorLocation, dependent string, dependencies []string) Dependencies {
	res := Dependencies{}
	for _, v := range dependencies {
		// A dependent cannot depend on itself or on an empty dependency
		if v != "" && v != dependent {
			depLocation := location.appendPath(v)
			dep := StackRef{
				ref:      v,
				location: depLocation,
				env:      env,
			}
			res.Content = append(res.Content, dep)
		}
	}
	return res
}
