package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/astutil"
)

const input = `package main

import "fmt"

func main() {
	x := 1
        fmt.Println(x)
}`

func main() {
	var config loader.Config
	config.Import(

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "input.go", input, 0)
	if err != nil {
		log.Fatal(err) // parse error
	}

	conf := types.Config{Importer: importer.Default()}
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	if _, err := conf.Check("cmd/hello", fset, []*ast.File{f}, &info); err != nil {
		log.Fatal(err) // type error
	}

	name := "x"
	for _, obj := range info.Defs {
		if obj != nil && obj.Name() == name {
			fmt.Printf("%s: %T %s\n", fset.Position(obj.Pos()), obj, obj)
		}
	}
}
