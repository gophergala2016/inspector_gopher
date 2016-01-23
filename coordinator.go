package inspector
import "github.com/libgit2/git2go"

type Coordinator struct {
	Repo string
}

func NewCoordinator(repo string) *Coordinator {
	return &Coordinator{
		Repo: repo,
	}
}

func (c *Coordinator) Heatmap() string {
	repo, _ := GetRepo(c.Repo)


	WalkCommits(repo, func(previous, current *git.Commit) bool {
		diff, err := GetDiff(repo, previous, current)
		if err != nil {
			return false
		}
		defer diff.Free()

		WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
		})

		return true
	})

	defer repo.Free()
	defer CleanTempDir()

	return c.Repo
}
