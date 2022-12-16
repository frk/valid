package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/xtypes"

	"github.com/frk/compare"
)

func TestChecker_orderedCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_ordered_Validator", err: nil,
	}, {
		name: "Test_ERR_ORDERED_TYPE_1_Validator",
		err: &Error{C: ERR_ORDERED_TYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"min:8"`,
				Object: &xtypes.Object{Type: T.Slice(T.string)},
				Var:    T._var,
			},
			ty: T.Slice(T.string),
			r: &Rule{
				Name: "min",
				Args: []*Arg{{Type: ARG_INT, Value: "8"}},
				Spec: GetSpec("min"),
			},
		},
	}, {
		name: "Test_ERR_ORDERED_TYPE_2_Validator",
		err: &Error{C: ERR_ORDERED_TYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"gt:8"`,
				Object: &xtypes.Object{Type: T.Slice(T.int)},
				Var:    T._var,
			},
			ty: T.Slice(T.int),
			r: &Rule{
				Name: "gt",
				Args: []*Arg{{Type: ARG_INT, Value: "8"}},
				Spec: GetSpec("gt"),
			},
		},
	}, {
		name: "Test_ERR_ORDERED_ARGTYPE_1_Validator",
		err: &Error{C: ERR_ORDERED_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"gte:0.8"`,
				Object: &xtypes.Object{Type: T.int},
				Var:    T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "gte",
				Args: []*Arg{{Type: ARG_FLOAT, Value: "0.8"}},
				Spec: GetSpec("gte"),
			},
			ra: &Arg{Type: ARG_FLOAT, Value: "0.8"},
		},
	}, {
		name: "Test_ERR_ORDERED_ARGTYPE_2_Validator",
		err: &Error{C: ERR_ORDERED_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"gte:foo"`,
				Object: &xtypes.Object{Type: T.float64},
				Var:    T._var,
			},
			ty: T.float64,
			r: &Rule{
				Name: "gte",
				Args: []*Arg{{Type: ARG_STRING, Value: "foo"}},
				Spec: GetSpec("gte"),
			},
			ra: &Arg{Type: ARG_STRING, Value: "foo"},
		},
	}, {
		name: "Test_ERR_ORDERED_ARGTYPE_3_Validator",
		err: &Error{C: ERR_ORDERED_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"lte:&S.F"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "lte",
				Args: []*Arg{{Type: ARG_FIELD_ABS, Value: "S.F"}},
				Spec: GetSpec("lte"),
			},
			ra:  &Arg{Type: ARG_FIELD_ABS, Value: "S.F"},
			raf: T._sf(),
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
