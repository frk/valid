package search

import (
	"fmt"
	"testing"

	"github.com/frk/compare"
)

func TestFindFunc(t *testing.T) {
	_err := fmt.Errorf("") // dummy to satisfy `cmp:"+"`

	tests := []struct {
		pkg    string
		name   string
		rawCfg []byte
		err    error
		show   bool
	}{{
		pkg:  "strings",
		name: "Contains",
	}, {
		pkg:  "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name: "IsFoo",
	}, {
		pkg:  "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name: "isBar",
		rawCfg: []byte(`
name: bar
args:
  - { default: "123" }
  - options:
    - { value: "1", alias: v1 }
    - { value: "2", alias: v2 }
    - { value: "3", alias: v3 }
error: { text: "bar is not valid" }
`),
	}, {
		pkg:  "strings",
		name: "Abracadabra",
		err:  &Error{C: ERR_FUNC_NOTFOUND, pkg: "strings", name: "Abracadabra"},
	}, {
		pkg:  "sgnirts",
		name: "Contains",
		err:  &Error{C: ERR_PKG_ERROR, pkg: "sgnirts", name: "Contains", err: _err},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for i, tt := range tests {
		fn, cd, err := FindFunc(tt.pkg, tt.name, &testast)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("Error: %v (%v)", e, err)
		} else if err == nil {
			if p := fn.Pkg(); p.Path() != tt.pkg || fn.Name() != tt.name {
				t.Errorf("#%d: want=%s.%s; got err=%v", i, tt.pkg, tt.name, fn)
			}
			rawCfg := cd.(configDecoder).rawYAML
			if e := compare.Compare(string(rawCfg), string(tt.rawCfg)); e != nil {
				t.Errorf("#%d: *rulecfg.RuleConfig: %v\n%q", i, e, rawCfg)
			}
		}
		if tt.show && tt.err != nil {
			fmt.Println(err)
		}
	}
}
