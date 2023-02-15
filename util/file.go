package util

import (
	"fmt"
	goast "go/ast"
	"go/token"
	"math"
	"strings"

	"github.com/sshelll/sinfra/ast"
)

const (
	testingPkgName = `testing`
	testifyPkgName = `github.com/stretchr/testify/suite`
)

type Func struct {
	goast.Node
	Name string
}

func Min[T ~int | ~float64 | ~float32](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Abs[T ~int | ~float64 | ~float32](a T) T {
	if a < 0 {
		return -a
	}

	return a
}

func FindNearstTestFunc(f *ast.File, pos token.Pos) string {
	goTests, testifyTests := ExtractTestFuncs(f), ExtractTestifySuiteTestMethods(f)
	testList := append(testifyTests, goTests...)

	for _, f := range testList {
		if pos >= f.Pos() && pos <= f.End() {
			return f.Name
		}
	}

	var nearstFunc *Func
	var minDistance token.Pos = math.MaxInt
	for _, f := range testList {
		dis := Min(Abs(pos-f.Pos()), Abs(pos-f.End()))
		if nearstFunc == nil {
			nearstFunc = f
			minDistance = dis
		}
		if dis < minDistance {
			nearstFunc = f
			minDistance = dis
		}
	}

	return nearstFunc.Name
}

func ExtractTestFuncs(f *ast.File) []*Func {
	fnList := make([]*Func, 0, len(f.FuncList))
	testingPkg := findTestingPkgName(f.ImportList)
	for _, fn := range f.FuncList {
		if ast.IsGoTestFunc(fn, &testingPkg) {
			fnList = append(fnList, &Func{
				Node: fn.AstDecl,
				Name: fn.Name,
			})
		}
	}
	return fnList
}

func ExtractTestifySuiteTestMethods(f *ast.File) []*Func {

	testingPkg := findTestingPkgName(f.ImportList)
	testifyPkg := findTestifyPkgName(f.ImportList)

	suiteEntryMap := make(map[string]string)
	for _, fn := range f.FuncList {
		suiteName, ok := ast.IsTestifySuiteEntryFunc(fn, &testingPkg, &testifyPkg)
		if ok {
			suiteEntryMap[suiteName] = fn.Name
		}
	}

	methodList := make([]*Func, 0, 16)
	for _, s := range f.StructList {
		entryName, ok := suiteEntryMap[s.Name]
		if !ok {
			continue
		}
		for _, m := range s.MethodList {
			if strings.HasPrefix(m.Name, "Test") {
				methodList = append(methodList, &Func{
					Node: m.AstDecl,
					Name: fmt.Sprintf("%s/%s", entryName, m.Name),
				})
			}
		}
	}

	return methodList

}

func findTestifyPkgName(importList []*ast.Import) string {
	alias := findPkgAlias(importList, testifyPkgName)
	if alias != nil {
		return *alias
	}
	return "suite"
}

func findTestingPkgName(importList []*ast.Import) string {
	alias := findPkgAlias(importList, testingPkgName)
	if alias != nil {
		return *alias
	}
	return "testing"
}

func findPkgAlias(importList []*ast.Import, pkg string) (alias *string) {
	for _, imp := range importList {
		if imp.Pkg == pkg {
			if imp.Alias == "" {
				return nil
			}
			return &imp.Alias
		}
	}
	return nil
}
