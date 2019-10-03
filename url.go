package model

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	pathIm "path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	//EkURL defines the url used into Ekara
	EkURL interface {
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
		ResolveReference(repo string) (EkURL, error)
		AddPathSuffix(s string)
		RemovePathSuffix(s string)
	}

	rootURL struct {
		url *url.URL
	}

	//FileURL defines a local url, typically a file url of something already downloaded into the platform
	FileURL struct {
		*rootURL
		filePath string
	}

	//RemoteURL defines a remote url,example over http; https, git...
	RemoteURL struct {
		*rootURL
	}
)

//MarshalYAML serialize the url content into YAML
func (fu FileURL) MarshalYAML() (interface{}, error) {
	res, err := yaml.Marshal(&struct {
		URL *url.URL
	}{
		URL: fu.url,
	})
	return string(res), err
}

//MarshalYAML serialize the url content into YAML
func (ru RemoteURL) MarshalYAML() (interface{}, error) {
	res, err := yaml.Marshal(&struct {
		URL *url.URL
	}{
		URL: ru.url,
	})
	return string(res), err
}

func (ru *rootURL) Path() string {
	return ru.url.Path
}

func (ru *rootURL) Host() string {
	return ru.url.Host
}

func (ru *rootURL) Scheme() string {
	return ru.url.Scheme
}

func (ru *rootURL) SetScheme(s string) {
	ru.url.Scheme = s
}

func (ru *rootURL) UpperScheme() string {
	return strings.ToUpper(ru.url.Scheme)
}

func (ru *rootURL) SetDefaultScheme() {
	// If no protocol, assume file
	if ru.Scheme() == SchemeUnknown {
		ru.SetScheme(strings.ToLower(SchemeFile))
	}
}

func (ru *rootURL) CheckSlashSuffix() {
	if !strings.HasSuffix(ru.Path(), "/") {
		ru.AddPathSuffix("/")
	}
}

func (ru *rootURL) AddPathSuffix(s string) {
	ru.url.Path = ru.url.Path + s
}

func (ru *rootURL) RemovePathSuffix(s string) {
	ru.url.Path = strings.TrimRight(ru.url.Path, s)
}

func (ru *rootURL) String() string {
	return ru.url.String()
}

//ResolveReference resolves the repository URI reference to an absolute URI from
// the RemoteURL as base URI
func (ru RemoteURL) ResolveReference(repository string) (EkURL, error) {
	repository = strings.TrimLeft(repository, "/")
	repoU, e := url.Parse(repository)
	if e != nil {
		return RemoteURL{}, e
	}
	return CreateUrl(ru.url.ResolveReference(repoU).String())
}

//ResolveReference resolves the repository URI reference to an absolute URI from
// the FileURL as base URI
func (fu FileURL) ResolveReference(repository string) (EkURL, error) {
	repository = strings.TrimLeft(repository, "/")
	repoU := fu.filePath + filepath.ToSlash(repository)
	return CreateUrl(repoU)
}

/////////////////////////////////////////////
func createFileURL(path string) (EkURL, error) {
	absPath, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	if DirExist(absPath) {
		absPath = absPath + string(filepath.Separator)
	}
	r := FileURL{filePath: absPath}

	absPath = filepath.ToSlash(absPath)
	absPath = strings.TrimLeft(absPath, "/")
	absPath = "file:///" + absPath
	u, e := url.Parse(absPath)
	if e != nil {
		return r, e
	}
	u.Path = pathIm.Clean(u.Path)
	r.rootURL = &rootURL{url: u}
	r.CheckSlashSuffix()
	return r, nil
}

func createRemoteUlr(path string) (EkURL, error) {
	r := RemoteURL{}
	u, e := url.Parse(path)
	if e != nil {
		return r, e
	}
	r.rootURL = &rootURL{url: u}
	r.CheckSlashSuffix()
	return r, nil
}

//CreateUrl creates an Ekara url for the given path. The provided path can be a path on
// a file system or a remote url over http, https, git...
func CreateUrl(path string) (EkURL, error) {
	var r EkURL
	var e error
	// If file exists locally, resolve its absolute path and convert it to an URL

	if b, _ := FileExist(path); b {
		r, e = createFileURL(path)
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

	_, e = url.Parse(r.String())
	if e != nil {
		return r, e
	}

	return r, nil
}

//ReadUrl reads the content referenced by the url
func (fu FileURL) ReadUrl() ([]byte, error) {
	location := fu.filePath

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

//AsFilePath return path corresponding to the file url
func (fu FileURL) AsFilePath() string {
	return fu.filePath
}

//ReadUrl reads the content referenced by the url
func (ru RemoteURL) ReadUrl() ([]byte, error) {
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

//AsFilePath return "" because it's a remote url
func (ru RemoteURL) AsFilePath() string {
	return ""
}

//GetCurrentDirectoryURL return the working directory as an url
func GetCurrentDirectoryURL(l *log.Logger) (EkURL, error) {
	wd, err := os.Getwd()
	if err != nil {
		l.Printf("Error getting the working directory: %s\n", err.Error())
		return FileURL{}, err
	}
	return CreateUrl(wd)
}
