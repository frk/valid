package checker

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"

	"github.com/frk/compare"
)

func Test_walk(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_E_RULE_ELEM_1_is_Validator",
		err: &Error{C: E_RULE_ELEM,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"[]email"`,
				Obj: &types.Obj{Type: T.string},
			},
			ty: T.string,
			tag: &rules.TagNode{
				Elem: &rules.TagNode{
					Rules: []*rules.Rule{{Name: "email"}},
				},
			},
		},
	}, {
		name: "Test_E_RULE_ELEM_2_is_Validator",
		err: &Error{C: E_RULE_ELEM,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"[][]email"`,
				Obj: &types.Obj{Type: T.Slice(T.string)},
			},
			ty: T.string,
			tag: &rules.TagNode{
				Elem: &rules.TagNode{
					Rules: []*rules.Rule{{Name: "email"}},
				},
			},
		},
	}, {
		name: "Test_E_RULE_ELEM_3_is_Validator",
		err: &Error{C: E_RULE_ELEM,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"[][]email,[]email"`,
				Obj: &types.Obj{Type: T.Slice(T.Map(T.string, T.string))},
			},
			ty: T.string,
			tag: &rules.TagNode{
				Rules: []*rules.Rule{{Name: "email"}},
				Elem:  &rules.TagNode{Rules: []*rules.Rule{{Name: "email"}}},
			},
		},
	}, {
		name: "Test_E_RULE_KEY_1_is_Validator",
		err: &Error{C: E_RULE_KEY,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"[email]"`,
				Obj: &types.Obj{Type: T.string},
			},
			ty:  T.string,
			tag: &rules.TagNode{Key: &rules.TagNode{Rules: []*rules.Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_E_RULE_KEY_2_is_Validator",
		err: &Error{C: E_RULE_KEY,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"[email]"`,
				Obj: &types.Obj{Type: T.Slice(T.string)},
			},
			ty:  T.Slice(T.string),
			tag: &rules.TagNode{Key: &rules.TagNode{Rules: []*rules.Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_E_RULE_KEY_3_is_Validator",
		err: &Error{C: E_RULE_KEY,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"[][email]"`,
				Obj: &types.Obj{Type: T.Map(T.string, T.Slice(T.string))},
			},
			ty:  T.Slice(T.string),
			tag: &rules.TagNode{Key: &rules.TagNode{Rules: []*rules.Rule{{Name: "email"}}}},
		},
	}, {
		name: "Test_E_RULE_ELEM_1_pre_Validator",
		err: &Error{C: E_RULE_ELEM,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"[]trim"`,
				Obj: &types.Obj{Type: T.string},
			},
			ty: T.string,
			tag: &rules.TagNode{
				Elem: &rules.TagNode{
					Rules: []*rules.Rule{{Name: "trim"}},
				},
			},
		},
	}, {
		name: "Test_E_RULE_ELEM_2_pre_Validator",
		err: &Error{C: E_RULE_ELEM,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"[][]trim"`,
				Obj: &types.Obj{Type: T.Slice(T.string)},
			},
			ty: T.string,
			tag: &rules.TagNode{
				Elem: &rules.TagNode{
					Rules: []*rules.Rule{{Name: "trim"}},
				},
			},
		},
	}, {
		name: "Test_E_RULE_ELEM_3_pre_Validator",
		err: &Error{C: E_RULE_ELEM,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"[][]trim,[]trim"`,
				Obj: &types.Obj{Type: T.Slice(T.Map(T.string, T.string))},
			},
			ty: T.string,
			tag: &rules.TagNode{
				Rules: []*rules.Rule{{Name: "trim"}},
				Elem: &rules.TagNode{
					Rules: []*rules.Rule{{Name: "trim"}},
				},
			},
		},
	}, {
		name: "Test_E_RULE_KEY_1_pre_Validator",
		err: &Error{C: E_RULE_KEY,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"[trim]"`,
				Obj: &types.Obj{Type: T.string},
			},
			ty:  T.string,
			tag: &rules.TagNode{Key: &rules.TagNode{Rules: []*rules.Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_E_RULE_KEY_2_pre_Validator",
		err: &Error{C: E_RULE_KEY,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"[trim]"`,
				Obj: &types.Obj{Type: T.Slice(T.string)},
			},
			ty:  T.Slice(T.string),
			tag: &rules.TagNode{Key: &rules.TagNode{Rules: []*rules.Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_E_RULE_KEY_3_pre_Validator",
		err: &Error{C: E_RULE_KEY,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"[][trim]"`,
				Obj: &types.Obj{Type: T.Map(T.string, T.Slice(T.string))},
			},
			ty:  T.Slice(T.string),
			tag: &rules.TagNode{Key: &rules.TagNode{Rules: []*rules.Rule{{Name: "trim"}}}},
		},
	}, {
		name: "Test_E_RULE_UNDEFINED_1_is_Validator",
		err: &Error{C: E_RULE_UNDEFINED,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"email,trim"`,
				Obj: &types.Obj{Type: T.string},
			},
			tag: &rules.TagNode{Rules: []*rules.Rule{
				{Name: "email", Spec: specs.Get("email")},
				{Name: "trim"},
			}},
			r: &rules.Rule{Name: "trim"},
		},
	}, {
		name: "Test_E_RULE_UNDEFINED_2_is_Validator",
		err: &Error{C: E_RULE_UNDEFINED,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `is:"email,trim"`,
				Obj: &types.Obj{Type: T.Ptr(T.string)},
			},
			tag: &rules.TagNode{Rules: []*rules.Rule{
				{Name: "email", Spec: specs.Get("email")},
				{Name: "trim"},
			}},
			r: &rules.Rule{Name: "trim"},
		},
	}, {
		name: "Test_E_RULE_UNDEFINED_1_pre_Validator",
		err: &Error{C: E_RULE_UNDEFINED,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"trim,email"`,
				Obj: &types.Obj{Type: T.string},
			},
			tag: &rules.TagNode{Rules: []*rules.Rule{
				{Name: "trim", Spec: specs.Get("pre:trim")},
				{Name: "email"},
			}},
			r: &rules.Rule{Name: "email"},
		},
	}, {
		name: "Test_E_RULE_UNDEFINED_2_pre_Validator",
		err: &Error{C: E_RULE_UNDEFINED,
			sf: &types.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag: `pre:"trim,email"`,
				Obj: &types.Obj{Type: T.Ptr(T.string)},
			},
			tag: &rules.TagNode{Rules: []*rules.Rule{
				{Name: "trim", Spec: specs.Get("pre:trim")},
				{Name: "email"},
			}},
			r: &rules.Rule{Name: "email"},
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

			cfg := Config{AST: &test_ast, FieldKey: fkCfg}
			err := Check(cfg, match, &Info{})

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
