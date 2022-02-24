package gotype

import (
	"go/types"
	"testing"

	"github.com/frk/compare"
)

func TestAnalyzer_Analyze(t *testing.T) {
	pkg0 := Pkg{
		Path:  "github.com/frk/valid/cmd/internal/gotype/testdata",
		Name:  "testdata",
		Local: "testdata",
	}
	pkg1 := Pkg{
		Path:  "github.com/frk/valid/cmd/internal/gotype/testdata/mypkg",
		Name:  "mypkg",
		Local: "mypkg",
	}

	tests := []struct {
		ty   types.Type
		want *Type
	}{{
		ty:   test_type("BoolType"),
		want: &Type{Pkg: pkg0, Name: "BoolType", Kind: K_BOOL, IsExported: true},
	}, {
		ty:   test_type("IntType"),
		want: &Type{Pkg: pkg0, Name: "IntType", Kind: K_INT, IsExported: true},
	}, {
		ty:   test_type("int8Type"),
		want: &Type{Pkg: pkg0, Name: "int8Type", Kind: K_INT8},
	}, {
		ty:   test_type("int16Type"),
		want: &Type{Pkg: pkg0, Name: "int16Type", Kind: K_INT16},
	}, {
		ty:   test_type("int32Type"),
		want: &Type{Pkg: pkg0, Name: "int32Type", Kind: K_INT32},
	}, {
		ty:   test_type("int64Type"),
		want: &Type{Pkg: pkg0, Name: "int64Type", Kind: K_INT64},
	}, {
		ty:   test_type("UintType"),
		want: &Type{Pkg: pkg0, Name: "UintType", Kind: K_UINT, IsExported: true},
	}, {
		ty:   test_type("uint8Type"),
		want: &Type{Pkg: pkg0, Name: "uint8Type", Kind: K_UINT8},
	}, {
		ty:   test_type("uint16Type"),
		want: &Type{Pkg: pkg0, Name: "uint16Type", Kind: K_UINT16},
	}, {
		ty:   test_type("uint32Type"),
		want: &Type{Pkg: pkg0, Name: "uint32Type", Kind: K_UINT32},
	}, {
		ty:   test_type("uint64Type"),
		want: &Type{Pkg: pkg0, Name: "uint64Type", Kind: K_UINT64},
	}, {
		ty:   test_type("UintptrType"),
		want: &Type{Pkg: pkg0, Name: "UintptrType", Kind: K_UINTPTR, IsExported: true},
	}, {
		ty:   test_type("Float32Type"),
		want: &Type{Pkg: pkg0, Name: "Float32Type", Kind: K_FLOAT32, IsExported: true},
	}, {
		ty:   test_type("float64Type"),
		want: &Type{Pkg: pkg0, Name: "float64Type", Kind: K_FLOAT64},
	}, {
		ty:   test_type("Complex64Type"),
		want: &Type{Pkg: pkg0, Name: "Complex64Type", Kind: K_COMPLEX64, IsExported: true},
	}, {
		ty:   test_type("complex128Type"),
		want: &Type{Pkg: pkg0, Name: "complex128Type", Kind: K_COMPLEX128},
	}, {
		ty:   test_type("StringType"),
		want: &Type{Pkg: pkg0, Name: "StringType", Kind: K_STRING, IsExported: true},
	}, {
		ty: test_type("IntSlice"),
		want: &Type{Pkg: pkg0, Name: "IntSlice", Kind: K_SLICE, IsExported: true,
			Elem: &Type{Pkg: pkg0, Name: "IntType", Kind: K_INT, IsExported: true}},
	}, {
		ty: test_type("strSlice"),
		want: &Type{Pkg: pkg0, Name: "strSlice", Kind: K_SLICE,
			Elem: &Type{Kind: K_STRING}},
	}, {
		ty: test_type("byteSlice"),
		want: &Type{Pkg: pkg0, Name: "byteSlice", Kind: K_SLICE,
			Elem: &Type{Kind: K_BYTE, IsByte: true}},
	}, {
		ty: test_type("runeSlice"),
		want: &Type{Pkg: pkg0, Name: "runeSlice", Kind: K_SLICE,
			Elem: &Type{Kind: K_RUNE, IsRune: true}},
	}, {
		ty: test_type("unsafePointerSlice"),
		want: &Type{Pkg: pkg0, Name: "unsafePointerSlice", Kind: K_SLICE,
			Elem: &Type{Kind: K_UNSAFEPOINTER}},
	}, {
		ty: test_type("IntArray"),
		want: &Type{Pkg: pkg0, Name: "IntArray", Kind: K_ARRAY, ArrayLen: 10, IsExported: true,
			Elem: &Type{Pkg: pkg0, Name: "IntType", Kind: K_INT, IsExported: true}},
	}, {
		ty: test_type("strArray"),
		want: &Type{Pkg: pkg0, Name: "strArray", Kind: K_ARRAY, ArrayLen: 4,
			Elem: &Type{Kind: K_STRING}},
	}, {
		ty: test_type("Str2BoolMap"),
		want: &Type{Pkg: pkg0, Name: "Str2BoolMap", Kind: K_MAP, IsExported: true,
			Key:  &Type{Pkg: pkg0, Name: "StringType", Kind: K_STRING, IsExported: true},
			Elem: &Type{Kind: K_BOOL}},
	}, {
		ty: test_type("str2BoolMap"),
		want: &Type{Pkg: pkg0, Name: "str2BoolMap", Kind: K_MAP,
			Key:  &Type{Kind: K_STRING},
			Elem: &Type{Pkg: pkg0, Name: "BoolType", Kind: K_BOOL, IsExported: true}},
	}, {
		ty: test_type("StrPointer"),
		want: &Type{Pkg: pkg0, Name: "StrPointer", Kind: K_PTR, IsExported: true,
			Elem: &Type{Pkg: pkg0, Name: "StringType", Kind: K_STRING, IsExported: true}},
	}, {
		ty: test_type("strPointer"),
		want: &Type{Pkg: pkg0, Name: "strPointer", Kind: K_PTR,
			Elem: &Type{Kind: K_STRING}},
	}, {
		ty: test_type("IfaceType"),
		want: &Type{Pkg: pkg0, Name: "IfaceType", Kind: K_INTERFACE, IsExported: true,
			Methods: []*Method{
				{Name: "Method1", Type: &Type{Kind: K_FUNC}, IsExported: true},
				{Name: "method2", Type: &Type{Kind: K_FUNC}, IsExported: false},
			}},
	}, {
		ty: test_type("ifaceType"),
		want: &Type{Pkg: pkg0, Name: "ifaceType", Kind: K_INTERFACE,
			Methods: []*Method{{
				Name: "Method", Type: &Type{
					Kind: K_FUNC,
					In:   []*Var{{Type: &Type{Kind: K_STRING}}},
					Out:  []*Var{{Type: &Type{Kind: K_BOOL}}},
				},
				IsExported: true,
			}}},
	}, {
		ty: test_type("FuncType"),
		want: &Type{Pkg: pkg0, Name: "FuncType", Kind: K_FUNC, IsExported: true,
			In:  []*Var{{Type: &Type{Kind: K_STRING}}},
			Out: []*Var{{Type: &Type{Kind: K_BOOL}}}},
	}, {
		ty: test_type("funcType"),
		want: &Type{Pkg: pkg0, Name: "funcType", Kind: K_FUNC,
			In: []*Var{
				{Type: &Type{Kind: K_INT}},
				{Type: &Type{Kind: K_STRING}},
				{Type: &Type{Pkg: pkg0, Name: "StringType", Kind: K_STRING, IsExported: true}},
			},
			Out: []*Var{
				{Type: &Type{Kind: K_INT}},
				{Type: &Type{
					Kind: K_INTERFACE,
					Name: "error",
					Methods: []*Method{{
						Name: "Error",
						Type: &Type{Kind: K_FUNC, Out: []*Var{
							{Type: &Type{Kind: K_STRING}},
						}},
						IsExported: true,
					}},
				}},
			}},
	}, {
		ty:   test_type("ChanType"),
		want: &Type{Pkg: pkg0, Name: "ChanType", Kind: K_CHAN, IsExported: true},
	}, {
		ty:   test_type("chanType"),
		want: &Type{Pkg: pkg0, Name: "chanType", Kind: K_CHAN},
	}, {
		ty: test_type("StructType"),
		want: &Type{Pkg: pkg0, Name: "StructType", Kind: K_STRUCT, IsExported: true,
			Fields: []*StructField{
				{Name: "F1", Type: &Type{Kind: K_STRING}, Tag: `is:"some_rule"`, IsExported: true, Var: &types.Var{}},
				{Name: "f2", Type: &Type{Kind: K_BOOL}, Tag: `foo bar baz`, IsExported: false, Var: &types.Var{}},
			}},
	}, {
		ty: test_type("structType"),
		want: &Type{Pkg: pkg0, Name: "structType", Kind: K_STRUCT,
			Fields: []*StructField{
				{Name: "f1", Type: &Type{
					Pkg:        pkg1,
					Name:       "ConstType1",
					Kind:       K_UINT,
					IsImported: true,
					IsExported: true,
				}, Tag: `is:"some_rule"`, Var: &types.Var{}},
				{Name: "f2", Type: &Type{
					Pkg:        pkg1,
					Name:       "Struct",
					Kind:       K_STRUCT,
					IsImported: true,
					IsExported: true,
					Fields: []*StructField{{
						Name:       "Field",
						Type:       &Type{Kind: K_STRING},
						Tag:        `key:"value"`,
						IsExported: true,
						Var:        &types.Var{},
					}},
				}, Tag: `is:"some_rule"`, Var: &types.Var{}},
			}},
	}, {
		ty: test_type("recursiveSlice"),
		want: func() *Type {
			t := &Type{Pkg: pkg0, Name: "recursiveSlice", Kind: K_SLICE}
			t.Elem = t
			return t
		}(),
	}, {
		ty: test_type("RecursiveStruct"),
		want: func() *Type {
			t := &Type{Pkg: pkg0, Name: "RecursiveStruct", Kind: K_STRUCT, IsExported: true}
			t.Fields = []*StructField{{Name: "F1", Type: &Type{Kind: K_PTR, Elem: t}, IsExported: true, Var: &types.Var{}}}
			return t
		}(),
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.ty.String(), func(t *testing.T) {
			a := NewAnalyzer(test_pkg.Type)
			got := a.Analyze(tt.ty)
			if err := compare.Compare(got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAnalyzer_Object(t *testing.T) {
	pkg0 := Pkg{
		Path:  "github.com/frk/valid/cmd/internal/gotype/testdata",
		Name:  "testdata",
		Local: "testdata",
	}

	tests := []struct {
		obj  types.Object
		want *Type
	}{{
		obj:  test_obj("Func1"),
		want: &Type{Pkg: pkg0, Kind: K_FUNC},
	}, {
		obj: test_obj("Func2"),
		want: &Type{Pkg: pkg0, Kind: K_FUNC,
			In:  []*Var{{Type: &Type{Kind: K_STRING}}},
			Out: []*Var{{Type: &Type{Kind: K_BOOL}}}},
	}, {
		obj: test_obj("Func3"),
		want: &Type{Pkg: pkg0, Kind: K_FUNC,
			In: []*Var{
				{Type: &Type{Kind: K_INT}},
				{Type: &Type{Kind: K_STRING}},
				{Type: &Type{Pkg: pkg0, Name: "StringType", Kind: K_STRING, IsExported: true}},
			},
			Out: []*Var{
				{Type: &Type{Kind: K_INT}},
				{Type: &Type{
					Kind: K_INTERFACE,
					Name: "error",
					Methods: []*Method{{
						Name: "Error",
						Type: &Type{Kind: K_FUNC, Out: []*Var{
							{Type: &Type{Kind: K_STRING}},
						}},
						IsExported: true,
					}},
				}},
			}},
	}}

	for _, tt := range tests {
		a := NewAnalyzer(test_pkg.Type)
		got := a.Object(tt.obj)
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
