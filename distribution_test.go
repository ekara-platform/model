package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultDistribution(t *testing.T) {
	ye := &yamlEnvironment{
		Ekara: yamlEkara{
			Distribution: yamlComponent{
				yamlAuth: yamlAuth{
					Auth: make(map[string]interface{}),
				},
			},
		},
	}
	ye.Ekara.Distribution.Auth["p1"] = "v1"
	ye.Ekara.Distribution.Auth["p2"] = "v2"
	assert.True(t, len(ye.Ekara.Distribution.Auth) == 2)

	b, e := CreateComponentBase(ye)
	assert.Nil(t, e)
	d, e := CreateDistribution(b, ye)
	assert.Nil(t, e)
	// When nothing is specified the distribution should be defaulted to the
	// ekara-platform/distribution on github
	assert.Equal(t, d.Repository.Url.String(), DefaultComponentBase+"/"+ekaraDistribution+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttps)
	// The defaulted distribution doesn't use authentication
	assert.True(t, len(d.Repository.Authentication) == 0)
}

func TestCreateDefaultDistributionOverDefinedBase(t *testing.T) {
	ye := &yamlEnvironment{
		Ekara: yamlEkara{
			Base: "project_base",
			Distribution: yamlComponent{
				yamlAuth: yamlAuth{
					Auth: make(map[string]interface{}),
				},
			},
		},
	}
	ye.Ekara.Distribution.Auth["p1"] = "v1"
	ye.Ekara.Distribution.Auth["p2"] = "v2"

	b, e := CreateComponentBase(ye)
	assert.Nil(t, e)
	d, e := CreateDistribution(b, ye)
	assert.Nil(t, e)
	// Even if the project defines its on base we need to get the defaulted ditribution
	// ekara-platform/distribution commint from the defaulted base on github
	assert.Equal(t, d.Repository.Url.String(), DefaultComponentBase+"/"+ekaraDistribution+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttps)
	// The defaulted distribution doesn't use authentication
	assert.True(t, len(d.Repository.Authentication) == 0)
}

func TestCreateDefinedDistributionOverDefinedBase(t *testing.T) {
	pbs := "http://project_base"
	ds := "projectOrganization/customDistribution"
	ye := &yamlEnvironment{
		Ekara: yamlEkara{
			Base: pbs,
			Distribution: yamlComponent{
				Repository: ds,
				yamlAuth: yamlAuth{
					Auth: make(map[string]interface{}),
				},
			},
		},
	}
	ye.Ekara.Distribution.Auth["p1"] = "v1"
	ye.Ekara.Distribution.Auth["p2"] = "v2"

	b, e := CreateComponentBase(ye)
	assert.Nil(t, e)
	d, e := CreateDistribution(b, ye)
	assert.Nil(t, e)
	assert.Equal(t, d.Repository.Url.String(), pbs+"/"+ds+GitExtension)
	assert.Equal(t, d.Repository.Url.UpperScheme(), SchemeHttp)
	// The projectdistribution usse authentication
	assert.True(t, len(d.Repository.Authentication) == 2)
}
