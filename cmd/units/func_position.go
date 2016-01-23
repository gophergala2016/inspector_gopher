package main

import (
	"io/ioutil"
	"go/ast"
	"go/token"
	"go/parser"
	"os"
	"fmt"
)

const (
	UNIT_TYPE_FUNCTION = iota
	UNIT_TYPE_STRUCT = iota
)

type File struct {
	Path         string
	NumberOfLine int
	Units        []Unit
}

type Unit struct {
	Name      string
	LineStart int
	LineEnd   int
	Type      int
}

func (f *File) AddUnit(u *Unit) {
	f.Units = append(f.Units, *u)
}

func (u *Unit) ContainsLine(line int) bool {
	return line >= u.LineStart && line <= u.LineEnd
}

func (u *Unit) InRange(lineStart int, lineEnd int) bool {
	return u.LineStart >= lineStart && u.LineStart <= lineEnd && u.LineEnd >= lineStart && u.LineEnd <= lineEnd
}

func (u *Unit) Size() int {
	return u.LineEnd - u.LineStart
}

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseFileContents(fileName string, contents string) []Unit {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, contents, 0)

	check(err)

	units := []Unit{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			u := Unit{
				Name:x.Name.Name,
				LineStart:fset.Position(x.Body.Lbrace).Line,
				LineEnd:fset.Position(x.Body.Rbrace).Line,
				Type:UNIT_TYPE_FUNCTION,
			}
			units = append(units, u)
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {
				u := Unit{
					Name:x.Name.Name,
					LineStart:fset.Position(x.Pos()).Line,
					LineEnd:fset.Position(x.End()).Line,
					Type:UNIT_TYPE_STRUCT,
				}
				units = append(units, u)
			}
		}

		return true
	})

	return units
}

func main() {
	dat, err := ioutil.ReadFile("cmd" + string(os.PathSeparator) + "units" + string(os.PathSeparator) + "to_parse.go")

	check(err)

	units := parseFileContents("to_parse.go", string(dat))

	for _, unit := range units {
		fmt.Printf("Type: %d Name: %s From: %d To: %d\n", unit.Type, unit.Name, unit.LineStart, unit.LineEnd)
	}
}
