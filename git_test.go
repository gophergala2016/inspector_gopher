package inspector

import (
	"github.com/libgit2/git2go"
	"testing"
)

const REPO_NAME = "lazartravica/Envy"

func TestGetRepo(t *testing.T) {
	repo, err := GetRepo(REPO_NAME)
	if err != nil {
		t.Error(err)
	}
	defer repo.Free()

	defer CleanTempDir()
}

func TestRepoWalking(t *testing.T) {
	repo, err := GetRepo(REPO_NAME)
	if err != nil {
		t.Error(err)
	}
	defer repo.Free()

	WalkCommits(repo, func(_, _ *git.Commit) bool {
		return true
	})

	defer CleanTempDir()
}

func TestGetDiff(t *testing.T) {
	repo, err := GetRepo(REPO_NAME)
	if err != nil {
		t.Error(err)
	}
	defer repo.Free()

	WalkCommits(repo, func(previous, current *git.Commit) bool {
		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			t.Error(err)
			return false
		}
		defer diff.Free()

		return true
	})

	defer CleanTempDir()
}

func TestWalkHunks(t *testing.T) {
	repo, err := GetRepo(REPO_NAME)
	if err != nil {
		t.Error(err)
	}
	defer repo.Free()

	WalkCommits(repo, func(previous, current *git.Commit) bool {
		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			t.Error(err)
			return false
		}
		defer diff.Free()

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
		})

		return true
	})

	defer CleanTempDir()
}
