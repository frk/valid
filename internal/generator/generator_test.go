package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/isvalid/internal/analysis"
	"github.com/frk/isvalid/internal/search"
)

func TestGenerator(t *testing.T) {
	tests := []string{
		"nothing_to_validate",
		"base_fields_with_rules",
		"nested_fields_with_rules",
		"base_fields_with_rules_and_nilguard",
		"base_fields_with_rules_and_notnil",
		"base_fields_with_rules_and_required",
		"nested_fields_with_rules_and_nilguard",
		"nested_fields_with_rules_and_notnil",
		"nested_fields_with_rules_and_required",
		/////////////

		"required",
		"notnil",
		"email",
		"url",
		"uri",
		"pan",
		"cvv",
		"ssn",
		"ein",
		"numeric",
		"hex",
		"hexcolor",
		"alphanum",
		"cidr",
		"phone",
		"zip",
		"uuid",
		"ip",
		"mac",
		"iso",
		"rfc",
		"re",
		"prefix",
		"suffix",
		"contains",
		"eq",
		"ne",
		"gt",
		"lt",
		"gte",
		"lte",
		"min",
		"max",
		"rng",
		"len",

		"error_constructor",
		"error_aggregator",
		"context_option",
		"references",
		"custom",
		"hooks",
		"isvalider",
		"enum",
		"slice",
	}

	anConf := analysis.Config{FieldKeyJoin: true, FieldKeySeparator: "."}

	var AST search.AST
	pkgs, err := search.Search("../testdata/generator", false, nil, &AST)
	if err != nil {
		t.Fatal(err)
	}
	pkg := pkgs[0]

	customrules := [][3]string{
		{"myrule", "github.com/frk/isvalid/internal/testdata/mypkg", "MyRule"},
		{"myrule2", "github.com/frk/isvalid/internal/testdata/mypkg", "MyRule2"},
		{"myrule3", "github.com/frk/isvalid/internal/testdata/mypkg", "MyRule3"},
	}
	for _, cr := range customrules {
		f, err := search.FindFunc(cr[1], cr[2], AST)
		if err != nil {
			t.Fatal(err)
		}
		if err := anConf.AddRuleFunc(cr[0], f); err != nil {
			t.Fatal(err)
		}
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			tinfos := []*TargetAnalysis{}
			fileprefix := "../testdata/generator/" + filename

			f, err := getFile(pkg, fileprefix+"_in.go")
			if err != nil {
				t.Fatal(err)
			}

			for _, match := range f.Matches {
				anInfo := &analysis.Info{}
				vs, err := anConf.Analyze(AST, match, anInfo)
				if err != nil {
					t.Error(err)
					return
				}

				tinfos = append(tinfos, &TargetAnalysis{vs, anInfo})
			}

			buf := new(bytes.Buffer)
			if err := Generate(buf, pkg.Name, tinfos); err != nil {
				t.Error(err)
				return
			}

			got := string(formatBytes(buf))

			out, err := ioutil.ReadFile(fileprefix + "_out.go")
			if err != nil {
				t.Fatal(err)
			}
			want := string(out)

			// compare
			if err := compare.Compare(got, want); err != nil {
				t.Error(err)
			}
		})
	}
}

// helper method...
func getFile(p *search.Package, filename string) (*search.File, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	for _, f := range p.Files {
		if f.Path == filename {
			return f, nil
		}
	}
	return nil, fmt.Errorf("file not found: %q", filename)
}

func formatBytes(buf *bytes.Buffer) []byte {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("format error: %s", err)
		return buf.Bytes()
	}
	return src
}
