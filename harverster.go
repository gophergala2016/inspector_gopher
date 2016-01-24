package inspector

import (
	"github.com/libgit2/git2go"
	"log"
	"time"
)

func Harvest(repoName string) string {
	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()

	count := 0
	WalkCommits(repo, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		diff, _ := GetDiff(repo, previousCommit, currentCommit)

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			count += 1
		})

		return true
	})

	log.Println(count)
	return "super"
}

func HarvestBenched(repoName string, depth int) float64 {
	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()

	start := time.Now()
	count := 0
	WalkDepthCommits(repo, depth, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		diff, _ := GetDiff(repo, previousCommit, currentCommit)

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			count += 1
		})

		return true
	})

	log.Println(count)

	return time.Since(start).Seconds()
}