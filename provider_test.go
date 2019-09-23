package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderDescType(t *testing.T) {
	p := Provider{}
	assert.Equal(t, p.DescType(), "Provider")
}

func TestProviderDescName(t *testing.T) {
	p := Provider{Name: "my_name"}
	assert.Equal(t, p.DescName(), "my_name")
}

func getProviderOrigin() *Provider {
	p1 := &Provider{Name: "my_name"}
	p1.EnvVars = make(map[string]string)
	p1.EnvVars["key1"] = "val1_target"
	p1.EnvVars["key2"] = "val2_target"
	p1.Parameters = make(map[string]interface{})
	p1.Parameters["key1"] = "val1_target"
	p1.Parameters["key2"] = "val2_target"

	p1.cRef = componentRef{
		ref: "cOriginal",
	}

	return p1
}

func getProviderOther(name string) *Provider {
	other := &Provider{Name: name}
	other.EnvVars = make(map[string]string)
	other.EnvVars["key2"] = "val2_other"
	other.EnvVars["key3"] = "val3_other"
	other.Parameters = make(map[string]interface{})
	other.Parameters["key2"] = "val2_other"
	other.Parameters["key3"] = "val3_other"

	other.cRef = componentRef{
		ref: "cOther",
	}

	return other
}

func checkProviderMerge(t *testing.T, p *Provider) {

	assert.Equal(t, p.cRef.ref, "cOther")

	if assert.Equal(t, 3, len(p.EnvVars)) {
		checkMap(t, p.EnvVars, "key1", "val1_target")
		checkMap(t, p.EnvVars, "key2", "val2_other")
		checkMap(t, p.EnvVars, "key3", "val3_other")
	}

	if assert.Equal(t, 3, len(p.Parameters)) {
		checkMapInterface(t, p.Parameters, "key1", "val1_target")
		checkMapInterface(t, p.Parameters, "key2", "val2_other")
		checkMapInterface(t, p.Parameters, "key3", "val3_other")
	}
}

func TestProviderMerge(t *testing.T) {

	o := getProviderOrigin()
	err := o.customize(*getProviderOther("my_name"))
	if assert.Nil(t, err) {
		checkProviderMerge(t, o)
	}
}

func TestMergeProviderItself(t *testing.T) {
	o := getProviderOrigin()
	oi := o
	err := o.customize(*o)
	if assert.Nil(t, err) {
		assert.Equal(t, oi, o)
	}
}

func TestProvidersMerge(t *testing.T) {
	origins := make(Providers)
	origins["myP"] = *getProviderOrigin()
	others := make(Providers)
	others["myP"] = *getProviderOther("my_name")

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		if assert.Len(t, customized, 1) {
			o := customized["myP"]
			checkProviderMerge(t, &o)
		}
	}
}

func TestProvidersMergeAddition(t *testing.T) {
	origins := make(Providers)
	origins["myP"] = *getProviderOrigin()
	others := make(Providers)
	others["myP"] = *getProviderOther("my_name")
	others["new"] = *getProviderOther("new")

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		assert.Len(t, customized, 2)
	}
}

func TestProvidersEmptyMerge(t *testing.T) {
	origins := make(Providers)
	o := getProviderOrigin()
	origins["myP"] = *o
	others := make(Providers)

	customized, err := origins.customize(&Environment{}, others)
	if assert.Nil(t, err) {
		if assert.Len(t, customized, 1) {
			oc := customized["myP"]
			assert.Equal(t, *o, oc)
		}
	}
}

func TestMergeProviderUnrelated(t *testing.T) {
	pro := Provider{
		Name: "Name",
	}

	o := Provider{
		Name: "Dummy",
	}

	err := pro.customize(o)
	if assert.NotNil(t, err) {
		assert.Equal(t, err.Error(), "cannot customize unrelated providers (Name != Dummy)")
	}
}
