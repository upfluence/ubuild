package frontenddeployer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/upfluence/ubuild/pkg/config"
	"github.com/upfluence/ubuild/pkg/context"
)

var errNoURL = errors.New("builder/frontenddeployer: No URL set")

type callbackData struct {
	Env        string `json:"env"`
	Repository string `json:"repository"`
	Commit     string `json:"commit_hash"`
}

func Build(ctx *context.Context, cfg *config.Configuration) error {
	if !ctx.Dist {
		return nil
	}

	var (
		d = cfg.GetDeployer()
		u = d.URL

		splittedRepo = strings.Split(cfg.GetRepo(), "/")
		repo         = splittedRepo[len(splittedRepo)-1]

		env = d.GetEnv(ctx.Version.Branch)
	)
	if u == "" {
		return errNoURL
	}

	buf, err := json.Marshal(
		&callbackData{Env: env, Repository: repo, Commit: ctx.Version.Commit[:7]},
	)

	if err != nil {
		return err
	}

	f := url.Values{}
	f.Add("env", env)
	f.Add("ref", "heads/"+ctx.Version.Branch)
	f.Add("repository", repo)
	f.Add("callback", u+"/activate")
	f.Add("callback_data", string(buf))

	req, err := http.NewRequest("POST", u+"/deploy", strings.NewReader(f.Encode()))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := d.Client().Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Deployment failed: [Status: %s]", res.Status)
	}

	return nil
}
