package search

import (
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
			Path: "github.com/frk/valid/cmd/internal/v2/search/testdata",
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
			Path: "github.com/frk/valid/cmd/internal/v2/search/testdata",
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
			Path: "github.com/frk/valid/cmd/internal/v2/search/testdata",
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
			Path: "github.com/frk/valid/cmd/internal/v2/search/testdata",
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
