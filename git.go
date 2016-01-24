package inspector

import (
	"errors"
	"github.com/libgit2/git2go"
	"io/ioutil"
	"log"
	"os"
)

const MAX_DEPTH = 500

var tempDir string
var repoName string

const tempDirLocation string = "/tmp"
func getRepoDir(repoName string) string {
	return tempDirLocation + string(os.PathSeparator) + repoName
	if tempDir == "" {
		tempDir, _ = ioutil.TempDir(tempDirLocation, "repo")
	}

	return tempDir + string(os.PathSeparator) + repoName
}

func CleanTempDir() {
//	os.RemoveAll(getRepoDir(repoName))
}

func GetRepo(rName string) (*git.Repository, error) {
	repoName = rName
	if _, err := os.Stat(getRepoDir(repoName)); err == nil {
		log.Printf("[START] OPEN REPO %s", getRepoDir(repoName))
		repo, err := git.OpenRepository(getRepoDir(repoName))

		if err != nil {
			return repo, err
		}
		log.Printf("[SUCCESS] OPEN REPO %s", repoName)

		log.Printf("[START] PULLING REPO")



		log.Printf("[SUCCESS] PULLING REPO")
		return repo, err
	}

	log.Printf("[START] CLONE REPO %s", repoName)
	defer log.Printf("[SUCCESS] CLONE REPO %s", repoName)
	return git.Clone("git://github.com/" + repoName + ".git", getRepoDir(repoName), &git.CloneOptions{})
}

func GetNumberOfCommits(repo *git.Repository) (count int, err error) {

	walker, err := repo.Walk()
	if err != nil {
		return 0, err
	}
	defer walker.Free()

	err = walker.PushHead()
	if err != nil {
		return 0, err
	}

	err = walker.Iterate(func(commit *git.Commit) bool {
		count++
		return true
	})

	return count, err
}

// Access commits via callback
type CommitWalkerFunc func(previousCommit *git.Commit, currentCommit *git.Commit) bool

func WalkCommits(repo *git.Repository, walkerFunc CommitWalkerFunc) error {
	if repo == nil {
		return errors.New("[FAIL] No repo supplied")
	}

	walker, err := repo.Walk()
	if err != nil {
		return err
	}
	defer walker.Free()

	walker.Sorting(git.SortTopological | git.SortReverse)
	err = walker.PushHead()
	if err != nil {
		return err
	}
	log.Println("[START] Walk commits")
	defer log.Printf("[SUCCESS] Walk commits")

	var previousCommit *git.Commit
	commitNumber := 0
	numberOfCommits, _ := GetNumberOfCommits(repo)

	err = walker.Iterate(func(commit *git.Commit) bool {
		if previousCommit == nil {
			previousCommit = commit
			return true
		}

		defer func() {
			previousCommit.Free()
			previousCommit = commit
			commitNumber += 1
		}()

		if (commitNumber + 1 < numberOfCommits - MAX_DEPTH) {
			return true
		}

		return walkerFunc(previousCommit, commit)
	})
	if err != nil {
		return err
	}

	if previousCommit != nil {
		previousCommit.Free()
	}

	return nil
}

func WalkDepthCommits(repo *git.Repository, depth int, walkerFunc CommitWalkerFunc) error {
	if repo == nil {
		return errors.New("[FAIL] No repo supplied")
	}

	walker, err := repo.Walk()
	if err != nil {
		return err
	}
	defer walker.Free()

	walker.Sorting(git.SortTopological | git.SortReverse)
	err = walker.PushHead()
	if err != nil {
		return err
	}
	log.Println("[START] Walk commits")
	defer log.Printf("[SUCCESS] Walk commits")

	var previousCommit *git.Commit
	commitNumber := 0
	numberOfCommits, _ := GetNumberOfCommits(repo)

	err = walker.Iterate(func(commit *git.Commit) bool {
		if previousCommit == nil {
			previousCommit = commit
			return true
		}

		defer func() {
			previousCommit.Free()
			previousCommit = commit
			commitNumber += 1
		}()

		if (commitNumber + 1 < numberOfCommits - depth) {
			return true
		}

		return walkerFunc(previousCommit, commit)
	})
	if err != nil {
		return err
	}

	if previousCommit != nil {
		previousCommit.Free()
	}

	return nil
}


func GetDiff(repo *git.Repository, previousCommit *git.Commit, currentCommit *git.Commit) (*git.Diff, error) {
	previousTree, err := previousCommit.Tree()
	defer previousTree.Free()
	if err != nil {
		return nil, err
	}

	var currentTree *git.Tree
	if currentCommit != nil {
		currentTree, err = currentCommit.Tree()
		defer currentTree.Free()
		if err != nil {
			return nil, err
		}
	} else {
		currentTree = nil
	}

	options, err := git.DefaultDiffOptions()
	if err != nil {
		return nil, err
	}

	options.ContextLines = uint32(0)

	return repo.DiffTreeToTree(previousTree, currentTree, &options)
}

type FileWalkerFunc func(file git.DiffDelta, process float64)

func WalkFiles(diff *git.Diff, walker FileWalkerFunc) error {
	err := diff.ForEach(func(file git.DiffDelta, process float64) (git.DiffForEachHunkCallback, error) {
		walker(file, process)
		return nil, nil
	}, git.DiffDetailFiles)

	return err
}

type HunkWalkerFunc func(file git.DiffDelta, hunk git.DiffHunk)

func WalkHunks(diff *git.Diff, walker HunkWalkerFunc) error {
	err := diff.ForEach(func(file git.DiffDelta, process float64) (git.DiffForEachHunkCallback, error) {
		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			walker(file, hunk)
			return nil, nil
		}, nil
	}, git.DiffDetailHunks)

	return err
}
