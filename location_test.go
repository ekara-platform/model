package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendPath(t *testing.T) {

	d1 := DescriptorLocation{
		Descriptor: "desc",
		Path:       "path",
	}

	d2 := d1.appendPath("suffix")

	assert.Equal(t, d2.Descriptor, d1.Descriptor)
	assert.Equal(t, d2.Path, d1.Path+"."+"suffix")

}

func TestAppendPathComplex(t *testing.T) {

	d := DescriptorLocation{
		Descriptor: "desc",
		Path:       "path1",
	}.appendPath("path2").appendPath("path3").appendPath("path4")

	assert.Equal(t, d.Descriptor, "desc")
	assert.Equal(t, d.Path, "path1.path2.path3.path4")

}

func TestAppendPathOnNoPath(t *testing.T) {

	d1 := DescriptorLocation{
		Descriptor: "desc",
	}

	d2 := d1.appendPath("suffix")

	assert.Equal(t, d2.Descriptor, d1.Descriptor)
	assert.Equal(t, d2.Path, "suffix")

}

func TestAppendPathOnNothing(t *testing.T) {

	d1 := DescriptorLocation{}

	d2 := d1.appendPath("suffix")

	assert.Equal(t, d2.Descriptor, "")
	assert.Equal(t, d2.Path, "suffix")

}

func TestDescriptorLocationEquals(t *testing.T) {
	d1 := DescriptorLocation{
		Descriptor: "desc",
		Path:       "path",
	}

	d2 := DescriptorLocation{
		Descriptor: "desc",
		Path:       "path",
	}
	assert.True(t, d1.equals(d2))

	d2.Descriptor = "dummy"
	assert.False(t, d1.equals(d2))

	d2.Descriptor = "desc"
	d2.Path = "dummy"
	assert.False(t, d1.equals(d2))

	d2.Descriptor = "dummy"
	d2.Path = "dummy"
	assert.False(t, d1.equals(d2))

}
