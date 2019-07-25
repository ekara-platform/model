package model

import (
	"strings"
)

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

	//Base represents the common location to all components defined into a single descriptor
	Base struct {
		// Url specifies the base location of a component
		Url EkURL
	}
)

//CreateBase a new Base for the provided url, if the url is not specified then
// it will be defaulted to DefaultComponentBase
func CreateBase(rawurl string) (Base, error) {
	b := Base{}
	if rawurl == "" {
		rawurl = DefaultComponentBase
	}
	u, err := CreateUrl(rawurl)
	if err != nil {
		return b, err
	}
	b.Url = u
	return b, nil
}

//CreateComponentBase returns a new Base for the url specified in the Ekara section of the
// provided environment/descriptor, if the url is not defined then
// it will be defaulted to DefaultComponentBase
func CreateComponentBase(yamlEkara yamlEkara) (Base, error) {
	if yamlEkara.Base == "" {
		return CreateBase(DefaultComponentBase)
	}
	return CreateBase(yamlEkara.Base)
}

//CreateBasedUrl creates a url under the base location
func (b Base) CreateBasedUrl(repo string) (EkURL, error) {
	return b.Url.ResolveReference(repo)
}

//Defaulted returns true if the base is the defaulted one
func (b Base) Defaulted() bool {
	return strings.TrimRight(DefaultComponentBase, "/ ") == strings.TrimRight(b.Url.String(), "/ ")
}
