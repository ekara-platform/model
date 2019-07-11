package model

import (
	"errors"
)

const (
	//EkaraParent The default repository for the ekara parent
	ekaraParent = "ekara-platform/distribution"
)

//Parent Represents the parent used to run Ekara
type Parent Component

//CreateParent creates the parent
//	Parameters
//
//		base: the base URL where to look for the parent
//		yamlEnv: the descriptor defining the parent
func CreateParent(base Base, yamlEnv *yamlEnvironment) (Parent, error) {
	defaulted := false
	repo := yamlEnv.Ekara.Parent.Repository
	if repo == "" {
		//If the parent is not specified we must look for the default Ekara one
		// even if the project has defined its own base.
		base, _ = CreateBase("")
		repo = ekaraParent
		defaulted = true
	}
	repoDist, e := CreateRepository(base, repo, yamlEnv.Ekara.Parent.Ref, "")
	if e != nil {
		return Parent{}, errors.New("invalid parent repository: " + e.Error())
	}
	if !defaulted {
		repoDist.setAuthentication(yamlEnv.Ekara.Parent)
	}
	c := CreateComponent(EkaraComponentId, repoDist)
	return Parent(c), nil
}

//Component returns the referenced component
func (p Parent) Component() (Component, error) {
	return Component(p), nil
}

//ComponentName returns the referenced component name
func (p Parent) ComponentName() string {
	return EkaraComponentId
}
