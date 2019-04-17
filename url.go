package model

import (
	"fmt"
	"io/ioutil"

	"net/http"
	"net/url"
	"os"
	pathIm "path"
	"path/filepath"
	"strings"
)

type (
	//EkUrl defines the url used into Ekara
	EkUrl interface {
		String() string
		//ReadUrl returns the content referenced by the url
		ReadUrl() ([]byte, error)
		Scheme() string
		SetScheme(s string)
		Path() string
		AsFilePath() string
		Host() string
		UpperScheme() string
		SetDefaultScheme()
		CheckSlashSuffix()
		ResolveReference(repo string) (EkUrl, error)
		AddPathSuffix(s string)
		RemovePathSuffix(s string)
	}

	rootUrl struct {
		url *url.URL
	}

	//FileUrl defines a local url, typically a file url of something already downloaded into the platform
	FileUrl struct {
		*rootUrl
		filePath string
	}

	//RemoteUrl defines a remote url,example over http; https, git...
	RemoteUrl struct {
		*rootUrl
	}
)

/////////////////////////////////////////////

func (ru *rootUrl) Path() string {
	return ru.url.Path
}

func (ru *rootUrl) Host() string {
	return ru.url.Host
}

func (ru *rootUrl) Scheme() string {
	return ru.url.Scheme
}

func (ru *rootUrl) SetScheme(s string) {
	ru.url.Scheme = s
}

func (ru *rootUrl) UpperScheme() string {
	return strings.ToUpper(ru.url.Scheme)
}

func (ru *rootUrl) SetDefaultScheme() {
	// If no protocol, assume file
	if ru.Scheme() == SchemeUnknown {
		ru.SetScheme(strings.ToLower(SchemeFile))
	}
}

func (ru *rootUrl) CheckSlashSuffix() {
	if !strings.HasSuffix(ru.Path(), "/") {
		ru.AddPathSuffix("/")
	}
}

func (ru *rootUrl) AddPathSuffix(s string) {
	ru.url.Path = ru.url.Path + s
}

func (ru *rootUrl) RemovePathSuffix(s string) {
	ru.url.Path = strings.TrimRight(ru.url.Path, s)
}

func (ru *rootUrl) String() string {
	return ru.url.String()
}

func (ru RemoteUrl) ResolveReference(repo string) (EkUrl, error) {
	repo = strings.TrimLeft(repo, "/")
	repoU, e := url.Parse(repo)
	if e != nil {
		return RemoteUrl{}, e
	}
	return CreateUrl(ru.url.ResolveReference(repoU).String())
}

func (ru FileUrl) ResolveReference(repo string) (EkUrl, error) {
	repo = strings.TrimLeft(repo, "/")
	repoU := ru.filePath + filepath.ToSlash(repo)
	return CreateUrl(repoU)
}

/////////////////////////////////////////////
func createFileUrl(path string) (EkUrl, error) {
	absPath, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	if DirExist(absPath) {
		absPath = absPath + string(filepath.Separator)
	}
	r := FileUrl{filePath: absPath}

	absPath = filepath.ToSlash(absPath)
	absPath = strings.TrimLeft(absPath, "/")
	absPath = "file:///" + absPath
	u, e := url.Parse(absPath)
	if e != nil {
		return r, e
	}
	u.Path = pathIm.Clean(u.Path)
	r.rootUrl = &rootUrl{url: u}
	r.CheckSlashSuffix()
	return r, nil
}

func createRemoteUlr(path string) (EkUrl, error) {
	r := RemoteUrl{}
	u, e := url.Parse(path)
	if e != nil {
		return r, e
	}
	r.rootUrl = &rootUrl{url: u}
	r.CheckSlashSuffix()
	return r, nil
}

func CreateUrl(path string) (EkUrl, error) {
	var r EkUrl
	var e error
	// If file exists locally, resolve its absolute path and convert it to an URL

	if b, _ := FileExist(path); b {
		r, e = createFileUrl(path)
		if e != nil {
			return r, e
		}
	} else {
		r, e = createRemoteUlr(path)
		if e != nil {
			return r, e
		}
	}
	r.SetDefaultScheme()
	return r, nil
}

func (ru FileUrl) ReadUrl() ([]byte, error) {
	location := ru.filePath

	file, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (ru FileUrl) AsFilePath() string {
	return ru.filePath
}

func (ru RemoteUrl) ReadUrl() ([]byte, error) {
	var response *http.Response
	response, err := http.Get(ru.url.String())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		err = fmt.Errorf("error reading URL %s, HTTP status %d", ru.url.String(), response.StatusCode)
		return nil, err
	}
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (ru RemoteUrl) AsFilePath() string {
	return ""
}
