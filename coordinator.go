package inspector

import (
	"github.com/libgit2/git2go"
)

const DEPTH = 100

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

	total := 0

//	files := createFilesFromPaths(ListFiles(repo.Path()))

	WalkCommits(repo, func(previous, current *git.Commit) bool {
		if total >= DEPTH {
			return false
		} else {
			total += 1
		}

		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			return false
		}
		defer diff.Free()

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {

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

//			units = append(units, *UnitFromHunk(file.NewFile, hunk))
		})

		return true
	})

	return c.RepoName
}