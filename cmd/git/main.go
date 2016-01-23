package main

import (
	"github.com/libgit2/git2go"
	"log"
	"os"
	"time"
)

// Snapshots collect data gathered from processing commits.
type Snapshot struct {

	Commit Commit

}

// Local representation of a commit.
type Commit struct {
	Name string
	Time time.Time
}

func main() {

	repoName := "libgit2/git2go"

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
	if _, err := os.Stat(repoName); err == nil {
		err = os.RemoveAll(repoName)
		if err != nil {
			return nil, err
		}
		log.Println("Cleaned up repo [" + repoName + "]")
	}

	repo, err := git.Clone("git://github.com/" + repoName + ".git", repoName, &git.CloneOptions{})
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
	err = walker.Iterate(func(c *git.Commit) bool {
		commit := Commit{
			Name: c.Message(),
			Time: c.Author().When,
		}

		snapshot := Snapshot{
			Commit: commit,
		}

		snapshots = append(snapshots, snapshot)

		return true
	})
	if err != nil {
		return snapshots, err
	}

	log.Printf("Finished processing repo, duration %d seconds.", time.Now().Unix() - start.Unix())

	return snapshots, nil
}
