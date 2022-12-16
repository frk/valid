package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/xtypes"

	"github.com/frk/compare"
)

func TestChecker_requiredCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_required_Validator", err: nil,
	}, {
		name: "Test_ERR_NOTNIL_TYPE_1_Validator",
		err: &Error{C: ERR_NOTNIL_TYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"notnil"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r:  &Rule{Name: "notnil", Spec: GetSpec("notnil")},
		},
	}, {
		name: "Test_ERR_NOTNIL_TYPE_2_Validator",
		err: &Error{C: ERR_NOTNIL_TYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"notnil"`,
				Object: &xtypes.Object{Type: T.bool},
				Var:    T._var,
			},
			ty: T.bool,
			r:  &Rule{Name: "notnil", Spec: GetSpec("notnil")},
		},
	}, {
		name: "Test_ERR_NOTNIL_TYPE_3_Validator",
		err: &Error{C: ERR_NOTNIL_TYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"notnil"`,
				Object: &xtypes.Object{Type: T.float64},
				Var:    T._var,
			},
			ty: T.float64,
			r:  &Rule{Name: "notnil", Spec: GetSpec("notnil")},
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
