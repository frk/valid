package types

import (
	"go/types"
	"testing"

	"github.com/frk/compare"
)

func TestAnalyzeValidator(t *testing.T) {
	pkg0 := Pkg{
		Path: "github.com/frk/valid/cmd/internal/v2/types/testdata",
		Name: "testdata",
	}

	errorType := &Type{
		Name: "error",
		Kind: INTERFACE,
		MethodSet: []*Method{{
			Name: "Error",
			Type: &Type{Kind: FUNC, Out: []*Var{
				{Type: &Type{Kind: STRING}},
			}},
			IsExported: true,
		}},
	}

	tests := []struct {
		named *types.Named
		want  *Validator
	}{{
		named: test_type("Test1Validator").(*types.Named),
		want:  &Validator{Type: &Type{Pkg: pkg0, Name: "Test1Validator", Kind: STRUCT, IsExported: true}},
	}, {
		named: test_type("Test2Validator").(*types.Named),
		want: &Validator{
			Type: &Type{
				Pkg:        pkg0,
				Name:       "Test2Validator",
				Kind:       STRUCT,
				IsExported: true,
				Fields: []*Field{{
					Pkg:  pkg0,
					Name: "err",
					Type: &Type{Pkg: pkg0, Name: "errorConstructor", Kind: STRUCT, MethodSet: []*Method{{
						Pkg:  pkg0,
						Name: "Error",
						Type: &Type{
							Kind:       FUNC,
							IsVariadic: true,
							In: []*Var{
								{Name: "key", Type: &Type{Kind: STRING}},
								{Name: "val", Type: &Type{Kind: INTERFACE}},
								{Name: "rule", Type: &Type{Kind: STRING}},
								{Name: "args", Type: &Type{Kind: SLICE, Elem: &Type{Kind: INTERFACE}}},
							},
							Out: []*Var{{Type: errorType}},
						},
						IsExported: true,
					}}},
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
				Kind:       STRUCT,
				IsExported: true,
				Fields: []*Field{{
					Pkg:  pkg0,
					Name: "err",
					Type: &Type{Pkg: pkg0, Name: "errorAggregator", Kind: STRUCT, MethodSet: []*Method{{
						Pkg:  pkg0,
						Name: "Error",
						Type: &Type{
							Kind:       FUNC,
							IsVariadic: true,
							In: []*Var{
								{Name: "key", Type: &Type{Kind: STRING}},
								{Name: "val", Type: &Type{Kind: INTERFACE}},
								{Name: "rule", Type: &Type{Kind: STRING}},
								{Name: "args", Type: &Type{Kind: SLICE, Elem: &Type{Kind: INTERFACE}}},
							},
						},
						IsExported: true,
					}, {
						Pkg:  pkg0,
						Name: "Out",
						Type: &Type{
							Kind: FUNC,
							Out:  []*Var{{Type: errorType}},
						},
						IsExported: true,
					}}},
				}},
			},
			ErrorHandlerField: &ErrorHandlerField{Name: "err", IsAggregator: true},
		},
	}, {
		named: test_type("Test5Validator").(*types.Named),
		want: &Validator{
			Type: &Type{
				Pkg:        pkg0,
				Name:       "Test5Validator",
				Kind:       STRUCT,
				IsExported: true,
				MethodSet: []*Method{{
					Pkg:  pkg0,
					Name: "AfterValidate", Type: &Type{Kind: FUNC, Out: []*Var{{Type: errorType}}},
					IsExported: true,
					IsPtr:      true,
				}, {
					Pkg:  pkg0,
					Name: "beforevalidate", Type: &Type{Kind: FUNC, Out: []*Var{{Type: errorType}}},
				}},
			},
			BeforeValidateMethod: &MethodInfo{Name: "beforevalidate"},
			AfterValidateMethod:  &MethodInfo{Name: "AfterValidate"}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for i, tt := range tests {
		t.Run(tt.named.String(), func(t *testing.T) {
			got := AnalyzeValidator(tt.named, &test_ast)
			if err := compare.Compare(got, tt.want); err != nil {
				t.Errorf("#%d: %v", i, err)
			}
		})
	}
}
