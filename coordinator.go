package inspector

import (
	"github.com/libgit2/git2go"
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

	units := []Unit {}

	WalkCommits(repo, func(previous, current *git.Commit) bool {
		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			return false
		}
		defer diff.Free()

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			//TODO:
			//	-	Checkout at the given commit
			//	-	Get AST to parse the given commit
			//	-	Run unitHunk.Intersects(*unitAst)
			//	-	Flag if true

			units = append(units, *UnitFromHunk(hunk))
		})

		return true
	})

	defer repo.Free()
	defer CleanTempDir()

	return c.RepoName
}
