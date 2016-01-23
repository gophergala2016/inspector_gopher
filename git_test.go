package inspector

import (
	"testing"
	"github.com/libgit2/git2go"
)

const REPO_NAME = "lazartravica/Envy"


func TestGitReadRepo(t *testing.T) {
	_, err := openRepo(REPO_NAME)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRepoWalking(t *testing.T) {
	WalkCommits(REPO_NAME, func(_ *git.Commit) bool {
		return true
	})
}
