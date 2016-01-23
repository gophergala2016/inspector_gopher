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

	//	defer repo.Free()
	//	type CommitWalkerFunc func(previousCommit *git.Commit, currentCommit *git.Commit) bool

	WalkCommits(repo, func(previous *git.Commit, current *git.Commit) bool {
		return true
	})

	//	gitUnits := (c.repo)
	//	astUnits := parseFileContents("/tmp/main.go")

	return c.Repo
}
