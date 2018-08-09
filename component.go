package model

import (
	"errors"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ScmType string

const (
	Git     ScmType = "GIT"
	Svn     ScmType = "SVN"
	Unknown ScmType = ""
)

type Component struct {
	Id         string
	Scm        ScmType
	Repository *url.URL
	Version    Version
}

type ComponentRef struct {
	component *Component
}

func (c ComponentRef) Resolve() Component {
	// just copy the component by value
	return *c.component
}

func createComponent(vErrs *ValidationErrors, location string, componentBase *url.URL, id string, repo string, version string) Component {
	repoUrl, e := ResolveRepositoryInfo(componentBase, repo)
	if e != nil {
		vErrs.AddError(e, location+".repository")
		return Component{}
	} else {
		return Component{
			Id:         id,
			Repository: repoUrl,
			Version:    createVersion(vErrs, location+".version", version),
			Scm:        resolveScm(vErrs, location+".repository", repoUrl)}
	}
}

func createComponentRef(vErrs *ValidationErrors, components map[string]Component, location string, componentRef string) ComponentRef {
	if len(componentRef) == 0 {
		vErrs.AddError(errors.New("empty component reference"), location)
	} else {
		if val, ok := components[componentRef]; ok {
			return ComponentRef{component: &val}
		} else {
			vErrs.AddError(errors.New("unknown component reference: "+componentRef), location)
		}
	}
	return ComponentRef{}
}

// ResolveRepository resolve a full URL from repository short-forms.
//
// - URLs starting with github.com or bitbucket.org are assumed as https://
// - URLs without protocol and matching org/repo are assumed as being prefixed with base
func ResolveRepositoryInfo(base *url.URL, repo string) (cUrl *url.URL, e error) {
	if repo == "" {
		e = errors.New("no repository specified")
		return
	}

	isSimpleRepo, _ := regexp.MatchString("^[_a-zA-Z0-9-]+/[_a-zA-Z0-9-]+$", repo)
	if isSimpleRepo {
		// Simple repositories are always resolved relatively to the base URL
		cUrl, e = url.Parse(repo)
		if e != nil {
			return
		}
		cUrl = base.ResolveReference(cUrl)
	} else {
		if _, e = os.Stat(repo); e == nil {
			// If it is a local file
			cUrl, e = PathToUrl(repo)
			if e != nil {
				return
			}
		} else {
			if hasPrefixIgnoringCase(repo, GitHubHost) || hasPrefixIgnoringCase(repo, BitBucketHost) {
				repo = "https://" + repo
			}
			cUrl, e = url.Parse(repo)
			if e != nil {
				return
			}
		}
	}

	// Normalize the URL
	cUrl, e = NormalizeUrl(cUrl)
	if e != nil {
		return
	}

	// If it's HTTP(S), assume it's GIT and add the suffix
	if (strings.ToUpper(cUrl.Scheme) == "HTTP" || strings.ToUpper(cUrl.Scheme) == "HTTPS") && !hasSuffixIgnoringCase(cUrl.Path, ".git") {
		cUrl.Path = cUrl.Path + ".git"
	}

	return
}

func resolveScm(vErrs *ValidationErrors, location string, url *url.URL) ScmType {
	switch strings.ToUpper(url.Scheme) {
	case "FILE":
		// TODO: for now assume git on local directories, later try to detect
		return Git
	case "GIT":
		return Git
	case "SVN":
		return Svn
	case "HTTP", "HTTPS":
		if hasSuffixIgnoringCase(url.Path, ".git") {
			return Git
		}
	}
	vErrs.AddError(errors.New("unknown fetch protocol: "+url.Scheme), location)
	return Unknown
}
