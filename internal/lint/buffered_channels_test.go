package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"testing"
)

func TestNoBufferedChannels(t *testing.T) {
	// List of files to check for buffered channels
	filesToCheck := []string{
		"../filecrawler/crawler.go",
		"../workerpool/pool.go",
	}

	for _, relPath := range filesToCheck {
		// Get absolute path
		absPath, err := filepath.Abs(relPath)
		if err != nil {
			t.Fatalf("Failed to get absolute path for %s: %v", relPath, err)
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, absPath, nil, parser.AllErrors)
		if err != nil {
			t.Fatalf("Failed to parse file %s: %v", absPath, err)
		}

		// Walk through the AST looking for channel declarations
		ast.Inspect(node, func(n ast.Node) bool {
			if makeExpr, ok := n.(*ast.CallExpr); ok {
				// Check if it's a make() call
				if ident, ok := makeExpr.Fun.(*ast.Ident); ok && ident.Name == "make" {
					// Check if it's making a channel
					if len(makeExpr.Args) >= 1 {
						if _, ok := makeExpr.Args[0].(*ast.ChanType); ok {
							// If there's more than one argument, it's a buffered channel
							// according to https://go.dev/ref/spec#Making_slices_maps_and_channels
							// if there is a second argument, it must be the size.
							if len(makeExpr.Args) > 1 {
								t.Errorf("File %s contains a buffered channel at position %v",
									relPath, fset.Position(makeExpr.Pos()))
							}
						}
					}
				}
			}
			return true
		})
	}
}
