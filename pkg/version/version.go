package version

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Masterminds/semver"
)

type Version struct {
	semver.Version
}

func (v *Version) RC() int64 {
	r := v.Prerelease()

	if r == "" {
		return 0
	}

	res, err := strconv.Atoi(strings.TrimPrefix(r, "rc"))

	if err != nil {
		return 0
	}

	return int64(res)
}

func (v *Version) String() string {
	return fmt.Sprintf("v%s", v.Version.String())
}

func (v *Version) Compare(v2 *Version) int {
	if v.Major() == v2.Major() && v.Minor() == v2.Minor() &&
		v.Patch() == v2.Patch() && v.RC() != v2.RC() {
		if v.RC()*v2.RC() == 0 {
			return int(v2.RC() - v.RC())
		}

		return int(v.RC() - v2.RC())
	}

	return v.Version.Compare(&v2.Version)
}

func (v *Version) IncMajor() {
	v.Version = v.Version.IncMajor()
}

func (v *Version) IncMinor() {
	v.Version = v.Version.IncMinor()
}

func (v *Version) IncPatch() {
	v.Version = v.Version.IncPatch()
}

func (v *Version) IncRC() {
	if v.RC() == 0 {
		v.IncPatch()
	}

	v.Version, _ = v.SetPrerelease(fmt.Sprintf("rc%d", v.RC()+1))
}

func IncrementVersionFromCommits(v *Version, messages []string) bool {
	var bumped bool

	for _, message := range messages {
		if strings.Contains(message, "bump-major") {
			v.IncMajor()
			bumped = true
		}

		if strings.Contains(message, "bump-minor") {
			v.IncMinor()
			bumped = true
		}
	}

	return bumped
}
