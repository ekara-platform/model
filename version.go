package descriptor

import (
	"strings"
	"strconv"
	"errors"
)

type Version struct {
	Major int
	Minor int
	Micro int
	Full  string
}

func createVersion(vErrs *ValidationErrors, location string, full string) Version {
	v := Version{Major: 0, Minor: 0, Micro: 0, Full: full}
	split := strings.Split(full, ".")

	if len(split) > 0 {
		major, err := strconv.Atoi(split[0])
		if err != nil {
			vErrs.AddError(err, location+".x")
		} else {
			v.Major = int(major)
		}
	} else {
		vErrs.AddError(errors.New("empty version was specified"), location)
	}

	if len(split) > 1 {
		minor, err := strconv.Atoi(split[1])
		if err != nil {
			vErrs.AddError(err, location+".y")
		} else {
			v.Minor = int(minor)
		}
	}
	if len(split) > 2 {
		minor, err := strconv.Atoi(split[2])
		if err != nil {
			vErrs.AddError(err, location+".z")
		} else {
			v.Micro = int(minor)
		}
	}
	if len(split) > 3 {
		vErrs.AddWarning("ignored extraneous characters after x.y.z in version "+full, location)
	}
	return v
}
