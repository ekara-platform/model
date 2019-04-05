package model

import (
	"errors"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type (
	// ScmType represents a type of Source Control Management system
	ScmType string

	//Component represents an element composing an ekara environment
	//
	//A component is always hosted into a source control management system.
	//
	//It can be for example a Provider or Software to deploy on the environment
	//
	Component struct {
		// Id specifies id of the component
		Id string
		// Scm specifies type of source sontrol management system holding the
		// component
		Scm ScmType
		// Repository specifies the repository Url where to fetch the component
		Repository *url.URL
		// The reference to the branch or tag to fetch. If not specified the default branch will be fetched
		Ref string
		// The authentication parameters to use if repository is not publicly accessible
		Authentication Parameters
		// Imports contains all the imports being declared within the component
		Imports []string
	}
)

const (
	//Git type of source control management system
	Git ScmType = "GIT"
	//Svn type of source control management system
	Svn ScmType = "SVN"
	//Unknown source control management system
	Unknown ScmType = ""

	//SchemeFile  scheme for a file
	SchemeFile string = "FILE"
	//SchemeGit  scheme for Git
	SchemeGit string = "GIT"
	//SchemeSvn  scheme for svn
	SchemeSvn string = "SVN"
	//SchemeHttp  scheme for http
	SchemeHttp string = "HTTP"
	//SchemeHttps  scheme for https
	SchemeHttps string = "HTTPS"
)

//CreateComponent creates a new component
//	Parameters
//
//		base: the base URL where to look for the component
//		id: the id of the component
//		repo: the repository Url where to fetch the component
//		ref: the ref to fetch, if the ref is not specified then the default branch will be fetched
//		imports: the imports located within the component
func CreateComponent(base *url.URL, id string, repo string, ref string, imports ...string) (Component, error) {
	repoUrl, e := resolveRepositoryInfo(base, repo)
	if e != nil {
		return Component{}, e
	}
	scmType, e := resolveScm(repoUrl)
	if e != nil {
		return Component{}, e
	}
	if len(imports) == 0 {
		imports = append(imports, DefaultDescriptorName)
	}
	return Component{Id: id, Repository: repoUrl, Ref: ref, Scm: scmType, Imports: imports}, nil
}

// resolveRepository resolves a full URL from repository short-forms.
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
	if (strings.ToUpper(cUrl.Scheme) == SchemeHttp || strings.ToUpper(cUrl.Scheme) == SchemeHttps) && !hasSuffixIgnoringCase(cUrl.Path, GitExtension) {
		cUrl.Path = cUrl.Path + GitExtension
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
		if hasSuffixIgnoringCase(url.Path, GitExtension) {
			return Git, nil
		}
	}
	return Unknown, errors.New("unknown fetch protocol: " + url.Scheme)
}
