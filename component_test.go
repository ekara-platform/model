package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildComponentFolderUrl(t *testing.T) {
	// GitHub assumed to be on https
	s := GitHubHost + "/blablabla"
	u, e := BuildComponentFolderUrl(s)
	assert.Nil(t, e)
	assert.Equal(t, u.Scheme, "https")

	// BitBucket assumed to be on https
	s = BitBucketHost + "/blablabla"
	u, e = BuildComponentFolderUrl(s)
	assert.Nil(t, e)
	assert.Equal(t, u.Scheme, "https")

	// oraganization/repo assumed located on GitHub on https
	s = "dummy_org/dummy_repo"
	u, e = BuildComponentFolderUrl(s)
	assert.Nil(t, e)
	assert.Equal(t, u.Scheme, "https")
	assert.Equal(t, u.Host, GitHubHost)

}

func TestBuildComponentGitUrl(t *testing.T) {

	// Github on http(s)
	u := url.URL{Scheme: "https", Host: GitHubHost, Path: "blablabla"}
	cUrl, e := BuildComponentGitUrl(u)
	assert.Nil(t, e)
	assert.Equal(t, u.Path+".git", cUrl.Path)

	// BitBucket on http(s)
	u = url.URL{Scheme: "https", Host: BitBucketHost, Path: "blablabla"}
	cUrl, e = BuildComponentGitUrl(u)
	assert.Nil(t, e)
	assert.Equal(t, u.Path+".git", cUrl.Path)

	// Bad scheme --> Not eligible
	u = url.URL{Scheme: "SVN", Host: "dummy", Path: "blablabla"}
	cUrl, e = BuildComponentGitUrl(u)
	assert.NotNil(t, e)

	// Already a GIT repo Url --> Not eligible
	u = url.URL{Scheme: "https", Host: GitHubHost, Path: "blablabla.git"}
	cUrl, e = BuildComponentGitUrl(u)
	assert.NotNil(t, e)

	// Bad host --> Not eligible
	u = url.URL{Scheme: "https", Host: "dummy", Path: "blablabla"}
	cUrl, e = BuildComponentGitUrl(u)
	assert.NotNil(t, e)
}
