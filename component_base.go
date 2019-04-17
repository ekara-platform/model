package model

const (
	//DefaultComponentBase specifies the default base URL where to look for a component
	//
	// For example if component is defined like this:
	//  components:
	//    aws:
	//      repository: ekara-platform/aws-provider
	//      version: 1.2.3
	//
	// We will assume that this is a Git component located in:
	//   https://github.com/: ekara-platform/aws-provider
	//
	DefaultComponentBase = "https://github.com"
)

type (

	//Repository represents a descriptor or component location
	Base struct {
		// Url specifies the base location of a component
		Url EkUrl
	}
)

func CreateBase(base string) (Base, error) {
	b := Base{}
	if base == "" {
		base = DefaultComponentBase
	}
	u, err := CreateUrl(base)
	if err != nil {
		return b, err
	}
	b.Url = u
	return b, nil
}

func CreateComponentBase(yamlEnv *yamlEnvironment) (Base, error) {
	if yamlEnv == nil && yamlEnv.Ekara.Base != "" {
		return CreateBase(DefaultComponentBase)
	} else {
		return CreateBase(yamlEnv.Ekara.Base)
	}

}

func (b Base) CreateBasedUrl(repo string) (EkUrl, error) {
	return b.Url.ResolveReference(repo)

}
