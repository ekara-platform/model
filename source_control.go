package model

import (
	"errors"
)

type (
	// SCMType represents a type of Source Control Management system
	SCMType string
)

const (
	//Git type of source control management system
	GitScm SCMType = SCMType(SchemeGits)
	//Svn type of source control management system
	SvnScm SCMType = SCMType(SchemeSvn)
	//Unknown source control management system
	UnknownScm SCMType = ""
)

func resolveSCMType(url EkUrl) (SCMType, error) {
	switch url.UpperScheme() {
	case SchemeFile:
		// TODO: for now assume git on local directories, later try to detect
		return GitScm, nil
	case SchemeGits:
		return GitScm, nil
	case SchemeSvn:
		return SvnScm, nil
	case SchemeHttp, SchemeHttps:
		if hasSuffixIgnoringCase(url.Path(), GitExtension) {
			return GitScm, nil
		}
	}
	return UnknownScm, errors.New("unknown fetch protocol: " + url.Scheme())
}
