package model

type (

	// Reference is the type used to identify a reference on a remote block whitin
	// the environment descriptor
	Reference struct {
		Id        string
		Type      string
		Mandatory bool
		Location  DescriptorLocation
		Repo      map[string]interface{}
	}

	ValidableReference interface {
		Reference() Reference
	}
)
