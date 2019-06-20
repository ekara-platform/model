package model

import (
	"errors"
)

const (
	//EkaraDistribution The default repository for the ekara distribution
	ekaraDistribution = "ekara-platform/distribution"
)

//Distribution Represents the distribution used to run Ekara
type Distribution Component

//CreateDistribution creates the distribution
//	Parameters
//
//		base: the base URL where to look for the distribution
//		yamlEnv: the descriptor defining the distribution
func CreateDistribution(base Base, yamlEnv *yamlEnvironment) (Distribution, error) {
	defaulted := false
	repo := yamlEnv.Ekara.Distribution.Repository
	if repo == "" {
		//If the distribution is not specified we must look for the default Ekara one
		// even if the project has defined its own base.
		base, _ = CreateBase("")
		repo = ekaraDistribution
		defaulted = true
	}
	repoDist, e := CreateRepository(base, repo, yamlEnv.Ekara.Distribution.Ref, "")
	if e != nil {
		return Distribution{}, errors.New("invalid distribution repository: " + e.Error())
	}
	if !defaulted {
		repoDist.setAuthentication(yamlEnv.Ekara.Distribution)
	}
	c := CreateComponent(EkaraComponentId, repoDist)
	return Distribution(c), nil
}

//Component returns the referenced component
func (r Distribution) Component() (Component, error) {
	return Component(r), nil
}

//ComponentName returns the referenced component name
func (r Distribution) ComponentName() string {
	return EkaraComponentId
}
