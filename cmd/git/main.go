package main

import (
	"github.com/libgit2/git2go"
	"log"
	"os"
	"time"
	"strconv"
)

// Path to which to clone repos to when analyzing them.
const clonePath string = "/tmp/inspector-gopher/"

// Snapshots collect data gathered from processing commits.
type Snapshot struct {
	Commit Commit
}

// Local representation of a commit.
type Commit struct {
	Hash      string
	Message   string
	Time      time.Time
	Developer Developer
}

type Developer struct {
	Name  string
	Email string
}

type File struct {
	Path       string
	Insertions []CodeBlock
	Deletions  []CodeBlock
}

type CodeBlock struct {
	Start, End int
}

func newFile(path string) *File {
	return &File{
		Path: path,
		Insertions: []CodeBlock{},
		Deletions: []CodeBlock{},
	}
}

func main() {

	repoName := "BencicAndrej/crAPI"
	//	repoName := "lazartravica/Envy"
	//	repoName := "libgit2/git2go"

	_, err := parse(repoName)
	if err != nil {
		log.Fatal(err)
	}

//	for index, snapshot := range snapshots {
//		log.Printf("Commit number: %d, Time: %s", index, snapshot.Commit.Time)
//	}
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
		//		err = os.RemoveAll(repoName)
		//		if err != nil {
		//			return nil, err
		//		}
		//		log.Println("Cleaned up repo [" + repoName + "]")

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

	var parentCommit *git.Commit

	err = walker.Iterate(func(c *git.Commit) bool {
		if parentCommit != nil {
			options, err := git.DefaultDiffOptions()
			if err != nil {
				return false
			}

			parentTree, _ := parentCommit.Tree()
			currentTree, _ := c.Tree()

			diff, err := repo.DiffTreeToTree(parentTree, currentTree, &options)
			if err != nil {
				return false
			}

			log.Println("COMMIT" + c.Message())
			err = diff.ForEach(func(file git.DiffDelta, progress float64) (git.DiffForEachHunkCallback, error) {
//				blob, _ := repo.LookupBlob(file.NewFile.Oid)
//
//				sum := string(blob.Contents())
//
//				log.Printf("Name: %s, %d", file.NewFile.Path, len([]byte(sum)))

				return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
					log.Printf(
						"Old range: %s, New range: %s",
						strconv.Itoa(hunk.OldStart) + "-" + strconv.Itoa(hunk.OldStart + hunk.OldLines),
						strconv.Itoa(hunk.NewStart) + "-" + strconv.Itoa(hunk.NewStart + hunk.NewLines),
					)
					return func(line git.DiffLine) error {
						return nil
					}, nil
				}, nil
//				newFile := newFile(file.NewFile.Path)
//
//				log.Printf("Filepath: %v", file.NewFile.Path)
//				return func(hunk git.DiffHunk) (git.DiffForEachLineCallback, error) {
//					newFile.Insertions = append(newFile.Insertions, CodeBlock{Start: hunk.NewStart + 3, End: hunk.NewStart + hunk.NewLines - 1 })
//					newFile.Deletions = append(newFile.Insertions, CodeBlock{Start: hunk.NewStart + 3, End: hunk.NewStart + hunk.NewLines - 1 })
//					log.Printf("Hunk: %v", hunk.Header)
//					return func(line git.DiffLine) error {
//						log.Printf("%s %d, NUMLINES: %d", line.Content, line.Origin, line.NumLines)
//						return nil
//					}, nil
//				}, nil
			}, git.DiffDetailHunks)

			if err != nil {
				panic(err)
			}
		}

		commit := Commit{
			Hash: c.Id().String(),
			Message: c.Message(),
			Time: c.Committer().When,
			Developer: Developer{
				Name: c.Committer().Name,
				Email: c.Committer().Email,
			},
		}

		snapshot := Snapshot{
			Commit: commit,
		}

		snapshots = append(snapshots, snapshot)

		parentCommit = c
		return true
	})
	if err != nil {
		return snapshots, err
	}

	log.Printf("Finished processing repo, duration %d seconds.", time.Now().Unix() - start.Unix())

	return snapshots, nil
}
