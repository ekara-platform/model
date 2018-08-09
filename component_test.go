package model

import (
	"net/url"
	"testing"

	"crypto/sha1"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"strings"
)

func TestBuildComponentInfo(t *testing.T) {
	// GitHub assumed to be on https
	s := GitHubHost + "/some/blablabla"
	u, e := ResolveRepositoryInfo(&url.URL{}, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)

	// BitBucket assumed to be on https
	s = BitBucketHost + "/some/blablabla"
	u, e = ResolveRepositoryInfo(&url.URL{}, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)

	// org/repo are prefixed with base
	s = "dummy_org/dummy_repo"
	baseUrl, _ := url.Parse("https://somebase.org")
	u, e = ResolveRepositoryInfo(baseUrl, s)
	assert.Nil(t, e)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "somebase.org", u.Host)
	assert.Equal(t, "/"+s+".git", u.Path)

	// local file
	u, e = ResolveRepositoryInfo(&url.URL{}, "testdata/dummy_org/dummy_repo")
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

func hashUrl(u *url.URL) string {
	hash := sha1.New()
	hash.Write([]byte(u.String()))
	return hex.EncodeToString(hash.Sum(nil))
}
