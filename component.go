package model

import (
	"strings"
	"net/url"
	"regexp"
	"errors"
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

func createComponentMap(vErrs *ValidationErrors, yamlEnv *yamlEnvironment) map[string]Version {
	res := map[string]Version{}
	for id, v := range yamlEnv.Components {
		res[buildComponentId(buildComponentUrl(vErrs, "components", id))] = createVersion(vErrs, "components."+id, v)
	}
	return res
}

func createComponent(vErrs *ValidationErrors, env *Environment, location string, repoUrl string, version string) Component {
	componentUrl := buildComponentUrl(vErrs, location+".repository", repoUrl)
	componentId := buildComponentId(componentUrl)
	var parsedVersion Version
	if len(version) > 0 {
		parsedVersion = createVersion(vErrs, location+".version", version)
	} else if managedVersion, ok := env.Components[componentId]; ok {
		parsedVersion = managedVersion
	} else {
		vErrs.AddError(errors.New("no version provided for component "+componentUrl.String()), location+".version")
	}
	return Component{
		Id:         buildComponentId(componentUrl),
		Repository: componentUrl.String(),
		Version:    parsedVersion,
		Scm:        resolveScm(vErrs, location+".repository", componentUrl)}
}

func buildComponentUrl(vErrs *ValidationErrors, location string, repoUrl string) url.URL {
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
		vErrs.AddError(e, location)
	}

	if hasPrefixIgnoringCase(parsedUrl.Scheme, "http") &&
		!hasSuffixIgnoringCase(parsedUrl.Path, ".git") &&
		(parsedUrl.Host == GitHubHost || parsedUrl.Host == BitBucketHost) {
		parsedUrl.Path = parsedUrl.Path + ".git"
	}

	return *parsedUrl
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
