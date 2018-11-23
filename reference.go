package model

type (

	// validatableReference is the type used to identify a reference on a remote component whitin
	// the environment descriptor
	validatableReference struct {
		//Id specifies the id/name of the referenced component
		Id string
		//Type specifies the type of the referenced component.
		//Example: NodeSet, Provider...
		Type string
		//Mandatory indicates if the reference is mandatory. A mandatory reference
		//without any referenced component will produce a validation error during the
		// environment validation
		Mandatory bool
		//Location indicates where the reference is located into the descriptor
		Location DescriptorLocation
		//Repo contains the list of components where to look for the one matching
		//the reference
		Repo map[string]interface{}
	}

	// validatableReferencer allows to get a validatable reference to a remote component
	// into the environment descriptor.
	//
	// If a structure containing a reference implements validatableReferencer then
	// this reference will be validated invoking:
	//  ErrorOnInvalid(myStruct)
	validatableReferencer interface {
		reference() validatableReference
	}
)
