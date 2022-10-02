package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_functionCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_function_Validator", err: nil,
	}, {
		name: "Test_ERR_FUNCTION_INTYPE_1_Validator",
		err: &Error{C: ERR_FUNCTION_INTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"contains:foo"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "contains",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("contains"),
			},
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGTYPE_1_Validator",
		err: &Error{C: ERR_FUNCTION_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"uuid:v6"`,
				Type: T.string,
				Var:  T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "uuid",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "v6"},
				},
				Spec: GetSpec("uuid"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "v6"},
			fp:  &gotype.Var{Name: "ver", Type: T.int},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_1_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"uuid:6"`,
				Type: T.string,
				Var:  T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "uuid",
				Args: []*Arg{
					{Type: ARG_INT, Value: "6"},
				},
				Spec: GetSpec("uuid"),
			},
			ra:  &Arg{Type: ARG_INT, Value: "6"},
			fp:  &gotype.Var{Name: "ver", Type: T.int},
			fpi: T.iptr(0),
		},
	}}

	cfg := loadConfig("testdata/configs/test_custom_rules.yaml")
	if err := initCustomSpecs(cfg, &test_ast); err != nil {
		t.Fatalf("loadConfig(testdata/configs/test_custom_rules.yaml) failed: %v", err)
	}

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
