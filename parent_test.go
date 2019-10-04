package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNoParent(t *testing.T) {
	ye := yamlEkara{
		Parent: yamlComponent{
			yamlAuth: yamlAuth{
				Auth: make(map[string]interface{}),
			},
		},
	}
	ye.Parent.Auth["p1"] = "v1"
	ye.Parent.Auth["p2"] = "v2"
	assert.True(t, len(ye.Parent.Auth) == 2)

	b, e := CreateComponentBase(ye)
	assert.Nil(t, e)
	_, bo, e := CreateParent(b, ye)
	assert.Nil(t, e)
	assert.False(t, bo)
}

func TestCreateDefinedParentOverDefinedBase(t *testing.T) {
	pbs := "http://project_base"
	ds := "projectOrganization/customParent"
	ye := yamlEkara{
		Base: pbs,
		Parent: yamlComponent{
			Repository: ds,
			yamlAuth: yamlAuth{
				Auth: make(map[string]interface{}),
			},
		},
	}
	ye.Parent.Auth["p1"] = "v1"
	ye.Parent.Auth["p2"] = "v2"

	b, e := CreateComponentBase(ye)
	assert.Nil(t, e)
	d, bo, e := CreateParent(b, ye)
	assert.Nil(t, e)
	assert.True(t, bo)
	assert.Equal(t, d.Repository.Url.String(), pbs+"/"+ds+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttp)
	// The project parent uses authentication
	assert.True(t, len(d.Repository.Authentication) == 2)
}
