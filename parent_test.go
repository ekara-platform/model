package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultParent(t *testing.T) {
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
	d, e := CreateParent(b, ye)
	assert.Nil(t, e)
	// When nothing is specified the parent should be defaulted to the
	// ekara-platform/distribution on github
	assert.Equal(t, d.Repository.Url.String(), DefaultComponentBase+"/"+ekaraParent+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttps)
	// The defaulted parent doesn't use authentication
	assert.True(t, len(d.Repository.Authentication) == 0)
}

func TestCreateDefaultParentOverDefinedBase(t *testing.T) {
	ye := yamlEkara{
		Base: "project_base",
		Parent: yamlComponent{
			yamlAuth: yamlAuth{
				Auth: make(map[string]interface{}),
			},
		},
	}
	ye.Parent.Auth["p1"] = "v1"
	ye.Parent.Auth["p2"] = "v2"

	b, e := CreateComponentBase(ye)
	assert.Nil(t, e)
	d, e := CreateParent(b, ye)
	assert.Nil(t, e)
	// Even if the project defines its on base we need to get the defaulted ditribution
	// ekara-platform/distribution coming from the defaulted base on github
	assert.Equal(t, d.Repository.Url.String(), DefaultComponentBase+"/"+ekaraParent+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttps)
	// The defaulted parent doesn't use authentication
	assert.True(t, len(d.Repository.Authentication) == 0)
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
	d, e := CreateParent(b, ye)
	assert.Nil(t, e)
	assert.Equal(t, d.Repository.Url.String(), pbs+"/"+ds+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttp)
	// The project parent uses authentication
	assert.True(t, len(d.Repository.Authentication) == 2)
}
