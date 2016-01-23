package inspector

import (
	"github.com/libgit2/git2go"
	"log"
	"os"
	"time"
	"errors"
)

// Repo base path.
const clonePath string = "/tmp/inspector-gopher/"

func GetRepo(repoName string) (*git.Repository, error) {
	if _, err := os.Stat(clonePath + repoName); err == nil {
		log.Println("Opened repo [" + repoName + "]")
		return git.OpenRepository(clonePath + repoName)
	}

	defer log.Println("Cloned repo [" + repoName + "]")
	return git.Clone("git://github.com/" + repoName + ".git", clonePath + repoName, &git.CloneOptions{})
}

// Access commits via callback
type CommitWalkerFunc func(previousCommit *git.Commit, currentCommit *git.Commit) bool

func WalkCommits(repo *git.Repository, walkerFunc CommitWalkerFunc) error {
	walker, err := repo.Walk()
	defer walker.Free()
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
		if previousCommit == nil {
			previousCommit = commit
			return true
		}

		walkForward := walkerFunc(previousCommit, commit)

		previousCommit.Free()

		previousCommit = commit

		return walkForward
	})
	if err != nil {
		return err
	}

	if previousCommit != nil {
		previousCommit.Free()
	}

	log.Printf("Finished processing repo, duration %d seconds.", time.Now().Unix() - start.Unix())

	return nil
}

func GetDiff(repo *git.Repository, previousCommit *git.Commit, currentCommit *git.Commit) (*git.Diff, error) {
	if (previousCommit == nil || currentCommit == nil) {
		return nil, errors.New("You must pass both commits to get the diff.")
	}

	previousTree, err := previousCommit.Tree()
	defer previousTree.Free()
	if err != nil {
		return nil, err
	}

	currentTree, err := currentCommit.Tree()
	defer currentTree.Free()
	if err != nil {
		return nil, err
	}

	options, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}

	return repo.DiffTreeToTree(previousTree, currentTree, &options)
}

//
//	err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
//		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
//
//			log.Printf("Hunk: %v", hunk.Header)
//			return nil, nil
//		}, nil
//	}, git.DiffDetailHunks)
//
//	if err != nil {
//		panic(err)
//	}
//
//	previousCommit = commit
//	return true
//}