package model

import (
	"errors"
	"net/url"
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
	Repository string
	Version    Version
}

func CreateDetachedComponent(repoUrl string, version string) (Component, error) {
	validationErrors := ValidationErrors{}
	c := createComponent(&validationErrors, nil, "<>", repoUrl, version)
	if validationErrors.HasErrors() {
		return Component{}, validationErrors
	}
	return c, nil
}

func createComponent(vErrs *ValidationErrors, env *Environment, location string, repoUrl string, version string) Component {
	//componentUrl := buildComponentUrl(vErrs, , repoUrl)
	cUrl, e := BuildComponentFolderUrl(repoUrl)
	if e != nil {
		vErrs.AddError(e, location+".repository")
	}
	cUrl, _ = BuildComponentGitUrl(cUrl)

	componentId := buildComponentId(cUrl)

	var parsedVersion Version
	if len(version) > 0 {
		parsedVersion = createVersion(vErrs, location+".version", version)
	} else {
		if managedVersion, ok := env.Components[componentId]; ok {
			parsedVersion = managedVersion
		} else {
			vErrs.AddError(errors.New("no version provided for component "+cUrl.String()), location+".version")
		}
	}

	return Component{
		Id:         buildComponentId(cUrl),
		Repository: cUrl.String(),
		Version:    parsedVersion,
		Scm:        resolveScm(vErrs, location+".repository", cUrl)}
}

func createComponentMap(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) map[string]Version {
	res := map[string]Version{}
	for id, v := range yamlEnv.Components {
		//res[buildComponentId(buildComponentUrl(vErrs, "components", id))] = createVersion(vErrs, "components."+id, v)
		cUrl, e := BuildComponentFolderUrl(id)
		if e != nil {
			vErrs.AddError(e, "components")
		}
		cUrl, _ = BuildComponentGitUrl(cUrl)
		res[buildComponentId(cUrl)] = createVersion(vErrs, "components."+id, v)
	}
	return res
}

// BuildComponentFolderUrl builds the complete URL of the folder containing all
// the component items.
//
// - URLs starting with github.com or bitbucket.org are assumed as https://
// - URLs without protocol and matching org/repo are assumed as https://github.com/...
func BuildComponentFolderUrl(repoUrl string) (url.URL, error) {
	// URL starting with github.com or bitbucket.org are assumed as https://
	if hasPrefixIgnoringCase(repoUrl, GitHubHost) || hasPrefixIgnoringCase(repoUrl, BitBucketHost) {
		repoUrl = "https://" + repoUrl
	}

	// URL without protocol and matching org/repo are assumed as https://github.com/...
	isSimpleRepo, _ := regexp.MatchString("^[_a-zA-Z0-9-]+/[_a-zA-Z0-9-]+$", repoUrl)
	if !hasPrefixIgnoringCase(repoUrl, "http") && isSimpleRepo {
		if hasPrefixIgnoringCase(repoUrl, "/") {
			repoUrl = "https://" + GitHubHost + repoUrl
		} else {
			repoUrl = "https://" + GitHubHost + "/" + repoUrl
		}
	}

	parsedUrl, e := url.Parse(repoUrl)
	if e != nil {
		return url.URL{}, e
	}
	return *parsedUrl, nil
}

// BuildComponentGitUrl builds the url of the git repository based on the
// url received has parameter
//
// If the received URL is not eligible to be converted into a GIT repository
// then an error will be returned and the unchanged  url will be returned;
func BuildComponentGitUrl(url url.URL) (url.URL, error) {
	if hasPrefixIgnoringCase(url.Scheme, "http") &&
		!hasSuffixIgnoringCase(url.Path, ".git") &&
		(url.Host == GitHubHost || url.Host == BitBucketHost) {
		url.Path = url.Path + ".git"
		return url, nil
	}
	return url, errors.New("the URL is not eligible to be a GIT repository")
}

func buildComponentId(componentUrl url.URL) string {
	id := componentUrl.Host
	if hasSuffixIgnoringCase(componentUrl.Path, ".git") {
		id += strings.Replace(componentUrl.Path[0:len(componentUrl.Path)-4], "/", "-", -1)
	} else {
		id += strings.Replace(componentUrl.Path, "/", "-", -1)
	}
	return id
}

func resolveScm(vErrs *ValidationErrors, location string, url url.URL) ScmType {
	switch strings.ToUpper(url.Scheme) {
	case "GIT":
		return Git
	case "SVN":
		return Svn
	case "HTTP", "HTTPS":
		if hasSuffixIgnoringCase(url.Path, ".git") {
			return Git
		}
	}
	vErrs.AddError(errors.New("unknown fetch protocol"), location)
	return Unknown
}

func hasPrefixIgnoringCase(s string, prefix string) bool {
	return strings.HasPrefix(strings.ToUpper(s), strings.ToUpper(prefix))
}

func hasSuffixIgnoringCase(s string, suffix string) bool {
	return strings.HasSuffix(strings.ToUpper(s), strings.ToUpper(suffix))
}
