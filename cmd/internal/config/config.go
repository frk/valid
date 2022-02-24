package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// The file from which the config was decoded. Used for error reporting.
	File String `yaml:"-"`
	// The directory in which the tool will search for files to process.
	// If not provided, the current working directory will be used by default.
	WorkDir String `yaml:"working_directory"`
	// If set to true, the tool will search the hierarchy of the working
	// directory for files to process.
	Recursive Bool `yaml:"recursive"`
	// List of files to be used as input for the tool.
	// The files must be located in the working directory.
	FileList StringSlice `yaml:"file_list"`
	// List of regular expressions to match input files that the tool
	// should process. The regular expressions must match files that
	// are located in the working directory.
	FilePatternList StringSlice `yaml:"file_pattern_list"`
	// A string containing the format to be used for generating
	// the name of the output files. The format string MUST end
	// with the ".go" file extension.
	//
	// Inside the format string the percent sign ("%") can be used as
	// the placeholder to be replaced by the tool with the base name of
	// the source file. If no placeholder is present in the format string
	// then the base name of the source file will simply be prefixed
	// to the format string.
	//
	// If not provided, the format will default to "%_valid.go".
	OutNameFormat String `yaml:"out_name_format"`
	// A string containing a regular expression that will be used by the tool
	// to identify the struct types for which to generate the validation code.
	//
	// If not provided, the pattern will default to "^(?i:\w*Validator)$".
	ValidatorNamePattern String `yaml:"validator_name_pattern"`
	// Configures the code generation of the handling of validation errors.
	ErrorHandling ErrorHandlingConfig `yaml:"error_handling"`
	// List of custom rules to be made available to the tool.
	Rules []RuleConfig `yaml:"rules"`

	// The compiled expressions of the FilePatternList slice.
	fileRegexpList []*regexp.Regexp
	// The compiled expression of the ValidatorNamePattern string.
	validatorNameRegexp *regexp.Regexp
}

type ErrorHandlingConfig struct {
	// Configures how field keys should be constructed.
	//
	// Field keys are used for error reporting by the generated code.
	// When a field fails validation the field's key, along with other
	// details, will be passed as an argument to the client's
	// implementation of the error handling code.
	FieldKey FieldKeyConfig `yaml:"field_key"`
	// The identifier of a function that the generated code should
	// use for constructing custom, application-specific errors.
	// The function's signature MUST be the following:
	//
	//     func(key string, val any, rule string, args ...any) error
	//
	Constructor ObjectIdent `yaml:"constructor"`
	// The identifier of a type that the generated code should use
	// for constructing & aggregation of custom, application-specific
	// errors. The type MUST implement the following interface:
	//
	//     interface {
	//         Error(key string, val any, rule string, args ...any)
	//         Out() error
	//     }
	//
	Aggregator ObjectIdent `yaml:"aggregator"`
}

type FieldKeyConfig struct {
	// If non-empty string, specifies the struct tag whose value will be
	// used for constructing the field keys. If explicitly set to an empty
	// string, the generator will default to use the fields' names for
	// constructing the field keys.
	//
	// A valid tag must begin with a letter (A-z) or an underscore (_),
	// subsequent characters in the tag can be letters, underscores,
	// and digits (0-9).
	//
	// If not provided, the tag "json" will be used by default.
	Tag String `yaml:"tag"`
	// When set to true, a nested struct field's key will be constructed by
	// joining it together with all of its parent fields. When false, a nested
	// struct field's key will be constructed only from that field's tag or name.
	//
	// If not provided, `true` will be used by default.
	Join Bool `yaml:"join"`
	// The separator to be used for joining fields' tags or names
	// when constructing field keys.
	//
	// The separator MUST be a single, one byte long, character.
	//
	// If not provided, the separator "." will be used by default.
	Separator String `yaml:"separator"`
}

type RuleConfig struct {
	// The function's name qualified with
	// the path of the function's package.
	Func ObjectIdent `yaml:"func"`
	// The spec for the function's rule. Optional if
	// the func's doc already has a valid config.
	Rule *RuleSpec `yaml:"rule"`
}

type RuleSpec struct {
	// The name of the rule.
	Name string `yaml:"name"`
	// The configuration for the rule's arguments.
	//
	// NOTE: If args is NOT empty, then it MUST contain a number
	// of elements that is compatible with the parameters of the
	// rule's function.
	Args []RuleArgConfig `yaml:"args"`
	// ArgMin can be used to enforce the number of arguments
	// that the rule must accept.
	//
	// - ArgMin is optional, if omitted its value will
	//   be inferred from the function's signature.
	// - When ArgMin is provided it will be used to enforce
	//   the valid number of arguments for variadic functions.
	// - When ArgMin is provided it MUST be compatible
	//   with the function's signature, if it isn't then
	//   the tool will exit with an error.
	ArgMin *uint `yaml:"arg_min"`
	// ArgMax can be used to override the upper limit of the number
	// of arguments that the rule should be allowed to accept.
	// A negative ArgMax can be used to indicate that there's no
	// upper limit to the number of arguments.
	//
	// - ArgMax is optional, if omitted its value will
	//   be inferred from the function's signature.
	// - When ArgMax is provided it will be used to enforce
	//   the valid number of arguments for variadic functions.
	// - When ArgMax is provided it MUST be compatible
	//   with the function's signature, if it isn't then
	//   the tool will exit with an error.
	ArgMax *int `yaml:"arg_max"`
	// The configuration for the error that should be generated for the rule.
	Error RuleErrorConfig `yaml:"error"`
	// The join operator that should be used to join
	// multiple instances of the rule into a single one.
	//
	// The value MUST be one of: "AND", "OR", or "NOT" (case insensitive).
	JoinOp JoinOp `yaml:"join_op"`
}

type RuleArgConfig struct {
	// The rule argument's default value. If nil, then the
	// rule argument's value MUST be provided in the struct tag.
	//
	// If not nil, the value must be a scalar value.
	Default *Scalar `yaml:"default"`
	// If options is empty, then ANY value can be provided for
	// the argument in the rule's struct tag.
	//
	// If options is NOT empty, then it is considered to represent, together
	// with the default value, the *complete* set of valid values that can be
	// provided as the argument in the rule's struct tag.
	Options []RuleArgOption `yaml:"options"`
}

type RuleArgOption struct {
	// Value specifies the value that the generator should supply
	// as the rule's argument in the generated code.
	Value Scalar `yaml:"value"`
	// Alias is an alternative identifier of the argument's value that
	// can be used within the rule's struct tag. This field is optional.
	Alias string `yaml:"alias"`
}

type RuleErrorConfig struct {
	// The text of the error message.
	Text string `yaml:"text,omitempty"`
	// If true the generated error message
	// will include the rule's arguments.
	WithArgs bool `yaml:"with_args,omitempty"`
	// The separator used to join the rule's
	// arguments for the error message.
	ArgSep string `yaml:"arg_sep,omitempty"`
	// The text to be appended after the list of arguments.
	ArgSuffix string `yaml:"arg_suffix,omitempty"`
}

// MergeAndCheck merges the receiver with a config file, if available.
func (c *Config) MergeAndCheck() error {
	if err := c.resolveFile(); err != nil {
		return &Error{C: ERR_CONFIG_FILE, file: c.File.Value, err: err}
	}
	if c.File.IsSet && len(c.File.Value) > 0 {
		if err := DecodeFile(c.File.Value, c); err != nil {
			return err
		}
	}

	dc := DefaultConfig()
	if !c.WorkDir.IsSet {
		c.WorkDir.Value = dc.WorkDir.Value
	}
	if !c.Recursive.IsSet {
		c.Recursive.Value = dc.Recursive.Value
	}
	if !c.OutNameFormat.IsSet {
		c.OutNameFormat.Value = dc.OutNameFormat.Value
	}
	if !c.ValidatorNamePattern.IsSet {
		c.ValidatorNamePattern.Value = dc.ValidatorNamePattern.Value
	}
	if !c.ErrorHandling.FieldKey.Tag.IsSet {
		c.ErrorHandling.FieldKey.Tag.Value = dc.ErrorHandling.FieldKey.Tag.Value
	}
	if !c.ErrorHandling.FieldKey.Join.IsSet {
		c.ErrorHandling.FieldKey.Join.Value = dc.ErrorHandling.FieldKey.Join.Value
	}
	if !c.ErrorHandling.FieldKey.Separator.IsSet {
		c.ErrorHandling.FieldKey.Separator.Value = dc.ErrorHandling.FieldKey.Separator.Value
	}

	return c.normalizeAndCheck()
}

func (c *Config) resolveFile() error {
	if c.File.IsSet && len(c.File.Value) > 0 {
		file, err := filepath.Abs(c.File.Value)
		if err != nil {
			return err
		}
		c.File.Value = file
		c.File.IsSet = true
		return nil
	}

	dir, err := filepath.Abs(c.WorkDir.Value)
	if err != nil {
		return err
	}

	var hasDotGit, hasConfig bool
	for len(dir) > 1 && dir[0] == '/' {
		hasDotGit, hasConfig, err = examineDir(dir)
		if err != nil {
			return err
		}
		if hasDotGit {
			break
		}

		// examine parent dir next
		dir = filepath.Dir(dir)
	}
	if hasConfig {
		file := filepath.Join(dir, ".valid.yaml")
		c.File.Value = file
		c.File.IsSet = true
	}
	return nil
}

var rxFKTag = regexp.MustCompile(`^(?:[A-Za-z_]\w*)?$`)

func (c *Config) normalizeAndCheck() (err error) {
	// update wd in case it wasn't abs
	wd, err := filepath.Abs(c.WorkDir.Value)
	if err != nil {
		return &Error{C: ERR_WORK_DIR, dir: c.WorkDir.Value, file: c.File.Value, err: err}
	}
	c.WorkDir.Value = wd
	c.WorkDir.IsSet = true

	// check that wd can be openned
	f, err := os.Open(c.WorkDir.Value)
	if err != nil {
		return &Error{C: ERR_WORK_DIR, dir: c.WorkDir.Value, file: c.File.Value, err: err}
	}
	f.Close()

	// update file paths to absolutes & make sure their are inside wd
	for i, f := range c.FileList.Value {
		file, err := filepath.Abs(f)
		if err != nil {
			return &Error{C: ERR_FILE_ITEM, dir: c.WorkDir.Value,
				file: c.File.Value, val: f, err: err}
		}
		if !strings.HasPrefix(file, c.WorkDir.Value) {
			err := fmt.Errorf("file not in working directory")
			return &Error{C: ERR_FILE_ITEM, dir: c.WorkDir.Value,
				file: c.File.Value, val: f, err: err}
		}
		c.FileList.Value[i] = file
	}

	// compile the input file regexeps
	if len(c.FilePatternList.Value) > 0 {
		c.fileRegexpList = make([]*regexp.Regexp, len(c.FilePatternList.Value))
		for i, expr := range c.FilePatternList.Value {
			rx, err := regexp.Compile(expr)
			if err != nil {
				return &Error{C: ERR_PATTERN, dir: c.WorkDir.Value,
					file: c.File.Value, key: "file_pattern",
					val: expr, err: err}
			}
			c.fileRegexpList[i] = rx
		}
	}

	// normalize & validate the output filename format
	if !strings.Contains(c.OutNameFormat.Value, "%") {
		c.OutNameFormat.Value = "%" + c.OutNameFormat.Value
	}
	if !strings.HasSuffix(c.OutNameFormat.Value, ".go") {
		return &Error{C: ERR_OUTNAME_FORMAT, dir: c.WorkDir.Value,
			file: c.File.Value, key: "out_name_format", val: c.OutNameFormat.Value}
	}

	// compile validator name regexp
	expr := c.ValidatorNamePattern.Value
	rx, err := regexp.Compile(expr)
	if err != nil {
		return &Error{C: ERR_PATTERN, dir: c.WorkDir.Value,
			file: c.File.Value, key: "validator_name_pattern",
			val: expr, err: err}
	}
	c.validatorNameRegexp = rx

	// check the field key configuration for errors
	if val := c.ErrorHandling.FieldKey.Tag.Value; !rxFKTag.MatchString(val) {
		return &Error{C: ERR_FKEY_TAG, dir: c.WorkDir.Value,
			file: c.File.Value, key: "field_key.tag",
			val: val, err: err}
	}
	if val := c.ErrorHandling.FieldKey.Separator.Value; len(val) > 1 {
		return &Error{C: ERR_FKEY_SEP, dir: c.WorkDir.Value,
			file: c.File.Value, key: "field_key.separator",
			val: val, err: err}
	}

	// check custom rules
	seen := make(map[string]bool) // to ensure uniqueness
	for i, rc := range c.Rules {
		if rc.Func.Name == "" {
			return &Error{C: ERR_RULE_NOFUNC, dir: c.WorkDir.Value,
				file: c.File.Value, val: strconv.Itoa(i)}
		}
		ruleKey := rc.Func.String()

		if rc.Rule != nil {
			if rc.Rule.Name == "" {
				return &Error{C: ERR_RULE_NONAME, dir: c.WorkDir.Value,
					file: c.File.Value, val: strconv.Itoa(i)}
			}
			ruleKey = rc.Rule.Name
		}

		// make sure rules are unique
		if seen[ruleKey] {
			return &Error{C: ERR_RULE_DUPNAME, dir: c.WorkDir.Value,
				file: c.File.Value, val: ruleKey}
		}
		seen[ruleKey] = true
	}

	return nil
}

func (c *Config) decodeFile() error {
	f, err := os.Open(c.File.Value)
	if err != nil {
		return &Error{C: ERR_CONFIG_FILE, file: c.File.Value, err: err}
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		if _, ok := err.(*Error); ok {
			return err
		}
		return &Error{C: ERR_YAML_FILE, file: c.File.Value, err: err}
	}
	return nil
}

func (c *Config) ValidatorRegexp() *regexp.Regexp {
	re := *c.validatorNameRegexp
	return &re
}

func (c *Config) FileFilterFunc() (filter func(filePath string) bool) {
	if len(c.FileList.Value) == 0 && len(c.FilePatternList.Value) == 0 {
		return nil
	}

	// copy file paths for matching
	allowFilePaths := make([]string, len(c.FileList.Value))
	copy(allowFilePaths, c.FileList.Value)

	// copy regular expressions for matching
	allowRegexps := make([]*regexp.Regexp, len(c.fileRegexpList))
	copy(allowRegexps, c.fileRegexpList)

	return func(filePath string) bool {
		for _, fp := range allowFilePaths {
			if filePath == fp {
				return true
			}
		}
		for _, rx := range allowRegexps {
			if rx.MatchString(filePath) {
				return true
			}
		}
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////
// helpers
////////////////////////////////////////////////////////////////////////////////

func DefaultConfig() Config {
	return Config{
		WorkDir:              String{Value: "."},
		Recursive:            Bool{Value: false},
		FileList:             StringSlice{},
		FilePatternList:      StringSlice{},
		OutNameFormat:        String{Value: "%_valid.go"},
		ValidatorNamePattern: String{Value: `^(?i:\w*Validator)$`},
		ErrorHandling: ErrorHandlingConfig{
			FieldKey: FieldKeyConfig{
				Tag:       String{Value: "json"},
				Join:      Bool{Value: true},
				Separator: String{Value: "."},
			},
		},
	}
}

// examineDir reports whether or not the directory at the given
// path is the root directory of a git project and whether or not
// it also contains a .valid.yaml config file.
func examineDir(path string) (hasDotGit, hasConfig bool, err error) {
	d, err := os.Open(path)
	if err != nil {
		return false, false, err
	}
	defer d.Close()

	infoList, err := d.Readdir(-1)
	if err != nil {
		return false, false, err
	}

	for _, info := range infoList {
		name := info.Name()
		if name == ".git" && info.IsDir() {
			hasDotGit = true
		}
		if name == ".valid.yaml" && !info.IsDir() {
			hasConfig = true
		}
	}

	return hasDotGit, hasConfig, nil
}

func DecodeFile(file string, c *Config) error {
	f, err := os.Open(file)
	if err != nil {
		return &Error{C: ERR_CONFIG_FILE,
			dir:  c.WorkDir.Value,
			file: file,
			err:  err}
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		if err, ok := err.(*Error); ok {
			err.dir = c.WorkDir.Value
			err.file = file
			return err
		}
		return &Error{C: ERR_YAML_FILE,
			dir:  c.WorkDir.Value,
			file: file,
			err:  err}
	}

	c.File.Value = file
	c.File.IsSet = true
	return nil
}
