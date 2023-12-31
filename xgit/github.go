// Package xgit
/*
 * @Date: 2023-07-20 10:22:13
 * @LastEditTime: 2023-07-20 10:50:26
 * @Description:
 */
package xgit

import (
	"context"

	"github.com/google/go-github/v51/github"
)

// GetGithubReleaseList ...
func GetGithubReleaseList(owner, repo string) ([]*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	releases, _, err := client.Repositories.ListReleases(context.Background(), owner, repo, nil)
	if err != nil {
		return nil, err
	}
	return releases, nil
}
