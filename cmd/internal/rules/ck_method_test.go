package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_methodCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_method_Validator", err: nil,
	}, {
		name: "Test_ERR_METHOD_TYPE_1_Validator",
		err: &Error{C: ERR_METHOD_TYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"isvalid"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r:  &Rule{Name: "isvalid", Spec: GetSpec("isvalid")},
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
