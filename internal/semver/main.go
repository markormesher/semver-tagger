package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

type SemVer struct {
	Prefix         string
	Major          int
	Minor          int
	Patch          int
	Rc             int
	CommitDistance int
}

var semVerRegex = regexp.MustCompile(`^([vV]?)(\d+)\.(\d+)\.(\d+)(?:\-rc(\d+))?(?:\-(\d+)\-g[[:xdigit:]]+)?$`)

func FromString(tag string) (SemVer, error) {
	groups := semVerRegex.FindStringSubmatch(tag)
	if len(groups) == 0 {
		return SemVer{}, fmt.Errorf("could not parse '%s' into a valid semver", tag)
	}

	prefix := groups[1]

	major, err := strconv.Atoi(groups[2])
	if err != nil {
		return SemVer{}, fmt.Errorf("could not parse '%s' into a valid semver", tag)
	}

	minor, err := strconv.Atoi(groups[3])
	if err != nil {
		return SemVer{}, fmt.Errorf("could not parse '%s' into a valid semver", tag)
	}

	patch, err := strconv.Atoi(groups[4])
	if err != nil {
		return SemVer{}, fmt.Errorf("could not parse '%s' into a valid semver", tag)
	}

	rc := 0
	if groups[5] != "" {
		rc, err = strconv.Atoi(groups[5])
		if err != nil {
			return SemVer{}, fmt.Errorf("could not parse '%s' into a valid semver", tag)
		}
	}

	commitDistance := 0
	if groups[6] != "" {
		commitDistance, err = strconv.Atoi(groups[6])
		if err != nil {
			return SemVer{}, fmt.Errorf("could not parse '%s' into a valid semver", tag)
		}
	}

	return SemVer{
		Prefix:         prefix,
		Major:          major,
		Minor:          minor,
		Patch:          patch,
		Rc:             rc,
		CommitDistance: commitDistance,
	}, nil
}

func (sv *SemVer) String() string {
	out := fmt.Sprintf("%s%d.%d.%d", sv.Prefix, sv.Major, sv.Minor, sv.Patch)
	if sv.Rc > 0 {
		out += fmt.Sprintf("-rc%d", sv.Rc)
	}
	return out
}

func (sv SemVer) Bump(major, minor, patch, rc, noRc bool) SemVer {
	switch {
	case major:
		sv.Major++
		sv.Minor = 0
		sv.Patch = 0
		sv.Rc = 0

	case minor:
		sv.Minor++
		sv.Patch = 0
		sv.Rc = 0

	case patch:
		sv.Patch++
		sv.Rc = 0
	}

	if rc {
		sv.Rc++
	}

	if noRc {
		sv.Rc = 0
	}

	return sv
}
