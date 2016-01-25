package inspector

import (
	"github.com/libgit2/git2go"
	"log"
	"time"
	"strings"
)

func Harvest(repoName string) *Everything {
	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()


	var files map[string]*File = make(map[string]*File)

	count := 0
	WalkCommits(repo, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		if (len(files) == 0) {
			//Initially populate the files map with the whole file structure at the specific commit
			log.Printf("[START] Create base tree")
			diff, _ := GetDiff(repo, previousCommit, nil)
			WalkFiles(diff, func(file git.DiffDelta, process float64) {
				if !strings.HasSuffix(strings.ToLower(file.OldFile.Path), ".go") {
					return
				}

				blob, err := repo.LookupBlob(file.OldFile.Oid)
				if err != nil {
					return
				}

				astFile := ParseFileContents(file.OldFile.Path, string(blob.Contents()))

				files[file.OldFile.Path] = astFile
			})

			log.Printf("[Success] Create base tree")
			return true
		}

		diff, _ := GetDiff(repo, previousCommit, currentCommit)

//		var commitFiles map[string]*File = make(map[string]*File)

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			if !strings.HasSuffix(strings.ToLower(file.OldFile.Path), ".go") {
				return
			}

			blob, err := repo.LookupBlob(file.OldFile.Oid)
			if err != nil {
				return
			}

//			var newFile *File
//			if commitFiles[file.NewFile.Path] != nil {
//				newFile = commitFiles[file.NewFile.Path]
//			} else {
//				newFile = ParseFileContents()
//			}

			astFile := ParseFileContents(file.OldFile.Path, string(blob.Contents()))

			files[file.OldFile.Path] = astFile
		})

		return true
	})


	log.Println(count)

	return &Everything{}
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