// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssautil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"golang.org/x/tools/go/ssa/ssautil"
	"golang.org/x/tools/go/types"
)

const hello = `package main

import "fmt"

func main() {
	fmt.Println("Hello, world")
}
`

func TestLoadPackage(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "hello.go", hello, 0)
	if err != nil {
		t.Fatal(err)
	}

	pkg := types.NewPackage("hello", "")
	ssapkg, _, err := ssautil.LoadPackage(new(types.Config), fset, pkg, []*ast.File{f}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Name() != "main" {
		t.Errorf("pkg.Name() = %s, want main", pkg.Name())
	}
	if ssapkg.Func("main") == nil {
		ssapkg.WriteTo(os.Stderr)
		t.Errorf("ssapkg has no main function")
	}
}

func TestLoadPackage_MissingImport(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "bad.go", `package bad; import "missing"`, 0)
	if err != nil {
		t.Fatal(err)
	}

	pkg := types.NewPackage("bad", "")
	ssapkg, _, err := ssautil.LoadPackage(new(types.Config), fset, pkg, []*ast.File{f}, 0)
	if err == nil || ssapkg != nil {
		t.Fatal("LoadPackage succeeded unexpectedly")
	}
}
