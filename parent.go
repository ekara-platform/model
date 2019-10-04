package model

import (
	"errors"
)

//Parent Represents the parent used to run Ekara
type Parent Component

//CreateParent creates the parent
func CreateParent(base Base, yamlEkara yamlEkara) (Parent, bool, error) {
	repo := yamlEkara.Parent.Repository
	if repo == "" {
		//If the parent is not specified we return an nil parent
		return Parent{}, false, nil
	}
	repoParent, e := CreateRepository(base, repo, yamlEkara.Parent.Ref, "")
	if e != nil {
		return Parent{}, false, errors.New("invalid parent repository: " + e.Error())
	}
	repoParent.setAuthentication(yamlEkara.Parent)
	c := CreateComponent(EkaraComponentId, repoParent)
	return Parent(c), true, nil
}

//Component returns the referenced component
func (p Parent) Component() (Component, error) {
	return Component(p), nil
}

//ComponentName returns the referenced component name
func (p Parent) ComponentName() string {
	return EkaraComponentId
}
