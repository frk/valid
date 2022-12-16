package xtypes

import (
	"testing"

	"github.com/frk/compare"
)

func TestConsts(t *testing.T) {
	pkg0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/xtypes/testdata",
		Name: "testdata",
	}
	pkg1 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/xtypes/testdata/mypkg",
		Name: "mypkg",
	}

	tests := []struct {
		typ  *Type
		want []Const
	}{{
		typ: &Type{Pkg: pkg0, Name: "ConstType1"},
		want: []Const{
			{Pkg: pkg0, Name: "CT1"},
			{Pkg: pkg0, Name: "CT2"},
			{Pkg: pkg0, Name: "CT4"},
			{Pkg: pkg0, Name: "ct5"},
		},
	}, {
		typ: &Type{Pkg: pkg0, Name: "ConstType2"},
		want: []Const{
			{Pkg: pkg0, Name: "ConstFoo"},
			{Pkg: pkg0, Name: "ConstBar"},
			{Pkg: pkg0, Name: "const_baz"},
		},
	}, {
		typ: &Type{Pkg: pkg0, Name: "constType3"},
		want: []Const{
			{Pkg: pkg0, Name: "kYES"},
			{Pkg: pkg0, Name: "kNO"},
		},
	}, {
		// make sure the unexported constant ct5 is
		// omitted if the const's type is imported.
		typ: &Type{Pkg: pkg1, Name: "ConstType1"},
		want: []Const{
			{Pkg: pkg1, Name: "CT1"},
			{Pkg: pkg1, Name: "CT2"},
			{Pkg: pkg1, Name: "CT4"},
		},
	}}

	a := NewAnalyzer(test_pkg.Type)
	for _, tt := range tests {
		got := a.Consts(tt.typ, &test_ast)
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
