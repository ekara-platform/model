package model

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveUnknownReference(t *testing.T) {

	//First we try to resolve un unknown reference
	cr := componentRef{
		env: &Environment{
			ekara: &Platform{
				Components: make(map[string]Component),
			},
		},
		ref: "unknown_ref",
	}

	_, e := cr.resolve()
	assert.NotNil(t, e)
	assert.Equal(t, e.Error(), fmt.Sprintf(unknownComponentRefError, cr.ref))

	//Now we add a known component and we will resolve it
	cr.env.ekara.Components["known_ref"] = Component{Id: "known_ref"}
	cr.ref = "known_ref"

	c, e := cr.resolve()
	assert.Nil(t, e)
	assert.Equal(t, cr.ref, c.Id)

}

func TestMergeComponentRefItself(t *testing.T) {
	cr := componentRef{
		ref:       "my_ref",
		mandatory: true,
		location: DescriptorLocation{
			Descriptor: "my_descriptor",
			Path:       "my_path",
		},
	}

	err := cr.customize(cr)
	assert.Nil(t, err)
	assert.Equal(t, cr.ref, "my_ref")
	assert.Equal(t, cr.location.Descriptor, "my_descriptor")
	assert.Equal(t, cr.location.Path, "my_path")
	assert.True(t, cr.mandatory)
}

func TestMergeComponentRefNoRef(t *testing.T) {
	cr := componentRef{
		ref:       "my_ref",
		mandatory: true,
		location: DescriptorLocation{
			Descriptor: "my_descriptor",
			Path:       "my_path",
		},
	}

	other := componentRef{
		ref:       "",
		mandatory: true,
	}

	err := cr.customize(other)
	assert.Nil(t, err)
	assert.Equal(t, cr.ref, "my_ref")
	assert.Equal(t, cr.location.Descriptor, "my_descriptor")
	assert.Equal(t, cr.location.Path, "my_path")
	assert.True(t, cr.mandatory)
}

func TestMergeComponentRefOptional(t *testing.T) {
	cr := componentRef{
		ref:       "my_ref",
		mandatory: true,
		location: DescriptorLocation{
			Descriptor: "my_descriptor",
			Path:       "my_path",
		},
	}

	other := componentRef{
		mandatory: false,
	}

	err := cr.customize(other)
	assert.Nil(t, err)
	assert.Equal(t, cr.ref, "my_ref")
	assert.Equal(t, cr.location.Descriptor, "my_descriptor")
	assert.Equal(t, cr.location.Path, "my_path")
	assert.False(t, cr.mandatory)
}

func TestMergeComponentRefLocation(t *testing.T) {
	cr := componentRef{
		ref:       "my_ref",
		mandatory: true,
		location: DescriptorLocation{
			Descriptor: "my_descriptor",
			Path:       "my_path",
		},
	}

	other := componentRef{
		location: DescriptorLocation{
			Descriptor: "other_descriptor",
			Path:       "other_path",
		},
	}

	err := cr.customize(other)
	assert.Nil(t, err)
	assert.Equal(t, cr.ref, "my_ref")
	assert.Equal(t, cr.location.Descriptor, "other_descriptor")
	assert.Equal(t, cr.location.Path, "other_path")
	assert.False(t, cr.mandatory)
}
