package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_lengthCheck(t *testing.T) {
	mypkg := gotype.Pkg{
		Path: "github.com/frk/valid/cmd/internal/rules/testdata/mypkg",
		Name: "mypkg",
	}
	timepkg := gotype.Pkg{
		Path: "time",
		Name: "time",
	}
	mytime := &gotype.Type{
		Pkg:  mypkg,
		Kind: gotype.K_STRUCT,
		Name: "Time",
		Fields: []*gotype.StructField{{
			Pkg:    timepkg,
			Name:   "wall",
			Object: &gotype.Object{Type: T.uint64},
			Var:    T._var,
		}, {
			Pkg:    timepkg,
			Name:   "ext",
			Object: &gotype.Object{Type: T.int64},
			Var:    T._var,
		}, {
			Pkg:  timepkg,
			Name: "loc",
			Object: &gotype.Object{Type: &gotype.Type{
				Kind: gotype.K_PTR,
				Elem: &gotype.Object{Type: T.loc},
			}},
			Var: T._var,
		}},
		IsExported: true,
	}

	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_length_Validator", err: nil,
	}, {
		name: "Test_ERR_LENGTH_NOLEN_1_Validator",
		err: &Error{C: ERR_LENGTH_NOLEN, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:10"`,
				Object: &gotype.Object{Type: T.uint},
				Var:    T._var,
			},
			ty: T.uint,
			r: &Rule{
				Name: "len",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("len"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NOLEN_2_Validator",
		err: &Error{C: ERR_LENGTH_NOLEN, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:10"`,
				Object: &gotype.Object{Type: mytime},
				Var:    T._var,
			},
			ty: mytime,
			r: &Rule{
				Name: "len",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("len"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NORUNE_1_Validator",
		err: &Error{C: ERR_LENGTH_NORUNE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:10"`,
				Object: &gotype.Object{Type: T.uint},
				Var:    T._var,
			},
			ty: T.uint,
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NORUNE_2_Validator",
		err: &Error{C: ERR_LENGTH_NORUNE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:10"`,
				Object: &gotype.Object{Type: mytime},
				Var:    T._var,
			},
			ty: mytime,
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NORUNE_3_Validator",
		err: &Error{C: ERR_LENGTH_NORUNE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:10"`,
				Object: &gotype.Object{Type: T.Slice(T.rune)},
				Var:    T._var,
			},
			ty: T.Slice(T.rune),
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NORUNE_4_Validator",
		err: &Error{C: ERR_LENGTH_NORUNE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:10"`,
				Object: &gotype.Object{Type: T.Slice(T.int32)},
				Var:    T._var,
			},
			ty: T.Slice(T.int32),
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NORUNE_5_Validator",
		err: &Error{C: ERR_LENGTH_NORUNE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:10"`,
				Object: &gotype.Object{Type: T.Slice(T.uint)},
				Var:    T._var,
			},
			ty: T.Slice(T.uint),
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NORUNE_6_Validator",
		err: &Error{C: ERR_LENGTH_NORUNE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:10"`,
				Object: &gotype.Object{Type: T.Array(21, T.byte)},
				Var:    T._var,
			},
			ty: T.Array(21, T.byte),
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NOARG_1_Validator",
		err: &Error{C: ERR_LENGTH_NOARG, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "len",
				Args: []*Arg{{}},
				Spec: GetSpec("len"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NOARG_2_Validator",
		err: &Error{C: ERR_LENGTH_NOARG, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len::"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "len",
				Args: []*Arg{{}, {}},
				Spec: GetSpec("len"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_NOARG_3_Validator",
		err: &Error{C: ERR_LENGTH_NOARG, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount::"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{}, {}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_BOUNDS_1_Validator",
		err: &Error{C: ERR_LENGTH_BOUNDS, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:10:5"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "len",
				Args: []*Arg{{Type: ARG_INT, Value: "10"}, {Type: ARG_INT, Value: "5"}},
				Spec: GetSpec("len"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_BOUNDS_2_Validator",
		err: &Error{C: ERR_LENGTH_BOUNDS, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:100:99"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_INT, Value: "100"}, {Type: ARG_INT, Value: "99"}},
				Spec: GetSpec("runecount"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_BOUNDS_3_Validator",
		err: &Error{C: ERR_LENGTH_BOUNDS, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:42:42"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			ty: T.string,
			r: &Rule{
				Name: "len",
				Args: []*Arg{{Type: ARG_INT, Value: "42"}, {Type: ARG_INT, Value: "42"}},
				Spec: GetSpec("len"),
			},
		},
	}, {
		name: "Test_ERR_LENGTH_ARGTYPE_1_Validator",
		err: &Error{C: ERR_LENGTH_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:foo"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			r: &Rule{
				Name: "len",
				Args: []*Arg{{Type: ARG_STRING, Value: "foo"}},
				Spec: GetSpec("len"),
			},
			ra: &Arg{Type: ARG_STRING, Value: "foo"},
		},
	}, {
		name: "Test_ERR_LENGTH_ARGTYPE_2_Validator",
		err: &Error{C: ERR_LENGTH_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:123.987"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_FLOAT, Value: "123.987"}},
				Spec: GetSpec("runecount"),
			},
			ra: &Arg{Type: ARG_FLOAT, Value: "123.987"},
		},
	}, {
		name: "Test_ERR_LENGTH_ARGTYPE_3_Validator",
		err: &Error{C: ERR_LENGTH_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"len:-123"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			r: &Rule{
				Name: "len",
				Args: []*Arg{{Type: ARG_INT, Value: "-123"}},
				Spec: GetSpec("len"),
			},
			ra: &Arg{Type: ARG_INT, Value: "-123"},
		},
	}, {
		name: "Test_ERR_LENGTH_ARGTYPE_4_Validator",
		err: &Error{C: ERR_LENGTH_ARGTYPE, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"runecount:&S.F"`,
				Object: &gotype.Object{Type: T.string},
				Var:    T._var,
			},
			r: &Rule{
				Name: "runecount",
				Args: []*Arg{{Type: ARG_FIELD_ABS, Value: "S.F"}},
				Spec: GetSpec("runecount"),
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
