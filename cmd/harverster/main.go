package main
import (
	"log"
	"github.com/gophergala2016/inspector_gopher"
	"os"
"github.com/libgit2/git2go"
)

func main() {
	repoName := "docker/docker"
	if len(os.Args) > 1 {
		repoName = os.Args[1]
	}

	repo, err := inspector.GetRepo(repoName)
	if err != nil {
		panic(err)
	}
	defer repo.Free()

	count := 0
	inspector.WalkCommits(repo, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		diff, _ := inspector.GetDiff(repo, previousCommit, currentCommit)

		inspector.WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			count++
		})

		return true
	})

	log.Printf("There are %d hunks in this repo.", count)
}
