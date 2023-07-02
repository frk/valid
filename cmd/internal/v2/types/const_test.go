package types

import (
	"testing"

	"github.com/frk/compare"
)

func TestFindConsts(t *testing.T) {
	pkg0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/v2/types/testdata",
		Name: "testdata",
	}
	pkg1 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/v2/types/testdata/mypkg",
		Name: "mypkg",
	}

	tests := []struct {
		typ  *Type
		want []Const
	}{{
		typ: &Type{Pkg: pkg0, Name: "ConstType1"},
		want: []Const{
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType1"}, Name: "CT1"},
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType1"}, Name: "CT2"},
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType1"}, Name: "CT4"},
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType1"}, Name: "ct5"},
		},
	}, {
		typ: &Type{Pkg: pkg0, Name: "ConstType2"},
		want: []Const{
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType2"}, Name: "ConstFoo"},
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType2"}, Name: "ConstBar"},
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "ConstType2"}, Name: "const_baz"},
		},
	}, {
		typ: &Type{Pkg: pkg0, Name: "constType3"},
		want: []Const{
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "constType3"}, Name: "kYES"},
			{Pkg: pkg0, Type: &Type{Pkg: pkg0, Name: "constType3"}, Name: "kNO"},
		},
	}, {
		typ: &Type{Pkg: pkg1, Name: "ConstType1"},
		want: []Const{
			{Pkg: pkg1, Type: &Type{Pkg: pkg1, Name: "ConstType1"}, Name: "CT1"},
			{Pkg: pkg1, Type: &Type{Pkg: pkg1, Name: "ConstType1"}, Name: "CT2"},
			{Pkg: pkg1, Type: &Type{Pkg: pkg1, Name: "ConstType1"}, Name: "CT4"},
			{Pkg: pkg1, Type: &Type{Pkg: pkg1, Name: "ConstType1"}, Name: "ct5"},
		},
	}}

	for _, tt := range tests {
		got := FindConsts(tt.typ, &test_src)
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
