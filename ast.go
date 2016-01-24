package inspector

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Extracts all of the functions and structures from the file
func ParseFileContents(filePath string, contents string) *File {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, filePath, contents, 0)

	check(err)

	file := File{
		Path: filePath,
		NumberOfLines: fset.Position(f.End()).Line,
		Units: []*Unit{},
		Commits: []*Commit{},
	}

	// We walk the syntax tree
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Functions have the position of their brackets
			u := Unit{
				Type:      UNIT_TYPE_FUNCTION,
				Name:      x.Name.Name,

				LineStart: fset.Position(x.Body.Lbrace).Line,
				LineEnd:   fset.Position(x.Body.Rbrace).Line,

				RationSum: 0,
				TimesChanged: 0,

				Commits: []*Commit{},
				File: &file,
			}
			file.Units = append(file.Units, &u)
			break
		case *ast.TypeSpec:
			if _, ok := x.Type.(*ast.StructType); ok {

				// Structures only have the position of their beginning and their end (unfortunately no bracket positions)
				u := Unit{
					Type:      UNIT_TYPE_STRUCT,
					Name:      x.Name.Name,

					LineStart: fset.Position(x.Pos()).Line,
					LineEnd:   fset.Position(x.End()).Line,

					RationSum: 0,
					TimesChanged: 0,

					Commits: []*Commit{},
					File: &file,
				}
				file.Units = append(file.Units, &u)
			}
		}

		return true
	})

	return &file
}
