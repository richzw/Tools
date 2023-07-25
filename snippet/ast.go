package snippet

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func GetFunctionLines(filename string, funcName string) (int, int, error) {
	// Parse the Go source file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return 0, 0, err
	}

	// Traverse the AST and look for function declarations
	var beginLine, endLine int
	ast.Inspect(node, func(n ast.Node) bool {
		// Check if the node is a function declaration with the given name
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == funcName {
			// Get the beginning and ending positions of the function
			beginPos := fset.Position(fn.Pos())
			endPos := fset.Position(fn.End())

			// Convert the positions to line numbers
			beginLine = beginPos.Line
			endLine = endPos.Line

			// Stop traversing the AST
			return false
		}

		// Continue traversing the AST
		return true
	})

	// Check if the function was found
	if beginLine == 0 && endLine == 0 {
		return 0, 0, fmt.Errorf("function %s not found in file %s", funcName, filename)
	}

	return beginLine, endLine, nil
}
