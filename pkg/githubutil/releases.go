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

const maxReleases = 30

func GetLatestTags(ctx context.Context, repo string) ([]string, error) {
	var (
		tags         = make([]string, 0, maxReleases)
		client       = buildClient().Repositories
		splittedRepo = strings.Split(repo, "/")
		rls, _, err  = client.ListReleases(
			ctx,
			splittedRepo[len(splittedRepo)-2],
			splittedRepo[len(splittedRepo)-1],
			&github.ListOptions{Page: 1, PerPage: maxReleases},
		)
	)

	if err != nil {
		return tags, err
	}

	for _, r := range rls {
		tags = append(tags, *r.TagName)
	}

	return tags, nil
}

func GetLastVersion(repo string) (*version.Version, error) {
	var tags, err = GetLatestTags(context.Background(), repo)

	if err != nil {
		return nil, err
	}

	ver, err := tagToVersion("v0.0.0")

	if err != nil {
		return nil, err
	}

	for _, t := range tags {
		if toCompare, err := tagToVersion(t); err != nil {
			return nil, err
		} else if ver.Compare(toCompare) < 0 {
			ver = toCompare
		}
	}

	return ver, err
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

func tagToVersion(tag string) (*version.Version, error) {
	var sv, err = semver.NewVersion(tag)

	if err != nil {
		return nil, err
	}

	return &version.Version{Version: *sv}, nil
}
