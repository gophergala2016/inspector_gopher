package inspector

import (
	"testing"
	"github.com/libgit2/git2go"
)

const REPO_NAME = "lazartravica/Envy"

func TestGetRepo(t *testing.T) {
	repo, err := GetRepo(REPO_NAME)
	defer repo.Free()
	if err != nil {
		t.Error(err)
	}
}

func TestRepoWalking(t *testing.T) {
	repo, err := GetRepo(REPO_NAME)
	defer repo.Free()
	if err != nil {
		t.Error(err)
	}

	WalkCommits(repo, func(_, _ *git.Commit) bool {
		return true
	})
}
