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
	"github.com/frk/valid/cmd/internal/rules/checker"
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/search"
)

func TestFile(t *testing.T) {
	tests := []string{
		"empty/single/v",
		"empty/multi/v",

		////////////////////////////////////////////////////////////////

		"is/required/v",
		"is/notnil/v",
		"is/optional/v",
		"is/omitnil/v",
		"is/noguard/v",
		"is/eq/v",
		"is/ne/v",
		"is/gt/v",
		"is/lt/v",
		"is/gte/v",
		"is/lte/v",
		"is/min/v",
		"is/max/v",
		"is/len/v",
		"is/runecount/v",
		"is/rng/v",
		"is/between/v",
		"is/enum/v",
		"is/contains/v",
		"is/prefix/v",
		"is/suffix/v",
		"is/isvalid/v",

		//////////////////////////////////////////////////////////////////////

		"pre/ceil/v",
		"pre/floor/v",
		"pre/htmlesc/v",
		"pre/htmlunesc/v",
		"pre/lower/v",
		"pre/ltrim/v",
		"pre/quote/v",
		"pre/quoteascii/v",
		"pre/quotegraphic/v",
		"pre/repeat/v",
		"pre/replace/v",
		"pre/round/v",
		"pre/rtrim/v",
		"pre/title/v",
		"pre/trim/v",
		"pre/trimprefix/v",
		"pre/trimsuffix/v",
		"pre/upper/v",
		"pre/urlqueryesc/v",
		"pre/urlpathesc/v",
		"pre/validutf8/v",

		//////////////////////////////////////////////////////////////////////

		"types/slice/v",
		//"types/map/v",

		//////////////////////////////////////////////////////////////////////

		//"included/re/v",
		//"included/ascii/v",
		//"included/alpha/v",
		//"included/alnum/v",
		//"included/bic/v",
		//"included/btc/v",
		//"included/base32/v",
		//"included/base58/v",
		//"included/base64/v",
		//"included/binary/v",
		//"included/bool/v",
		//"included/cidr/v",
		//"included/cvv/v",
		//"included/ccy/v",
		//"included/datauri/v",
		//"included/decimal/v",
		//"included/digits/v",
		//"included/ean/v",
		//// TODO "included/ein/v",
		//"included/eth/v",
		//"included/email/v",
		//"included/fqdn/v",
		//"included/float/v",
		//"included/hsl/v",
		//"included/hash/v",
		//"included/hex/v",
		//"included/hexcolor/v",
		//"included/iban/v",
		//// TODO "included/ic/v",
		//"included/imei/v",
		//"included/ip/v",
		//"included/iprange/v",
		//"included/isbn/v",
		//"included/isin/v",
		//"included/iso639/v",
		//"included/iso31661a/v",
		//"included/iso4217/v",
		//"included/isrc/v",
		//"included/issn/v",
		//"included/in/v",
		//"included/int/v",
		//"included/json/v",
		//"included/jwt/v",
		//"included/latlong/v",
		//"included/locale/v",
		//"included/lower/v",
		//"included/mac/v",
		//"included/md5/v",
		//"included/mime/v",
		//"included/magneturi/v",
		//"included/mongoid/v",
		//"included/numeric/v",
		//"included/octal/v",
		//"included/pan/v",
		//"included/phone/v",
		//"included/port/v",
		//"included/rgb/v",
		//"included/ssn/v",
		//"included/semver/v",
		//"included/slug/v",
		//"included/strongpass/v",
		//// TODO "included/url/v",
		//"included/uuid/v",
		//"included/uint/v",
		//"included/upper/v",
		//"included/vat/v",
		//"included/zip/v",

		//////////////////////////////////////////////////////////////////////

		//"error/error_returning_rule_func/v",

		//////////////////////////////////////////////////////////////////////

		//"example/basic/v",
		//"example/preproc_only/v",
		//"example/preproc_only2/v",
		//"example/preproc_basic/v",
	}

	var _debug string
	//_debug = "error/error_returning_rule_func/v"

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
	if err := specs.Load(cfg, &AST); err != nil {
		t.Fatal(err)
	}

	fkCfg := &config.FieldKeyConfig{
		Join:      config.Bool{Value: true, IsSet: true},
		Separator: config.String{Value: ".", IsSet: true},
	}

	for _, filename := range tests {
		if _debug != "" && _debug != filename {
			continue
		}
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

			infos := make([]*checker.Info, len(f.Matches))
			for k, match := range f.Matches {
				cfg := checker.Config{
					AST:      &AST,
					FieldKey: fkCfg,
				}
				info := new(checker.Info)
				if err := checker.Check(cfg, match, info); err != nil {
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
