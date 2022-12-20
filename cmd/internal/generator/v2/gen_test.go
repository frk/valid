package generate

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/check"
	"github.com/frk/valid/cmd/internal/rules/spec"
	"github.com/frk/valid/cmd/internal/search"
)

func TestFile(t *testing.T) {
	tests := []string{
		"empty/single/v",
		"empty/multi/v",

		"is/required/v",
	}

	var AST search.AST
	pkgs, err := search.Search(
		"../testdata/",
		true,
		nil,
		nil,
		&AST,
	)
	if err != nil {
		t.Fatal(err)
	}

	cfg := loadConfig("../testdata/config.yaml")
	if err := spec.Load(cfg, &AST); err != nil {
		t.Fatal(err)
	}

	fkCfg := &config.FieldKeyConfig{
		Join:      config.Bool{Value: true, IsSet: true},
		Separator: config.String{Value: ".", IsSet: true},
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			fileprefix := "../testdata/" + filename
			f, pkg, err := getFile(pkgs, fileprefix+"_in.go")
			if err != nil {
				t.Fatal(err)
			}

			// if init, ok := globals[filename]; ok {
			// 	init(t)
			// 	defer types.Globals.Unset()
			// }

			infos := make([]*check.Info, len(f.Matches))
			for k, match := range f.Matches {
				cfg := check.Config{
					AST:      &AST,
					FieldKey: fkCfg,
				}
				info := new(check.Info)
				if err := check.Check(cfg, match, info); err != nil {
					t.Fatal(err)
				}
				infos[k] = info
			}

			code, err := File(pkg.Pkg(), infos)
			if err != nil {
				t.Error(err)
				return
			}

			got := string(formatBytes(code))
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

////////////////////////////////////////////////////////////////////////////////
// helpers

func getFile(pkgs []*search.Package, filename string) (*search.File, *search.Package, error) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		return nil, nil, err
	}

	for _, p := range pkgs {
		for _, f := range p.Files {
			if f.Path == filename {
				return f, p, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("file not found: %q", filename)
}

func formatBytes(code []byte) []byte {
	src, err := format.Source(code)
	if err != nil {
		log.Printf("format error: %s", err)
		return code
	}
	return src
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
