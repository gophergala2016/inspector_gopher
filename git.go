package inspector

import (
	"github.com/libgit2/git2go"
	"log"
	"os"
	"time"
)

// Repo base path.
const clonePath string = "/tmp/inspector-gopher/"

type CommitWalkerFunc func(commit *git.Commit) bool

func WalkCommits(repoName string, walker CommitWalkerFunc) error {
	repo, err := openRepo(repoName)
	if err != nil {
		return err
	}

	return walkRepo(repo, walker)
}

// Clones repo to disc, if already exists, deletes the existing files first.
func openRepo(repoName string) (*git.Repository, error) {
	if _, err := os.Stat(clonePath + repoName); err == nil {
		log.Println("Opened repo [" + repoName + "]")
		return git.OpenRepository(clonePath + repoName)
	}

	repo, err := git.Clone("git://github.com/" + repoName + ".git", clonePath + repoName, &git.CloneOptions{})
	if err != nil {
		return nil, err
	}
	log.Println("Cloned repo [" + repoName + "]")

	return repo, nil
}

func walkRepo(repo *git.Repository, walkerFunc CommitWalkerFunc) error {
	walker, err := repo.Walk()
	if err != nil {
		return err
	}

	walker.Sorting(git.SortTopological | git.SortReverse)
	err = walker.PushHead()
	if err != nil {
		return err
	}
	log.Println("Started Processing repo")

	start := time.Now() // Log time passed for processing.

	var previousCommit *git.Commit

	err = walker.Iterate(func(commit *git.Commit) bool {
		if (previousCommit == nil) {
			previousCommit = commit
			return true
		}

		options, err := git.DefaultDiffOptions()
		if err != nil {
			return false
		}

		previousTree, _ := previousCommit.Tree()
		currentTree, _ := commit.Tree()

		diff, err := repo.DiffTreeToTree(previousTree, currentTree, &options)
		if err != nil {
			return false
		}

		err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
			return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {

				log.Printf("Hunk: %v", hunk.Header)
				return nil, nil
			}, nil
		}, git.DiffDetailHunks)

		if err != nil {
			panic(err)
		}

		previousCommit = commit
		return true
	})
	if err != nil {
		return err
	}

	log.Printf("Finished processing repo, duration %d seconds.", time.Now().Unix() - start.Unix())

	return nil
}
