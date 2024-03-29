package search

import (
	"fmt"
	"go/token"
	"go/types"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/frk/compare"
)

var testast AST

func TestMain(m *testing.M) {
	if _, err := Search("testdata/", false, nil, nil, &testast); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestSearch(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name      string
		dir       string
		recursive bool
		rx        *regexp.Regexp
		filter    func(string) bool
		want      []*Package
		err       error
	}{{
		name:      "should find one match in both files",
		dir:       "testdata/",
		recursive: true,
		rx:        regexp.MustCompile(`(?i:validator)$`),
		filter:    nil,
		want: []*Package{{
			Name: "testdata",
			Path: "github.com/frk/valid/cmd/internal/search/testdata",
			Fset: &token.FileSet{},
			Type: &types.Package{},
			Info: &types.Info{},
			Files: []*File{{
				Path:    wd + "/testdata/validator_1.go",
				Package: &Package{},
				Matches: []*Match{{
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}},
			}, {
				Path:    wd + "/testdata/validator_2.go",
				Package: &Package{},
				Matches: []*Match{{
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}},
			}},
		}},
	}, {
		name: "should find one match in both files",
		dir:  "testdata/",
		rx:   regexp.MustCompile(`(?i:Params)$`),
		want: []*Package{{
			Name: "testdata",
			Path: "github.com/frk/valid/cmd/internal/search/testdata",
			Fset: &token.FileSet{},
			Type: &types.Package{},
			Info: &types.Info{},
			Files: []*File{{
				Path:    wd + "/testdata/validator_1.go",
				Package: &Package{},
				Matches: []*Match{{
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}},
			}, {
				Path:    wd + "/testdata/validator_2.go",
				Package: &Package{},
				Matches: []*Match{{
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}},
			}},
		}},
	}, {
		name: "should find two matches in file xxx_1.go",
		dir:  "testdata/",
		rx:   regexp.MustCompile(`(?i:validator)|(?i:params)$`),
		filter: func(filepath string) bool {
			return strings.HasSuffix(filepath, "_1.go")
		},
		want: []*Package{{
			Name: "testdata",
			Path: "github.com/frk/valid/cmd/internal/search/testdata",
			Fset: &token.FileSet{},
			Type: &types.Package{},
			Info: &types.Info{},
			Files: []*File{{
				Path:    wd + "/testdata/validator_1.go",
				Package: &Package{},
				Matches: []*Match{{
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}, {
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}},
			}},
		}},
	}, {
		name: "should find one match in file xxx_2.go",
		dir:  "testdata/",
		rx:   regexp.MustCompile(`(?i:params)$`),
		filter: func(filepath string) bool {
			return strings.HasSuffix(filepath, "_2.go")
		},
		want: []*Package{{
			Name: "testdata",
			Path: "github.com/frk/valid/cmd/internal/search/testdata",
			Fset: &token.FileSet{},
			Type: &types.Package{},
			Info: &types.Info{},
			Files: []*File{{
				Path:    wd + "/testdata/validator_2.go",
				Package: &Package{},
				Matches: []*Match{{
					Named: &types.Named{},
					Fset:  &token.FileSet{},
					Pos:   token.Pos(1),
				}},
			}},
		}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := Search(tt.dir, tt.recursive, tt.rx, tt.filter, nil)
			if e := compare.Compare(out, tt.want); e != nil {
				t.Errorf("Out: %v", e)
			}
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("Error: %v (%v)", e, err)
			}
			//fn, rawCfg, err := FindFunc(tt.pkg, tt.name, &testast)
			//if e := compare.Compare(err, tt.err); e != nil {
			//	t.Errorf("Error: %v (%v)", e, err)
			//} else if err == nil {
			//	if p := fn.Pkg(); p.Path() != tt.pkg || fn.Name() != tt.name {
			//		t.Errorf("#%d: want=%s.%s; got err=%v", i, tt.pkg, tt.name, fn)
			//	}
			//	if e := compare.Compare(rawCfg, tt.rawCfg); e != nil {
			//		t.Errorf("#%d: *rulecfg.RuleConfig: %v\n%q", i, e, rawCfg)
			//	}
			//}
			//if tt.show && tt.err != nil {
			//	fmt.Println(err)
			//}
		})
	}
}

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
		pkg:  "github.com/frk/valid/cmd/internal/search/testdata",
		name: "IsFoo",
	}, {
		pkg:  "github.com/frk/valid/cmd/internal/search/testdata",
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
		fn, rawCfg, err := FindFunc(tt.pkg, tt.name, &testast)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("Error: %v (%v)", e, err)
		} else if err == nil {
			if p := fn.Pkg(); p.Path() != tt.pkg || fn.Name() != tt.name {
				t.Errorf("#%d: want=%s.%s; got err=%v", i, tt.pkg, tt.name, fn)
			}
			if e := compare.Compare(rawCfg, tt.rawCfg); e != nil {
				t.Errorf("#%d: *rulecfg.RuleConfig: %v\n%q", i, e, rawCfg)
			}
		}
		if tt.show && tt.err != nil {
			fmt.Println(err)
		}
	}
}

func TestFindConstantsByType(t *testing.T) {
	type konst struct {
		name string
		pkg  string
		typ  string
	}

	searchpkg := "github.com/frk/valid/cmd/internal/search/testdata"

	tests := []struct {
		pkg     string
		typname string
		want    []konst
	}{{
		pkg: searchpkg, typname: "K1",
		want: []konst{
			{name: "k11", pkg: searchpkg, typ: "K1"},
			{name: "k12", pkg: searchpkg, typ: "K1"},
			{name: "k13", pkg: searchpkg, typ: "K1"},
			{name: "k14", pkg: searchpkg, typ: "K1"},
			{name: "k1111", pkg: searchpkg, typ: "K1"},
		},
	}, {
		pkg: searchpkg, typname: "K2",
		want: []konst{
			{name: "k21", pkg: searchpkg, typ: "K2"},
			{name: "k22", pkg: searchpkg, typ: "K2"},
			{name: "k23", pkg: searchpkg, typ: "K2"},
		},
	}, {
		pkg: searchpkg, typname: "K3",
		want: []konst{
			{name: "k3foo", pkg: searchpkg, typ: "K3"},
			{name: "k3bar", pkg: searchpkg, typ: "K3"},
			{name: "k3baz", pkg: searchpkg, typ: "K3"},
		},
	}, {
		pkg: "go/parser", typname: "Mode",
		want: []konst{
			{name: "CustomMode", pkg: searchpkg, typ: "Mode"},
			{name: "PackageClauseOnly", pkg: "go/parser", typ: "Mode"},
			{name: "ImportsOnly", pkg: "go/parser", typ: "Mode"},
			{name: "ParseComments", pkg: "go/parser", typ: "Mode"},
			{name: "Trace", pkg: "go/parser", typ: "Mode"},
			{name: "DeclarationErrors", pkg: "go/parser", typ: "Mode"},
			{name: "SpuriousErrors", pkg: "go/parser", typ: "Mode"},
			{name: "SkipObjectResolution", pkg: "go/parser", typ: "Mode"}, // this constant was added in 1.17
			{name: "AllErrors", pkg: "go/parser", typ: "Mode"},
		},
	}}

	for _, tt := range tests {
		got := FindConstantsByType(tt.pkg, tt.typname, &testast)
		if len(got) != len(tt.want) {
			t.Errorf("%s.%s: len got=%d; want=%d", tt.pkg, tt.typname, len(got), len(tt.want))
		} else {
		kloop:
			for _, k := range tt.want {
				for _, c := range got {
					if typeNamed, ok := c.Type().(*types.Named); ok {
						name := c.Name()
						pkg := c.Pkg().Path()
						typ := typeNamed.Obj().Name()
						if k.name == name && k.pkg == pkg && k.typ == typ {
							continue kloop
						}
					}
				}

				t.Errorf("%s.%s: no match for %+v", tt.pkg, tt.typname, k)
			}
		}
	}
}

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
		pkgpath: "github.com/frk/valid/cmd/internal/search/testdata",
		name:    "FuncObject",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/search/testdata",
		name:    "TypeObject1",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/search/testdata",
		name:    "typeObject2",
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/search/testdata",
		name:    "typeObject3",
		err: &Error{C: ERR_OBJECT_NOTFOUND,
			pkg:  "github.com/frk/valid/cmd/internal/search/testdata",
			name: "typeObject3",
		},
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/search/testdata",
		name:    "VarObject",
		err: &Error{C: ERR_OBJECT_NOTFOUND,
			pkg:  "github.com/frk/valid/cmd/internal/search/testdata",
			name: "VarObject",
		},
	}, {
		pkgpath: "github.com/frk/valid/cmd/internal/search/testdata",
		name:    "ConstObject",
		err: &Error{C: ERR_OBJECT_NOTFOUND,
			pkg:  "github.com/frk/valid/cmd/internal/search/testdata",
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
