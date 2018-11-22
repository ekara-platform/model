package model

import (
	"fmt"

	"encoding/json"
	"regexp"
	"strconv"
)

var semanticVersioningPattern = regexp.MustCompile("^(?P<major>0|[1-9]\\d*)(\\.(?P<minor>0|[1-9]\\d*))?(\\.(?P<patch>0|[1-9]\\d*))?(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$")

//Version represents the version of a component
type Version struct {
	Major int
	Minor int
	Micro int
	full  string
}

// MarshalJSON returns the serialized content of the version as JSON
func (r Version) MarshalJSON() ([]byte, error) {
	t := struct {
		Major string `json:",omitempty"`
		Minor string `json:",omitempty"`
		Micro string `json:",omitempty"`
		Full  string `json:",omitempty"`
	}{}

	if r.Major > -1 {
		t.Major = strconv.Itoa(r.Major)
		t.Minor = strconv.Itoa(r.Minor)
		t.Micro = strconv.Itoa(r.Micro)
	} else {
		t.Full = r.full
	}

	return json.Marshal(t)
}

func createVersion(full string) (Version, error) {
	v := Version{Major: -1, Minor: -1, Micro: -1, full: full}

	if semanticVersioningPattern.MatchString(full) {
		match := semanticVersioningPattern.FindStringSubmatch(full)
		result := make(map[string]string)
		for i, name := range semanticVersioningPattern.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		if result["major"] != "" {
			major, err := strconv.Atoi(result["major"])
			if err != nil {
				return Version{}, err
			} else {
				v.Major = int(major)
			}
		}
		if result["minor"] != "" {
			minor, err := strconv.Atoi(result["minor"])
			if err != nil {
				return Version{}, err
			} else {
				v.Minor = int(minor)
			}
		}
		if result["patch"] != "" {
			patch, err := strconv.Atoi(result["patch"])
			if err != nil {
				return Version{}, err
			} else {
				v.Micro = int(patch)
			}
		}
		// TODO take prerelease tag into account
	}
	return v, nil
}

func (v Version) IncludesVersion(other Version) bool {
	if v.Major >= 0 {
		if v.Major != other.Major {
			return false
		}
		if v.Minor >= 0 && v.Minor != other.Minor {
			return false
		}
		if v.Micro >= 0 && v.Micro != other.Micro {
			return false
		}
		return true
	} else {
		return v.full == other.full
	}
}

// String returns the string representation of the version
func (v Version) String() string {
	if v.Major >= 0 {
		return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Micro)
	} else {
		return v.full
	}
}
