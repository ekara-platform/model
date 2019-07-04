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

func TestProviderMerge(t *testing.T) {
	p1 := Provider{Name: "my_name"}
	p1.EnvVars = make(map[string]string)
	p1.EnvVars["key1"] = "val1"
	p1.EnvVars["key2"] = "val2"
	p1.Parameters = make(map[string]interface{})
	p1.Parameters["key1"] = "val1"
	p1.Parameters["key2"] = "val2"

	other := Provider{Name: "my_name"}
	other.EnvVars = make(map[string]string)
	other.EnvVars["key2"] = "val2_not_overwritten"
	other.EnvVars["key3"] = "val3"
	other.Parameters = make(map[string]interface{})
	other.Parameters["key2"] = "val2_not_overwritten"
	other.Parameters["key3"] = "val3"

	err := p1.merge(other)
	if assert.Nil(t, err) {
		if assert.Equal(t, 3, len(p1.EnvVars)) {
			checkMap(t, p1.EnvVars, "key1", "val1")
			checkMap(t, p1.EnvVars, "key2", "val2")
			checkMap(t, p1.EnvVars, "key3", "val3")
		}

		if assert.Equal(t, 3, len(p1.Parameters)) {
			checkMapInterface(t, p1.Parameters, "key1", "val1")
			checkMapInterface(t, p1.Parameters, "key2", "val2")
			checkMapInterface(t, p1.Parameters, "key3", "val3")
		}
	}
}

func TestProvidersMerge(t *testing.T) {
	p1 := Provider{Name: "my_name"}
	p1.EnvVars = make(map[string]string)
	p1.EnvVars["key1"] = "val1"
	p1.EnvVars["key2"] = "val2"
	p1.Parameters = make(map[string]interface{})
	p1.Parameters["key1"] = "val1"
	p1.Parameters["key2"] = "val2"

	other := Provider{Name: "my_name"}
	other.EnvVars = make(map[string]string)
	other.EnvVars["key2"] = "val2_not_overwritten"
	other.EnvVars["key3"] = "val3"
	other.Parameters = make(map[string]interface{})
	other.Parameters["key2"] = "val2_not_overwritten"
	other.Parameters["key3"] = "val3"

	p1s := make(map[string]Provider)
	p1s["myP"] = p1
	others := make(map[string]Provider)
	others["myP"] = other

	pms, err := Providers(p1s).merge(&Environment{}, Providers(others))
	if assert.Nil(t, err) {
		if assert.Equal(t, 1, len(pms)) {
			p := pms["myP"]
			if assert.Equal(t, 3, len(p.EnvVars)) {
				checkMap(t, p.EnvVars, "key1", "val1")
				checkMap(t, p.EnvVars, "key2", "val2")
				checkMap(t, p.EnvVars, "key3", "val3")
			}

			if assert.Equal(t, 3, len(p.Parameters)) {
				checkMapInterface(t, p.Parameters, "key1", "val1")
				checkMapInterface(t, p.Parameters, "key2", "val2")
				checkMapInterface(t, p.Parameters, "key3", "val3")
			}
		}
	}
}

func TestProvidersEmptyMerge(t *testing.T) {
	other := Provider{Name: "my_name"}
	other.EnvVars = make(map[string]string)
	other.EnvVars["key2"] = "val2"
	other.EnvVars["key3"] = "val3"
	other.Parameters = make(map[string]interface{})
	other.Parameters["key2"] = "val2"
	other.Parameters["key3"] = "val3"

	p1s := make(map[string]Provider)

	others := make(map[string]Provider)
	others["myP"] = other

	pms, err := Providers(p1s).merge(&Environment{}, Providers(others))
	if assert.Nil(t, err) {
		if assert.Equal(t, 1, len(pms)) {
			p := pms["myP"]
			if assert.Equal(t, 2, len(p.EnvVars)) {
				checkMap(t, p.EnvVars, "key2", "val2")
				checkMap(t, p.EnvVars, "key3", "val3")
			}

			if assert.Equal(t, 2, len(p.Parameters)) {
				checkMapInterface(t, p.Parameters, "key2", "val2")
				checkMapInterface(t, p.Parameters, "key3", "val3")
			}
		}
	}
}
