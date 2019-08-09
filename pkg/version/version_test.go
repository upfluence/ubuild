package version

import (
	"testing"

	"github.com/Masterminds/semver"
)

func TestIncPatch(t *testing.T) {
	for _, tCase := range []struct {
		in  string
		out string
	}{
		{in: "v0.0.0", out: "v0.0.1"},
		{in: "v1.0.0-rc1", out: "v1.0.0"},
	} {
		v, err := semver.NewVersion(tCase.in)

		if err != nil {
			t.Errorf("Can't parse: %v: %v", tCase.in, err)
		}

		r := &Version{Version: *v}
		r.IncPatch()

		if res := r.String(); res != tCase.out {
			t.Errorf("Wrong version computed: %v instead of: %v", res, tCase.out)
		}
	}
}

func TestIncRC(t *testing.T) {
	for _, tCase := range []struct {
		in  string
		out string
	}{
		{in: "v0.0.0", out: "v0.0.1-rc1"},
		{in: "v1.0.0", out: "v1.0.1-rc1"},
		{in: "v1.0.0-rc1", out: "v1.0.0-rc2"},
	} {
		v, err := semver.NewVersion(tCase.in)

		if err != nil {
			t.Errorf("Can't parse: %v: %v", tCase.in, err)
		}

		r := &Version{Version: *v}
		r.IncRC()

		if res := r.String(); res != tCase.out {
			t.Errorf("Wrong version computed: %v instead of: %v", res, tCase.out)
		}
	}
}

func TestCompare(t *testing.T) {
	for _, tCase := range []struct {
		from string
		to   string
		out  int
	}{
		{from: "v0.0.0", to: "v0.0.0-rc1", out: 1},
		{from: "v0.0.0-rc1", to: "v0.0.0", out: -1},
		{from: "v1.0.0", to: "v1.0.1-rc1", out: -1},
		{from: "v1.0.1-rc2", to: "v1.0.1-rc1", out: 1},
		{from: "v1.0.1-rc1", to: "v1.0.1-rc2", out: -1},
		{from: "v1.0.1-rc9", to: "v1.0.1-rc10", out: -1},
		{from: "v1.0.1-rc1", to: "v1.0.1", out: -1},
		{from: "v1.0.0", to: "v1.0.0", out: 0},
	} {
		from, errf := semver.NewVersion(tCase.from)
		to, errt := semver.NewVersion(tCase.to)

		if errf != nil {
			t.Errorf("Can't parse: %v: %v", tCase.from, errf)
		}

		if errt != nil {
			t.Errorf("Can't parse: %v: %v", tCase.to, errt)
		}

		r := &Version{Version: *from}
		vt := &Version{Version: *to}

		if res := r.Compare(vt); res != tCase.out {
			t.Errorf(
				"Wrong compare [from: %q to %q]: %v instead of: %v",
				tCase.from,
				tCase.to,
				res,
				tCase.out,
			)
		}
	}
}
