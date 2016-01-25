package inspector

import (
	"testing"
	"github.com/libgit2/git2go"
)

const REPO_ENVY = "lazartravica/Envy"
const REPO_HEXY = "lazartravica/Hexy"
const REPO_GIT2GO = "libgit2/git2go"
const REPO_MUX = "gorilla/mux"
const REPO_CONSUL = "hashicorp/consul"
const REPO_DOCKER = "docker/docker"

func TestGetRepo(t *testing.T) {
	repo, err := GetRepo(REPO_ENVY)
	if err != nil {
		t.Error(err)
	}
	defer repo.Free()

	defer CleanTempDir()
}

func TestNumberOfCommits(t *testing.T) {
	repo, err := GetRepo(REPO_HEXY)
	if err != nil {
		t.Error(err)
	}

	defer repo.Free()
	defer CleanTempDir()

	numberOfCommits, err :=GetNumberOfCommits(repo)
	if err != nil {
		t.Error(err)
	}

	if numberOfCommits != 41 {
		t.Errorf("Number of commits is not 41, got: %d", numberOfCommits)
	}
}

func TestRepoWalking(t *testing.T) {
	repo, err := GetRepo(REPO_ENVY)
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
	repo, err := GetRepo(REPO_ENVY)
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
	repo, err := GetRepo(REPO_DOCKER)
	if err != nil {
		t.Error(err)
	}

	defer repo.Free()
	defer CleanTempDir()

	WalkCommits(repo, func(previous, current *git.Commit) bool {
		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			t.Error(err)
			return false
		}
		defer diff.Free()

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
//			time.Sleep(time.Millisecond)
		})

		return true
	})
}

func TestListFiles(t *testing.T) {
	repo, err := GetRepo(REPO_DOCKER)
	if err != nil {
		t.Error(err)
	}

	defer repo.Free()
	defer CleanTempDir()

	ListFilesInWorkingDir(REPO_DOCKER)
}
