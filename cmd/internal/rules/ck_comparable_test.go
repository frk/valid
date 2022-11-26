package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_comparableCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_comparable_Validator", err: nil,
	}, {
		name: "Test_ERR_ARG_BADCMP_1_Validator",
		err: &Error{C: ERR_ARG_BADCMP, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"eq:42:64:foo:-22"`,
				Object: &gotype.Object{Type: T.int},
				Var:    T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "eq",
				Args: []*Arg{
					{Type: ARG_INT, Value: "42"},
					{Type: ARG_INT, Value: "64"},
					{Type: ARG_STRING, Value: "foo"},
					{Type: ARG_INT, Value: "-22"},
				},
				Spec: GetSpec("eq"),
			},
			ra: &Arg{Type: ARG_STRING, Value: "foo"},
		},
	}, {
		name: "Test_ERR_ARG_BADCMP_2_Validator",
		err: &Error{C: ERR_ARG_BADCMP, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"eq:123:&S.G"`,
				Object: &gotype.Object{Type: T.int},
				Var:    T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "eq",
				Args: []*Arg{
					{Type: ARG_INT, Value: "123"},
					{Type: ARG_FIELD_ABS, Value: "S.G"},
				},
				Spec: GetSpec("eq"),
			},
			ra:  &Arg{Type: ARG_FIELD_ABS, Value: "S.G"},
			raf: T._sf(),
		},
	}, {
		name: "Test_ERR_ARG_BADCMP_3_Validator",
		err: &Error{C: ERR_ARG_BADCMP, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"eq:0.03"`,
				Object: &gotype.Object{Type: T.int},
				Var:    T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "eq",
				Args: []*Arg{{Type: ARG_FLOAT, Value: "0.03"}},
				Spec: GetSpec("eq"),
			},
			ra: &Arg{Type: ARG_FLOAT, Value: "0.03"},
		},
	}, {
		name: "Test_ERR_ARG_BADCMP_4_Validator",
		err: &Error{C: ERR_ARG_BADCMP, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"ne:-345"`,
				Object: &gotype.Object{Type: T.uint},
				Var:    T._var,
			},
			ty: T.uint,
			r: &Rule{
				Name: "ne",
				Args: []*Arg{{Type: ARG_INT, Value: "-345"}},
				Spec: GetSpec("ne"),
			},
			ra: &Arg{Type: ARG_INT, Value: "-345"},
		},
	}, {
		name: "Test_ERR_ARG_BADCMP_5_Validator",
		err: &Error{C: ERR_ARG_BADCMP, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"eq:1"`,
				Object: &gotype.Object{Type: T.bool},
				Var:    T._var,
			},
			ty: T.bool,
			r: &Rule{
				Name: "eq",
				Args: []*Arg{{Type: ARG_INT, Value: "1"}},
				Spec: GetSpec("eq"),
			},
			ra: &Arg{Type: ARG_INT, Value: "1"},
		},
	}, {
		name: "Test_ERR_ARG_BADCMP_6_Validator",
		err: &Error{C: ERR_ARG_BADCMP, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"eq:true"`,
				Object: &gotype.Object{Type: T.int},
				Var:    T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "eq",
				Args: []*Arg{{Type: ARG_BOOL, Value: "true"}},
				Spec: GetSpec("eq"),
			},
			ra: &Arg{Type: ARG_BOOL, Value: "true"},
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
