package descriptor

import (
	"strings"
	"strconv"
)

type Version struct {
	Major int
	Minor int
	Micro int
	Full  string
}

func createVersion(full string) (Version, error) {
	result := Version{Full: full}
	split := strings.Split(full, ".")
	if len(split) > 0 {
		major, err := strconv.Atoi(split[0])
		if err != nil {
			return result, err
		}
		result.Major = int(major)
	}
	if len(split) > 1 {
		minor, err := strconv.Atoi(split[1])
		if err != nil {
			return result, err
		}
		result.Minor = int(minor)
	}
	if len(split) > 2 {
		minor, err := strconv.Atoi(split[2])
		if err != nil {
			return result, err
		}
		result.Micro = int(minor)
	}
	return result, nil
}
