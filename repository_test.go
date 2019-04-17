package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepositoryWrongExtension(t *testing.T) {
	b, e := CreateBase("")
	assert.Nil(t, e)
	_, e = CreateRepository(b, "repo", "master", "descriptor.txt")
	if assert.NotNil(t, e) {
		assert.Equal(t, unsupportedFileExtension, e.Error())
	}
}

func TestCreateRepositoryYamlExtension(t *testing.T) {
	b, e := CreateBase("")
	assert.Nil(t, e)
	_, e = CreateRepository(b, "repo", "master", "descriptor.yaml")
	assert.Nil(t, e)
}

func TestCreateRepositoryYmlExtension(t *testing.T) {
	b, e := CreateBase("")
	assert.Nil(t, e)
	_, e = CreateRepository(b, "repo", "master", "descriptor.yml")
	assert.Nil(t, e)
}

func TestCreateRepositoryDefault(t *testing.T) {
	b, e := CreateBase("")
	assert.Nil(t, e)
	r, e := CreateRepository(b, "repo", "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, DefaultDescriptorName, r.DescriptorName)
	}
}

func TestCreateSimpleRepository(t *testing.T) {
	repo := "organisation/repo"
	b, e := CreateBase("")
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, GitScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeHttps)
		// Simple repositories are supposed to be prefixed by a base. In this case
		// The default base DefaultComponentBase
		assert.Equal(t, r.Url.String(), DefaultComponentBase+"/"+repo+GitExtension)
	}
}

func TestCreateExplicitRepository(t *testing.T) {
	repo := "http://github.my_company.com/organisation/repo"
	b, e := CreateBase("")
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, GitScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeHttp)
		// Explicit repositories are not supposed to be prefixed by any base
		assert.Equal(t, r.Url.String(), repo+GitExtension)
	}
}

func TestCreateBasedSimpleRepository(t *testing.T) {
	base := "http://github.my_company.com"
	repo := "organisation/repo"
	b, e := CreateBase(base)
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, GitScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeHttp)
		assert.Equal(t, r.Url.String(), base+"/"+repo+GitExtension)
	}
}

func TestCreateExplicitSvnRepository(t *testing.T) {
	repo := "svn://github.my_company.com/organisation/repo"
	b, e := CreateBase("")
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, SvnScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeSvn)
		// Explicit repositories are not supposed to be prefixed by any base
		assert.Equal(t, r.Url.String(), repo+"/")
	}
}

func TestCreateBasedSimpleSvnRepository(t *testing.T) {
	base := "svn://svnblabla.my_company.com"
	repo := "organisation/repo"
	b, e := CreateBase(base)
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, SvnScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeSvn)
		assert.Equal(t, r.Url.String(), base+"/"+repo+"/")
	}
}

func TestCreateExplicitGitRepository(t *testing.T) {
	repo := "git://github.my_company.com/organisation/repo"
	b, e := CreateBase("")
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, GitScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeGits)
		// Explicit repositories are not supposed to be prefixed by any base
		assert.Equal(t, r.Url.String(), repo+GitExtension)
	}
}

func TestCreateBasedSimpleGitRepository(t *testing.T) {
	base := "git://svnblabla.my_company.com"
	repo := "organisation/repo"
	b, e := CreateBase(base)
	assert.Nil(t, e)
	r, e := CreateRepository(b, repo, "master", "")
	if assert.Nil(t, e) {
		assert.Equal(t, r.Ref, "master")
		assert.Equal(t, r.Scm, GitScm)
		assert.Equal(t, r.Url.UpperScheme(), SchemeGits)
		assert.Equal(t, r.Url.String(), base+"/"+repo+GitExtension)
	}
}
