package inspector

import (
	"github.com/libgit2/git2go"
)

func Harvest(repoName string) string {
	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()

	count := 0
	WalkCommits(repo, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		diff, _ := GetDiff(repo, previousCommit, currentCommit)

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			count++
		})

		return true
	})

	return "super"
}