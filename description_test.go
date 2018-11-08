package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {

	n1 := NodeSet{Name: "NameNode1"}
	n2 := NodeSet{Name: "NameNode2"}
	p1 := Provider{Name: "NameProvider1"}
	p2 := Provider{Name: "NameProvider2"}

	chained := ChainDescribable(n1, n2, p1, p2)

	assert.Equal(t, "NameNode1-NameNode2-NameProvider1-NameProvider2", chained.DescName())
	assert.Equal(t, "NodeSet-NodeSet-Provider-Provider", chained.DescType())
}
