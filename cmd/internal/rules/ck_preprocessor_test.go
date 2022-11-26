package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func init() {
	_custom["pre:p1"] = &Spec{
		Name: "p1", Kind: PREPROC,
		FName: "MyPreproc1",
		FType: &gotype.Type{
			Pkg: gotype.Pkg{
				Name: "mypkg",
				Path: "github.com/frk/valid/cmd/internal/rules/testdata/mypkg",
			},
			In:   []*gotype.Var{{Name: "v", Type: &gotype.Type{Kind: gotype.K_FLOAT64}}},
			Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
			Kind: gotype.K_FUNC,
		},
	}
	_custom["pre:p2"] = &Spec{
		Name: "p2", Kind: PREPROC,
		FName: "MyPreproc2",
		FType: &gotype.Type{
			Pkg: gotype.Pkg{
				Name: "mypkg",
				Path: "github.com/frk/valid/cmd/internal/rules/testdata/mypkg",
			},
			In: []*gotype.Var{
				{Name: "v", Type: &gotype.Type{Kind: gotype.K_STRING}},
				{Name: "opt", Type: &gotype.Type{Kind: gotype.K_UINT}},
			},
			Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_FLOAT64}}},
			Kind: gotype.K_FUNC,
		},
	}
	_custom["pre:p4"] = &Spec{
		Name: "p4", Kind: PREPROC,
		FName: "MyPreproc4",
		FType: &gotype.Type{
			Pkg: gotype.Pkg{
				Name: "mypkg",
				Path: "github.com/frk/valid/cmd/internal/rules/testdata/mypkg",
			},
			In: []*gotype.Var{
				{Name: "v", Type: &gotype.Type{Kind: gotype.K_STRING}},
				{Name: "opt", Type: &gotype.Type{Kind: gotype.K_UINT}},
			},
			Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
			Kind: gotype.K_FUNC,
		},
		ArgMin: 1,
		ArgMax: 1,
	}

	// non-PREPROC, for a test that ensures the code checks the Kind==PREPROC
	_custom["pre:p0"] = &Spec{
		Name: "p0", Kind: FUNCTION,
		FName: "MyPreproc0",
		FType: &gotype.Type{
			Pkg: gotype.Pkg{
				Name: "mypkg",
				Path: "github.com/frk/valid/cmd/internal/rules/testdata/mypkg",
			},
			In:   []*gotype.Var{{Name: "v", Type: &gotype.Type{Kind: gotype.K_STRING}}},
			Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
			Kind: gotype.K_FUNC,
		},
		ArgMin: 1,
		ArgMax: 1,
	}
}

func TestChecker_preprocessorCheck(t *testing.T) {

	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_preproc_Validator", err: nil,
	}, {
		name: "Test_ERR_PREPROC_INTYPE_1_Validator",
		err: &Error{C: ERR_PREPROC_INTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"p1"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "p1",
				Spec: GetSpec("pre:p1"),
			},
		},
	}, {
		name: "Test_ERR_PREPROC_OUTTYPE_1_Validator",
		err: &Error{C: ERR_PREPROC_OUTTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"p2"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "p2",
				Spec: GetSpec("pre:p2"),
			},
		},
	}, {
		name: "Test_ERR_PREPROC_ARGTYPE_1_Validator",
		err: &Error{C: ERR_PREPROC_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"p4:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "p4",
				Args: []*Arg{{Type: ARG_STRING, Value: "foo"}},
				Spec: GetSpec("pre:p4"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "opt", Type: T.uint},
			fpi: T.iptr(0),
		},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	fkCfg := &config.FieldKeyConfig{
		Join:      config.Bool{Value: true, IsSet: true},
		Separator: config.String{Value: ".", IsSet: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := testMatch(t, tt.name)

			info := new(Info)
			checker := NewChecker(&test_ast, test_pkg.Pkg(), fkCfg, info)
			err := checker.Check(match)

			got := _ttError(err)
			want := _ttError(tt.err)
			if e := compare.Compare(got, want); e != nil {
				t.Errorf("Error: %v", e)
			}

			if tt.show && err != nil {
				fmt.Println(err)
			}
		})
	}
}
