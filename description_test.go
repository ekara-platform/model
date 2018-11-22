package model

import (
	"fmt"
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

func ExampleChainDescribable() {
	p := Provider{Name: "MyProviderName"}
	n := NodeSet{Name: "MyNodesetName"}

	c := ChainDescribable(p, n)
	fmt.Printf("Chained types :%s", c.DescType())
	fmt.Printf("Chained names :%s", c.DescName())
	// Output: Chained types :Provider-NodeSet
	// Chained names :MyProviderName-MyNodesetName
}
