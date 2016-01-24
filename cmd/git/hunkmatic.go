package main
import (
"github.com/gophergala2016/inspector_gopher"
"github.com/libgit2/git2go"
	"log"
)

func main() {

	repo, _ := inspector.GetRepo("bogdanhabic/testrepo")

	inspector.WalkCommits(repo, func(previousCommit *git.Commit, currentCommit *git.Commit) bool {
		diff, _ := inspector.GetDiff(repo, previousCommit, currentCommit)

		log.Printf("--COMMIT-- \n%s\n--ENDMESSAGE--", currentCommit.Message())

		stats, _ := diff.Stats()
		log.Printf("INSERTIONS: %d DELETIONS: %d", stats.Insertions(), stats.Deletions())

		inspector.WalkHunks(diff, func(file git.DiffDelta, hunk git.DiffHunk) {
			log.Printf("%s: FROM: %d TO: %d | FROM: %d TO: %d",
				file.NewFile.Path,
				hunk.OldStart, hunk.OldStart + hunk.OldLines,
				hunk.NewStart, hunk.NewStart + hunk.NewLines,
			)
		})

		return true
	})
}
