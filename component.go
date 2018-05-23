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
		Id:         cId,
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
func ResolveRepositoryInfo(base *url.URL, repo string) (cId string, cUrl *url.URL, e error) {
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
			repo, e = filepath.Abs(repo)
			if e != nil {
				return
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
		}
		cUrl, e = url.Parse(repo)
		if e != nil {
			return
		}
	}

	// If it's HTTP(S), assume it's GIT and add the suffix
	if (strings.ToUpper(cUrl.Scheme) == "HTTP" || strings.ToUpper(cUrl.Scheme) == "HTTPS") && !hasSuffixIgnoringCase(cUrl.Path, ".git") {
		cUrl.Path = cUrl.Path + ".git"
	}

	// Compute the last segment in path (without extension) + hash of full url
	splitPath := strings.Split(cUrl.Path, "/")
	cId = splitPath[len(splitPath)-1]
	if strings.Contains(cId, ".") {
		cId = cId[:strings.LastIndex(cId, ".")]
	}
	hash := sha1.New()
	hash.Write([]byte(cUrl.String()))
	cId = cId + "-" + hex.EncodeToString(hash.Sum(nil))

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