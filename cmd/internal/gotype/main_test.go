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
		if ok {
			return named
		}
		alias, ok := typeName.Type().(*types.Alias)
		if ok {
			return alias
		}
		continue
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

type test_values struct {
	string  *Type
	int     *Type
	int32   *Type
	int64   *Type
	uint    *Type
	uint8   *Type
	uint16  *Type
	uint64  *Type
	bool    *Type
	float64 *Type
	rune    *Type
	byte    *Type
	bytes   *Type
	error   *Type
	pkg     Pkg
}

var T = test_values{
	string:  &Type{Kind: K_STRING},
	int:     &Type{Kind: K_INT},
	int32:   &Type{Kind: K_INT32},
	int64:   &Type{Kind: K_INT64},
	uint:    &Type{Kind: K_UINT},
	uint8:   &Type{Kind: K_UINT8},
	uint16:  &Type{Kind: K_UINT16},
	uint64:  &Type{Kind: K_UINT64},
	float64: &Type{Kind: K_FLOAT64},
	bool:    &Type{Kind: K_BOOL},
	rune:    &Type{Kind: K_INT32, IsRune: true},
	byte:    &Type{Kind: K_UINT8, IsByte: true},
	bytes:   &Type{Kind: K_SLICE, Elem: &Type{Kind: K_UINT8, IsByte: true}},
	error: &Type{Kind: K_INTERFACE, Name: "error",
		Methods: []*Method{{Name: "Error", IsExported: true,
			Type: &Type{Kind: K_FUNC, Out: []*Var{{Type: &Type{Kind: K_STRING}}}},
		}},
	},

	pkg: Pkg{
		Path: "github.com/frk/valid/cmd/internal/gotype/testdata",
		Name: "testdata",
	},
}
