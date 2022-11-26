package gotype

import (
	"go/types"
	"testing"

	"github.com/frk/compare"
)

func TestAnalyzer_Validator(t *testing.T) {
	pkg0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/gotype/testdata",
		Name: "testdata",
	}

	errorType := &Type{
		Name: "error",
		Kind: K_INTERFACE,
		Methods: []*Method{{
			Name: "Error",
			Type: &Type{Kind: K_FUNC, Out: []*Var{
				{Type: &Type{Kind: K_STRING}},
			}},
			IsExported: true,
		}},
	}

	tests := []struct {
		named *types.Named
		want  *Validator
	}{{
		named: test_type("Test1Validator").(*types.Named),
		want:  &Validator{Type: &Type{Pkg: pkg0, Name: "Test1Validator", Kind: K_STRUCT, IsExported: true}},
	}, {
		named: test_type("Test2Validator").(*types.Named),
		want: &Validator{
			Type: &Type{
				Pkg:        pkg0,
				Name:       "Test2Validator",
				Kind:       K_STRUCT,
				IsExported: true,
				Fields: []*StructField{{
					Pkg:  pkg0,
					Name: "err",
					Object: &Object{Type: &Type{Pkg: pkg0, Name: "errorConstructor", Kind: K_STRUCT, Methods: []*Method{{
						Pkg:  pkg0,
						Name: "Error",
						Type: &Type{
							Kind:       K_FUNC,
							IsVariadic: true,
							In: []*Var{
								{Name: "key", Type: &Type{Kind: K_STRING}},
								{Name: "val", Type: &Type{Kind: K_INTERFACE}},
								{Name: "rule", Type: &Type{Kind: K_STRING}},
								{Name: "args", Type: &Type{Kind: K_SLICE, Elem: &Object{Type: &Type{Kind: K_INTERFACE}}}},
							},
							Out: []*Var{{Type: errorType}},
						},
						IsExported: true,
					}}}},
					Var: &types.Var{},
				}},
			},
			ErrorHandlerField: &ErrorHandlerField{Name: "err"},
		},
	}, {
		named: test_type("Test3Validator").(*types.Named),
		want: &Validator{
			Type: &Type{
				Pkg:        pkg0,
				Name:       "Test3Validator",
				Kind:       K_STRUCT,
				IsExported: true,
				Fields: []*StructField{{
					Pkg:  pkg0,
					Name: "err",
					Object: &Object{Type: &Type{Pkg: pkg0, Name: "errorAggregator", Kind: K_STRUCT, Methods: []*Method{{
						Pkg:  pkg0,
						Name: "Error",
						Type: &Type{
							Kind:       K_FUNC,
							IsVariadic: true,
							In: []*Var{
								{Name: "key", Type: &Type{Kind: K_STRING}},
								{Name: "val", Type: &Type{Kind: K_INTERFACE}},
								{Name: "rule", Type: &Type{Kind: K_STRING}},
								{Name: "args", Type: &Type{Kind: K_SLICE, Elem: &Object{Type: &Type{Kind: K_INTERFACE}}}},
							},
						},
						IsExported: true,
					}, {
						Pkg:  pkg0,
						Name: "Out",
						Type: &Type{
							Kind: K_FUNC,
							Out:  []*Var{{Type: errorType}},
						},
						IsExported: true,
					}}}},
					Var: &types.Var{},
				}},
			},
			ErrorHandlerField: &ErrorHandlerField{Name: "err", IsAggregator: true},
		},
	}, {
		// (*gotype.Validator).Type.Methods[0].Type.Out[0].Type.Methods[0].Pkg.Name
		named: test_type("Test5Validator").(*types.Named),
		want: &Validator{
			Type: &Type{
				Pkg:        pkg0,
				Name:       "Test5Validator",
				Kind:       K_STRUCT,
				IsExported: true,
				Methods: []*Method{{
					Pkg:  pkg0,
					Name: "beforevalidate", Type: &Type{Kind: K_FUNC, Out: []*Var{{Type: errorType}}},
				}, {
					Pkg:  pkg0,
					Name: "AfterValidate", Type: &Type{Kind: K_FUNC, Out: []*Var{{Type: errorType}}},
					IsExported: true,
				}},
			},
			BeforeValidateMethod: &MethodInfo{Name: "beforevalidate"},
			AfterValidateMethod:  &MethodInfo{Name: "AfterValidate"}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for i, tt := range tests {
		t.Run(tt.named.String(), func(t *testing.T) {
			a := NewAnalyzer(test_pkg.Type)
			got := a.Validator(tt.named)
			if err := compare.Compare(got, tt.want); err != nil {
				t.Errorf("#%d: %v", i, err)
			}
		})
	}
}
