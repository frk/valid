package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_optionalCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_optional_Validator", err: nil,
	}, {
		name: "Test_ERR_OPTIONAL_CONFLICT_1_Validator",
		err: &Error{C: ERR_OPTIONAL_CONFLICT, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"optional,required"`,
				Type: T.Ptr(T.string),
				Var:  T._var,
			},
			ty: T.Ptr(T.string),
			r:  &Rule{Name: "optional", Spec: GetSpec("optional")},
			r2: &Rule{Name: "required", Spec: GetSpec("required")},
		},
	}, {
		name: "Test_ERR_OPTIONAL_CONFLICT_2_Validator",
		err: &Error{C: ERR_OPTIONAL_CONFLICT, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"required,optional"`,
				Type: T.Ptr(T.string),
				Var:  T._var,
			},
			ty: T.Ptr(T.string),
			r:  &Rule{Name: "optional", Spec: GetSpec("optional")},
			r2: &Rule{Name: "required", Spec: GetSpec("required")},
		},
	}, {
		name: "Test_ERR_OPTIONAL_CONFLICT_3_Validator",
		err: &Error{C: ERR_OPTIONAL_CONFLICT, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"notnil,optional"`,
				Type: T.Ptr(T.string),
				Var:  T._var,
			},
			ty: T.Ptr(T.string),
			r:  &Rule{Name: "optional", Spec: GetSpec("optional")},
			r2: &Rule{Name: "notnil", Spec: GetSpec("notnil")},
		},
	}, {
		name: "Test_ERR_OPTIONAL_CONFLICT_4_Validator",
		err: &Error{C: ERR_OPTIONAL_CONFLICT, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"omitnil,notnil"`,
				Type: T.Ptr(T.string),
				Var:  T._var,
			},
			ty: T.Ptr(T.string),
			r:  &Rule{Name: "omitnil", Spec: GetSpec("omitnil")},
			r2: &Rule{Name: "notnil", Spec: GetSpec("notnil")},
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
