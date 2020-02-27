package model

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefaultComponentBase(t *testing.T) {
	ekara := yamlEkara{
		Base: "",
	}
	b, e := CreateComponentBase(ekara)
	assert.Nil(t, e)
	assert.Equal(t, SchemeHttps, b.Url.UpperScheme())
	assert.True(t, reflect.TypeOf(b.Url) == reflect.TypeOf(RemoteURL{}))
	// The url shoud end with a slash
	assert.True(t, hasSuffixIgnoringCase(b.Url.String(), "/"))
	// The base should be defaulted to DefaultComponentBase
	assert.Equal(t, DefaultComponentBase+"/", b.Url.String())

}

func TestCreateHttpComponentBase(t *testing.T) {
	bs := "http://www.google.com/my_path"
	ekara := yamlEkara{
		Base: bs,
	}
	b, e := CreateComponentBase(ekara)
	assert.Nil(t, e)
	assert.Equal(t, SchemeHttp, b.Url.UpperScheme())
	assert.True(t, reflect.TypeOf(b.Url) == reflect.TypeOf(RemoteURL{}))
	assert.True(t, hasSuffixIgnoringCase(b.Url.String(), "/"))
	assert.Equal(t, bs+"/", b.Url.String())
}

func TestCreateFileComponentBase(t *testing.T) {
	wd, e := os.Getwd()
	assert.Nil(t, e)
	var p string

	if os.PathSeparator == '/' {
		p = wd + "/some/path"
	} else if os.PathSeparator == '\\' {
		p = wd + "\\some\\path"
	}

	e = os.MkdirAll(p, 0777)
	assert.Nil(t, e)
	assert.True(t, DirExist(p))

	ekara := yamlEkara{
		Base: p,
	}
	b, e := CreateComponentBase(ekara)
	assert.Nil(t, e)

	assert.Equal(t, SchemeFile, b.Url.UpperScheme())
	assert.True(t, reflect.TypeOf(b.Url) == reflect.TypeOf(FileURL{}))
	assert.True(t, hasSuffixIgnoringCase(b.Url.String(), "/"))
	defer func() {
		var e error
		if os.PathSeparator == '/' {
			e = os.RemoveAll("./some/")
		} else if os.PathSeparator == '\\' {
			e = os.RemoveAll(".\\some\\")
		}
		assert.Nil(t, e)
	}()
}

func TestCreateComponentBaseError(t *testing.T) {
	ekara := yamlEkara{
		Base: "://missing/scheme/should/generate/an/error",
	}
	_, e := CreateComponentBase(ekara)
	assert.NotNil(t, e)
}

func TestDefaulted(t *testing.T) {
	b, e := CreateBase("")
	assert.Nil(t, e)
	assert.True(t, b.Defaulted())

	b, e = CreateBase(DefaultComponentBase)
	assert.Nil(t, e)
	assert.True(t, b.Defaulted())

	b, e = CreateBase(DefaultComponentBase + "/")
	assert.Nil(t, e)
	assert.True(t, b.Defaulted())

	b, e = CreateBase("http://project_base")
	assert.Nil(t, e)
	assert.False(t, b.Defaulted())

}
