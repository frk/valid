package generator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/xtypes"
)

func TestGenerator(t *testing.T) {
	tests := []string{
		"validation/00_nothing_to_validate",
		"validation/01_required",
		"validation/02_notnil",
		"validation/03_eq",
		"validation/04_ne",
		"validation/05_gt",
		"validation/06_lt",
		"validation/07_gte",
		"validation/08_lte",
		"validation/09_min",
		"validation/10_max",
		"validation/11_rng",
		"validation/12_len",
		"validation/13_runecount",
		"validation/14_enum",
		"validation/15_optional",
		"validation/16_isvalid",
		"validation/17_-isvalid",
		"validation/18_prefix",
		"validation/19_suffix",
		"validation/20_contains",
		"validation/21_email",
		"validation/22_alnum",
		"validation/23_base64",
		"validation/24_cidr",
		"validation/25_cvv",
		//"validation/26_ein", (currently not implemented by github.com/frk/valid)
		"validation/27_fqdn",
		"validation/28_hex",
		"validation/29_hexcolor",
		"validation/30_ip",
		"validation/31_mac",
		"validation/32_numeric",
		"validation/33_pan",
		"validation/34_phone",
		"validation/35_ssn",
		"validation/36_strongpass",
		"validation/37_upper",
		//"validation/38_url", (currently not implemented by github.com/frk/valid)
		"validation/39_uuid",
		"validation/40_zip",
		"validation/41_error_constructor",
		"validation/42_error_aggregator",
		"validation/43_custom_funcs",
		"validation/44_regexp",
		"validation/45_hooks",
		"validation/46_references",
		"validation/47_nilguard",
		"validation/48_notnil",
		"validation/49_required",
		"validation/50_multiple_rules",
		"validation/51_slice",
		"validation/52_slice",
		"validation/53_map",
		"validation/54_map_and_slice",
		"validation/56_struct_nilguard",
		"validation/57_struct_notnil",
		"validation/58_struct_required",
		"validation/59_noguard",
		"validation/60_embed",
		"validation/61_omitkey",
		"validation/62_omitnil",

		// TODO test IsValid with `is:"-"`
		// TODO test IsValid combined with other rules
		// TODO test IsValid with `is:"-isvalid"` and combined with other rules

		"preproc/01_preproc_only",
		"preproc/02_preproc_only",
		"preproc/03_basic",
		"preproc/04_required",
		"preproc/05_slice",
		"preproc/06_map",

		// global
		"global/01_global_error_constructor",
		"global/02_global_error_aggregator",
		"global/03_global_error_agg_has_priority_over_ctor",
		"global/04_local_has_priority_over_global",

		// builtin/stdlib validation
		"is/required/v",
		"is/optional/v",
		"is/omitnil/v",
		"is/noguard/v",
		"is/eq/v",
		"is/ne/v",
		"is/gt/v",
		"is/lt/v",
		"is/gte/v",
		"is/min/v",
		"is/lte/v",
		"is/max/v",
		"is/rng/v",
		"is/between/v",
		"is/len/v",
		"is/runecount/v",
		"is/enum/v",
		"is/isvalid/v",
		"is/contains/v",
		"is/prefix/v",
		"is/suffix/v",

		// builtin/stdlib preprocessors
		"pre/lower/v",
		"pre/upper/v",
		"pre/title/v",
		"pre/trim/v",
		"pre/ltrim/v",
		"pre/rtrim/v",
		"pre/trimprefix/v",
		"pre/trimsuffix/v",
		"pre/repeat/v",
		"pre/replace/v",
		"pre/validutf8/v",
		"pre/quote/v",
		"pre/quoteascii/v",
		"pre/quotegraphic/v",
		"pre/urlqueryesc/v",
		"pre/urlpathesc/v",
		"pre/htmlesc/v",
		"pre/htmlunesc/v",
		"pre/round/v",
		"pre/ceil/v",
		"pre/floor/v",

		// included validation
		"included/re/v",
		"included/ascii/v",
		"included/alpha/v",
		"included/alnum/v",
		"included/bic/v",
		"included/btc/v",
		"included/base32/v",
		"included/base58/v",
		"included/base64/v",
		"included/binary/v",
		"included/bool/v",
		"included/cidr/v",
		"included/cvv/v",
		"included/ccy/v",
		"included/datauri/v",
		"included/decimal/v",
		"included/digits/v",
		"included/ean/v",
		// TODO "included/ein/v",
		"included/eth/v",
		"included/email/v",
		"included/fqdn/v",
		"included/float/v",
		"included/hsl/v",
		"included/hash/v",
		"included/hex/v",
		"included/hexcolor/v",
		"included/iban/v",
		// TODO "included/ic/v",
		"included/imei/v",
		"included/ip/v",
		"included/iprange/v",
		"included/isbn/v",
		"included/isin/v",
		"included/iso639/v",
		"included/iso31661a/v",
		"included/iso4217/v",
		"included/isrc/v",
		"included/issn/v",
		"included/in/v",
		"included/int/v",
		"included/json/v",
		"included/jwt/v",
		"included/latlong/v",
		"included/locale/v",
		"included/lower/v",
		"included/mac/v",
		"included/md5/v",
		"included/mime/v",
		"included/magneturi/v",
		"included/mongoid/v",
		"included/numeric/v",
		"included/octal/v",
		"included/pan/v",
		"included/phone/v",
		"included/port/v",
		"included/rgb/v",
		"included/ssn/v",
		"included/semver/v",
		"included/slug/v",
		"included/strongpass/v",
		// TODO "included/url/v",
		"included/uuid/v",
		"included/uint/v",
		"included/upper/v",
		"included/vat/v",
		"included/zip/v",

		// nested
		"nested/basic/v",
		"nested/complex/v",
		"nested/fieldabs/v",
		"nested/fieldrel/v",

		// error
		"error/error_returning_rule_func/v",

		// examples
		"example/basic/v",
	}

	var AST search.AST
	pkgs, err := search.Search(
		"testdata/",
		true,
		nil,
		nil,
		&AST,
	)
	if err != nil {
		t.Fatal(err)
	}

	mypkg := "github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
	globals := map[string]func(t *testing.T){
		"global/01_global_error_constructor": func(t *testing.T) {
			var cfg config.Config
			cfg.ErrorHandling.Constructor = config.ObjectIdent{mypkg, "NewError", true}
			if err := xtypes.Globals.Init(cfg, &AST); err != nil {
				t.Fatal(err)
			}
		},
		"global/02_global_error_aggregator": func(t *testing.T) {
			var cfg config.Config
			cfg.ErrorHandling.Aggregator = config.ObjectIdent{mypkg, "ErrorList", true}
			if err := xtypes.Globals.Init(cfg, &AST); err != nil {
				t.Fatal(err)
			}
		},
		"global/03_global_error_agg_has_priority_over_ctor": func(t *testing.T) {
			var cfg config.Config
			cfg.ErrorHandling.Constructor = config.ObjectIdent{mypkg, "NewError", true}
			cfg.ErrorHandling.Aggregator = config.ObjectIdent{mypkg, "ErrorList", true}
			if err := xtypes.Globals.Init(cfg, &AST); err != nil {
				t.Fatal(err)
			}
		},
		"global/04_local_has_priority_over_global": func(t *testing.T) {
			var cfg config.Config
			cfg.ErrorHandling.Constructor = config.ObjectIdent{mypkg, "NewError", true}
			cfg.ErrorHandling.Aggregator = config.ObjectIdent{mypkg, "ErrorList", true}
			if err := xtypes.Globals.Init(cfg, &AST); err != nil {
				t.Fatal(err)
			}
		},
	}

	cfg := loadConfig("testdata/config.yaml")
	if err := rules.InitSpecs(cfg, &AST); err != nil {
		t.Fatal(err)
	}

	fkCfg := &config.FieldKeyConfig{
		Join:      config.Bool{Value: true, IsSet: true},
		Separator: config.String{Value: ".", IsSet: true},
	}

	for _, filename := range tests {
		t.Run(filename, func(t *testing.T) {
			fileprefix := "testdata/" + filename
			f, pkg, err := getFile(pkgs, fileprefix+"_in.go")
			if err != nil {
				t.Fatal(err)
			}

			if init, ok := globals[filename]; ok {
				init(t)
				defer xtypes.Globals.Unset()
			}

			infos := make([]*rules.Info, len(f.Matches))
			for k, match := range f.Matches {
				info := new(rules.Info)
				checker := rules.NewChecker(&AST, pkg.Pkg(), fkCfg, info)
				if err := checker.Check(match); err != nil {
					t.Fatal(err)
				}
				infos[k] = info
			}

			code, err := Generate(pkg.Pkg(), infos)
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

// helper method...
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
