package model

type (
	//Component represents an element composing an ekara environment
	//
	//A component is always hosted into a source control management system.
	//
	//It can be for example a Provider or Software to deploy on the environment
	//
	Component struct {
		// Id specifies id of the component
		Id string
		// Repository specifies the location where to look for the component
		Repository Repository
		//Templates Defines the content to template for the component
		Templates Patterns
	}

	//ComponentReferencer allows to access to a component through its reference
	ComponentReferencer interface {
		//Component returns the referenced component
		Component() (Component, error)
		//ComponentName returns the referenced component name
		ComponentName() string
	}
)

//CreateComponent creates a new component
//	Parameters
//
//		id: the id of the component
//		repo: the repository where to fetch the component
func CreateComponent(id string, repo Repository) Component {
	return Component{
		//Id is the id of the component
		Id: id,
		//Repository is the repository where to look for the component
		Repository: repo,
	}
}

//Templatable indicates if the component contains templates
func (c Component) Templatable() (bool, Patterns) {
	return len(c.Templates) > 0, c.Templates
}

func (c Component) Component() (Component, error) {
	return c, nil
}

func (c Component) ComponentName() string {
	return c.Id
}
