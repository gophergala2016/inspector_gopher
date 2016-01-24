package inspector

import (
	"github.com/libgit2/git2go"
	"log"
	"time"
)

//type Contributor struct {
//	Name string
//	Email string
//	Files []*File
//	Units []*Unit
//	Commits []*Commit
//}
//
//type Commit struct {
//	Hash string
//	Contributor *Contributor
//	Message string
//	Time time.Time
//	Files []*File
//	Unit []*Unit
//}
//
//type File struct {
//	Path string
//	Units []*Unit
//	Commits []*Commit
//}
//
//type Unit struct {
//	UnitType string
//	Name string
//	Signature string
//	File *File
//}
//
//var files map[string]File
//
//func init() {
//	files = make(map[string]File)
//}

func Harvest(repoName string) string {
	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()

	count := 0
	WalkCommits(repo, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		if (len(files) == 0) {
			diff, _ := GetDiff(repo, previousCommit, nil)
			WalkFiles(diff, func(file git.DiffDelta, process float64) {
				log.Printf("FILE: %s", file.OldFile.Path)
			})
			return false
		}

		diff, _ := GetDiff(repo, previousCommit, currentCommit)

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {

			count += 1
		})

		return true
	})

	log.Println(count)
	return "super"
}

func HarvestBenched(repoName string, depth int) float64 {
	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()

	start := time.Now()
	count := 0
	WalkDepthCommits(repo, depth, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		diff, _ := GetDiff(repo, previousCommit, currentCommit)

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			count += 1
		})

		return true
	})

	log.Println(count)

	return time.Since(start).Seconds()
}