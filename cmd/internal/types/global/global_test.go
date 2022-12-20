package global

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/search"

	"github.com/frk/compare"
)

var test_ast search.AST
var test_pkg search.Package

func TestMain(m *testing.M) {
	pkgs, err := search.Search(
		"../testdata/",
		true,
		regexp.MustCompile(`(?i:Validator)$`),
		nil,
		&test_ast,
	)
	if err != nil {
		log.Fatal(err)
	}

	test_pkg = *(pkgs[0])

	os.Exit(m.Run())
}

func TestInit(t *testing.T) {
	_pkg := "github.com/frk/valid/cmd/internal/types/testdata"
	_err := fmt.Errorf("") // dummy to satisfy `cmp:"+"`

	var objId = func(name string) config.ObjectIdent {
		return config.ObjectIdent{
			Pkg:   _pkg,
			Name:  name,
			IsSet: true,
		}
	}

	var findObj = func(name string) types.Object {
		obj, err := search.FindObject(_pkg, name, &test_ast)
		if err != nil {
			return nil
		}
		return obj
	}

	type obj struct {
		id   config.ObjectIdent
		want bool
	}

	tests := []struct {
		ctor obj
		agg  obj
		err  error
		show bool
	}{{
		ctor: obj{id: objId("CustomErrorConstructor"), want: true},
		agg:  obj{id: objId("CustomErrorAggregator"), want: true},
	}, {
		ctor: obj{id: objId("Undefined"), want: false},
		err:  &Error{C: E_OBJECT_SEARCH, oid: objId("Undefined"), err: _err},
	}, {
		agg: obj{id: objId("Undefined"), want: false},
		err: &Error{C: E_OBJECT_SEARCH, oid: objId("Undefined"), err: _err},
	}, {
		ctor: obj{id: objId("NotAFuncObject"), want: false},
		err: &Error{
			C:   E_ERROR_CONSTRUCTOR_OBJECT,
			oid: objId("NotAFuncObject"),
			obj: findObj("NotAFuncObject"),
		},
	}, {
		agg: obj{id: objId("NotANamedTypeObject"), want: false},
		err: &Error{
			C:   E_ERROR_AGGREGATOR_OBJECT,
			oid: objId("NotANamedTypeObject"),
			obj: findObj("NotANamedTypeObject"),
		},
	}, {
		ctor: obj{id: objId("ErrorConstructorWithBadSignature"), want: false},
		err: &Error{
			C:   E_ERROR_CONSTRUCTOR_TYPE,
			oid: objId("ErrorConstructorWithBadSignature"),
			obj: findObj("ErrorConstructorWithBadSignature"),
		},
	}, {
		agg: obj{id: objId("ErrorAggregatorWithBadImpl"), want: false},
		err: &Error{
			C:   E_ERROR_AGGREGATOR_TYPE,
			oid: objId("ErrorAggregatorWithBadImpl"),
			obj: findObj("ErrorAggregatorWithBadImpl"),
		},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		cfg := config.Config{}
		cfg.ErrorHandling.Constructor = tt.ctor.id
		cfg.ErrorHandling.Aggregator = tt.agg.id

		err := Init(cfg, &test_ast)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("Error: %v (%v)", e, err)
		} else if err == nil {
			if got, want := (ErrorConstructor != nil), tt.ctor.want; got != want {
				t.Errorf("ErrorConstructor got=%t want=%t", got, want)
			}
			if got, want := (ErrorAggregator != nil), tt.agg.want; got != want {
				t.Errorf("ErrorAggregator got=%t want=%t", got, want)
			}
		}
		if tt.show && tt.err != nil {
			fmt.Println(err)
		}
	}
}
