package inspector

import (
	"github.com/libgit2/git2go"
	"log"
	"strings"
	"net/http"
)

func Harvest(repoName string) *Everything {
	resp, _ := http.Get("http://github.com/" + repoName)
	if resp.StatusCode == 404 {
		return nil
	}

	repo, _ := GetRepo(repoName)
	defer repo.Free()
	defer CleanTempDir()

	var everything *Everything = &Everything{}
	var files map[string]*File = make(map[string]*File)
	var commits map[string]*Commit = make(map[string]*Commit)
	everything.Files = files
	everything.Commits = commits

	numberOfCommits, _ := GetNumberOfCommits(repo)

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
		count++

		diff, _ := GetDiff(repo, previousCommit, currentCommit)

		newestCommit := &Commit{
			Hash: currentCommit.Id().String(),
			Contributor: nil,
			Message: currentCommit.Message(),
			Time: currentCommit.Author().When,
		}

		commits[newestCommit.Hash] = newestCommit
		var commitFiles map[string]*File = make(map[string]*File)

		var checked map[string]bool = map[string]bool{}
		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			if !strings.HasSuffix(strings.ToLower(file.OldFile.Path), ".go") {
				return
			}

			var newFile *File
			if commitFiles[file.NewFile.Path] != nil {
				newFile = commitFiles[file.NewFile.Path]
			} else {
				blob, err := repo.LookupBlob(file.NewFile.Oid)
				if err != nil {
					return
				}

				newFile = ParseFileContents(file.NewFile.Path, string(blob.Contents()))
			}

			if files[newFile.Path] != nil {
				files[file.OldFile.Path] = newFile
			}

			if !checked[newFile.Path] {
				checked[newFile.Path] = true
				if files[newFile.Path] != nil {
					remove := []int{}
					for i, u := range files[newFile.Path].Units {
						for _, un := range newFile.Units {
							if u.Type == un.Type && u.Name == un.Name {
								remove = append(remove, i)
							}
						}
					}
					newUnits := []*Unit{}
					for i, unit := range files[newFile.Path].Units {
						ff := false
						for _, j := range remove {
							if i == j {
								ff = true
							}
						}
						if !ff {
							newUnits = append(newUnits, unit)
						}
						files[newFile.Path].Units = newUnits
					}
				} else {
					files[newFile.Path] = newFile
				}
			}

			for _, unit := range newFile.Units {
				unit.RatioSum += 0.5
				unit.TimesChanged++
				if unit.IntersectsInt(hunk.NewStart, hunk.NewStart + hunk.NewLines) {
					magicNumber := unit.numberOfLinesIntersected(hunk.NewStart, hunk.NewStart + hunk.NewLines) * len(commits)
					unit.RatioSum += float64(magicNumber) * float64(count) / float64(numberOfCommits)
					unit.Commits = append(unit.Commits, newestCommit)
					if files[file.NewFile.Path] != nil {
						for _, u := range files[file.NewFile.Path].Units {
							if u.Type == unit.Type && u.Name == unit.Name {
								u.LineStart = unit.LineStart
								u.LineEnd = unit.LineEnd
							}
						}
					} else {
						files[file.NewFile.Path] = newFile
					}
				} else {

				}
			}
		})

		return true
	})

	return everything
}