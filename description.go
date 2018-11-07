package model

// A Describable is a type that can describe itself with a type and a name
type Describable interface {
	DescType() string
	DescName() string
}
