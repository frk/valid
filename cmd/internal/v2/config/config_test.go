package config

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/frk/compare"
	"gopkg.in/yaml.v3"
)

func _ptr[T any](v T) *T { return &v }

func Test_Config_MergeAndCheck(t *testing.T) {
	_err_ := fmt.Errorf("dummy error")

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dc := DefaultConfig()

	tests := []struct {
		c    string
		wd   string
		want *Config
		err  error
		show bool
	}{{
		c:   "testdata/config-file-does-not-exist.yaml",
		err: &Error{C: ERR_CONFIG_FILE, file: wd + "/testdata/config-file-does-not-exist.yaml", err: _err_},
	}, {
		wd:  "testdata/working-dir-does-not-exist/",
		err: &Error{C: ERR_CONFIG_FILE, file: "", err: _err_},
	}, {
		c:   "testdata/bad_config_test/empty_config.yaml",
		err: &Error{C: ERR_YAML_FILE, file: wd + "/testdata/bad_config_test/empty_config.yaml", err: _err_},
	}, {
		c:   "testdata/bad_config_test/not_yaml_config.xml",
		err: &Error{C: ERR_YAML_FILE, file: wd + "/testdata/bad_config_test/not_yaml_config.xml", err: _err_},
	}, {
		c:   "testdata/bad_config_test/bad_syntax_config.yaml",
		err: &Error{C: ERR_YAML_FILE, file: wd + "/testdata/bad_config_test/bad_syntax_config.yaml", err: _err_},
	}, {
		c: "testdata/bad_config_test/bad_yaml_type.yaml",
		err: &Error{C: ERR_YAML_ERROR, file: wd + "/testdata/bad_config_test/bad_yaml_type.yaml",
			tt: &Bool{}, node: &yaml.Node{}, err: &yaml.TypeError{}},
	}, {
		c:   "testdata/bad_config_test/bad_wd_config.yaml",
		err: &Error{C: ERR_WORK_DIR, dir: "/foo/bar/baz", file: wd + "/testdata/bad_config_test/bad_wd_config.yaml", err: _err_},
	}, {
		c: "testdata/bad_config_test/file_list_item_not_in_wd.yaml",
		err: &Error{C: ERR_FILE_ITEM,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/file_list_item_not_in_wd.yaml",
			val:  "/foo/bar/file1.go", err: _err_},
	}, {
		c: "testdata/bad_config_test/file_pattern_list_item_not_valid_regexp.yaml",
		err: &Error{C: ERR_PATTERN,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/file_pattern_list_item_not_valid_regexp.yaml",
			key:  "file_pattern", val: "[a-z)", err: _err_},
	}, {
		c: "testdata/bad_config_test/out_name_format_not_valid.yaml",
		err: &Error{C: ERR_OUTNAME_FORMAT,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/out_name_format_not_valid.yaml",
			key:  "out_name_format", val: "%_out"},
	}, {
		c: "testdata/bad_config_test/validator_name_pattern_not_valid_regexp.yaml",
		err: &Error{C: ERR_PATTERN,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/validator_name_pattern_not_valid_regexp.yaml",
			key:  "validator_name_pattern", val: "[a-z)", err: _err_},
	}, {
		c: "testdata/bad_config_test/bad_field_key_tag.yaml",
		err: &Error{C: ERR_FKEY_TAG,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/bad_field_key_tag.yaml",
			key:  "field_key.tag", val: "foo bar"},
	}, {
		c: "testdata/bad_config_test/bad_field_key_separator.yaml",
		err: &Error{C: ERR_FKEY_SEP,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/bad_field_key_separator.yaml",
			key:  "field_key.separator", val: "..."},
	}, {
		c: "testdata/bad_config_test/rule_with_no_name.yaml",
		err: &Error{C: ERR_RULE_NONAME,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/rule_with_no_name.yaml",
			val:  "1"},
	}, {
		c: "testdata/bad_config_test/rule_with_no_func.yaml",
		err: &Error{C: ERR_RULE_NOFUNC,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/rule_with_no_func.yaml",
			val:  "0"},
	}, {
		c: "testdata/bad_config_test/rule_with_duplicate_name.yaml",
		err: &Error{C: ERR_RULE_DUPNAME,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/rule_with_duplicate_name.yaml",
			val:  "bar"},
	}, {
		c: "testdata/bad_config_test/type_with_no_name.yaml",
		err: &Error{C: ERR_TYPE_NONAME,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/type_with_no_name.yaml",
			val:  "2"},
	}, {
		c: "testdata/bad_config_test/type_with_duplicate_name.yaml",
		err: &Error{C: ERR_TYPE_DUPNAME,
			dir:  wd + "/testdata",
			file: wd + "/testdata/bad_config_test/type_with_duplicate_name.yaml",
			val:  "example.com/test.Test1"},
	}, {
		wd: "testdata/no_config_test/",
		want: &Config{
			WorkDir:              String{Value: wd + "/testdata/no_config_test", IsSet: true},
			OutNameFormat:        dc.OutNameFormat,
			ValidatorNamePattern: dc.ValidatorNamePattern,
			ErrorHandling:        dc.ErrorHandling,
			validatorNameRegexp:  regexp.MustCompile(dc.ValidatorNamePattern.Value),
		},
	}, {
		wd: "testdata/implicit_config_test/foo/bar/baz",
		want: &Config{
			File:                 String{Value: wd + "/testdata/implicit_config_test/.valid.yaml", IsSet: true},
			WorkDir:              String{Value: wd + "/testdata/implicit_config_test/foo/bar/baz", IsSet: true},
			OutNameFormat:        dc.OutNameFormat,
			ValidatorNamePattern: dc.ValidatorNamePattern,
			ErrorHandling:        dc.ErrorHandling,
			validatorNameRegexp:  regexp.MustCompile(dc.ValidatorNamePattern.Value),
			Rules: []FuncConfig{{
				Func: ObjectIdent{Pkg: "mod/pkg", Name: "Test", IsSet: true},
				Rule: &RuleConfig{Name: "test"},
			}},
		},
	}, {
		c: "testdata/good_config_test/explicit_config.yaml",
		want: &Config{
			File:                 String{Value: wd + "/testdata/good_config_test/explicit_config.yaml", IsSet: true},
			WorkDir:              String{Value: wd, IsSet: true},
			OutNameFormat:        dc.OutNameFormat,
			ValidatorNamePattern: dc.ValidatorNamePattern,
			ErrorHandling:        dc.ErrorHandling,
			validatorNameRegexp:  regexp.MustCompile(dc.ValidatorNamePattern.Value),
		},
	}, {
		c: "testdata/good_config_test/full_config.yaml",
		want: &Config{
			File:      String{Value: wd + "/testdata/good_config_test/full_config.yaml", IsSet: true},
			WorkDir:   String{Value: wd + "/testdata", IsSet: true},
			Recursive: Bool{Value: true, IsSet: true},
			FileList: StringSlice{
				Value: []string{
					wd + "/testdata/good_config_test/foo/bar/baz/file1.go",
					wd + "/testdata/good_config_test/foo/bar/baz/file_2.go",
				},
				IsSet: true,
			},
			FilePatternList: StringSlice{
				Value: []string{
					"^\\/path\\/to\\/my\\/\\w+_foo.go$",
					"^\\/path\\/to\\/my\\/\\w+_bar.go$",
				},
				IsSet: true,
			},
			OutNameFormat:        String{Value: "%_out.go", IsSet: true},
			ValidatorNamePattern: String{Value: "^\\w+Input$", IsSet: true},
			ErrorHandling: ErrorHandlingConfig{
				FieldKey: FieldKeyConfig{
					Tag:       String{Value: "json", IsSet: true},
					Join:      Bool{Value: true, IsSet: true},
					Separator: String{Value: ".", IsSet: true},
				},
				Constructor: ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "NewError",
					IsSet: true,
				},
				Aggregator: ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "MyErrorAggregator",
					IsSet: true,
				},
			},
			Rules: []FuncConfig{{
				Func: ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "IsFoobar",
					IsSet: true,
				},
				Rule: &RuleConfig{
					Name: "foobar",
					Args: []RuleArg{{
						Default: &Scalar{Type: NIL},
						Options: []RuleArgOption{{
							Value: Scalar{Type: INT, Value: "123"},
							Alias: "x",
						}},
					}},
					ArgMin: _ptr[uint](1),
					ArgMax: _ptr[int](2),
					Error: RuleErrMesg{
						Text:      "invalid foobar",
						WithArgs:  true,
						ArgSep:    ", ",
						ArgSuffix: " (bazzz)",
					},
					JoinOp: JOIN_OR,
				},
			}},
			Types: []TypeConfig{{
				Type: ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "MyType1",
					IsSet: true,
				},
				Value: TypeValue{
					Get: String{Value: "Value", IsSet: true},
					Set: String{Value: "Value", IsSet: true},
				},
				RequiredCheck: String{Value: "IsZero", IsSet: true},
				NotnilCheck:   String{Value: "IsNil", IsSet: true},
				OptionalCheck: String{Value: "IsNotZero", IsSet: true},
				OmitnilCheck:  String{Value: "IsNotNil", IsSet: true},
			}, {
				Type: ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "MyType2",
					IsSet: true,
				},
				Value: TypeValue{
					Get: String{Value: "GetValue", IsSet: true},
					Set: String{Value: "SetValue", IsSet: true},
				},
				RequiredCheck: String{Value: "IsZero", IsSet: true},
				NotnilCheck:   String{Value: "IsNil", IsSet: true},
				OptionalCheck: String{Value: "IsNotZero", IsSet: true},
				OmitnilCheck:  String{Value: "IsNotNil", IsSet: true},
			}},
			fileRegexpList: []*regexp.Regexp{
				regexp.MustCompile("^\\/path\\/to\\/my\\/\\w+_foo.go$"),
				regexp.MustCompile("^\\/path\\/to\\/my\\/\\w+_bar.go$"),
			},
			validatorNameRegexp: regexp.MustCompile("^\\w+Input$"),
		},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		name := tt.c
		if name == "" {
			name = tt.wd
		}

		t.Run(name, func(t *testing.T) {
			got := Config{
				File:    String{Value: tt.c, IsSet: tt.c != ""},
				WorkDir: String{Value: tt.wd, IsSet: tt.wd != ""},
			}

			err := got.MergeAndCheck()
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if tt.want != nil {
				if e := compare.Compare(got, *tt.want); e != nil {
					t.Error(e)
				}
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}
