package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_rangeCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_range_Validator", err: nil,
	}, {
		name: "Test_between_Validator", err: nil,
	}, {
		name: "Test_ERR_RANGE_TYPE_1_Validator",
		err: &Error{C: ERR_RANGE_TYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:10:20"`,
				Type: T.string,
				Var:  T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_INT, Value: "10"},
					{Type: ARG_INT, Value: "20"},
				},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_TYPE_2_Validator",
		err: &Error{C: ERR_RANGE_TYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:10:20"`,
				Type: T.Slice(T.int),
				Var:  T._var,
			},
			ty: T.Slice(T.int),
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_INT, Value: "10"},
					{Type: ARG_INT, Value: "20"},
				},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_NOARG_1_Validator",
		err: &Error{C: ERR_RANGE_NOARG, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng::"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{{}, {}},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_NOARG_2_Validator",
		err: &Error{C: ERR_RANGE_NOARG, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:10:"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}, {}},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_NOARG_3_Validator",
		err: &Error{C: ERR_RANGE_NOARG, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng::20"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{{}, {Type: ARG_INT, Value: "20"}},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_ARGTYPE_1_Validator",
		err: &Error{C: ERR_RANGE_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:foo:20"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
					{Type: ARG_INT, Value: "20"},
				},
				Spec: GetSpec("rng"),
			},
			ra: &Arg{Type: ARG_STRING, Value: "foo"},
		},
	}, {
		name: "Test_ERR_RANGE_ARGTYPE_2_Validator",
		err: &Error{C: ERR_RANGE_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:3.14:20"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_FLOAT, Value: "3.14"},
					{Type: ARG_INT, Value: "20"},
				},
				Spec: GetSpec("rng"),
			},
			ra: &Arg{Type: ARG_FLOAT, Value: "3.14"},
		},
	}, {
		name: "Test_ERR_RANGE_ARGTYPE_3_Validator",
		err: &Error{C: ERR_RANGE_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:-10:20"`,
				Type: T.uint16,
				Var:  T._var,
			},
			ty: T.uint16,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_INT, Value: "-10"},
					{Type: ARG_INT, Value: "20"},
				},
				Spec: GetSpec("rng"),
			},
			ra: &Arg{Type: ARG_INT, Value: "-10"},
		},
	}, {
		name: "Test_ERR_RANGE_ARGTYPE_4_Validator",
		err: &Error{C: ERR_RANGE_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:&S.F:20"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_FIELD, Value: "S.F"},
					{Type: ARG_INT, Value: "20"},
				},
				Spec: GetSpec("rng"),
			},
			ra:  &Arg{Type: ARG_FIELD, Value: "S.F"},
			raf: T._sf(),
		},
	}, {
		name: "Test_ERR_RANGE_BOUNDS_1_Validator",
		err: &Error{C: ERR_RANGE_BOUNDS, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:20:10"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_INT, Value: "20"},
					{Type: ARG_INT, Value: "10"},
				},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_BOUNDS_2_Validator",
		err: &Error{C: ERR_RANGE_BOUNDS, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:10:10"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_INT, Value: "10"},
					{Type: ARG_INT, Value: "10"},
				},
				Spec: GetSpec("rng"),
			},
		},
	}, {
		name: "Test_ERR_RANGE_BOUNDS_3_Validator",
		err: &Error{C: ERR_RANGE_BOUNDS, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:  `is:"rng:10:-20"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "rng",
				Args: []*Arg{
					{Type: ARG_INT, Value: "10"},
					{Type: ARG_INT, Value: "-20"},
				},
				Spec: GetSpec("rng"),
			},
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
