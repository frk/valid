package rules

import (
	"fmt"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"

	"github.com/frk/compare"
)

func TestChecker_enumCheck(t *testing.T) {
	tests := []struct {
		name string
		err  error
		show bool
	}{{
		name: "Test_enum_Validator", err: nil,
	}, {
		name: "Test_ERR_ENUM_NONAME_1_Validator",
		err: &Error{C: ERR_ENUM_NONAME, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"enum"`,
				Object: &gotype.Object{Type: T.uint},
				Var:    T._var,
			},
			ty: T.uint,
			r:  &Rule{Name: "enum", Spec: GetSpec("enum")},
		},
	}, {
		name: "Test_ERR_ENUM_KIND_1_Validator",
		err: &Error{C: ERR_ENUM_KIND, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"enum"`,
				Object: &gotype.Object{Type: &gotype.Type{Kind: gotype.K_STRUCT, Name: "enum_kind", Pkg: T.pkg}},
				Var:    T._var,
			},
			ty: &gotype.Type{Kind: gotype.K_STRUCT, Name: "enum_kind", Pkg: T.pkg},
			r:  &Rule{Name: "enum", Spec: GetSpec("enum")},
		},
	}, {
		name: "Test_ERR_ENUM_NOCONST_1_Validator",
		err: &Error{C: ERR_ENUM_NOCONST, a: T._ast, sfv: T._var,
			sf: &gotype.StructField{
				Pkg:  T.pkg,
				Name: "F", IsExported: true,
				Tag:    `is:"enum"`,
				Object: &gotype.Object{Type: &gotype.Type{Kind: gotype.K_UINT, Name: "enum_noconst", Pkg: T.pkg}},
				Var:    T._var,
			},
			ty: &gotype.Type{Kind: gotype.K_UINT, Name: "enum_noconst", Pkg: T.pkg},
			r:  &Rule{Name: "enum", Spec: GetSpec("enum")},
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
