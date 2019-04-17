package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScmGitOnFile(t *testing.T) {
	us := "file:///blablabla"
	checkScm(t, us, GitScm)
}

func TestScmGitOnGit(t *testing.T) {
	us := "git:///blablabla"
	checkScm(t, us, GitScm)
}

func TestScmSvnOnSvn(t *testing.T) {
	us := "svn:///blablabla"
	checkScm(t, us, SvnScm)
}

func TestScmGitOnHttp(t *testing.T) {
	us := "http:///blablabla/my_repo.git"
	checkScm(t, us, GitScm)
}

func TestScmGitOnHttps(t *testing.T) {
	us := "https:///blablabla/my_repo.git"
	checkScm(t, us, GitScm)
}

func checkScm(t *testing.T, rawurl string, wanted SCMType) {
	u, e := url.Parse(rawurl)
	assert.Nil(t, e)

	s, e := resolveSCMType(FileUrl{rootUrl: &rootUrl{url: u}})
	if assert.Nil(t, e) {
		assert.Equal(t, s, wanted)
	}
}

func TestScmUnknown(t *testing.T) {
	us := "dummy:///blablabla/my_repo"
	u, e := url.Parse(us)
	assert.Nil(t, e)

	s, e := resolveSCMType(FileUrl{rootUrl: &rootUrl{url: u}})
	if assert.NotNil(t, e) {
		assert.Equal(t, s, UnknownScm)
	}
}
