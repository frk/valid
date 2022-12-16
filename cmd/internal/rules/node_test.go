package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/xtypes"

	"github.com/frk/compare"
)

func TestChecker_makeNode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		echo bool
	}{{
		name: "Test_ERR_RULE_ELEM_1_is_Validator",
		err: &Error{C: ERR_RULE_ELEM, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"[]email"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty:  T.string,
			tag: &Tag{Elem: &Tag{Rules: []*Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_ERR_RULE_ELEM_2_is_Validator",
		err: &Error{C: ERR_RULE_ELEM, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"[][]email"`,
				Object: &xtypes.Object{Type: T.Slice(T.string)},
				Var:    T._var,
			},
			ty:  T.string,
			tag: &Tag{Elem: &Tag{Rules: []*Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_ERR_RULE_ELEM_3_is_Validator",
		err: &Error{C: ERR_RULE_ELEM, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"[][]email,[]email"`,
				Object: &xtypes.Object{Type: T.Slice(T.Map(T.string, T.string))},
				Var:    T._var,
			},
			ty: T.string,
			tag: &Tag{
				Rules: []*Rule{{Name: "email", Spec: GetSpec("email")}},
				Elem:  &Tag{Rules: []*Rule{{Name: "email"}}},
			},
		},
	}, {
		name: "Test_ERR_RULE_KEY_1_is_Validator",
		err: &Error{C: ERR_RULE_KEY, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"[email]"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty:  T.string,
			tag: &Tag{Key: &Tag{Rules: []*Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_ERR_RULE_KEY_2_is_Validator",
		err: &Error{C: ERR_RULE_KEY, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"[email]"`,
				Object: &xtypes.Object{Type: T.Slice(T.string)},
				Var:    T._var,
			},
			ty:  T.Slice(T.string),
			tag: &Tag{Key: &Tag{Rules: []*Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_ERR_RULE_KEY_3_is_Validator",
		err: &Error{C: ERR_RULE_KEY, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"[][email]"`,
				Object: &xtypes.Object{Type: T.Map(T.string, T.Slice(T.string))},
				Var:    T._var,
			},
			ty:  T.Slice(T.string),
			tag: &Tag{Key: &Tag{Rules: []*Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_ERR_RULE_ELEM_1_pre_Validator",
		err: &Error{C: ERR_RULE_ELEM, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"[]trim"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty:  T.string,
			tag: &Tag{Elem: &Tag{Rules: []*Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_ERR_RULE_ELEM_2_pre_Validator",
		err: &Error{C: ERR_RULE_ELEM, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"[][]trim"`,
				Object: &xtypes.Object{Type: T.Slice(T.string)},
				Var:    T._var,
			},
			ty:  T.string,
			tag: &Tag{Elem: &Tag{Rules: []*Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_ERR_RULE_ELEM_3_pre_Validator",
		err: &Error{C: ERR_RULE_ELEM, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"[][]trim,[]trim"`,
				Object: &xtypes.Object{Type: T.Slice(T.Map(T.string, T.string))},
				Var:    T._var,
			},
			ty: T.string,
			tag: &Tag{
				Rules: []*Rule{{Name: "trim", Spec: GetSpec("pre:trim")}},
				Elem:  &Tag{Rules: []*Rule{{Name: "trim"}}},
			},
		},
	}, {
		name: "Test_ERR_RULE_KEY_1_pre_Validator",
		err: &Error{C: ERR_RULE_KEY, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"[trim]"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty:  T.string,
			tag: &Tag{Key: &Tag{Rules: []*Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_ERR_RULE_KEY_2_pre_Validator",
		err: &Error{C: ERR_RULE_KEY, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"[trim]"`,
				Object: &xtypes.Object{Type: T.Slice(T.string)},
				Var:    T._var,
			},
			ty:  T.Slice(T.string),
			tag: &Tag{Key: &Tag{Rules: []*Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_ERR_RULE_KEY_3_pre_Validator",
		err: &Error{C: ERR_RULE_KEY, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"[][trim]"`,
				Object: &xtypes.Object{Type: T.Map(T.string, T.Slice(T.string))},
				Var:    T._var,
			},
			ty:  T.Slice(T.string),
			tag: &Tag{Key: &Tag{Rules: []*Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_ERR_RULE_UNDEFINED_1_is_Validator",
		err: &Error{C: ERR_RULE_UNDEFINED, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"email,trim"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			tag: &Tag{Rules: []*Rule{
				{Name: "email", Spec: GetSpec("email")},
				{Name: "trim"},
			}},
			r: &Rule{Name: "trim"},
		},
	}, {
		name: "Test_ERR_RULE_UNDEFINED_2_is_Validator",
		err: &Error{C: ERR_RULE_UNDEFINED, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"email,trim"`,
				Object: &xtypes.Object{Type: T.Ptr(T.string)},
				Var:    T._var,
			},
			ty: T.Ptr(T.string),
			tag: &Tag{Rules: []*Rule{
				{Name: "email", Spec: GetSpec("email")},
				{Name: "trim"},
			}},
			r: &Rule{Name: "trim"},
		},
	}, {
		name: "Test_ERR_RULE_UNDEFINED_1_pre_Validator",
		err: &Error{C: ERR_RULE_UNDEFINED, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"trim,email"`,
				Object: &xtypes.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			tag: &Tag{Rules: []*Rule{
				{Name: "trim", Spec: GetSpec("pre:trim")},
				{Name: "email"},
			}},
			r: &Rule{Name: "email"},
		},
	}, {
		name: "Test_ERR_RULE_UNDEFINED_2_pre_Validator",
		err: &Error{C: ERR_RULE_UNDEFINED, a: T._ast, sfv: T._var,
			sf: &xtypes.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `pre:"trim,email"`,
				Object: &xtypes.Object{Type: T.Ptr(T.string)},
				Var:    T._var,
			},
			ty: T.string,
			tag: &Tag{Rules: []*Rule{
				{Name: "trim", Spec: GetSpec("pre:trim")},
				{Name: "email"},
			}},
			r: &Rule{Name: "email"},
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

			got := (*ttError)(err.(*Error))
			want := (*ttError)(tt.err.(*Error))
			if e := compare.Compare(got, want); e != nil {
				t.Errorf("Error: %v", e)
			}

			if tt.echo {
				fmt.Println(err)
			}
		})
	}
}
