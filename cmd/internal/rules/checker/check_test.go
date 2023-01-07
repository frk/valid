package checker

import (
	"fmt"
	stdtypes "go/types"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"

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

	T.loc = types.MustGetType("time", "Location", &test_ast)

	if err := specs.Load(config.Config{}, &test_ast); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestCheckerCheck(t *testing.T) {
	mypkg := types.Pkg{
		Path: "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
		Name: "mypkg",
	}
	timepkg := types.Pkg{
		Path: "time",
		Name: "time",
	}
	mytime := &types.Type{
		Pkg:  mypkg,
		Kind: types.STRUCT,
		Name: "Time",
		Fields: []*types.StructField{{
			Pkg:  timepkg,
			Name: "wall",
			Obj:  &types.Obj{Type: T.uint64},
		}, {
			Pkg:  timepkg,
			Name: "ext",
			Obj:  &types.Obj{Type: T.int64},
		}, {
			Pkg:  timepkg,
			Name: "loc",
			Obj: &types.Obj{Type: &types.Type{
				Kind: types.PTR,
				Elem: &types.Obj{Type: T.loc},
			}},
		}},
		IsExported: true,
	}

	type testCase struct {
		name string
		err  error
		show bool
	}

	groups := []struct {
		name  string
		tests []testCase
	}{{
		name: "general",
		tests: []testCase{{
			name: "Test_checker_Validator",
		}, {
			name: "Test_E_FIELD_UNKNOWN_1_Validator",
			err: &Error{
				C: E_FIELD_UNKNOWN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gt:&num"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "gt",
					Args: []*rules.Arg{
						{Type: rules.ARG_FIELD_ABS, Value: "num"},
					},
					Spec: specs.Get("gt"),
				},
				ra: &rules.Arg{Type: rules.ARG_FIELD_ABS, Value: "num"},
			},
		}, {
			name: "Test_E_FIELD_UNKNOWN_2_Validator",
			err: &Error{
				C: E_FIELD_UNKNOWN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p4:&num"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "p4",
					Args: []*rules.Arg{{Type: rules.ARG_FIELD_ABS, Value: "num"}},
					Spec: specs.Get("pre:p4"),
				},
				ra: &rules.Arg{Type: rules.ARG_FIELD_ABS, Value: "num"},
			},
		}, {
			name: "Test_E_FIELD_UNKNOWN_3_Validator",
			err: &Error{
				C: E_FIELD_UNKNOWN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gt:.num"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "gt",
					Args: []*rules.Arg{
						{Type: rules.ARG_FIELD_REL, Value: "num"},
					},
					Spec: specs.Get("gt"),
				},
				ra: &rules.Arg{Type: rules.ARG_FIELD_REL, Value: "num"},
			},
		}, {
			name: "Test_E_FIELD_UNKNOWN_4_Validator",
			err: &Error{
				C: E_FIELD_UNKNOWN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p4:.num"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "p4",
					Args: []*rules.Arg{{Type: rules.ARG_FIELD_REL, Value: "num"}},
					Spec: specs.Get("pre:p4"),
				},
				ra: &rules.Arg{Type: rules.ARG_FIELD_REL, Value: "num"},
			},
		}, {
			name: "Test_E_RULE_ARGMIN_1_Validator",
			err: &Error{
				C: E_RULE_ARGMIN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gt"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r:  &rules.Rule{Name: "gt", Spec: specs.Get("gt")},
			},
		}, {
			name: "Test_E_RULE_ARGMIN_2_Validator",
			err: &Error{
				C: E_RULE_ARGMIN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p4"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r:  &rules.Rule{Name: "p4", Spec: specs.Get("pre:p4")},
			},
		}, {
			name: "Test_E_RULE_ARGMAX_1_Validator",
			err: &Error{
				C: E_RULE_ARGMAX,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gt:4:5"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "gt",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "4"},
						{Type: rules.ARG_INT, Value: "5"},
					},
					Spec: specs.Get("gt"),
				},
			},
		}, {
			name: "Test_E_RULE_ARGMAX_2_Validator",
			err: &Error{
				C: E_RULE_ARGMAX,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p4:1:2:3"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "p4",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "1"},
						{Type: rules.ARG_INT, Value: "2"},
						{Type: rules.ARG_INT, Value: "3"},
					},
					Spec: specs.Get("pre:p4"),
				},
			},
		}, {
			name: "Test_E_PREPROC_INVALID_1_Validator",
			err: &Error{
				C: E_PREPROC_INVALID,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p0"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r:  &rules.Rule{Name: "p0", Spec: specs.Get("pre:p0")},
			},
		}},
	}, {
		name: "required",
		tests: []testCase{{
			name: "Test_required_Validator",
		}, {
			name: "Test_E_NOTNIL_TYPE_1_Validator",
			err: &Error{C: E_NOTNIL_TYPE,
				sf: &types.StructField{
					Pkg:        T.pkg,
					Name:       "F",
					IsExported: true,
					Tag:        `is:"notnil"`,
					Obj:        &types.Obj{Type: T.string},
				},
				ty: T.string,
				r:  &rules.Rule{Name: "notnil", Spec: specs.Get("notnil")},
			},
		}, {
			name: "Test_E_NOTNIL_TYPE_2_Validator",
			err: &Error{C: E_NOTNIL_TYPE,
				sf: &types.StructField{
					Pkg:        T.pkg,
					Name:       "F",
					IsExported: true,
					Tag:        `is:"notnil"`,
					Obj:        &types.Obj{Type: T.bool},
				},
				ty: T.bool,
				r:  &rules.Rule{Name: "notnil", Spec: specs.Get("notnil")},
			},
		}, {
			name: "Test_E_NOTNIL_TYPE_3_Validator",
			err: &Error{C: E_NOTNIL_TYPE,
				sf: &types.StructField{
					Pkg:        T.pkg,
					Name:       "F",
					IsExported: true,
					Tag:        `is:"notnil"`,
					Obj:        &types.Obj{Type: T.float64},
				},
				ty: T.float64,
				r:  &rules.Rule{Name: "notnil", Spec: specs.Get("notnil")},
			},
		}},
	}, {
		name: "comparable",
		tests: []testCase{{
			name: "Test_comparable_Validator",
		}, {
			name: "Test_E_ARG_BADCMP_1_Validator",
			err: &Error{C: E_ARG_BADCMP,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"eq:42:64:foo:-22"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "eq",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "42"},
						{Type: rules.ARG_INT, Value: "64"},
						{Type: rules.ARG_STRING, Value: "foo"},
						{Type: rules.ARG_INT, Value: "-22"},
					},
					Spec: specs.Get("eq"),
				},
				ra: &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
			},
		}, {
			name: "Test_E_ARG_BADCMP_2_Validator",
			err: &Error{C: E_ARG_BADCMP,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"eq:123:&S.G"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "eq",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "123"},
						{Type: rules.ARG_FIELD_ABS, Value: "S.G"},
					},
					Spec: specs.Get("eq"),
				},
				ra:  &rules.Arg{Type: rules.ARG_FIELD_ABS, Value: "S.G"},
				raf: T._sf(),
			},
		}, {
			name: "Test_E_ARG_BADCMP_3_Validator",
			err: &Error{C: E_ARG_BADCMP,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"eq:0.03"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "eq",
					Args: []*rules.Arg{{Type: rules.ARG_FLOAT, Value: "0.03"}},
					Spec: specs.Get("eq"),
				},
				ra: &rules.Arg{Type: rules.ARG_FLOAT, Value: "0.03"},
			},
		}, {
			name: "Test_E_ARG_BADCMP_4_Validator",
			err: &Error{C: E_ARG_BADCMP,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"ne:-345"`,
					Obj: &types.Obj{Type: T.uint},
				},
				ty: T.uint,
				r: &rules.Rule{
					Name: "ne",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "-345"}},
					Spec: specs.Get("ne"),
				},
				ra: &rules.Arg{Type: rules.ARG_INT, Value: "-345"},
			},
		}, {
			name: "Test_E_ARG_BADCMP_5_Validator",
			err: &Error{C: E_ARG_BADCMP,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"eq:1"`,
					Obj: &types.Obj{Type: T.bool},
				},
				ty: T.bool,
				r: &rules.Rule{
					Name: "eq",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "1"}},
					Spec: specs.Get("eq"),
				},
				ra: &rules.Arg{Type: rules.ARG_INT, Value: "1"},
			},
		}, {
			name: "Test_E_ARG_BADCMP_6_Validator",
			err: &Error{C: E_ARG_BADCMP,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"eq:true"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "eq",
					Args: []*rules.Arg{{Type: rules.ARG_BOOL, Value: "true"}},
					Spec: specs.Get("eq"),
				},
				ra: &rules.Arg{Type: rules.ARG_BOOL, Value: "true"},
			},
		}},
	}, {
		name: "ordered",
		tests: []testCase{{
			name: "Test_ordered_Validator", err: nil,
		}, {
			name: "Test_E_ORDERED_TYPE_1_Validator",
			err: &Error{C: E_ORDERED_TYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"min:8"`,
					Obj: &types.Obj{Type: T.Slice(T.string)},
				},
				ty: T.Slice(T.string),
				r: &rules.Rule{
					Name: "min",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "8"}},
					Spec: specs.Get("min"),
				},
			},
		}, {
			name: "Test_E_ORDERED_TYPE_2_Validator",
			err: &Error{C: E_ORDERED_TYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gt:8"`,
					Obj: &types.Obj{Type: T.Slice(T.int)},
				},
				ty: T.Slice(T.int),
				r: &rules.Rule{
					Name: "gt",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "8"}},
					Spec: specs.Get("gt"),
				},
			},
		}, {
			name: "Test_E_ORDERED_ARGTYPE_1_Validator",
			err: &Error{C: E_ORDERED_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gte:0.8"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "gte",
					Args: []*rules.Arg{{Type: rules.ARG_FLOAT, Value: "0.8"}},
					Spec: specs.Get("gte"),
				},
				ra: &rules.Arg{Type: rules.ARG_FLOAT, Value: "0.8"},
			},
		}, {
			name: "Test_E_ORDERED_ARGTYPE_2_Validator",
			err: &Error{C: E_ORDERED_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"gte:foo"`,
					Obj: &types.Obj{Type: T.float64},
				},
				ty: T.float64,
				r: &rules.Rule{
					Name: "gte",
					Args: []*rules.Arg{{Type: rules.ARG_STRING, Value: "foo"}},
					Spec: specs.Get("gte"),
				},
				ra: &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
			},
		}, {
			name: "Test_E_ORDERED_ARGTYPE_3_Validator",
			err: &Error{C: E_ORDERED_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"lte:&S.F"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "lte",
					Args: []*rules.Arg{{Type: rules.ARG_FIELD_ABS, Value: "S.F"}},
					Spec: specs.Get("lte"),
				},
				ra:  &rules.Arg{Type: rules.ARG_FIELD_ABS, Value: "S.F"},
				raf: T._sf(),
			},
		}},
	}, {
		name: "length",
		tests: []testCase{{
			name: "Test_length_Validator", err: nil,
		}, {
			name: "Test_E_LENGTH_NOLEN_1_Validator",
			err: &Error{C: E_LENGTH_NOLEN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:10"`,
					Obj: &types.Obj{Type: T.uint},
				},
				ty: T.uint,
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("len"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NOLEN_2_Validator",
			err: &Error{C: E_LENGTH_NOLEN,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:10"`,
					Obj: &types.Obj{Type: mytime},
				},
				ty: mytime,
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("len"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NORUNE_1_Validator",
			err: &Error{C: E_LENGTH_NORUNE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:10"`,
					Obj: &types.Obj{Type: T.uint},
				},
				ty: T.uint,
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NORUNE_2_Validator",
			err: &Error{C: E_LENGTH_NORUNE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:10"`,
					Obj: &types.Obj{Type: mytime},
				},
				ty: mytime,
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NORUNE_3_Validator",
			err: &Error{C: E_LENGTH_NORUNE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:10"`,
					Obj: &types.Obj{Type: T.rune},
				},
				ty: T.rune,
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NORUNE_4_Validator",
			err: &Error{C: E_LENGTH_NORUNE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:10"`,
					Obj: &types.Obj{Type: T.Slice(T.int16)},
				},
				ty: T.Slice(T.int16),
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NORUNE_5_Validator",
			err: &Error{C: E_LENGTH_NORUNE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:10"`,
					Obj: &types.Obj{Type: T.Slice(T.uint)},
				},
				ty: T.Slice(T.uint),
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NORUNE_6_Validator",
			err: &Error{C: E_LENGTH_NORUNE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:10"`,
					Obj: &types.Obj{Type: T.Array(21, T.byte)},
				},
				ty: T.Array(21, T.byte),
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NOARG_1_Validator",
			err: &Error{C: E_LENGTH_NOARG,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{}},
					Spec: specs.Get("len"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NOARG_2_Validator",
			err: &Error{C: E_LENGTH_NOARG,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len::"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{}, {}},
					Spec: specs.Get("len"),
				},
			},
		}, {
			name: "Test_E_LENGTH_NOARG_3_Validator",
			err: &Error{C: E_LENGTH_NOARG,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount::"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{}, {}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_BOUNDS_1_Validator",
			err: &Error{C: E_LENGTH_BOUNDS,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:10:5"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}, {Type: rules.ARG_INT, Value: "5"}},
					Spec: specs.Get("len"),
				},
			},
		}, {
			name: "Test_E_LENGTH_BOUNDS_2_Validator",
			err: &Error{C: E_LENGTH_BOUNDS,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:100:99"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "100"}, {Type: rules.ARG_INT, Value: "99"}},
					Spec: specs.Get("runecount"),
				},
			},
		}, {
			name: "Test_E_LENGTH_BOUNDS_3_Validator",
			err: &Error{C: E_LENGTH_BOUNDS,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:42:42"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "42"}, {Type: rules.ARG_INT, Value: "42"}},
					Spec: specs.Get("len"),
				},
			},
		}, {
			name: "Test_E_LENGTH_ARGTYPE_1_Validator",
			err: &Error{C: E_LENGTH_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{Type: rules.ARG_STRING, Value: "foo"}},
					Spec: specs.Get("len"),
				},
				ra: &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
			},
		}, {
			name: "Test_E_LENGTH_ARGTYPE_2_Validator",
			err: &Error{C: E_LENGTH_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:123.987"`,
					Obj: &types.Obj{Type: T.string},
				},
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_FLOAT, Value: "123.987"}},
					Spec: specs.Get("runecount"),
				},
				ra: &rules.Arg{Type: rules.ARG_FLOAT, Value: "123.987"},
			},
		}, {
			name: "Test_E_LENGTH_ARGTYPE_3_Validator",
			err: &Error{C: E_LENGTH_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"len:-123"`,
					Obj: &types.Obj{Type: T.string},
				},
				r: &rules.Rule{
					Name: "len",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "-123"}},
					Spec: specs.Get("len"),
				},
				ra: &rules.Arg{Type: rules.ARG_INT, Value: "-123"},
			},
		}, {
			name: "Test_E_LENGTH_ARGTYPE_4_Validator",
			err: &Error{C: E_LENGTH_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"runecount:&S.F"`,
					Obj: &types.Obj{Type: T.string},
				},
				r: &rules.Rule{
					Name: "runecount",
					Args: []*rules.Arg{{Type: rules.ARG_FIELD_ABS, Value: "S.F"}},
					Spec: specs.Get("runecount"),
				},
				ra:  &rules.Arg{Type: rules.ARG_FIELD_ABS, Value: "S.F"},
				raf: T._sf(),
			},
		}},
	}, {
		name: "range",
		tests: []testCase{{
			name: "Test_range_Validator", err: nil,
		}, {
			name: "Test_between_Validator", err: nil,
		}, {
			name: "Test_E_RANGE_TYPE_1_Validator",
			err: &Error{C: E_RANGE_TYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:10:20"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "10"},
						{Type: rules.ARG_INT, Value: "20"},
					},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_TYPE_2_Validator",
			err: &Error{C: E_RANGE_TYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:10:20"`,
					Obj: &types.Obj{Type: T.Slice(T.int)},
				},
				ty: T.Slice(T.int),
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "10"},
						{Type: rules.ARG_INT, Value: "20"},
					},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_NOARG_1_Validator",
			err: &Error{C: E_RANGE_NOARG,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng::"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{{}, {}},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_NOARG_2_Validator",
			err: &Error{C: E_RANGE_NOARG,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:10:"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{{Type: rules.ARG_INT, Value: "10"}, {}},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_NOARG_3_Validator",
			err: &Error{C: E_RANGE_NOARG,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng::20"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{{}, {Type: rules.ARG_INT, Value: "20"}},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_ARGTYPE_1_Validator",
			err: &Error{C: E_RANGE_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:foo:20"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
						{Type: rules.ARG_INT, Value: "20"},
					},
					Spec: specs.Get("rng"),
				},
				ra: &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
			},
		}, {
			name: "Test_E_RANGE_ARGTYPE_2_Validator",
			err: &Error{C: E_RANGE_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:3.14:20"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_FLOAT, Value: "3.14"},
						{Type: rules.ARG_INT, Value: "20"},
					},
					Spec: specs.Get("rng"),
				},
				ra: &rules.Arg{Type: rules.ARG_FLOAT, Value: "3.14"},
			},
		}, {
			name: "Test_E_RANGE_ARGTYPE_3_Validator",
			err: &Error{C: E_RANGE_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:-10:20"`,
					Obj: &types.Obj{Type: T.uint16},
				},
				ty: T.uint16,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "-10"},
						{Type: rules.ARG_INT, Value: "20"},
					},
					Spec: specs.Get("rng"),
				},
				ra: &rules.Arg{Type: rules.ARG_INT, Value: "-10"},
			},
		}, {
			name: "Test_E_RANGE_ARGTYPE_4_Validator",
			err: &Error{C: E_RANGE_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:&S.F:20"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_FIELD_ABS, Value: "S.F"},
						{Type: rules.ARG_INT, Value: "20"},
					},
					Spec: specs.Get("rng"),
				},
				ra:  &rules.Arg{Type: rules.ARG_FIELD_ABS, Value: "S.F"},
				raf: T._sf(),
			},
		}, {
			name: "Test_E_RANGE_BOUNDS_1_Validator",
			err: &Error{C: E_RANGE_BOUNDS,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:20:10"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "20"},
						{Type: rules.ARG_INT, Value: "10"},
					},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_BOUNDS_2_Validator",
			err: &Error{C: E_RANGE_BOUNDS,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:10:10"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "10"},
						{Type: rules.ARG_INT, Value: "10"},
					},
					Spec: specs.Get("rng"),
				},
			},
		}, {
			name: "Test_E_RANGE_BOUNDS_3_Validator",
			err: &Error{C: E_RANGE_BOUNDS,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"rng:10:-20"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r: &rules.Rule{
					Name: "rng",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "10"},
						{Type: rules.ARG_INT, Value: "-20"},
					},
					Spec: specs.Get("rng"),
				},
			},
		}},
	}, {
		name: "enum",
		tests: []testCase{{
			name: "Test_enum_Validator", err: nil,
		}, {
			name: "Test_E_ENUM_NONAME_1_Validator",
			err: &Error{C: E_ENUM_NONAME,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"enum"`,
					Obj: &types.Obj{Type: T.uint},
				},
				ty: T.uint,
				r:  &rules.Rule{Name: "enum", Spec: specs.Get("enum")},
			},
		}, {
			name: "Test_E_ENUM_KIND_1_Validator",
			err: &Error{C: E_ENUM_KIND,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"enum"`,
					Obj: &types.Obj{Type: &types.Type{Kind: types.STRUCT, Name: "enum_kind", Pkg: T.pkg}},
				},
				ty: &types.Type{Kind: types.STRUCT, Name: "enum_kind", Pkg: T.pkg},
				r:  &rules.Rule{Name: "enum", Spec: specs.Get("enum")},
			},
		}, {
			name: "Test_E_ENUM_NOCONST_1_Validator",
			err: &Error{C: E_ENUM_NOCONST,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"enum"`,
					Obj: &types.Obj{Type: &types.Type{Kind: types.UINT, Name: "enum_noconst", Pkg: T.pkg}},
				},
				ty: &types.Type{Kind: types.UINT, Name: "enum_noconst", Pkg: T.pkg},
				r:  &rules.Rule{Name: "enum", Spec: specs.Get("enum")},
			},
		}},
	}, {
		name: "function",
		tests: []testCase{{
			name: "Test_function_Validator", err: nil,
			show: true,
		}, {
			name: "Test_E_FUNCTION_INTYPE_1_Validator",
			err: &Error{C: E_FUNCTION_INTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"contains:foo"`,
					Obj: &types.Obj{Type: T.float32},
				},
				ty: T.float32,
				r: &rules.Rule{
					Name: "contains",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("contains"),
				},
			},
		}, {
			name: "Test_E_FUNCTION_ARGTYPE_1_Validator",
			err: &Error{C: E_FUNCTION_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"uuid:v6"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "uuid",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "v6"},
					},
					Spec: specs.Get("uuid"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "v6"},
				fp:  &types.Var{Name: "ver", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_1_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"alpha:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "alpha",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("alpha"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "lang", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_2_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"alnum:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "alnum",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("alnum"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "lang", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_3_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"ccy:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "ccy",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
						{},
					},
					Spec: specs.Get("ccy"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "code", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_4_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"decimal:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "decimal",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("decimal"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "locale", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_5_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"hash:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "hash",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("hash"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "algo", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_6_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"ip:5"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "ip",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "5"},
					},
					Spec: specs.Get("ip"),
				},
				ra:  &rules.Arg{Type: rules.ARG_INT, Value: "5"},
				fp:  &types.Var{Name: "ver", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_7_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"isbn:12"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "isbn",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "12"},
					},
					Spec: specs.Get("isbn"),
				},
				ra:  &rules.Arg{Type: rules.ARG_INT, Value: "12"},
				fp:  &types.Var{Name: "ver", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_8_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"iso639:3"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "iso639",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "3"},
					},
					Spec: specs.Get("iso639"),
				},
				ra:  &rules.Arg{Type: rules.ARG_INT, Value: "3"},
				fp:  &types.Var{Name: "num", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_9_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"iso31661a:1"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "iso31661a",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "1"},
					},
					Spec: specs.Get("iso31661a"),
				},
				ra:  &rules.Arg{Type: rules.ARG_INT, Value: "1"},
				fp:  &types.Var{Name: "num", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_10_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"mac:7"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "mac",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "7"},
					},
					Spec: specs.Get("mac"),
				},
				ra:  &rules.Arg{Type: rules.ARG_INT, Value: "7"},
				fp:  &types.Var{Name: "space", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_11_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"re:[0-9)?"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "re",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "[0-9)?"},
					},
					Spec: specs.Get("re"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "[0-9)?"},
				fp:  &types.Var{Name: "expr", Type: T.string},
				fpi: T.iptr(0),
				err: T._err,
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_12_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"uuid:6"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "uuid",
					Args: []*rules.Arg{
						{Type: rules.ARG_INT, Value: "6"},
					},
					Spec: specs.Get("uuid"),
				},
				ra:  &rules.Arg{Type: rules.ARG_INT, Value: "6"},
				fp:  &types.Var{Name: "ver", Type: T.int},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_13_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"phone:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "phone",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("phone"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "cc", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_14_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"vat:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "vat",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("vat"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "cc", Type: T.string},
				fpi: T.iptr(0),
			},
		}, {
			name: "Test_E_FUNCTION_ARGVALUE_15_Validator",
			err: &Error{C: E_FUNCTION_ARGVALUE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"zip:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "zip",
					Args: []*rules.Arg{
						{Type: rules.ARG_STRING, Value: "foo"},
					},
					Spec: specs.Get("zip"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "cc", Type: T.string},
				fpi: T.iptr(0),
			},
		}},
	}, {
		name: "method",
		tests: []testCase{{
			name: "Test_method_Validator", err: nil,
		}, {
			name: "Test_E_METHOD_TYPE_1_Validator",
			err: &Error{C: E_METHOD_TYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"isvalid"`,
					Obj: &types.Obj{Type: T.int},
				},
				ty: T.int,
				r:  &rules.Rule{Name: "isvalid", Spec: specs.Get("isvalid")},
			},
		}},
	}, {
		name: "optional",
		tests: []testCase{{
			name: "Test_optional_Validator", err: nil,
		}, {
			name: "Test_E_OPTIONAL_CONFLICT_1_Validator",
			err: &Error{C: E_OPTIONAL_CONFLICT,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"optional,required"`,
					Obj: &types.Obj{Type: T.Ptr(T.string)},
				},
				ty: T.Ptr(T.string),
				r:  &rules.Rule{Name: "optional", Spec: specs.Get("optional")},
				r2: &rules.Rule{Name: "required", Spec: specs.Get("required")},
			},
		}, {
			name: "Test_E_OPTIONAL_CONFLICT_2_Validator",
			err: &Error{C: E_OPTIONAL_CONFLICT,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"required,optional"`,
					Obj: &types.Obj{Type: T.Ptr(T.string)},
				},
				ty: T.Ptr(T.string),
				r:  &rules.Rule{Name: "optional", Spec: specs.Get("optional")},
				r2: &rules.Rule{Name: "required", Spec: specs.Get("required")},
			},
		}, {
			name: "Test_E_OPTIONAL_CONFLICT_3_Validator",
			err: &Error{C: E_OPTIONAL_CONFLICT,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"notnil,optional"`,
					Obj: &types.Obj{Type: T.Ptr(T.string)},
				},
				ty: T.Ptr(T.string),
				r:  &rules.Rule{Name: "optional", Spec: specs.Get("optional")},
				r2: &rules.Rule{Name: "notnil", Spec: specs.Get("notnil")},
			},
		}, {
			name: "Test_E_OPTIONAL_CONFLICT_4_Validator",
			err: &Error{C: E_OPTIONAL_CONFLICT,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `is:"omitnil,notnil"`,
					Obj: &types.Obj{Type: T.Ptr(T.string)},
				},
				ty: T.Ptr(T.string),
				r:  &rules.Rule{Name: "omitnil", Spec: specs.Get("omitnil")},
				r2: &rules.Rule{Name: "notnil", Spec: specs.Get("notnil")},
			},
		}},
	}, {
		name: "preprocessor",
		tests: []testCase{{
			name: "Test_preproc_Validator", err: nil,
		}, {
			name: "Test_E_PREPROC_INTYPE_1_Validator",
			err: &Error{C: E_PREPROC_INTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p1"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r:  &rules.Rule{Name: "p1", Spec: specs.Get("pre:p1")},
			},
		}, {
			name: "Test_E_PREPROC_OUTTYPE_1_Validator",
			err: &Error{C: E_PREPROC_OUTTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p2"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r:  &rules.Rule{Name: "p2", Spec: specs.Get("pre:p2")},
			},
		}, {
			name: "Test_E_PREPROC_ARGTYPE_1_Validator",
			err: &Error{C: E_PREPROC_ARGTYPE,
				sf: &types.StructField{
					Pkg:  T.pkg,
					Name: "F", IsExported: true,
					Tag: `pre:"p4:foo"`,
					Obj: &types.Obj{Type: T.string},
				},
				ty: T.string,
				r: &rules.Rule{
					Name: "p4",
					Args: []*rules.Arg{{Type: rules.ARG_STRING, Value: "foo"}},
					Spec: specs.Get("pre:p4"),
				},
				ra:  &rules.Arg{Type: rules.ARG_STRING, Value: "foo"},
				fp:  &types.Var{Name: "opt", Type: T.uint},
				fpi: T.iptr(0),
			},
		}},
	}}

	cfg := loadConfig("testdata/configs/test_custom_rules.yaml")
	if err := specs.LoadCustomSpecs(cfg, &test_ast); err != nil {
		t.Fatalf("loadConfig(testdata/configs/test_custom_rules.yaml) failed: %v", err)
	}
	fkCfg := &config.FieldKeyConfig{
		Join:      config.Bool{Value: true, IsSet: true},
		Separator: config.String{Value: ".", IsSet: true},
	}
	compare := compare.Config{ObserveFieldTag: "cmp"}

	for _, gg := range groups {
		t.Run(gg.name, func(t *testing.T) {
			for _, tt := range gg.tests {
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
	int16   *types.Type
	int32   *types.Type
	int64   *types.Type
	uint    *types.Type
	uint8   *types.Type
	uint16  *types.Type
	uint64  *types.Type
	bool    *types.Type
	float32 *types.Type
	float64 *types.Type
	rune    *types.Type
	byte    *types.Type
	loc     *types.Type

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
	int16:   &types.Type{Kind: types.INT16},
	int32:   &types.Type{Kind: types.INT32},
	int64:   &types.Type{Kind: types.INT64},
	uint:    &types.Type{Kind: types.UINT},
	uint8:   &types.Type{Kind: types.UINT8},
	uint16:  &types.Type{Kind: types.UINT16},
	uint64:  &types.Type{Kind: types.UINT64},
	float32: &types.Type{Kind: types.FLOAT32},
	float64: &types.Type{Kind: types.FLOAT64},
	bool:    &types.Type{Kind: types.BOOL},
	rune:    &types.Type{Kind: types.INT32, IsRune: true},
	byte:    &types.Type{Kind: types.UINT8, IsByte: true},
	loc:     nil,

	pkg: types.Pkg{
		Path: "github.com/frk/valid/cmd/internal/rules/checker/testdata",
		Name: "testdata",
	},
}

func init() {
	specs.AddCustomSpec("pre:p1", &rules.Spec{
		Name: "p1", Kind: rules.PREPROC,
		Func: &rules.FuncIdent{
			Pkg:  "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
			Name: "MyPreproc1",
		}},
		&types.Func{
			Name: "MyPreproc1",
			Type: &types.Type{
				Pkg: types.Pkg{
					Name: "mypkg",
					Path: "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
				},
				In:   []*types.Var{{Name: "v", Type: &types.Type{Kind: types.FLOAT64}}},
				Out:  []*types.Var{{Type: &types.Type{Kind: types.STRING}}},
				Kind: types.FUNC,
			},
		})

	specs.AddCustomSpec("pre:p2", &rules.Spec{
		Name: "p2", Kind: rules.PREPROC,
		Func: &rules.FuncIdent{
			Pkg:  "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
			Name: "MyPreproc2",
		}},
		&types.Func{
			Name: "MyPreproc2",
			Type: &types.Type{
				Pkg: types.Pkg{
					Name: "mypkg",
					Path: "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
				},
				In: []*types.Var{
					{Name: "v", Type: &types.Type{Kind: types.STRING}},
					{Name: "opt", Type: &types.Type{Kind: types.UINT}},
				},
				Out:  []*types.Var{{Type: &types.Type{Kind: types.FLOAT64}}},
				Kind: types.FUNC,
			},
		})

	specs.AddCustomSpec("pre:p4", &rules.Spec{
		Name: "p4", Kind: rules.PREPROC,
		ArgMin: 1, ArgMax: 1,
		Func: &rules.FuncIdent{
			Pkg:  "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
			Name: "MyPreproc4",
		}},
		&types.Func{
			Name: "MyPreproc4",
			Type: &types.Type{
				Pkg: types.Pkg{
					Name: "mypkg",
					Path: "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
				},
				In: []*types.Var{
					{Name: "v", Type: &types.Type{Kind: types.STRING}},
					{Name: "opt", Type: &types.Type{Kind: types.UINT}},
				},
				Out:  []*types.Var{{Type: &types.Type{Kind: types.STRING}}},
				Kind: types.FUNC,
			},
		})

	// non-PREPROC, for a test that ensures the code checks the Kind==PREPROC
	specs.AddCustomSpec("pre:p0", &rules.Spec{
		Name: "p0", Kind: rules.FUNCTION,
		ArgMin: 1, ArgMax: 1,
		Func: &rules.FuncIdent{
			Pkg:  "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
			Name: "MyPreproc0",
		}},
		&types.Func{
			Name: "MyPreproc0",
			Type: &types.Type{
				Pkg: types.Pkg{
					Name: "mypkg",
					Path: "github.com/frk/valid/cmd/internal/rules/checker/testdata/mypkg",
				},
				In:   []*types.Var{{Name: "v", Type: &types.Type{Kind: types.STRING}}},
				Out:  []*types.Var{{Type: &types.Type{Kind: types.STRING}}},
				Kind: types.FUNC,
			},
		})
}
