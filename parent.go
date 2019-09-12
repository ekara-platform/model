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
func CreateParent(base Base, yamlEkara yamlEkara) (Parent, error) {
	defaulted := false
	repo := yamlEkara.Parent.Repository
	if repo == "" {
		//If the parent is not specified we must look for the default Ekara one
		// even if the project has defined its own base.
		base, _ = CreateBase("")
		repo = ekaraParent
		defaulted = true
	}
	repoParent, e := CreateRepository(base, repo, yamlEkara.Parent.Ref, "")
	if e != nil {
		return Parent{}, errors.New("invalid parent repository: " + e.Error())
	}
	if !defaulted {
		repoParent.setAuthentication(yamlEkara.Parent)
	}
	c := CreateComponent(EkaraComponentId, repoParent)
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
