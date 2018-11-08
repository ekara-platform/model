package model

import (
	"strings"
)

// A Describable is a type that can describe itself with a type and a name
type Describable interface {
	DescType() string
	DescName() string
}

type chained struct {
	descTypes []string
	descNames []string
}

func (c chained) DescType() string {
	return strings.Join(c.descTypes, "-")
}

func (c chained) DescName() string {
	return strings.Join(c.descNames, "-")
}

func ChainDescribable(descs ...Describable) chained {
	r := chained{}
	for _, v := range descs {
		r.descTypes = append(r.descTypes, v.DescType())
		r.descNames = append(r.descNames, v.DescName())
	}
	return r
}
