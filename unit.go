package inspector

import (
	"time"
	"fmt"
)

// Types of unit.
const (
	UNIT_TYPE_FUNCTION = iota
	UNIT_TYPE_STRUCT
	UNIT_TYPE_HUNK
)

type Everything struct {
	Files map[string]*File
	Commits map[string]*Commit
}

// File represents a single .go file.
type File struct {
	Path          string
	NumberOfLines int
	Units         []*Unit
	Commits       []*Commit
}

// Unit is either a function or a data structure.
// Holds enough information for later determining
// whether it changed or not.
type Unit struct {
	Type         int
	Name         string

	LineStart    int
	LineEnd      int

	RatioSum     int
	TimesChanged int

	Commits      []*Commit
	File         *File
}

type Commit struct {
	Contributor *Contributor
	Hash        string
	Message     string
	Time        time.Time
	Files       []*File
	Units       []*Unit
}

type Contributor struct {
	Name    string
	Email   string
	Files   []*File
	Units   []*Unit
	Commits []*Commit
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

func (f File) String() string {
	units := "Units[\n"
	for _, u := range f.Units {
		units += u.String() + "\n"
	}
	units += "]"

	return fmt.Sprintf(
		"FILE {Path: %s, NumberOfLines: %s, Units: %s, TimesChanged: %d}",
		f.Path,
		f.NumberOfLines,
		units,
		len(f.Commits),
	)
}

func (u Unit) String() string {
	return fmt.Sprintf(
		"UNIT {Type: %v, Name: %s, LineStart: %d, LineEnd: %d, RatioSum: %d, TimesChanged: %d, File: %s}",
		u.Type,
		u.Name,
		u.LineStart,
		u.LineEnd,
		u.RatioSum,
		u.TimesChanged,
		u.File.Path,
	)
}
