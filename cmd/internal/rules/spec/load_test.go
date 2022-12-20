package spec

import (
	"fmt"
	stdtypes "go/types"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"
)

var test_ast search.AST

func TestMain(m *testing.M) {
	_, err := search.Search(
		"testdata/",
		true,
		nil,
		nil,
		&test_ast,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := LoadIncludedSpecs(&test_ast); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestLoadCustomSpecs(t *testing.T) {
	tests := []struct {
		name string
		want error
		show bool
	}{{
		name: "func does not exist",
		want: &Error{C: E_CONFIG_FUNCSEARCH, a: T._ast, c: T._cfg, err: &search.Error{}},
	}, {
		name: "package does not exist",
		want: &Error{C: E_CONFIG_FUNCSEARCH, a: T._ast, c: T._cfg, err: &search.Error{}},
	}, {
		name: "bad config in comment",
		want: &Error{C: E_CONFIG_INVALID, a: T._ast, c: T._cfg, ft: T._func, err: T._err},
	}, {
		name: "missing config",
		want: &Error{C: E_CONFIG_MISSING, a: T._ast, c: T._cfg, ft: T._func},
	}, {
		name: "no rule name",
		want: &Error{C: E_CONFIG_NONAME, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{},
		},
	}, {
		name: "no rule name 2",
		want: &Error{C: E_CONFIG_NONAME, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "required"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "notnil"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "omitnil"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "optional"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "noguard"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "isvalid"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "-isvalid"},
		},
	}, {
		name: "rule name is reserved",
		want: &Error{C: E_CONFIG_RESERVED, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "enum"},
		},
	}, {
		name: "has invalid signature",
		want: &Error{C: E_CONFIG_FUNCTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "myrule"},
		},
	}, {
		name: "has invalid signature #2",
		want: &Error{C: E_CONFIG_FUNCTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "myrule"},
		},
	}, {
		name: "incompatible number of args",
		want: &Error{C: E_CONFIG_ARGNUM, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
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
		want: &Error{C: E_CONFIG_ARGNUM, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name: "test", Args: []config.RuleArgConfig{
				{Default: &config.Scalar{Type: config.INT, Value: "123"}},
			},
		}},
	}, {
		name: "incompatible bounds",
		want: &Error{C: E_CONFIG_ARGBOUNDS, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name:   "test",
			ArgMin: T.uptr(8),
		}},
	}, {
		name: "incompatible bounds",
		want: &Error{C: E_CONFIG_ARGBOUNDS, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name:   "test",
			ArgMax: T.iptr(7),
		}},
	}, {
		name: "incompatible bounds",
		want: &Error{C: E_CONFIG_ARGBOUNDS, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name:   "test",
			ArgMin: T.uptr(8),
			ArgMax: T.iptr(7),
		}},
	}, {
		name: "pre: bad signature for PREPROC rule",
		want: &Error{C: E_CONFIG_PREFUNCTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "pre:rule"},
		},
	}, {
		name: "pre: illegal use of 'join_op' for PREPROC rule",
		want: &Error{C: E_CONFIG_PREPROCJOIN, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{Name: "pre:rule", JoinOp: config.JOIN_AND},
		},
	}, {
		name: "pre: illegal use of 'err' for PREPROC rule",
		want: &Error{C: E_CONFIG_PREPROCERROR, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{
				Name: "pre:rule",
				Error: config.RuleErrorConfig{
					Text: "foo bar",
				},
			},
		},
	}, {
		name: "pre: incompatible number of args",
		want: &Error{C: E_CONFIG_ARGNUM, a: T._ast, c: T._cfg, ft: T._func, rs: &config.RuleSpec{
			Name: "pre:foo_bar", Args: []config.RuleArgConfig{
				{Default: &config.Scalar{Type: config.INT, Value: "123"}},
			},
		}},
	}, {
		name: "pre: incompatible arg types",
		want: &Error{C: E_CONFIG_ARGTYPE, a: T._ast, c: T._cfg, ft: T._func,
			rs: &config.RuleSpec{
				Name: "pre:foo_bar", Args: []config.RuleArgConfig{
					{Default: &config.Scalar{Type: config.STRING, Value: "foo"}},
				},
			},
			rca:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
			rcai: T.iptr(0),
			rcak: T.sptr(""),
			fp:   &types.Var{Name: "opt", Type: T.bool},
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
			err := LoadCustomSpecs(cfg, &test_ast)

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

type test_values struct {
	_ast  *search.AST
	_cfg  *config.Config
	_var  *stdtypes.Var
	_func *stdtypes.Func
	_err  error

	string  *types.Type
	int     *types.Type
	int32   *types.Type
	int64   *types.Type
	uint    *types.Type
	uint8   *types.Type
	uint16  *types.Type
	uint64  *types.Type
	bool    *types.Type
	float64 *types.Type
	rune    *types.Type
	byte    *types.Type

	pkg types.Pkg
}

func (test_values) _sf() *types.StructField {
	return &types.StructField{}
}

func (test_values) sptr(s string) *string {
	return &s
}

func (test_values) uptr(u uint) *uint {
	return &u
}

func (test_values) iptr(i int) *int {
	return &i
}

func (test_values) Slice(e *types.Type) *types.Type {
	return &types.Type{Kind: types.SLICE, Elem: &types.Obj{Type: e}}
}

func (test_values) Array(n int64, e *types.Type) *types.Type {
	return &types.Type{Kind: types.ARRAY, ArrayLen: n, Elem: &types.Obj{Type: e}}
}

func (test_values) Ptr(e *types.Type) *types.Type {
	return &types.Type{Kind: types.PTR, Elem: &types.Obj{Type: e}}
}

func (test_values) Map(k, e *types.Type) *types.Type {
	return &types.Type{Kind: types.MAP, Key: &types.Obj{Type: k}, Elem: &types.Obj{Type: e}}
}

var T = test_values{
	_ast:  &search.AST{},
	_cfg:  &config.Config{},
	_var:  &stdtypes.Var{},
	_func: &stdtypes.Func{},
	_err:  fmt.Errorf(""),

	string:  &types.Type{Kind: types.STRING},
	int:     &types.Type{Kind: types.INT},
	int32:   &types.Type{Kind: types.INT32},
	int64:   &types.Type{Kind: types.INT64},
	uint:    &types.Type{Kind: types.UINT},
	uint8:   &types.Type{Kind: types.UINT8},
	uint16:  &types.Type{Kind: types.UINT16},
	uint64:  &types.Type{Kind: types.UINT64},
	float64: &types.Type{Kind: types.FLOAT64},
	bool:    &types.Type{Kind: types.BOOL},
	rune:    &types.Type{Kind: types.INT32, IsRune: true},
	byte:    &types.Type{Kind: types.UINT8, IsByte: true},

	pkg: types.Pkg{
		Path: "github.com/frk/valid/cmd/internal/rules/spec/testdata",
		Name: "testdata",
	},
}
