package model

import (
	"errors"
	"net/url"
	"os"
	"regexp"
	"strings"
)

// ScmType is the type used to identify the Source Control Management system
type (
	ScmType string

	Component struct {
		Id         string
		Scm        ScmType
		Repository *url.URL
		Version    Version
		Imports    []string
	}

	ComponentRef struct {
		ref       string
		mandatory bool

		env      *Environment
		location DescriptorLocation
	}
)

const (
	Git     ScmType = "GIT"
	Svn     ScmType = "SVN"
	Unknown ScmType = ""

	SchemeFile  string = "FILE"
	SchemeGit   string = "GIT"
	SchemeSvn   string = "SVN"
	SchemeHttp  string = "HTTP"
	SchemeHttps string = "HTTPS"
)

func CreateComponent(base *url.URL, id string, repo string, version string, imports ...string) (Component, error) {
	repoUrl, e := resolveRepositoryInfo(base, repo)
	if e != nil {
		return Component{}, e
	}
	scmType, e := resolveScm(repoUrl)
	if e != nil {
		return Component{}, e
	}
	parsedVersion, e := createVersion(version)
	if e != nil {
		return Component{}, e
	}
	if len(imports) == 0 {
		imports = append(imports, DefaultDescriptorName)
	}
	return Component{Id: id, Repository: repoUrl, Version: parsedVersion, Scm: scmType, Imports: imports}, nil
}

func createComponentRef(env *Environment, location DescriptorLocation, componentRef string, mandatory bool) ComponentRef {
	return ComponentRef{env: env, location: location, ref: componentRef, mandatory: mandatory}
}

func (r ComponentRef) validate() ValidationErrors {
	validationErrors := ValidationErrors{}
	if r.ref == "" {
		if r.mandatory {
			validationErrors.addError(errors.New("empty component reference"), r.location)
		}
	} else {
		if _, ok := r.env.Ekara.Components[r.ref]; !ok {
			validationErrors.addError(errors.New("reference to unknown component: "+r.ref), r.location)
		}
	}
	return validationErrors
}

func (r *ComponentRef) merge(other ComponentRef) {
	if r.ref == "" {
		r.ref = other.ref
	}
}

func (r ComponentRef) Resolve() Component {
	validationErrors := r.validate()
	if validationErrors.HasErrors() {
		panic(validationErrors)
	}
	return r.env.Ekara.Components[r.ref]
}

// ResolveRepository resolve a full URL from repository short-forms.
//
// URLs without protocol and matching org/repo are assumed as being prefixed with base
func resolveRepositoryInfo(base *url.URL, repo string) (cUrl *url.URL, e error) {
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
	if (strings.ToUpper(cUrl.Scheme) == SchemeHttp || strings.ToUpper(cUrl.Scheme) == SchemeHttps) && !hasSuffixIgnoringCase(cUrl.Path, ".git") {
		cUrl.Path = cUrl.Path + ".git"
	}

	return
}

func resolveScm(url *url.URL) (ScmType, error) {
	switch strings.ToUpper(url.Scheme) {
	case SchemeFile:
		// TODO: for now assume git on local directories, later try to detect
		return Git, nil
	case SchemeGit:
		return Git, nil
	case SchemeSvn:
		return Svn, nil
	case SchemeHttp, SchemeHttps:
		if hasSuffixIgnoringCase(url.Path, ".git") {
			return Git, nil
		}
	}
	return Unknown, errors.New("unknown fetch protocol: " + url.Scheme)
}
