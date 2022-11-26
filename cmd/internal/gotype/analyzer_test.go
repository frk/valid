package gotype

import (
	"go/types"
	"testing"

	"github.com/frk/compare"
)

func TestAnalyzer_Analyze(t *testing.T) {
	p0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/gotype/testdata",
		Name: "testdata",
	}
	p1 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/gotype/testdata/mypkg",
		Name: "mypkg",
	}

	tests := []struct {
		ty   types.Type
		want *Type
	}{{
		ty:   test_type("BoolType"),
		want: &Type{Pkg: p0, Name: "BoolType", Kind: K_BOOL, IsExported: true},
	}, {
		ty:   test_type("IntType"),
		want: &Type{Pkg: p0, Name: "IntType", Kind: K_INT, IsExported: true},
	}, {
		ty:   test_type("int8Type"),
		want: &Type{Pkg: p0, Name: "int8Type", Kind: K_INT8},
	}, {
		ty:   test_type("int16Type"),
		want: &Type{Pkg: p0, Name: "int16Type", Kind: K_INT16},
	}, {
		ty:   test_type("int32Type"),
		want: &Type{Pkg: p0, Name: "int32Type", Kind: K_INT32},
	}, {
		ty:   test_type("int64Type"),
		want: &Type{Pkg: p0, Name: "int64Type", Kind: K_INT64},
	}, {
		ty:   test_type("UintType"),
		want: &Type{Pkg: p0, Name: "UintType", Kind: K_UINT, IsExported: true},
	}, {
		ty:   test_type("uint8Type"),
		want: &Type{Pkg: p0, Name: "uint8Type", Kind: K_UINT8},
	}, {
		ty:   test_type("uint16Type"),
		want: &Type{Pkg: p0, Name: "uint16Type", Kind: K_UINT16},
	}, {
		ty:   test_type("uint32Type"),
		want: &Type{Pkg: p0, Name: "uint32Type", Kind: K_UINT32},
	}, {
		ty:   test_type("uint64Type"),
		want: &Type{Pkg: p0, Name: "uint64Type", Kind: K_UINT64},
	}, {
		ty:   test_type("UintptrType"),
		want: &Type{Pkg: p0, Name: "UintptrType", Kind: K_UINTPTR, IsExported: true},
	}, {
		ty:   test_type("Float32Type"),
		want: &Type{Pkg: p0, Name: "Float32Type", Kind: K_FLOAT32, IsExported: true},
	}, {
		ty:   test_type("float64Type"),
		want: &Type{Pkg: p0, Name: "float64Type", Kind: K_FLOAT64},
	}, {
		ty:   test_type("Complex64Type"),
		want: &Type{Pkg: p0, Name: "Complex64Type", Kind: K_COMPLEX64, IsExported: true},
	}, {
		ty:   test_type("complex128Type"),
		want: &Type{Pkg: p0, Name: "complex128Type", Kind: K_COMPLEX128},
	}, {
		ty:   test_type("StringType"),
		want: &Type{Pkg: p0, Name: "StringType", Kind: K_STRING, IsExported: true},
	}, {
		ty: test_type("IntSlice"),
		want: &Type{Pkg: p0, Name: "IntSlice", Kind: K_SLICE, IsExported: true,
			Elem: &Object{Type: &Type{Pkg: p0, Name: "IntType", Kind: K_INT, IsExported: true}}},
	}, {
		ty: test_type("strSlice"),
		want: &Type{Pkg: p0, Name: "strSlice", Kind: K_SLICE,
			Elem: &Object{Type: T.string}},
	}, {
		ty: test_type("byteSlice"),
		want: &Type{Pkg: p0, Name: "byteSlice", Kind: K_SLICE,
			Elem: &Object{Type: &Type{Kind: K_BYTE, IsByte: true}}},
	}, {
		ty: test_type("runeSlice"),
		want: &Type{Pkg: p0, Name: "runeSlice", Kind: K_SLICE,
			Elem: &Object{Type: &Type{Kind: K_RUNE, IsRune: true}}},
	}, {
		ty: test_type("unsafePointerSlice"),
		want: &Type{Pkg: p0, Name: "unsafePointerSlice", Kind: K_SLICE,
			Elem: &Object{Type: &Type{Kind: K_UNSAFEPOINTER}}},
	}, {
		ty: test_type("IntArray"),
		want: &Type{Pkg: p0, Name: "IntArray", Kind: K_ARRAY, ArrayLen: 10, IsExported: true,
			Elem: &Object{Type: &Type{Pkg: p0, Name: "IntType", Kind: K_INT, IsExported: true}}},
	}, {
		ty: test_type("strArray"),
		want: &Type{Pkg: p0, Name: "strArray", Kind: K_ARRAY, ArrayLen: 4,
			Elem: &Object{Type: T.string}},
	}, {
		ty: test_type("Str2BoolMap"),
		want: &Type{Pkg: p0, Name: "Str2BoolMap", Kind: K_MAP, IsExported: true,
			Key:  &Object{Type: &Type{Pkg: p0, Name: "StringType", Kind: K_STRING, IsExported: true}},
			Elem: &Object{Type: T.bool}},
	}, {
		ty: test_type("str2BoolMap"),
		want: &Type{Pkg: p0, Name: "str2BoolMap", Kind: K_MAP,
			Key:  &Object{Type: T.string},
			Elem: &Object{Type: &Type{Pkg: p0, Name: "BoolType", Kind: K_BOOL, IsExported: true}}},
	}, {
		ty: test_type("StrPointer"),
		want: &Type{Pkg: p0, Name: "StrPointer", Kind: K_PTR, IsExported: true,
			Elem: &Object{Type: &Type{Pkg: p0, Name: "StringType", Kind: K_STRING, IsExported: true}}},
	}, {
		ty: test_type("strPointer"),
		want: &Type{Pkg: p0, Name: "strPointer", Kind: K_PTR,
			Elem: &Object{Type: T.string}},
	}, {
		ty: test_type("IfaceType"),
		want: &Type{Pkg: p0, Name: "IfaceType", Kind: K_INTERFACE, IsExported: true,
			Methods: []*Method{
				{Pkg: p0, Name: "Method1", Type: &Type{Kind: K_FUNC}, IsExported: true},
				{Pkg: p0, Name: "method2", Type: &Type{Kind: K_FUNC}, IsExported: false},
			}},
	}, {
		ty: test_type("ifaceType"),
		want: &Type{Pkg: p0, Name: "ifaceType", Kind: K_INTERFACE,
			Methods: []*Method{{
				Pkg:  p0,
				Name: "Method", Type: &Type{
					Kind: K_FUNC,
					In:   []*Var{{Type: T.string}},
					Out:  []*Var{{Type: T.bool}},
				},
				IsExported: true,
			}}},
	}, {
		ty: test_type("FuncType"),
		want: &Type{Pkg: p0, Name: "FuncType", Kind: K_FUNC, IsExported: true,
			In:  []*Var{{Type: T.string}},
			Out: []*Var{{Type: T.bool}}},
	}, {
		ty: test_type("funcType"),
		want: &Type{Pkg: p0, Name: "funcType", Kind: K_FUNC,
			In: []*Var{
				{Type: T.int},
				{Type: T.string},
				{Type: &Type{Pkg: p0, Name: "StringType", Kind: K_STRING, IsExported: true}},
			},
			Out: []*Var{
				{Type: T.int},
				{Type: T.error},
			}},
	}, {
		ty:   test_type("ChanType"),
		want: &Type{Pkg: p0, Name: "ChanType", Kind: K_CHAN, IsExported: true},
	}, {
		ty:   test_type("chanType"),
		want: &Type{Pkg: p0, Name: "chanType", Kind: K_CHAN},
	}, {
		ty: test_type("StructType"),
		want: &Type{Pkg: p0, Name: "StructType", Kind: K_STRUCT, IsExported: true,
			Fields: []*StructField{
				{Pkg: p0, Name: "F1", Object: &Object{Type: T.string}, Tag: `is:"some_rule"`, IsExported: true, Var: &types.Var{}},
				{Pkg: p0, Name: "f2", Object: &Object{Type: T.bool}, Tag: `foo bar baz`, IsExported: false, Var: &types.Var{}},
			}},
	}, {
		ty: test_type("structType"),
		want: &Type{Pkg: p0, Name: "structType", Kind: K_STRUCT,
			Fields: []*StructField{
				{Pkg: p0, Name: "f1", Object: &Object{Type: &Type{
					Pkg:        p1,
					Name:       "ConstType1",
					Kind:       K_UINT,
					IsExported: true,
				}}, Tag: `is:"some_rule"`, Var: &types.Var{}},
				{Pkg: p0, Name: "f2", Object: &Object{Type: &Type{
					Pkg:        p1,
					Name:       "Struct",
					Kind:       K_STRUCT,
					IsExported: true,
					Fields: []*StructField{{
						Pkg:        p1,
						Name:       "Field",
						Object:     &Object{Type: T.string},
						Tag:        `key:"value"`,
						IsExported: true,
						Var:        &types.Var{},
					}},
				}}, Tag: `is:"some_rule"`, Var: &types.Var{}},
			}},
	}, {
		ty: test_type("recursiveSlice"),
		want: func() *Type {
			t := &Type{Pkg: p0, Name: "recursiveSlice", Kind: K_SLICE}
			t.Elem = &Object{Type: t}
			return t
		}(),
	}, {
		ty: test_type("RecursiveStruct"),
		want: func() *Type {
			t := &Type{Pkg: p0, Name: "RecursiveStruct", Kind: K_STRUCT, IsExported: true}
			t.Fields = []*StructField{{Pkg: p0, Name: "F1",
				Object:     &Object{Type: &Type{Kind: K_PTR, Elem: &Object{Type: t}}},
				IsExported: true, Var: &types.Var{}}}
			return t
		}(),
	}, {
		ty: test_type("genMap1"),
		want: &Type{Pkg: p0, Name: "genMap1", Kind: K_MAP,
			Key:  &Object{Type: &Type{Kind: K_INTERFACE, Name: "comparable"}},
			Elem: &Object{Type: &Type{Kind: K_INTERFACE}},
			TypeParams: []*TypeParam{{
				Pkg: p0, Name: "K",
				Constraint: &Type{Kind: K_INTERFACE, Name: "comparable"},
			}, {
				Pkg: p0, Name: "V",
				Constraint: &Type{Kind: K_INTERFACE},
			}},
		},
	}, {
		ty: test_type("genMap2"),
		want: &Type{Pkg: p0, Name: "genMap2", Kind: K_MAP,
			Key: &Object{Type: &Type{Kind: K_INTERFACE, Embeddeds: []*Type{
				{Kind: K_UNION, Terms: []*Term{
					{Tilde: true, Type: T.string},
				}},
			}}},
			Elem: &Object{Type: &Type{Kind: K_INTERFACE}},
			TypeParams: []*TypeParam{{
				Pkg: p0, Name: "K",
				Constraint: &Type{Kind: K_INTERFACE, Embeddeds: []*Type{
					{Kind: K_UNION, Terms: []*Term{
						{Tilde: true, Type: T.string},
					}},
				}},
			}, {
				Pkg: p0, Name: "V",
				Constraint: &Type{Kind: K_INTERFACE},
			}},
		},
	}, {
		ty: test_type("genStruct1"),
		want: &Type{Pkg: p0, Name: "genStruct1", Kind: K_STRUCT,
			TypeParams: []*TypeParam{{
				Pkg: p0, Name: "C", Constraint: &Type{
					Kind: K_INTERFACE,
					Methods: []*Method{{
						Pkg:  Pkg{Name: "io", Path: "io"},
						Name: "Read",
						Type: &Type{
							Kind: K_FUNC,
							In:   []*Var{{Name: "p", Type: T.bytes}},
							Out: []*Var{
								{Name: "n", Type: T.int},
								{Name: "err", Type: T.error},
							},
						},
						IsExported: true,
					}},
					Embeddeds: []*Type{{
						Pkg:        Pkg{Path: "io", Name: "io"},
						Name:       "Reader",
						Kind:       K_INTERFACE,
						IsExported: true,
						Methods: []*Method{{
							Pkg:  Pkg{Name: "io", Path: "io"},
							Name: "Read",
							Type: &Type{
								Kind: K_FUNC,
								In:   []*Var{{Name: "p", Type: T.bytes}},
								Out: []*Var{
									{Name: "n", Type: T.int},
									{Name: "err", Type: T.error},
								},
							},
							IsExported: true,
						}},
					}},
				},
			}},
		},
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
	p0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/gotype/testdata",
		Name: "testdata",
	}

	tests := []struct {
		obj  types.Object
		want *Type
	}{{
		obj:  test_obj("Func1"),
		want: &Type{Pkg: p0, Kind: K_FUNC},
	}, {
		obj: test_obj("Func2"),
		want: &Type{Pkg: p0, Kind: K_FUNC,
			In:  []*Var{{Type: T.string}},
			Out: []*Var{{Type: T.bool}}},
	}, {
		obj: test_obj("Func3"),
		want: &Type{Pkg: p0, Kind: K_FUNC,
			In: []*Var{
				{Type: T.int},
				{Type: T.string},
				{Type: &Type{Pkg: p0, Name: "StringType", Kind: K_STRING, IsExported: true}},
			},
			Out: []*Var{
				{Type: T.int},
				{Type: &Type{
					Kind: K_INTERFACE,
					Name: "error",
					Methods: []*Method{{
						Name: "Error",
						Type: &Type{Kind: K_FUNC, Out: []*Var{
							{Type: T.string},
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
