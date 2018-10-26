package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
)

func TestBuildComponentInfo(t *testing.T) {
	// org/repo are prefixed with base
	s := "dummy_org/dummy_repo"
	baseUrl, _ := url.Parse("https://somebase.org")
	u, e := resolveRepositoryInfo(baseUrl, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "somebase.org", u.Host)
	assert.Equal(t, "/"+s+".git", u.Path)

	// local file
	u, e = resolveRepositoryInfo(&url.URL{}, "testdata/dummy_org/dummy_repo")
	assert.Nil(t, e)
	assert.Equal(t, "file", u.Scheme)
	assert.Equal(t, "", u.Host)
	absPath, _ := filepath.Abs("testdata/dummy_org/dummy_repo")
	absPath = filepath.ToSlash(absPath)
	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}
	assert.Equal(t, absPath+"/", u.Path)
}
