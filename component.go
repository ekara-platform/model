package model

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net/url"
	"os"
	"path/filepath"
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

func createComponent(vErrs *ValidationErrors, env *Environment, location string, repo string, version string) Component {
	cId, cUrl, e := ResolveRepositoryInfo(env.Settings.ComponentBase, repo)
	if e != nil {
		vErrs.AddError(e, location+".repository")
	}

	var parsedVersion Version
	if len(version) > 0 {
		parsedVersion = createVersion(vErrs, location+".version", version)
	} else {
		if managedVersion, ok := env.Components[cId]; ok {
			parsedVersion = managedVersion
		} else {
			vErrs.AddError(errors.New("no version provided for component "+cUrl.String()), location+".version")
		}
	}

	return Component{
		Id:         BuildComponentId(cUrl),
		Repository: cUrl,
		Version:    parsedVersion,
		Scm:        resolveScm(vErrs, location+".repository", cUrl)}
}

func createComponentMap(vErrs *ValidationErrors, env *Environment, yamlEnv *yamlEnvironment) map[string]Version {
	res := map[string]Version{}
	for repo, v := range yamlEnv.Components {
		cId, _, e := ResolveRepositoryInfo(env.Settings.ComponentBase, repo)
		if e != nil {
			vErrs.AddError(e, "components")
		}
		res[cId] = createVersion(vErrs, "components."+repo, v)
	}
	return res
}

// ResolveRepository resolve a full URL from repository short-forms.
//
// - URLs starting with github.com or bitbucket.org are assumed as https://
// - URLs without protocol and matching org/repo are assumed as being prefixed with base
func ResolveRepositoryInfo(base *url.URL, repo string) (string, *url.URL, error) {
	isSimpleRepo := false

	if _, e := os.Stat(repo); e == nil {
		// If it is a local file
		repo, e = filepath.Abs(repo)
		if e != nil {
			return "", nil, e
		}
		repo = filepath.ToSlash(repo)
		if strings.HasPrefix(repo, "/") {
			repo = "file://" + repo
		} else {
			repo = "file:///" + repo
		}
	} else if hasPrefixIgnoringCase(repo, GitHubHost) || hasPrefixIgnoringCase(repo, BitBucketHost) {
		// If not check if it begins with a known source provider (github, bitbucket, ...)
		repo = "https://" + repo
	} else {
		// If it is a simple form (org/repo), resolve it according to the base URL
		isSimpleRepo, _ = regexp.MatchString("^[_a-zA-Z0-9-]+/[_a-zA-Z0-9-]+$", repo)
	}

	// Parse the resulting URL
	cUrl, e := url.Parse(repo)
	if e != nil {
		return "", nil, e
	}

	// If it was a simple repo, resolve the parsed URL relatively to the base
	if isSimpleRepo {
		cUrl = base.ResolveReference(cUrl)
	}

	// If it's HTTP(S), assume it's GIT and add the suffix
	if (strings.ToUpper(cUrl.Scheme) == "HTTP" || strings.ToUpper(cUrl.Scheme) == "HTTPS") && !hasSuffixIgnoringCase(cUrl.Path, ".git") {
		cUrl.Path = cUrl.Path + ".git"
	}

	// Compute the last segment in path (without extension) + hash of full url
	splitPath := strings.Split(cUrl.Path, "/")
	cId := splitPath[len(splitPath)-1]
	if strings.Contains(cId, ".") {
		cId = cId[:strings.LastIndex(cId, ".")]
	}
	hash := sha1.New()
	hash.Write([]byte(cUrl.String()))
	cId = cId + "-" + hex.EncodeToString(hash.Sum(nil))

	return cId, cUrl, nil
}

func BuildComponentId(componentUrl *url.URL) string {
	id := componentUrl.Host
	if hasSuffixIgnoringCase(componentUrl.Path, ".git") {
		id += strings.Replace(componentUrl.Path[0:len(componentUrl.Path)-4], "/", "-", -1)
	} else {
		id += strings.Replace(componentUrl.Path, "/", "-", -1)
	}
	return id
}

func resolveScm(vErrs *ValidationErrors, location string, url *url.URL) ScmType {
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
