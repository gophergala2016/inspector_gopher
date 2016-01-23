package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"github.com/gophergala2016/inspector_gopher/common"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseFileContents(filePath string, contents string) common.FileRevision {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, contents, 0)

	check(err)

	units := []common.Unit{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			u := common.Unit{
				Name:      x.Name.Name,
				LineStart: fset.Position(x.Body.Lbrace).Line,
				LineEnd:   fset.Position(x.Body.Rbrace).Line,
				Type:      common.UNIT_TYPE_FUNCTION,
			}
			units = append(units, u)
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {
				u := common.Unit{
					Name:      x.Name.Name,
					LineStart: fset.Position(x.Pos()).Line,
					LineEnd:   fset.Position(x.End()).Line,
					Type:      common.UNIT_TYPE_STRUCT,
				}
				units = append(units, u)
			}
		}

		return true
	})

	file := common.FileRevision{NumberOfLines: fset.Position(f.End()).Line, Units: units}

	return file
}
