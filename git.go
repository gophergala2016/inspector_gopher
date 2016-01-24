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

var tempDir string

const tempDirLocation string = "/tmp"
func getRepoDir(repoName string) string {
	if tempDir == "" {
		tempDir, _ = ioutil.TempDir(tempDirLocation, "repo")
	}

	return tempDir + string(os.PathSeparator) + repoName
}

func CleanTempDir() {
	os.RemoveAll(getRepoDir())
}

func GetRepo(repoName string) (*git.Repository, error) {
	if _, err := os.Stat(); err == nil {
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

	return nil
}

func GetDiff(repo *git.Repository, previousCommit *git.Commit, currentCommit *git.Commit) (*git.Diff, error) {
	if previousCommit == nil || currentCommit == nil {
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

	options.ContextLines = uint32(0)

	return repo.DiffTreeToTree(previousTree, currentTree, &options)
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

func ListFiles(repoPath string) []string {
	files := []string{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		stat, err := os.Stat(path)
		if err != nil {
			return err
		}

		if stat.Name() == ".git" {
			return filepath.SkipDir
		} else if !stat.IsDir() {
			files = append(files, strings.Replace(path, repoPath, "", 1))
		}

		if err != nil {
			return err
		}
		return nil
	}
	err := filepath.Walk(repoPath, walkFn)
	if err != nil {
		log.Fatal(err)
	}

	for _, val := range files {
		log.Println(val)
	}

	return files
}

func Pull(repo *git.Repository) error {
	// Get the name
	name := "master"

	// Locate remote
	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	// Fetch changes from remote
	if err := remote.Fetch([]string{}, nil, ""); err != nil {
		return err
	}

	// Get remote master
	remoteBranch, err := repo.References.Lookup("refs/remotes/origin/"+name)
	if err != nil {
		return err
	}

	remoteBranchID := remoteBranch.Target()
	// Get annotated commit
	annotatedCommit, err := repo.AnnotatedCommitFromRef(remoteBranch)
	if err != nil {
		return err
	}

	// Do the merge analysis
	mergeHeads := make([]*git.AnnotatedCommit, 1)
	mergeHeads[0] = annotatedCommit
	analysis, _, err := repo.MergeAnalysis(mergeHeads)
	if err != nil {
		return err
	}

	// Get repo head
	head, err := repo.Head()
	if err != nil {
		return err
	}

	if analysis & git.MergeAnalysisUpToDate != 0 {
		return nil
	}  else if analysis & git.MergeAnalysisNormal != 0 {
		// Just merge changes
		if err := repo.Merge([]*git.AnnotatedCommit{annotatedCommit}, nil, nil); err != nil {
			return err
		}
		// Check for conflicts
		index, err := repo.Index()
		if err != nil {
			return err
		}

		if index.HasConflicts() {
			return errors.New("Conflicts encountered. Please resolve them.")
		}

		// Make the merge commit
		sig, err := repo.DefaultSignature()
		if err != nil {
			return err
		}

		// Get Write Tree
		treeId, err := index.WriteTree()
		if err != nil {
			return err
		}

		tree, err := repo.LookupTree(treeId)
		if err != nil {
			return err
		}

		localCommit, err := repo.LookupCommit(head.Target())
		if err != nil {
			return err
		}

		remoteCommit, err := repo.LookupCommit(remoteBranchID)
		if err != nil {
			return err
		}

		repo.CreateCommit("HEAD", sig, sig, "", tree, localCommit, remoteCommit)

		// Clean up
		repo.StateCleanup()
	} else if analysis & git.MergeAnalysisFastForward != 0 {
		// Fast-forward changes
		// Get remote tree
		remoteTree, err := repo.LookupTree(remoteBranchID)
		if err != nil {
			return err
		}

		// Checkout
		if err := repo.CheckoutTree(remoteTree, nil); err != nil {
			return err
		}

		branchRef, err := repo.References.Lookup("refs/heads/"+name)
		if err != nil {
			return err
		}

		// Point branch to the object
		branchRef.SetTarget(remoteBranchID, "")
		if _, err := head.SetTarget(remoteBranchID, ""); err != nil {
			return err
		}

	} else {
		return errors.New("Unexpected merge analysis result")
	}

	return nil
}