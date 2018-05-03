package version

import (
	"testing"

	"github.com/Masterminds/semver"
)

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
