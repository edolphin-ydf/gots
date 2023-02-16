package main

import (
	"flag"
	"fmt"
	"go/token"
	"os"

	"github.com/edolphin-ydf/gots/util"
	"github.com/sshelll/sinfra/ast"
)

var (
	file = flag.String("f", "", "file path")
	pos = flag.Int("p", 0, "position in file")
)

func main() {
	flag.Parse()
	if *file == "" {
		fmt.Fprintln(os.Stderr, "file name is required")
		os.Exit(1)
	}

	fInfo, err := ast.Parse(*file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[gott] ast parse failed: %v\n", err.Error())
	}

	funcName := util.FindNearstTestFunc(fInfo, token.Pos(*pos))
	funcName = "^" + funcName + "$"

	fmt.Fprintln(os.Stdout, funcName)
}
