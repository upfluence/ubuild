package githubutil

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/upfluence/ubuild/pkg/version"
)

func GetLastVersion(repo string) (*version.Version, error) {
	splittedRepo := strings.Split(repo, "/")
	r, _, err := buildClient().Repositories.GetLatestRelease(
		context.Background(),
		splittedRepo[len(splittedRepo)-2],
		splittedRepo[len(splittedRepo)-1],
	)

	if err != nil {
		if e, ok := err.(*github.ErrorResponse); ok && e.Response.StatusCode == 404 {
			return &version.Version{Version: *semver.MustParse("v0.0.0")}, nil
		}

		return nil, err
	}

	ver, err := semver.NewVersion(*r.TagName)

	if err != nil {
		return nil, err
	}

	return &version.Version{Version: *ver}, err
}

func CreateRelease(repo string, sha string, v *version.Version) error {
	splittedRepo := strings.Split(repo, "/")
	vStr := v.String()
	t := "commit"
	ref := fmt.Sprintf("refs/tags/%s", v.String())

	if _, _, err := buildClient().Git.CreateTag(
		context.Background(),
		splittedRepo[len(splittedRepo)-2],
		splittedRepo[len(splittedRepo)-1],
		&github.Tag{
			Tag:     &vStr,
			Message: &vStr,
			Object: &github.GitObject{
				SHA:  &sha,
				Type: &t,
			},
		},
	); err != nil {
		return err
	}

	if _, _, err := buildClient().Git.CreateRef(
		context.Background(),
		splittedRepo[len(splittedRepo)-2],
		splittedRepo[len(splittedRepo)-1],
		&github.Reference{
			Ref: &ref,
			Object: &github.GitObject{
				SHA:  &sha,
				Type: &t,
			},
		},
	); err != nil {
		return err
	}

	pre := false

	_, _, err := buildClient().Repositories.CreateRelease(
		context.Background(),
		splittedRepo[len(splittedRepo)-2],
		splittedRepo[len(splittedRepo)-1],
		&github.RepositoryRelease{
			TagName:    &vStr,
			Name:       &vStr,
			Prerelease: &pre,
		},
	)

	return err
}

func FetchCommits(repo, from, to string) ([]string, error) {
	splittedRepo := strings.Split(repo, "/")
	res := []string{}

	r, _, err := buildClient().Repositories.CompareCommits(
		context.Background(),
		splittedRepo[len(splittedRepo)-2],
		splittedRepo[len(splittedRepo)-1],
		from,
		to,
	)

	if err != nil {
		if e, ok := err.(*github.ErrorResponse); ok && e.Response.StatusCode == 404 {
			return res, nil
		}

		return nil, err
	}

	for _, commit := range r.Commits {
		res = append(res, *commit.Commit.Message)
	}

	return res, nil
}

func buildClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
