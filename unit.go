package main

const (
	UNIT_TYPE_FUNCTION = iota
	UNIT_TYPE_STRUCT
)

type File struct {
	Path    string
	Changes []FileRevision
}

type FileRevision struct {
	NumberOfLines int
	Units         []Unit
}

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
