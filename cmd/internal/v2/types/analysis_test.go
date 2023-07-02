package types

import (
	"go/types"
	"testing"

	"github.com/frk/compare"
)

func TestAnalyze(t *testing.T) {
	p0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/v2/types/testdata",
		Name: "testdata",
	}
	p1 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/v2/types/testdata/mypkg",
		Name: "mypkg",
	}

	tests := []struct {
		skip bool
		ty   types.Type
		want *Type
	}{{
		ty:   test_type("BoolType"),
		want: &Type{Pkg: p0, Name: "BoolType", Kind: BOOL, IsExported: true},
	}, {
		ty:   test_type("IntType"),
		want: &Type{Pkg: p0, Name: "IntType", Kind: INT, IsExported: true},
	}, {
		ty:   test_type("int8Type"),
		want: &Type{Pkg: p0, Name: "int8Type", Kind: INT8},
	}, {
		ty:   test_type("int16Type"),
		want: &Type{Pkg: p0, Name: "int16Type", Kind: INT16},
	}, {
		ty:   test_type("int32Type"),
		want: &Type{Pkg: p0, Name: "int32Type", Kind: INT32},
	}, {
		ty:   test_type("int64Type"),
		want: &Type{Pkg: p0, Name: "int64Type", Kind: INT64},
	}, {
		ty:   test_type("UintType"),
		want: &Type{Pkg: p0, Name: "UintType", Kind: UINT, IsExported: true},
	}, {
		ty:   test_type("uint8Type"),
		want: &Type{Pkg: p0, Name: "uint8Type", Kind: UINT8},
	}, {
		ty:   test_type("uint16Type"),
		want: &Type{Pkg: p0, Name: "uint16Type", Kind: UINT16},
	}, {
		ty:   test_type("uint32Type"),
		want: &Type{Pkg: p0, Name: "uint32Type", Kind: UINT32},
	}, {
		ty:   test_type("uint64Type"),
		want: &Type{Pkg: p0, Name: "uint64Type", Kind: UINT64},
	}, {
		ty:   test_type("UintptrType"),
		want: &Type{Pkg: p0, Name: "UintptrType", Kind: UINTPTR, IsExported: true},
	}, {
		ty:   test_type("Float32Type"),
		want: &Type{Pkg: p0, Name: "Float32Type", Kind: FLOAT32, IsExported: true},
	}, {
		ty:   test_type("float64Type"),
		want: &Type{Pkg: p0, Name: "float64Type", Kind: FLOAT64},
	}, {
		ty:   test_type("Complex64Type"),
		want: &Type{Pkg: p0, Name: "Complex64Type", Kind: COMPLEX64, IsExported: true},
	}, {
		ty:   test_type("complex128Type"),
		want: &Type{Pkg: p0, Name: "complex128Type", Kind: COMPLEX128},
	}, {
		ty:   test_type("StringType"),
		want: &Type{Pkg: p0, Name: "StringType", Kind: STRING, IsExported: true},
	}, {
		ty: test_type("IntSlice"),
		want: &Type{Pkg: p0, Name: "IntSlice", Kind: SLICE, IsExported: true,
			Elem: &Type{Pkg: p0, Name: "IntType", Kind: INT, IsExported: true}},
	}, {
		ty: test_type("strSlice"),
		want: &Type{Pkg: p0, Name: "strSlice", Kind: SLICE,
			Elem: T.string},
	}, {
		ty: test_type("byteSlice"),
		want: &Type{Pkg: p0, Name: "byteSlice", Kind: SLICE,
			Elem: &Type{Kind: BYTE, IsByte: true}},
	}, {
		ty: test_type("runeSlice"),
		want: &Type{Pkg: p0, Name: "runeSlice", Kind: SLICE,
			Elem: &Type{Kind: RUNE, IsRune: true}},
	}, {
		ty: test_type("unsafePointerSlice"),
		want: &Type{Pkg: p0, Name: "unsafePointerSlice", Kind: SLICE,
			Elem: &Type{Kind: UNSAFEPOINTER}},
	}, {
		ty: test_type("IntArray"),
		want: &Type{Pkg: p0, Name: "IntArray", Kind: ARRAY, ArrayLen: 10, IsExported: true,
			Elem: &Type{Pkg: p0, Name: "IntType", Kind: INT, IsExported: true}},
	}, {
		ty: test_type("strArray"),
		want: &Type{Pkg: p0, Name: "strArray", Kind: ARRAY, ArrayLen: 4,
			Elem: T.string},
	}, {
		ty: test_type("Str2BoolMap"),
		want: &Type{Pkg: p0, Name: "Str2BoolMap", Kind: MAP, IsExported: true,
			Key:  &Type{Pkg: p0, Name: "StringType", Kind: STRING, IsExported: true},
			Elem: T.bool},
	}, {
		ty: test_type("str2BoolMap"),
		want: &Type{Pkg: p0, Name: "str2BoolMap", Kind: MAP,
			Key:  T.string,
			Elem: &Type{Pkg: p0, Name: "BoolType", Kind: BOOL, IsExported: true}},
	}, {
		ty: test_type("StrPointer"),
		want: &Type{Pkg: p0, Name: "StrPointer", Kind: PTR, IsExported: true,
			Elem: &Type{Pkg: p0, Name: "StringType", Kind: STRING, IsExported: true}},
	}, {
		ty: test_type("strPointer"),
		want: &Type{Pkg: p0, Name: "strPointer", Kind: PTR,
			Elem: T.string},
	}, {
		ty: test_type("IfaceType"),
		want: &Type{Pkg: p0, Name: "IfaceType", Kind: INTERFACE, IsExported: true,
			MethodSet: []*Method{
				{Pkg: p0, Name: "Method1", Type: &Type{Kind: FUNC}, IsExported: true},
				{Pkg: p0, Name: "method2", Type: &Type{Kind: FUNC}, IsExported: false},
			}},
	}, {
		ty: test_type("ifaceType"),
		want: &Type{Pkg: p0, Name: "ifaceType", Kind: INTERFACE,
			MethodSet: []*Method{{
				Pkg:  p0,
				Name: "Method", Type: &Type{
					Kind: FUNC,
					In:   []*Var{{Type: T.string}},
					Out:  []*Var{{Type: T.bool}},
				},
				IsExported: true,
			}}},
	}, {
		ty: test_type("FuncType"),
		want: &Type{Pkg: p0, Name: "FuncType", Kind: FUNC, IsExported: true,
			In:  []*Var{{Type: T.string}},
			Out: []*Var{{Type: T.bool}}},
	}, {
		ty: test_type("funcType"),
		want: &Type{Pkg: p0, Name: "funcType", Kind: FUNC,
			In: []*Var{
				{Type: T.int},
				{Type: T.string},
				{Type: &Type{Pkg: p0, Name: "StringType", Kind: STRING, IsExported: true}},
			},
			Out: []*Var{
				{Type: T.int},
				{Type: T.error},
			}},
	}, {
		ty:   test_type("ChanType"),
		want: &Type{Pkg: p0, Name: "ChanType", Kind: CHAN, IsExported: true},
	}, {
		ty:   test_type("chanType"),
		want: &Type{Pkg: p0, Name: "chanType", Kind: CHAN},
	}, {
		ty: test_type("StructType"),
		want: &Type{Pkg: p0, Name: "StructType", Kind: STRUCT, IsExported: true,
			Fields: []*Field{
				{Pkg: p0, Name: "F1", Type: T.string, Tag: `is:"some_rule"`, IsExported: true},
				{Pkg: p0, Name: "f2", Type: T.bool, Tag: `foo bar baz`, IsExported: false},
			}},
	}, {
		ty: test_type("structType"),
		want: &Type{Pkg: p0, Name: "structType", Kind: STRUCT,
			Fields: []*Field{
				{Pkg: p0, Name: "f1", Type: &Type{
					Pkg:        p1,
					Name:       "ConstType1",
					Kind:       UINT,
					IsExported: true,
				}, Tag: `is:"some_rule"`},
				{Pkg: p0, Name: "f2", Type: &Type{
					Pkg:        p1,
					Name:       "Struct",
					Kind:       STRUCT,
					IsExported: true,
					Fields: []*Field{{
						Pkg:        p1,
						Name:       "Field",
						Type:       T.string,
						Tag:        `key:"value"`,
						IsExported: true,
					}},
				}, Tag: `is:"some_rule"`},
			}},
	}, {
		ty: test_type("recursiveSlice"),
		want: func() *Type {
			t := &Type{Pkg: p0, Name: "recursiveSlice", Kind: SLICE, IsRecursive: true}
			t.Elem = t
			return t
		}(),
	}, {
		ty: test_type("RecursiveStruct"),
		want: func() *Type {
			t := &Type{Pkg: p0, Name: "RecursiveStruct", Kind: STRUCT, IsExported: true, IsRecursive: true}
			t.Fields = []*Field{{
				Pkg:        p0,
				Name:       "F1",
				Type:       &Type{Kind: PTR, Elem: t},
				IsExported: true,
			}}
			return t
		}(),
	}, {
		ty: test_type("genMap1"),
		want: &Type{Pkg: p0, Name: "genMap1", Kind: MAP,
			Key:  &Type{Kind: INTERFACE, Name: "comparable"},
			Elem: &Type{Kind: INTERFACE},
			TypeParams: []*TypeParam{{
				Pkg: p0, Name: "K",
				Constraint: &Type{Kind: INTERFACE, Name: "comparable"},
			}, {
				Pkg: p0, Name: "V",
				Constraint: &Type{Kind: INTERFACE},
			}},
		},
	}, {
		ty: test_type("genMap2"),
		want: &Type{Pkg: p0, Name: "genMap2", Kind: MAP,
			Key: &Type{Kind: INTERFACE, Embeddeds: []*Type{
				{Kind: UNION, Terms: []*Term{
					{Tilde: true, Type: T.string},
				}},
			}},
			Elem: &Type{Kind: INTERFACE},
			TypeParams: []*TypeParam{{
				Pkg: p0, Name: "K",
				Constraint: &Type{Kind: INTERFACE, Embeddeds: []*Type{
					{Kind: UNION, Terms: []*Term{
						{Tilde: true, Type: T.string},
					}},
				}},
			}, {
				Pkg: p0, Name: "V",
				Constraint: &Type{Kind: INTERFACE},
			}},
		},
	}, {
		ty: test_type("genStruct1"),
		want: &Type{Pkg: p0, Name: "genStruct1", Kind: STRUCT,
			TypeParams: []*TypeParam{{
				Pkg: p0, Name: "C", Constraint: &Type{
					Kind: INTERFACE,
					MethodSet: []*Method{{
						Pkg:  Pkg{Name: "io", Path: "io"},
						Name: "Read",
						Type: &Type{
							Kind: FUNC,
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
						Kind:       INTERFACE,
						IsExported: true,
						MethodSet: []*Method{{
							Pkg:  Pkg{Name: "io", Path: "io"},
							Name: "Read",
							Type: &Type{
								Kind: FUNC,
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
		if tt.skip {
			continue
		}
		t.Run(tt.ty.String(), func(t *testing.T) {
			got := Analyze(tt.ty, &test_src)
			if err := compare.Compare(got, tt.want); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAnalyzeObject(t *testing.T) {
	p0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/v2/types/testdata",
		Name: "testdata",
	}

	tests := []struct {
		obj  types.Object
		want *Type
	}{{
		obj:  test_obj("Func1"),
		want: &Type{Pkg: p0, Kind: FUNC},
	}, {
		obj: test_obj("Func2"),
		want: &Type{Pkg: p0, Kind: FUNC,
			In:  []*Var{{Type: T.string}},
			Out: []*Var{{Type: T.bool}}},
	}, {
		obj: test_obj("Func3"),
		want: &Type{Pkg: p0, Kind: FUNC,
			In: []*Var{
				{Type: T.int},
				{Type: T.string},
				{Type: &Type{Pkg: p0, Name: "StringType", Kind: STRING, IsExported: true}},
			},
			Out: []*Var{
				{Type: T.int},
				{Type: &Type{
					Kind: INTERFACE,
					Name: "error",
					MethodSet: []*Method{{
						Name: "Error",
						Type: &Type{Kind: FUNC, Out: []*Var{
							{Type: T.string},
						}},
						IsExported: true,
					}},
				}},
			}},
	}}

	for _, tt := range tests {
		got := AnalyzeObject(tt.obj, &test_src)
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
