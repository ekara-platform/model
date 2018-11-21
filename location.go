package model

import (
	"reflect"
)

type (
	// Type used to describe the location of any element within the environment
	// descriptor
	DescriptorLocation struct {
		Descriptor string
		Path       string
	}
)

func (r DescriptorLocation) equals(o DescriptorLocation) bool {
	return reflect.DeepEqual(r, o)
}

func (r DescriptorLocation) appendPath(suffix string) DescriptorLocation {
	newLoc := DescriptorLocation{Path: r.Path, Descriptor: r.Descriptor}
	if newLoc.Path == "" {
		newLoc.Path = suffix
	} else {
		newLoc.Path = newLoc.Path + "." + suffix
	}
	return newLoc
}
