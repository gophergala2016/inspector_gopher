package main

import (
	"go/ast"
	"go/token"
	"go/parser"
	"io/ioutil"
	"os"
	"fmt"
)


// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseFileContents(filePath string, contents string) FileRevision {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, contents, 0)

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

	file := FileRevision{NumberOfLines:fset.Position(f.End()).Line, Units:units}

	return file
}

func main() {
	dat, err := ioutil.ReadFile("cmd" + string(os.PathSeparator) + "units" + string(os.PathSeparator) + "to_parse.go")

	check(err)

	parsedFile := parseFileContents("to_parse.go", string(dat))

	for _, unit := range parsedFile.Units {
		fmt.Printf("Type: %d Name: %s From: %d To: %d\n", unit.Type, unit.Name, unit.LineStart, unit.LineEnd)
	}
}

