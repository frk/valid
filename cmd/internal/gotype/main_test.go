package gotype

import (
	"go/types"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/frk/valid/cmd/internal/search"
)

var test_ast search.AST
var test_pkg search.Package

func TestMain(m *testing.M) {
	pkgs, err := search.Search(
		"testdata/",
		true,
		regexp.MustCompile(`(?i:Validator)$`),
		nil,
		&test_ast,
	)
	if err != nil {
		log.Fatal(err)
	}

	test_pkg = *(pkgs[0])

	os.Exit(m.Run())
}

func test_type(name string) types.Type {
	for id, obj := range test_pkg.Info.Defs {
		if id.Name != name {
			continue
		}

		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		named, ok := typeName.Type().(*types.Named)
		if !ok {
			continue
		}
		return named
	}
	return nil
}

func test_obj(name string) types.Object {
	for id, obj := range test_pkg.Info.Defs {
		if id.Name == name {
			return obj
		}
	}
	return nil
}
