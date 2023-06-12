package search

import (
	"go/types"
	"testing"
)

func TestFindConstantsByType(t *testing.T) {
	type konst struct {
		name string
		pkg  string
		typ  string
	}

	searchpkg := "github.com/frk/valid/cmd/internal/v2/search/testdata"

	tests := []struct {
		pkg     string
		typname string
		want    []konst
	}{{
		pkg: searchpkg, typname: "K1",
		want: []konst{
			{name: "k11", pkg: searchpkg, typ: "K1"},
			{name: "k12", pkg: searchpkg, typ: "K1"},
			{name: "k13", pkg: searchpkg, typ: "K1"},
			{name: "k14", pkg: searchpkg, typ: "K1"},
			{name: "k1111", pkg: searchpkg, typ: "K1"},
		},
	}, {
		pkg: searchpkg, typname: "K2",
		want: []konst{
			{name: "k21", pkg: searchpkg, typ: "K2"},
			{name: "k22", pkg: searchpkg, typ: "K2"},
			{name: "k23", pkg: searchpkg, typ: "K2"},
		},
	}, {
		pkg: searchpkg, typname: "K3",
		want: []konst{
			{name: "k3foo", pkg: searchpkg, typ: "K3"},
			{name: "k3bar", pkg: searchpkg, typ: "K3"},
			{name: "k3baz", pkg: searchpkg, typ: "K3"},
		},
	}, {
		pkg: "go/parser", typname: "Mode",
		want: []konst{
			{name: "CustomMode", pkg: searchpkg, typ: "Mode"},
			{name: "PackageClauseOnly", pkg: "go/parser", typ: "Mode"},
			{name: "ImportsOnly", pkg: "go/parser", typ: "Mode"},
			{name: "ParseComments", pkg: "go/parser", typ: "Mode"},
			{name: "Trace", pkg: "go/parser", typ: "Mode"},
			{name: "DeclarationErrors", pkg: "go/parser", typ: "Mode"},
			{name: "SpuriousErrors", pkg: "go/parser", typ: "Mode"},
			{name: "SkipObjectResolution", pkg: "go/parser", typ: "Mode"}, // this constant was added in 1.17
			{name: "AllErrors", pkg: "go/parser", typ: "Mode"},
		},
	}}

	for _, tt := range tests {
		got := FindConstantsByType(tt.pkg, tt.typname, &testast)
		if len(got) != len(tt.want) {
			t.Errorf("%s.%s: len got=%d; want=%d", tt.pkg, tt.typname, len(got), len(tt.want))
		} else {
		kloop:
			for _, k := range tt.want {
				for _, c := range got {
					if typeNamed, ok := c.Type().(*types.Named); ok {
						name := c.Name()
						pkg := c.Pkg().Path()
						typ := typeNamed.Obj().Name()
						if k.name == name && k.pkg == pkg && k.typ == typ {
							continue kloop
						}
					}
				}

				t.Errorf("%s.%s: no match for %+v", tt.pkg, tt.typname, k)
			}
		}
	}
}
