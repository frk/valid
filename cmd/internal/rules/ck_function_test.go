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
				Tag:    `is:"contains:foo"`,
				Object: &gotype.Object{Type: T.int},
				Var:    T._var,
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
				Tag:    `is:"uuid:v6"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
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
				Tag:    `is:"alpha:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "alpha",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("alpha"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "lang", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_2_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"alnum:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "alnum",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("alnum"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "lang", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_3_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"ccy:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "ccy",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
					{},
				},
				Spec: GetSpec("ccy"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "code", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_4_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"decimal:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "decimal",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("decimal"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "locale", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_5_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"hash:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "hash",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("hash"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "algo", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_6_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"ip:5"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "ip",
				Args: []*Arg{
					{Type: ARG_INT, Value: "5"},
				},
				Spec: GetSpec("ip"),
			},
			ra:  &Arg{Type: ARG_INT, Value: "5"},
			fp:  &gotype.Var{Name: "ver", Type: T.int},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_7_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"isbn:12"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "isbn",
				Args: []*Arg{
					{Type: ARG_INT, Value: "12"},
				},
				Spec: GetSpec("isbn"),
			},
			ra:  &Arg{Type: ARG_INT, Value: "12"},
			fp:  &gotype.Var{Name: "ver", Type: T.int},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_8_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"iso639:3"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "iso639",
				Args: []*Arg{
					{Type: ARG_INT, Value: "3"},
				},
				Spec: GetSpec("iso639"),
			},
			ra:  &Arg{Type: ARG_INT, Value: "3"},
			fp:  &gotype.Var{Name: "num", Type: T.int},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_9_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"iso31661a:1"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "iso31661a",
				Args: []*Arg{
					{Type: ARG_INT, Value: "1"},
				},
				Spec: GetSpec("iso31661a"),
			},
			ra:  &Arg{Type: ARG_INT, Value: "1"},
			fp:  &gotype.Var{Name: "num", Type: T.int},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_10_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"mac:7"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "mac",
				Args: []*Arg{
					{Type: ARG_INT, Value: "7"},
				},
				Spec: GetSpec("mac"),
			},
			ra:  &Arg{Type: ARG_INT, Value: "7"},
			fp:  &gotype.Var{Name: "space", Type: T.int},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_11_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"re:[0-9)?"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "re",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "[0-9)?"},
				},
				Spec: GetSpec("re"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "[0-9)?"},
			fp:  &gotype.Var{Name: "expr", Type: T.string},
			fpi: T.iptr(0),
			err: T._err,
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_12_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"uuid:6"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
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
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_13_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"phone:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "phone",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("phone"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "cc", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_14_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"vat:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "vat",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("vat"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "cc", Type: T.string},
			fpi: T.iptr(0),
		},
	}, {
		name: "Test_ERR_FUNCTION_ARGVALUE_15_Validator",
		err: &Error{C: ERR_FUNCTION_ARGVALUE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"zip:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "zip",
				Args: []*Arg{
					{Type: ARG_STRING, Value: "foo"},
				},
				Spec: GetSpec("zip"),
			},
			ra:  &Arg{Type: ARG_STRING, Value: "foo"},
			fp:  &gotype.Var{Name: "cc", Type: T.string},
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
