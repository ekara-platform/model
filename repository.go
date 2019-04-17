package model

import (
	"errors"
	"regexp"
)

const (

	//DefaultDescriptorName specifies the default name of the environment descriptor
	//
	//When the environment descriptor is not specified, for example into a use
	// component then we will look for a default descriptor name "ekara.yaml"
	DefaultDescriptorName = "ekara.yaml"

	unsupportedFileExtension = "unsupported file extension, only .yaml and .yml are supported"
)

type (

	//Repository represents a descriptor or component location
	Repository struct {
		// Scm specifies type of source sontrol management system holding the
		// component
		Scm SCMType
		// Url specifies the repository Url where to fetch the component
		Url EkUrl
		// The reference to the branch or tag to fetch. If not specified the default branch will be fetched
		Ref string
		//DescriptorName specifies the name of the descriptor
		DescriptorName string
		// The authentication parameters to use if repository is not publicly accessible
		Authentication Parameters
	}
)

//CreateRepository creates a repository
//	Parameters
//
//		base: the base URL where to look for the component
//		repo: the repository Url where to fetch the component
//		ref: the ref to fetch, if the ref is not specified then the default branch will be fetched
//		descriptor: the name of the descriptor, if not specified then it will be defaulted
func CreateRepository(base Base, repo string, ref string, descriptor string) (Repository, error) {

	if descriptor == "" {
		descriptor = DefaultDescriptorName
	}

	r := Repository{
		Ref:            ref,
		DescriptorName: descriptor,
	}

	if !hasSuffixIgnoringCase(descriptor, ".yaml") && !hasSuffixIgnoringCase(descriptor, ".yml") {
		return r, errors.New(unsupportedFileExtension)
	}

	repoUrl, e := resolveRepositoryInfo(base, repo)
	if e != nil {
		return r, e
	}
	r.Url = repoUrl
	scmType, e := resolveSCMType(repoUrl)
	if e != nil {
		return r, e
	}
	r.Scm = scmType

	return r, e
}

// resolveRepository resolves a full URL from repository short-forms.
//
// URLs without protocol and matching org/repo are assumed as being prefixed with base
func resolveRepositoryInfo(base Base, repo string) (cUrl EkUrl, e error) {
	if repo == "" {
		e = errors.New("no repository specified")
		return
	}

	isSimpleRepo, _ := regexp.MatchString("^[_a-zA-Z0-9-]+/[_a-zA-Z0-9-]+$", repo)
	if isSimpleRepo {
		cUrl, e = base.CreateBasedUrl(repo)
		if e != nil {
			return
		}
	} else {
		cUrl, e = CreateUrl(repo)
		if e != nil {
			return
		}
	}

	// If it's HTTP(S), assume it's GIT and add the suffix
	if (cUrl.UpperScheme() == SchemeHttp || cUrl.UpperScheme() == SchemeHttps || cUrl.UpperScheme() == SchemeGits) && !hasSuffixIgnoringCase(cUrl.Path(), GitExtension) {
		if hasSuffixIgnoringCase(cUrl.Path(), "/") {
			cUrl.RemovePathSuffix("/")
		}
		cUrl.AddPathSuffix(GitExtension)
	}
	return
}
