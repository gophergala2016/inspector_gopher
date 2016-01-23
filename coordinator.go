package inspector

import (
	"github.com/libgit2/git2go"
	"log"
)

type Coordinator struct {
	RepoName string
}

func NewCoordinator(repo string) *Coordinator {
	return &Coordinator{
		RepoName: repo,
	}
}

func (c *Coordinator) Heatmap() string {
	repo, _ := GetRepo(c.RepoName)
	defer repo.Free()
	defer CleanTempDir()

	units := []Unit{}

	WalkCommits(repo, func(previous, current *git.Commit) bool {
		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			return false
		}
		defer diff.Free()

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			if file.NewFile.Oid.IsZero() {
//				log.Printf("FILE DELETED")
				//File deleted
				return
			}
			if file.OldFile.Oid.IsZero() {
//				log.Printf("FILE CREATED")
				//File created
				return
			}

			//File modified
//			blob, err := repo.LookupBlob(file.NewFile.Oid)
//			if err != nil {
//				panic("Cannot lookup blob, error: " + err.Error())
//			}
//
//			totalBytes += len(blob.Contents())

//			blob.Free()

			//TODO:
			//	-	Checkout at the given commit
			//	-	Get AST to parse the given commit
			//	-	Run unitHunk.Intersects(*unitAst)
			//	-	Flag if true

			units = append(units, *UnitFromHunk(file.NewFile, hunk))
		})

		return true
	})

	defer repo.Free()
	defer CleanTempDir()

	return c.RepoName
}
