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
	u = NormalizeUrl(u)

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
	} else if hasPrefixIgnoringCase(u.Scheme, "file") {
		logger.Println("loading local URL", u.String())

		// Fetch the content
		location, err := filepath.Abs(filepath.FromSlash(u.EscapedPath()))
		if err != nil {
			return nil, nil, err
		}
		file, err := os.Open(filepath.FromSlash(u.EscapedPath()))
		if err != nil {
			return nil, nil, err
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, nil, err
		}

		// Compute the base
		base, err := url.Parse("file://" + filepath.ToSlash(filepath.Dir(location)) + "/")
		if err != nil {
			return nil, nil, err
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
