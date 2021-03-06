package inspector

import (
	"errors"
	"github.com/libgit2/git2go"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const MAX_DEPTH = 50

var tempDir string
var repoName string
var filesInWorkingTree []string

var tempDirLocation string = "/tmp"
func getRepoDir(repoName string) string {
	return tempDirLocation + string(os.PathSeparator) + repoName // Uncomment to prefix with rand.
	if tempDir == "" {
		tempDir, _ = ioutil.TempDir(tempDirLocation, "repo")
	}

	return tempDir + string(os.PathSeparator) + repoName
}

func SetTempDir(newLocation string) {
	tempDirLocation = newLocation
	log.Println("Location set to " + tempDirLocation)
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

		filesInWorkingTree, err = ListFilesInWorkingDir(repoName)

		return repo, err
	}

	log.Printf("[START] CLONE REPO %s", repoName)
	repo, err := git.Clone("git://github.com/" + repoName + ".git", getRepoDir(repoName), &git.CloneOptions{})
	log.Printf("[SUCCESS] CLONE REPO %s", repoName)

	filesInWorkingTree, err = ListFilesInWorkingDir(repoName)

	return repo, err
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

		// Check if the file is in the master HEAD master.
		found := false
		for _, fileInWorkingTree := range filesInWorkingTree {
			if file.NewFile.Path == fileInWorkingTree {
				found = true
				break;
			}
		}

		// If the file is not found, we don't want to analyze it.
		if !found {
			return nil, nil
		}

		return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
			walker(file, hunk)
			return nil, nil
		}, nil
	}, git.DiffDetailHunks)

	return err
}

func ListFilesInWorkingDir(repoName string) (files []string, err error) {

	walkFn := func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.Name() == ".git" {
			return filepath.SkipDir
		} else if !stat.IsDir() {
			files = append(files, strings.Replace(path, getRepoDir(repoName) + "/", "", 1))
		}

		if err != nil {
			return err
		}
		return nil
	}

	err = filepath.Walk(getRepoDir(repoName), walkFn)
	if err != nil {
		return files, err
	}

	return files, nil
}