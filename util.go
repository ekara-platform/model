package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func PathToUrl(path string) (*url.URL, error) {
	absPath, e := filepath.Abs(path)
	if e != nil {
		return nil, e
	}
	if fi, e := os.Stat(absPath); e == nil && fi.IsDir() {
		absPath = absPath + string(filepath.Separator)
	}
	absPath = filepath.ToSlash(absPath)
	if strings.HasPrefix(absPath, "/") {
		path = "file://" + filepath.ToSlash(absPath)
	} else {
		path = "file:///" + filepath.ToSlash(absPath)
	}
	u, e := url.Parse(path)
	if e != nil {
		return nil, e
	}
	return u, nil
}

func UrlToPath(u *url.URL) (string, error) {
	if strings.ToUpper(u.Scheme) != "FILE" {
		return "", errors.New("not a valid local URL: " + u.String())
	}
	p := filepath.FromSlash(u.Path)
	if strings.HasPrefix(p, "\\") {
		// windows paths should be stripped from first character
		p = p[1:]
	}
	p = filepath.Clean(filepath.FromSlash(p))
	return p, nil
}

func EnsurePathSuffix(u *url.URL, suffix string) *url.URL {
	if strings.HasSuffix(u.Path, suffix) {
		return u
	} else {
		if strings.HasSuffix(u.Path, "/") {
			u.Path = u.Path + suffix
		} else {
			u.Path = u.Path + "/" + suffix
		}
	}
	return u
}

func NormalizeUrl(u *url.URL) (*url.URL, error) {
	res := *u
	if res.Scheme == "" {
		res.Scheme = "file"
	}
	if strings.ToUpper(res.Scheme) == "FILE" {
		p, e := UrlToPath(&res)
		if e != nil {
			return nil, e
		}
		return PathToUrl(p)
	} else {
		res.Path = path.Clean(res.Path)
	}
	return &res, nil
}

func ReadUrl(u *url.URL) ([]byte, error) {
	if hasPrefixIgnoringCase(u.Scheme, "http") {
		// Fetch the content
		var response *http.Response
		response, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			err = fmt.Errorf("error reading URL "+u.String()+", HTTP status %d", response.StatusCode)
			return nil, err
		}
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		return content, nil
	} else if strings.ToUpper(u.Scheme) == "FILE" {
		// Fetch the content
		location, err := UrlToPath(u)
		if err != nil {
			return nil, err
		}
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
	} else {
		return nil, errors.New("unsupported protocol <" + u.Scheme + ">")
	}
}

func hasPrefixIgnoringCase(s string, prefix string) bool {
	return strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(prefix))
}

func hasSuffixIgnoringCase(s string, suffix string) bool {
	return strings.HasSuffix(strings.ToUpper(s), strings.ToUpper(suffix))
}
