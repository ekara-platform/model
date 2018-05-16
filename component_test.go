package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
)

func TestBuildComponentUrl(t *testing.T) {
	// GitHub assumed to be on https
	s := GitHubHost + "/blablabla"
	u, e := ResolveRepositoryUrl(&url.URL{}, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)

	// BitBucket assumed to be on https
	s = BitBucketHost + "/blablabla"
	u, e = ResolveRepositoryUrl(&url.URL{}, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)

	// org/repo are prefixed with base
	s = "dummy_org/dummy_repo"
	baseUrl, _ := url.Parse("https://somebase.org")
	u, e = ResolveRepositoryUrl(baseUrl, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "somebase.org", u.Host)
	assert.Equal(t, "/"+s+".git", u.Path)

	// local file
	u, e = ResolveRepositoryUrl(&url.URL{}, "testdata/dummy_repo")
	assert.Nil(t, e)
	assert.Equal(t, "file", u.Scheme)
	assert.Equal(t, "", u.Host)
	absPath, _ := filepath.Abs("testdata/dummy_repo")
	absPath = filepath.ToSlash(absPath)
	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}
	assert.Equal(t, absPath, u.Path)

}
