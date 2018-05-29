package model

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"testing"
)

func TestRelativePathToUrl(t *testing.T) {
	absPath, e := filepath.Abs("testdata/components")
	assert.Nil(t, e)
	url, e := PathToUrl("testdata/components")
	assert.Nil(t, e)

	assert.Equal(t, "file", url.Scheme)
	if runtime.GOOS == "windows" {
		assert.Equal(t, "/"+filepath.ToSlash(absPath), url.Path)
	} else {
		assert.Equal(t, filepath.ToSlash(absPath), url.Path)
	}
}

func TestUrlToPath(t *testing.T) {
	absPath, e := filepath.Abs("testdata/components")
	assert.Nil(t, e)
	url, e := PathToUrl("testdata/components")
	assert.Nil(t, e)
	path, e := UrlToPath(url)
	assert.Nil(t, e)

	assert.Equal(t, absPath, path)
}
