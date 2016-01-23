package main

import (
	"io/ioutil"
	"go/ast"
	"go/token"
	"go/parser"
	"os"
	"container/list"
	"fmt"
)

const (
	UNIT_TYPE_FUNCTION = iota
	UNIT_TYPE_STRUCT = iota
)

type Unit struct {
	Name      string
	LineStart int
	LineEnd   int
	Type      int
}

func (u *Unit) ContainsLine(line int) bool {
	return line >= u.LineStart && line <= u.LineEnd
}

func (u *Unit) InRange(lineStart int, lineEnd int) bool {
	return u.LineStart >= lineStart && u.LineStart <= lineEnd && u.LineEnd >= lineStart && u.LineEnd <= lineEnd
}


// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseFileContents(fileName string, contents string) *list.List {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, fileName, contents, 0)

	check(err)

	units := list.New()

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			u := Unit{
				Name:x.Name.Name,
				LineStart:fset.Position(x.Body.Lbrace).Line,
				LineEnd:fset.Position(x.Body.Rbrace).Line,
				Type:UNIT_TYPE_FUNCTION,
			}
			units.PushFront(u)
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {
				u := Unit{
					Name:x.Name.Name,
					LineStart:fset.Position(x.Pos()).Line,
					LineEnd:fset.Position(x.End()).Line,
					Type:UNIT_TYPE_STRUCT,
				}
				units.PushFront(u)
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

	for e := units.Front(); e != nil; e = e.Next() {
		unit := e.Value.(Unit)
		fmt.Printf("Type: %d Name: %s From: %d To: %d\n", unit.Type, unit.Name, unit.LineStart, unit.LineEnd)
	}
}
