package search

import (
	"fmt"
	"testing"

	"github.com/frk/compare"
)

func TestFindObject(t *testing.T) {
	_err := fmt.Errorf("") // dummy to satisfy `cmp:"+"`

	tests := []struct {
		pkgpath string
		name    string
		err     error
		show    bool
	}{{
		pkgpath: "strings",
		name:    "Contains",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name:    "FuncObject",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name:    "TypeObject1",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name:    "typeObject2",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name:    "typeObject3",
		err: &Error{C: ERR_OBJECT_NOTFOUND,
			pkg:  "github.com/frk/valid/cmd/internal/v2/search/testdata",
			name: "typeObject3",
		},
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name:    "VarObject",
		err: &Error{C: ERR_OBJECT_NOTFOUND,
			pkg:  "github.com/frk/valid/cmd/internal/v2/search/testdata",
			name: "VarObject",
		},
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/v2/search/testdata",
		name:    "ConstObject",
		err: &Error{C: ERR_OBJECT_NOTFOUND,
			pkg:  "github.com/frk/valid/cmd/internal/v2/search/testdata",
			name: "ConstObject",
		},
	}, {
		pkgpath: "strings",
		name:    "Abracadabra",
		err:     &Error{C: ERR_OBJECT_NOTFOUND, pkg: "strings", name: "Abracadabra"},
	}, {
		pkgpath: "sgnirts",
		name:    "Contains",
		err:     &Error{C: ERR_PKG_ERROR, pkg: "sgnirts", name: "Contains", err: _err},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for i, tt := range tests {
		obj, err := FindObject(tt.pkgpath, tt.name, &testast)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("Error: %v (%v)", e, err)
		} else if err == nil {
			if p := obj.Pkg(); p.Path() != tt.pkgpath || obj.Name() != tt.name {
				t.Errorf("#%d: want=%s.%s; got err=%v", i, tt.pkgpath, tt.name, obj)
			}
		}
		if tt.show && tt.err != nil {
			fmt.Println(err)
		}
	}
}
