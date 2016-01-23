package inspector

// Types of unit.
const (
	UNIT_TYPE_FUNCTION = iota
	UNIT_TYPE_STRUCT
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

func (f *FileRevision) AddUnit(u *Unit) {
	f.Units = append(f.Units, *u)
}

func (u *Unit) ContainsLine(line int) bool {
	return line >= u.LineStart && line <= u.LineEnd
}

func (u *Unit) InRange(lineStart int, lineEnd int) bool {
	return u.LineStart >= lineStart && u.LineStart <= lineEnd && u.LineEnd >= lineStart && u.LineEnd <= lineEnd
}

func (u *Unit) Intersects(lineStart int, lineEnd int) bool {
	return (u.LineStart >= lineStart && u.LineStart <= lineEnd) || (u.LineEnd >= lineStart && u.LineEnd <= lineEnd)
}

func (u *Unit) Size() int {
	return u.LineEnd - u.LineStart
}
