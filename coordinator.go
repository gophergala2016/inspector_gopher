package inspector

type Coordinator struct {
	repo string 
}

func NewCoordinator(repo string) *Coordinator {
	return &Coordinator{
		repo: repo,
	}
}

func (c *Coordinator) Heatmap() string {
//	gitUnits := (c.repo)
//	astUnits := parseFileContents("/tmp/main.go")

	return "some files"
}
