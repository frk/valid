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
	"github.com/frk/isvalid/internal/parser"
)

func TestGenerator(t *testing.T) {
	tests := []string{
		"required",
		"notnil",
		"email",
		"url",
		"uri",
	}

	pkgs, err := parser.Parse("../testdata/generator", false, nil)
	if err != nil {
		t.Fatal(err)
	}
	pkg := pkgs[0]

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			tinfos := []*TargetInfo{}
			fileprefix := "../testdata/generator/" + filename

			f, err := getFile(pkg, fileprefix+"_in.go")
			if err != nil {
				t.Fatal(err)
			}

			for _, target := range f.Targets {
				ainfo := &analysis.Info{}
				vs, err := analysis.Run(pkg.Fset, target.Named, target.Pos, ainfo)
				if err != nil {
					t.Error(err)
					return
				}

				tinfos = append(tinfos, &TargetInfo{vs, ainfo})
			}

			buf := new(bytes.Buffer)
			conf := Config{}
			if err := Write(buf, pkg.Name, tinfos, conf); err != nil {
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
func getFile(p *parser.Package, filename string) (*parser.File, error) {
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
