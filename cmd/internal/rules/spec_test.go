package rules

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/xtypes"
)

func Test_initCustomSpecs(t *testing.T) {
	tests := []struct {
		name string
		want error
		show bool
	}{{
		name: "func does not exist",
		want: &Error{C: ERR_CONFIG_FUNCSEARCH, a: T._ast, c: T._cfg, err: &search.Error{}},
	}, {
		name: "package does not exist",
		want: &Error{C: ERR_CONFIG_FUNCSEARCH, a: T._ast, c: T._cfg, err: &search.Error{}},
	}, {
		name: "bad config in comment",
		want: &Error{C: ERR_CONFIG_INVALID, a: T._ast, c: T._cfg, ft: T._func, err: T._err},
	}, {
		name: "missing config",
		want: &Error{C: ERR_CONFIG_MISSING, a: T._ast, c: T._cfg, ft: T._func},
	}, {
		name: "no rule name",
		want: &Error{C: ERR_CONFIG_NONAME, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{},
		},
	}, {
		name: "no rule name 2",
		want: &Error{C: ERR_CONFIG_NONAME, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "required"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "notnil"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "omitnil"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "optional"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "noguard"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "isvalid"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "-isvalid"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: ERR_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "enum"},
		},
	}, {
		name: "has invalid signature",
		want: &Error{C: ERR_CONFIG_FUNCTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "myrule"},
		},
	}, {
		name: "has invalid signature #2",
		want: &Error{C: ERR_CONFIG_FUNCTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "myrule"},
		},
	}, {
		name: "incompatible number of args",
		want: &Error{C: ERR_CONFIG_ARGNUM, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name: "test", Args: []config.RuleArgConfig{
				{Default: &config.Scalar{Type: config.INT, Value: "123"}},
				{Default: &config.Scalar{Type: config.BOOL, Value: "true"}},
				{Options: []config.RuleArgOption{
					{Value: config.Scalar{Type: config.STRING, Value: "foo"}},
					{Value: config.Scalar{Type: config.STRING, Value: "bar"}},
				}},
			},
		}},
	}, {
		name: "incompatible number of args",
		want: &Error{C: ERR_CONFIG_ARGNUM, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name: "test", Args: []config.RuleArgConfig{
				{Default: &config.Scalar{Type: config.INT, Value: "123"}},
			},
		}},
	}, {
		name: "incompatible bounds",
		want: &Error{C: ERR_CONFIG_ARGBOUNDS, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name:   "test",
			ArgMin: T.uptr(8),
		}},
	}, {
		name: "incompatible bounds",
		want: &Error{C: ERR_CONFIG_ARGBOUNDS, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name:   "test",
			ArgMax: T.iptr(7),
		}},
	}, {
		name: "incompatible bounds",
		want: &Error{C: ERR_CONFIG_ARGBOUNDS, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name:   "test",
			ArgMin: T.uptr(8),
			ArgMax: T.iptr(7),
		}},
	}, {
		name: "pre: bad signature for PREPROC rule",
		want: &Error{C: ERR_CONFIG_PREFUNCTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "pre:rule"},
		},
	}, {
		name: "pre: illegal use of 'join_op' for PREPROC rule",
		want: &Error{C: ERR_CONFIG_PREPROCJOIN, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "pre:rule", JoinOp: config.JOIN_AND},
		},
	}, {
		name: "pre: illegal use of 'err' for PREPROC rule",
		want: &Error{C: ERR_CONFIG_PREPROCERROR, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{
				Name: "pre:rule",
				Error: config.RuleErrorConfig{
					Text: "foo bar",
				},
			},
		},
	}, {
		name: "pre: incompatible number of args",
		want: &Error{C: ERR_CONFIG_ARGNUM, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name: "pre:foo_bar", Args: []config.RuleArgConfig{
				{Default: &config.Scalar{Type: config.INT, Value: "123"}},
			},
		}},
	}, {
		name: "pre: incompatible arg types",
		want: &Error{C: ERR_CONFIG_ARGTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{
				Name: "pre:foo_bar", Args: []config.RuleArgConfig{
					{Default: &config.Scalar{Type: config.STRING, Value: "foo"}},
				},
			},
			rca:  &Arg{Type: ARG_STRING, Value: "foo"},
			rcai: T.iptr(0),
			rcak: T.sptr(""),
			fp:   &xtypes.Var{Name: "opt", Type: T.bool},
			fpi:  T.iptr(0),
		},
	}}

	cfg := loadConfig("testdata/configs/bad_custom_rules.yaml")
	if len(cfg.Rules) == 0 {
		t.Error("loadConfig: failed to load rules...")
	}

	compare := compare.Config{ObserveFieldTag: "cmp"}

	for i, c := range cfg.Rules {
		tt := tests[i]

		cfg := cfg
		cfg.Rules = []config.RuleConfig{c}

		t.Run(tt.name, func(t *testing.T) {
			err := initCustomSpecs(cfg, &test_ast)

			got := (*ttError)(err.(*Error))
			want := (*ttError)(tt.want.(*Error))
			if e := compare.Compare(got, want); e != nil {
				t.Error(e)
			}
			if tt.show && err != nil {
				fmt.Println(err)
			}
		})
	}
}

func loadConfig(file string) (c config.Config) {
	file, err := filepath.Abs(file)
	if err != nil {
		log.Fatal(err)
	}
	if err := config.DecodeFile(file, &c); err != nil {
		log.Fatal(err)
	}
	return c
}
