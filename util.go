package model

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	return NormalizeUrl(u), nil
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
	res := NormalizeUrl(u)
	if strings.HasSuffix(res.Path, suffix) {
		return res
	} else {
		if strings.HasSuffix(res.Path, "/") {
			res.Path = res.Path + suffix
		} else {
			res.Path = res.Path + "/" + suffix
		}
	}
	return res
}

func NormalizeUrl(u *url.URL) *url.URL {
	res := *u
	if res.Scheme == "" {
		res.Scheme = "file"
	}
	if strings.ToUpper(res.Scheme) == "FILE" {
		res.Path = filepath.ToSlash(filepath.Clean(res.Path))
	}
	res.Path = path.Clean(res.Path)
	return &res
}

func ReadUrl(logger *log.Logger, u *url.URL) (*url.URL, []byte, error) {
	if hasPrefixIgnoringCase(u.Scheme, "http") {
		logger.Println("loading remote URL", u.String())

		// Fetch the content
		var response *http.Response
		response, err := http.Get(u.String())
		if err != nil {
			return nil, nil, err
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			err = fmt.Errorf("error reading URL "+u.String()+", HTTP status %d", response.StatusCode)
			return nil, nil, err
		}
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}

		// Compute the base
		i := strings.LastIndex(u.EscapedPath(), "/")
		if i == -1 {
			err = errors.New("cannot determine base URL for " + u.String())
			return nil, nil, err
		}
		base := *u
		base.Path = base.Path[0 : i+1]

		return &base, content, nil
	} else if strings.ToUpper(u.Scheme) == "FILE" {
		logger.Println("loading local URL", u.String())

		// Fetch the content
		location, err := UrlToPath(u)
		if err != nil {
			return nil, nil, err
		}
		file, err := os.Open(location)
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, nil, err
		}

		// Compute the base
		base, err := PathToUrl(filepath.Dir(location))
		if err != nil {
			return nil, nil, err
		}
		if !strings.HasSuffix(base.Path, "/") {
			base.Path = base.Path + "/"
		}

		return base, content, nil
	} else {
		return nil, nil, errors.New("unsupported protocol <" + u.Scheme + ">")
	}
}

func hasPrefixIgnoringCase(s string, prefix string) bool {
	return strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(prefix))
}

func hasSuffixIgnoringCase(s string, suffix string) bool {
	return strings.HasSuffix(strings.ToUpper(s), strings.ToUpper(suffix))
}
