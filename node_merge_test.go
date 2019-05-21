package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeSetsMerge(t *testing.T) {
	env := &Environment{}
	n1 := NodeSet{Name: "n1"}
	n2 := NodeSet{Name: "n2"}

	no1 := NodeSet{Name: "n1"}
	no2 := NodeSet{Name: "n2"}
	no3 := NodeSet{Name: "n3"}

	origin := make(map[string]NodeSet)
	origin[n1.Name] = n1
	origin[n2.Name] = n2

	other := make(map[string]NodeSet)
	other[no1.Name] = no1
	other[no2.Name] = no2
	other[no3.Name] = no3

	origin, err := NodeSets(origin).merge(env, NodeSets(other))
	assert.Nil(t, err)
	assert.Equal(t, len(origin), 3)

	n, ok := origin["n1"]
	assert.True(t, ok)
	assert.Equal(t, n.Name, "n1")

	n, ok = origin["n2"]
	assert.True(t, ok)
	assert.Equal(t, n.Name, "n2")

	n, ok = origin["n3"]
	assert.True(t, ok)
	assert.Equal(t, n.Name, "n3")
}
