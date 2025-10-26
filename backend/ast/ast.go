package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// ExtractNodeTypesFromDir returns a map[nodeType]count for all Go files in dir
func ExtractNodeTypesFromDir(dir string) (map[string]int, error) {
	nodeFreq := make(map[string]int)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Println("Skipping file (parse error):", path)
			return nil
		}

		ast.Inspect(f, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			nodeType := fmt.Sprintf("%T", n)
			nodeFreq[nodeType]++
			return true
		})
		return nil
	})

	return nodeFreq, err
}

// WeightedJaccard computes similarity between two node frequency maps
func WeightedJaccard(a, b map[string]int) float64 {
	var intersection, union float64
	seen := make(map[string]struct{})

	for node, fa := range a {
		fb := b[node]
		intersection += float64(min(fa, fb))
		union += float64(max(fa, fb))
		seen[node] = struct{}{}
	}

	for node, fb := range b {
		if _, ok := seen[node]; !ok {
			union += float64(fb)
		}
	}

	if union == 0 {
		return 0.0
	}
	return intersection / union
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
