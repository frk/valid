package search

import (
	"go/types"
	"log"
	"os"
	"testing"

	"github.com/frk/compare"
)

var testast AST

func TestMain(m *testing.M) {
	if _, err := Search("../testdata/search", false, nil, &testast); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestFindFunc(t *testing.T) {
	tests := []struct {
		pkgpath string
		name    string
		err     error
	}{
		{"strings", "Contains", nil},
		{"github.com/frk/isvalid/internal/testdata/search", "IsFoo", nil},

		{"strings", "Abracadabra", findFuncError{"strings", "Abracadabra"}},
		{"sgnirts", "Contains", pkgLoadError{"sgnirts", "Contains", nil}},
	}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for i, tt := range tests {
		f, err := FindFunc(tt.pkgpath, tt.name, testast)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("Error: %v", e)
		} else if err == nil {
			if p := f.Pkg(); p.Path() != tt.pkgpath || f.Name() != tt.name {
				t.Errorf("#%d: want=%s.%s; got err=%v", i, tt.pkgpath, tt.name, f)
			}
		}
	}
}

func TestFindConstantsByType(t *testing.T) {
	type konst struct {
		name string
		pkg  string
		typ  string
	}

	searchpkg := "github.com/frk/isvalid/internal/testdata/search"

	tests := []struct {
		pkgpath string
		typname string
		want    []konst
	}{{
		pkgpath: searchpkg, typname: "K1",
		want: []konst{
			{name: "k11", pkg: searchpkg, typ: "K1"},
			{name: "k12", pkg: searchpkg, typ: "K1"},
			{name: "k13", pkg: searchpkg, typ: "K1"},
			{name: "k14", pkg: searchpkg, typ: "K1"},
			{name: "k1111", pkg: searchpkg, typ: "K1"},
		},
	}, {
		pkgpath: searchpkg, typname: "K2",
		want: []konst{
			{name: "k21", pkg: searchpkg, typ: "K2"},
			{name: "k22", pkg: searchpkg, typ: "K2"},
			{name: "k23", pkg: searchpkg, typ: "K2"},
		},
	}, {
		pkgpath: searchpkg, typname: "K3",
		want: []konst{
			{name: "k3foo", pkg: searchpkg, typ: "K3"},
			{name: "k3bar", pkg: searchpkg, typ: "K3"},
			{name: "k3baz", pkg: searchpkg, typ: "K3"},
		},
	}, {
		pkgpath: "go/parser", typname: "Mode",
		want: []konst{
			{name: "CustomMode", pkg: searchpkg, typ: "Mode"},
			{name: "PackageClauseOnly", pkg: "go/parser", typ: "Mode"},
			{name: "ImportsOnly", pkg: "go/parser", typ: "Mode"},
			{name: "ParseComments", pkg: "go/parser", typ: "Mode"},
			{name: "Trace", pkg: "go/parser", typ: "Mode"},
			{name: "DeclarationErrors", pkg: "go/parser", typ: "Mode"},
			{name: "SpuriousErrors", pkg: "go/parser", typ: "Mode"},
			{name: "AllErrors", pkg: "go/parser", typ: "Mode"},
		},
	}}

	for _, tt := range tests {
		got := FindConstantsByType(tt.pkgpath, tt.typname, testast)
		if len(got) != len(tt.want) {
			t.Errorf("%s.%s: len got=%d; want=%d", tt.pkgpath, tt.typname, len(got), len(tt.want))
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

				t.Errorf("%s.%s: no match for %+v", tt.pkgpath, tt.typname, k)
			}
		}
	}
}
