package main

import (
	"github.com/libgit2/git2go"
	"log"
	"os"
	"time"
)

// Path to which to clone repos to when analyzing them.
const clonePath string = "/tmp/inspector-gopher/"

// Snapshots collect data gathered from processing commits.
type Snapshot struct {

	Commit Commit

}

// Local representation of a commit.
type Commit struct {
	Name string
	Time time.Time
	Commit git.Commit
}

func main() {

	repoName := "lazartravica/Envy"

	snapshots, err := parse(repoName)
	if err != nil {
		log.Fatal(err)
	}

	for index, snapshot := range snapshots {
		log.Printf("Commit number: %d, Time: %s", index, snapshot.Commit.Time)
	}
}

func parse(repoName string) ([]Snapshot, error) {
	repo, err := cloneRepo(repoName)
	if err != nil {
		return nil, err
	}

	return walkRepo(repo)
}

// Clones repo to disc, if already exists, deletes the existing files first.
func cloneRepo(repoName string) (*git.Repository, error) {
	if _, err := os.Stat(clonePath + repoName); err == nil {
		err = os.RemoveAll(clonePath + repoName)
		if err != nil {
			return nil, err
		}
		log.Println("Cleaned up repo [" + repoName + "]")
	}

	repo, err := git.Clone("git://github.com/" + repoName + ".git", clonePath + repoName, &git.CloneOptions{})
	if err != nil {
		return nil, err
	}
	log.Println("Cloned repo [" + repoName + "]")

	return repo, nil
}

func walkRepo(repo *git.Repository) ([]Snapshot, error) {
	snapshots := []Snapshot{}

	walker, err := repo.Walk()
	if err != nil {
		return snapshots, err
	}
	walker.Sorting(git.SortTopological | git.SortReverse)
	err = walker.PushHead()
	if err != nil {
		return snapshots, err
	}
	log.Println("Started Processing repo")
	start := time.Now()

	var prevCommit *git.Commit

	err = walker.Iterate(func(c *git.Commit) bool {

		commit := Commit{
			Name: c.Message(),
			Time: c.Author().When,
		}

		log.Printf("%d", len(snapshots))

		if prevCommit != nil {
			currTree, _ := c.Tree()
			prevTree, _ := prevCommit.Tree()

			diffOptions, _ := git.DefaultDiffOptions()
			diff, _ := repo.DiffTreeToTree(currTree, prevTree, &diffOptions)

			err := diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
				return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
					log.Printf("Hunk: %v", hunk.Header)
					return func(line git.DiffLine) error {
						log.Printf("%s %d", line.Content, line.Origin)
						return nil
					}, nil
				}, nil
			}, git.DiffDetailHunks)

			if err != nil {
				panic(err)
			}
		}

		snapshot := Snapshot{
			Commit: commit,
		}

		snapshots = append(snapshots, snapshot)

		prevCommit = c
		return true
	})
	if err != nil {
		return snapshots, err
	}

	log.Printf("Finished processing repo, duration %d seconds.", time.Now().Unix() - start.Unix())

	return snapshots, nil
}
