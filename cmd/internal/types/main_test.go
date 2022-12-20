package types

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

func test_named(name string) *types.Named {
	for id, obj := range test_pkg.Info.Defs {
		if id.Name != name {
			continue
		}
		if tn, ok := obj.(*types.TypeName); ok {
			if named, ok := tn.Type().(*types.Named); ok {
				return named
			}
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
	string:  &Type{Kind: STRING},
	int:     &Type{Kind: INT},
	int32:   &Type{Kind: INT32},
	int64:   &Type{Kind: INT64},
	uint:    &Type{Kind: UINT},
	uint8:   &Type{Kind: UINT8},
	uint16:  &Type{Kind: UINT16},
	uint64:  &Type{Kind: UINT64},
	float64: &Type{Kind: FLOAT64},
	bool:    &Type{Kind: BOOL},
	rune:    &Type{Kind: INT32, IsRune: true},
	byte:    &Type{Kind: UINT8, IsByte: true},
	bytes:   &Type{Kind: SLICE, Elem: &Obj{Type: &Type{Kind: UINT8, IsByte: true}}},
	error: &Type{Kind: INTERFACE, Name: "error",
		MethodSet: []*Method{{Name: "Error", IsExported: true,
			Type: &Type{Kind: FUNC, Out: []*Var{{Type: &Type{Kind: STRING}}}},
		}},
	},

	pkg: Pkg{
		Path: "github.com/frk/valid/cmd/internal/types/testdata",
		Name: "testdata",
	},
}
