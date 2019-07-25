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
