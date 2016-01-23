package inspector
import "github.com/libgit2/git2go"

// Types of unit.
const (
	UNIT_TYPE_FUNCTION = iota
	UNIT_TYPE_STRUCT
	UNIT_TYPE_HUNK
)

// File represents a single .go file and holds all revisions of itself
type File struct {
	Path    string
	Changes []FileRevision
}

// File revision represents the files state in a single commit.
type FileRevision struct {
	NumberOfLines int
	Units         []Unit
}

// Unit is either a function or a data structure.
// Holds enough information for later determining
// whether it changed or not.
type Unit struct {
	Name      string
	LineStart int
	LineEnd   int
	Type      int
}

func UnitFromHunk(file git.DiffFile, hunk git.DiffHunk) (*Unit) {
	start := 0

	if hunk.NewStart > 3 {
		start = 3
	}

	start += hunk.NewStart
	end := hunk.NewStart + hunk.NewLines

	if file.Size - end >= 4 {
		end -= 4
	}

	return &Unit {
		Name: "Hunk",
		LineStart: start,
		LineEnd: end,
		Type: UNIT_TYPE_HUNK,
	}
}

func (f *FileRevision) AddUnit(u *Unit) {
	f.Units = append(f.Units, *u)
}

func (u *Unit) ContainsLine(line int) bool {
	return line >= u.LineStart && line <= u.LineEnd
}

func (u1 *Unit) InRange(u2 *Unit) bool {
	return u1.LineStart >= u2.LineStart && u1.LineStart <= u2.LineEnd && u1.LineEnd >= u2.LineStart && u1.LineEnd <= u2.LineEnd
}

// Checks if either the beginning or the end line are contained in the hunk
func (u1 *Unit) Intersects(u2 *Unit) bool {
	return (u1.LineStart >= u2.LineStart && u1.LineStart <= u2.LineEnd) || (u1.LineEnd >= u2.LineStart && u1.LineEnd <= u2.LineEnd)
}

// Number of lines that the Unit has
func (u *Unit) Size() int {
	return u.LineEnd - u.LineStart
}
