package main
import (
"github.com/gophergala2016/inspector_gopher"
"github.com/libgit2/git2go"
	"log"
	"os"
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
		count++
		return true
	})

	log.Printf("There are %d commits in this repo.", count)
}
