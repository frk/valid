package rules

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/search"

	"github.com/frk/compare"
)

var test_ast search.AST
var test_pkg search.Package

func TestMain(m *testing.M) {
	pkgs, err := search.Search(
		"testdata/",
		true,
		nil,
		nil,
		&test_ast,
	)
	if err != nil {
		log.Fatal(err)
	}
	test_pkg = *(pkgs[0])
	//fmt.Println(pkgs)

	T.loc = gotype.MustGetType("time", "Location", &test_ast)

	if err := loadBuiltinSpecs(&test_ast); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestCheckerCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_ERR_CONTEXT_NOFIELD_1_Validator",
		err: &Error{C: ERR_CONTEXT_NOFIELD, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Name: "F", IsExported: true,
				Tag:  `is:"eq:42:@foo"`,
				Type: T.int,
				Var:  T._var,
			},
			ty: T.int,
			r: &Rule{
				Name: "eq",
				Args: []*Arg{{Type: ARG_INT, Value: "42"}},
				Spec: GetSpec("eq"),
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

func testMatch(t *testing.T, name string) *search.Match {
	for _, file := range test_pkg.Files {
		for _, match := range file.Matches {
			if match.Named.Obj().Name() == name {
				return match
			}
		}
	}

	t.Fatal(name, " not found")
	return nil
}

type test_values struct {
	_ast  *search.AST
	_cfg  *config.Config
	_var  *types.Var
	_func *types.Func
	_err  error

	string  *gotype.Type
	int     *gotype.Type
	int32   *gotype.Type
	int64   *gotype.Type
	uint    *gotype.Type
	uint8   *gotype.Type
	uint16  *gotype.Type
	uint64  *gotype.Type
	bool    *gotype.Type
	float64 *gotype.Type
	rune    *gotype.Type
	byte    *gotype.Type
	loc     *gotype.Type

	pkg gotype.Pkg
}

func (test_values) _sf() *gotype.StructField {
	return &gotype.StructField{}
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

func (test_values) Slice(e *gotype.Type) *gotype.Type {
	return &gotype.Type{Kind: gotype.K_SLICE, Elem: e}
}

func (test_values) Array(n int64, e *gotype.Type) *gotype.Type {
	return &gotype.Type{Kind: gotype.K_ARRAY, ArrayLen: n, Elem: e}
}

func (test_values) Ptr(e *gotype.Type) *gotype.Type {
	return &gotype.Type{Kind: gotype.K_PTR, Elem: e}
}

func (test_values) Map(k, e *gotype.Type) *gotype.Type {
	return &gotype.Type{Kind: gotype.K_MAP, Key: k, Elem: e}
}

var T = test_values{
	_ast:  &search.AST{},
	_cfg:  &config.Config{},
	_var:  &types.Var{},
	_func: &types.Func{},
	_err:  fmt.Errorf(""),

	string:  &gotype.Type{Kind: gotype.K_STRING},
	int:     &gotype.Type{Kind: gotype.K_INT},
	int32:   &gotype.Type{Kind: gotype.K_INT32},
	int64:   &gotype.Type{Kind: gotype.K_INT64},
	uint:    &gotype.Type{Kind: gotype.K_UINT},
	uint8:   &gotype.Type{Kind: gotype.K_UINT8},
	uint16:  &gotype.Type{Kind: gotype.K_UINT16},
	uint64:  &gotype.Type{Kind: gotype.K_UINT64},
	float64: &gotype.Type{Kind: gotype.K_FLOAT64},
	bool:    &gotype.Type{Kind: gotype.K_BOOL},
	rune:    &gotype.Type{Kind: gotype.K_INT32, IsRune: true},
	byte:    &gotype.Type{Kind: gotype.K_UINT8, IsByte: true},
	loc:     nil,

	pkg: gotype.Pkg{
		Path: "github.com/frk/valid/cmd/internal/rules/testdata",
		Name: "testdata",
	},
}