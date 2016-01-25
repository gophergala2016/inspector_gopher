package inspector

import (
	"testing"
)

func TestHarvest(t *testing.T) {

	//	repo, _ := GetRepo("lazartravica/envy")
	//	head, _ := repo.Head()
	//	commit, _ := repo.LookupCommit(head.Target())
	//
	//	diff, _ := GetDiff(repo, commit, nil)
	//	WalkFiles(diff, func(file git.DiffDelta, process float64) {
	//		log.Printf("FILE: %s", file.OldFile.Path)
	//	})

	Harvest("asdfasdf/asdfsd")
}