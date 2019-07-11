package model

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	assert.Equal(t, u.Path, ru.Path())
}

func TestScheme(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	assert.Equal(t, u.Scheme, ru.Scheme())
	assert.Equal(t, strings.ToUpper(SchemeHttp), ru.UpperScheme())
	assert.Equal(t, SchemeHttp, ru.UpperScheme())
}

func TestSetScheme(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	assert.Equal(t, SchemeHttp, ru.UpperScheme())
	ru.SetScheme(SchemeHttps)
	assert.Equal(t, SchemeHttps, ru.UpperScheme())
}

func TestDefaultScheme(t *testing.T) {
	u, e := url.Parse("//www.google.com/my_path")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	assert.Equal(t, SchemeUnknown, ru.Scheme())
	ru.SetDefaultScheme()
	assert.Equal(t, SchemeFile, ru.UpperScheme())
}

func TestCheckSlashSuffixHttpOk(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	ru.CheckSlashSuffix()
	assert.True(t, hasSuffixIgnoringCase(ru.Path(), "/my_path/"))
}

func TestCheckSlashSuffixHttpKo(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path/")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	ru.CheckSlashSuffix()
	assert.True(t, hasSuffixIgnoringCase(ru.Path(), "/my_path/"))
}

func TestAddSuffix(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path/")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	ru.AddPathSuffix("dummySuffix")
	assert.True(t, hasSuffixIgnoringCase(ru.Path(), "/my_path/dummySuffix"))
}

func TestRemoveSuffix(t *testing.T) {
	u, e := url.Parse("http://www.google.com/my_path/")
	assert.Nil(t, e)
	ru := rootURL{url: u}
	ru.RemovePathSuffix("/")
	assert.True(t, hasSuffixIgnoringCase(ru.Path(), "/my_path"))
}

func TestCreateRemoteUlr(t *testing.T) {
	p := "http://www.google.com/my_path"
	u, e := createRemoteUlr(p)
	assert.Nil(t, e)
	assert.Equal(t, SchemeHttp, u.UpperScheme())
	assert.True(t, hasSuffixIgnoringCase(u.String(), "/"))
	assert.Equal(t, "/my_path/", u.Path())

	val, ok := u.(RemoteURL)
	assert.True(t, ok)
	// The slash suffix is almays added
	assert.Equal(t, val.rootURL.String(), p+"/")
}

func TestCreateFileUlr(t *testing.T) {

	var p string
	wd, e := os.Getwd()
	assert.Nil(t, e)

	if os.PathSeparator == '/' {
		p = wd + "/some/base"
	} else if os.PathSeparator == '\\' {
		p = wd + "\\some\\base"
	}

	u, e := createFileURL(p)
	assert.Nil(t, e)

	os.Remove(p)
	assert.Nil(t, e)

	assert.Equal(t, SchemeFile, u.UpperScheme())
	assert.True(t, hasSuffixIgnoringCase(u.String(), "/"))
	val, ok := u.(FileURL)
	assert.True(t, ok)
	if os.PathSeparator == '/' {
		assert.Equal(t, p, val.filePath)
	} else if os.PathSeparator == '\\' {
		assert.Equal(t, p, val.filePath)
	}

	ps := filepath.ToSlash(p)

	if hasPrefixIgnoringCase(wd, "/") {
		assert.Equal(t, val.rootURL.String(), "file://"+ps+"/")
	} else {
		assert.Equal(t, val.rootURL.String(), "file:///"+ps+"/")
	}

}

func TestCreateBasedRemoteURL(t *testing.T) {

	p1 := "http://github.com"
	p2 := "organisation/repository"

	env := &yamlEnvironment{
		Ekara: yamlEkara{Base: p1},
	}
	b, e := CreateComponentBase(env)
	assert.Nil(t, e)
	u, e := b.CreateBasedUrl(p2)
	assert.Nil(t, e)
	assert.Equal(t, SchemeHttp, u.UpperScheme())

	val, ok := u.(RemoteURL)
	assert.True(t, ok)
	assert.Equal(t, val.rootURL.String(), p1+"/"+p2+"/")
}

func TestCreateBasedRemoteURL2(t *testing.T) {

	p1 := "http://github.com/"
	p2 := "/organisation/repository/"

	env := &yamlEnvironment{
		Ekara: yamlEkara{Base: p1},
	}
	b, e := CreateComponentBase(env)
	assert.Nil(t, e)
	u, e := b.CreateBasedUrl(p2)
	assert.Nil(t, e)
	assert.Equal(t, SchemeHttp, u.UpperScheme())

	val, ok := u.(RemoteURL)
	assert.True(t, ok)
	assert.Equal(t, val.rootURL.String(), p1+p2[1:])
}

func TestCreateBasedLocalUrl(t *testing.T) {

	var p1 string
	var p2 string
	wd, e := os.Getwd()
	assert.Nil(t, e)

	if os.PathSeparator == '/' {
		p1 = wd + "/some/base/"
		p2 = "organisation/repository"
	} else if os.PathSeparator == '\\' {
		p1 = wd + "\\some\\base\\"
		p2 = "organisation\\repository"
	}

	e = os.MkdirAll(p1+p2, 0777)
	assert.Nil(t, e)
	assert.True(t, DirExist(p1+p2))

	env := &yamlEnvironment{
		Ekara: yamlEkara{Base: p1},
	}

	b, e := CreateComponentBase(env)
	assert.Nil(t, e)
	assert.Equal(t, SchemeFile, b.Url.UpperScheme())
	u, e := b.CreateBasedUrl(p2)
	assert.Nil(t, e)
	assert.Equal(t, SchemeFile, u.UpperScheme())

	val, ok := u.(FileURL)
	assert.True(t, ok)

	if os.PathSeparator == '/' {
		assert.Equal(t, p1+p2+"/", val.filePath)
	} else if os.PathSeparator == '\\' {
		assert.Equal(t, p1+p2+"\\", val.filePath)
	}

	ps := filepath.ToSlash(p1 + p2 + "/")

	if hasPrefixIgnoringCase(wd, "/") {
		assert.Equal(t, val.rootURL.String(), "file://"+ps)
	} else {
		assert.Equal(t, val.rootURL.String(), "file:///"+ps)
	}

	defer func(wd string) {
		var e error
		if os.PathSeparator == '/' {
			e = os.RemoveAll(wd + "/some/")
		} else if os.PathSeparator == '\\' {
			e = os.RemoveAll(wd + "\\some\\")
		}
		assert.Nil(t, e)
	}(wd)
}
