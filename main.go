// Copyright 2013 Rocky Bernstein.

// +build ignore

package main

import (
	"os"
	"go/parser"
	"code.google.com/p/go-gnureadline"
	"code.google.com/p/go.tools/importer"
	// "github.com/rocky/ssa-interp/gub"
	// "code.google.com/p/go.tools/go/types"
	// "code.google.com/p/go.tools/go/exact"
	"go/ast"
	"fmt"
)

func setup() (*importer.Importer, *importer.PackageInfo) {
	test := `
package main

import "os"

type testEntry struct {
	src string
    num int
}

var testTypes = []testEntry{ {"a", 1}, {"b", 2}}

func main() {
    intary := []int{1,2,3,4}
    os.Exit(intary[2])
}
`
	imp := importer.New(new(importer.Config)) // no Loader; uses GC importer

	f, err := parser.ParseFile(imp.Fset, "<input>", test, parser.DeclarationErrors)
	if err != nil {
		fmt.Println("parse error: %s\n", err)
		os.Exit(1)
	}

	info := imp.CreateSourcePackage("main", []*ast.File{f})
	if info.Err != nil {
		fmt.Println(info.Err.Error())
		os.Exit(2)
	}
	return imp, info
}

func main() {
	imp, info := setup()
	fmt.Println(*imp)
	fmt.Println("==============================")
	fmt.Println(*info)
	fmt.Println(info.Scope)
	var line string
	var err error
	for i:=1; err == nil && line != "quit"; i++ {
		line, err = gnureadline.Readline(fmt.Sprintf("Enter something [%d]: ", i), true)
		switch line {
		case "quit":
			break;
		}
		fmt.Printf("You typed: %s\n", line)
		expr, err := parser.ParseExpr(line)
		if err != nil {
			fmt.Printf("Error parsing %s: %s\n", line, err.Error())
			continue
		}
		ast.Print(nil, expr)
     }
}
