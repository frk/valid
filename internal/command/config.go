package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	// The directory in which the tool will search for files to process.
	// If not provided, the current working directory will be used by default.
	WorkingDirectory String `json:"working_directory"`
	// If set to true, the tool will search the hierarchy of the working
	// directory for files to process.
	Recursive Bool `json:"recursive"`
	// List of files to be used as input for the tool.
	// The files must be located in the working directory.
	InputFiles StringSlice `json:"input_files"`
	// List of regular expressions to match input files that the tool should
	// process. The regular expressions must match files that are located in
	// the working directory.
	InputFileRegexps StringSlice `json:"input_file_regexps"`
	// The format used for generating the name of the output files.
	//
	// The format can contain one (and only one) "%s" placeholder which the
	// tool will replace with the input file's base name, if no placeholder is
	// present then the input file's base name will be prefixed to the format.
	//
	// If not provided, the format "%s_isvalid.go" will be used by default.
	OutputFileNameFormat String `json:"output_file_name_format"`
	// If set to a non-empty string, it specifies the struct tag to be used
	// for constructing the field keys that will be used by the generator for
	// error reporting. A valid tag must begin with a letter (A-z) or an
	// underscore (_), subsequent characters in the tag can be letters,
	// underscores, and digits (0-9). If set to "" (empty string), the generator
	// will default to use the field names instead of struct tags to construct
	// the field keys.
	//
	// If not provided, the tag "json" will be used by default.
	FieldKeyTag String `json:"field_key_tag"`
	// If set, instructs the generator to use only the base of a tag/field
	// chain to construct the field keys.
	//
	// If not provided, `false` will be used by default.
	FieldKeyBase Bool `json:"field_key_base"`
	// The separator to be used to join a chain of tag/field values for
	// constructing the field keys. The separator can be at most one byte long.
	//
	// If not provided, the separator "." will be used by default.
	FieldKeySeparator String `json:"field_key_separator"`

	// TODO add documentation
	CustomRules []*RuleConfig `json:"custom_rules"`

	// holds the compiled expressions of the InputFileRegexps slice.
	compiledInputFileRegexps []*regexp.Regexp
}

type RuleConfig struct {
	Name string `json:"name"`
	Func string `json:"func"`

	funcPkg  string `json:"-"`
	funcName string `json:"-"`
}

var DefaultConfig = Config{
	WorkingDirectory:     String{Value: "."},
	Recursive:            Bool{Value: false},
	InputFiles:           StringSlice{},
	InputFileRegexps:     StringSlice{},
	OutputFileNameFormat: String{Value: "%s_isvalid.go"},
	FieldKeyTag:          String{Value: "json"},
	FieldKeyBase:         Bool{Value: false},
	FieldKeySeparator:    String{Value: "."},
}

// ParseFlags unmarshals the cli flags into the receiver.
func (c *Config) ParseFlags() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = printUsage
	fs.Var(&c.WorkingDirectory, "wd", "")
	fs.Var(&c.Recursive, "r", "")
	fs.Var(&c.InputFiles, "f", "")
	fs.Var(&c.InputFileRegexps, "rx", "")
	fs.Var(&c.OutputFileNameFormat, "o", "")
	fs.Var(&c.FieldKeyTag, "fktag", "")
	fs.Var(&c.FieldKeyBase, "fkbase", "")
	fs.Var(&c.FieldKeySeparator, "fksep", "")
	_ = fs.Parse(os.Args[1:])
}

// ParseFile looks for an isvalid config file in the git project's root of the receiver's
// working directory, if it finds such a file it will then unmarshal it into the receiver.
func (c *Config) ParseFile() error {
	dir, err := filepath.Abs(c.WorkingDirectory.Value)
	if err != nil {
		return err
	}

	var isRoot bool
	var confName string
	for len(dir) > 1 && dir[0] == '/' {
		isRoot, confName, err = examineDir(dir)
		if err != nil {
			return err
		}
		if isRoot {
			break
		}
		dir = filepath.Dir(dir) // parent dir will be examined next
	}

	// if found, unamrshal the config file
	if confName != "" {
		confpath := filepath.Join(dir, confName)
		f, err := os.Open(confpath)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := json.NewDecoder(f).Decode(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) FileFilterFunc() (filter func(filePath string) bool) {
	if len(c.InputFiles.Value) == 0 && len(c.InputFileRegexps.Value) == 0 {
		return nil
	}

	// copy file paths for matching
	allowFilePaths := make([]string, len(c.InputFiles.Value))
	copy(allowFilePaths, c.InputFiles.Value)

	// copy regular expressions for matching
	allowRegexps := make([]*regexp.Regexp, len(c.compiledInputFileRegexps))
	copy(allowRegexps, c.compiledInputFileRegexps)

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

// validate checks the config for errors and updates some of the values to a more "normalized" format.
func (c *Config) validate() (err error) {
	// check that the working directory can be openned
	f, err := os.Open(c.WorkingDirectory.Value)
	if err != nil {
		return fmt.Errorf("failed to open working directory: %q -- %v", c.WorkingDirectory.Value, err)
	}
	f.Close()

	// update file paths to absolutes
	for i, fp := range c.InputFiles.Value {
		abs, err := filepath.Abs(fp)
		if err != nil {
			return fmt.Errorf("error resolving absolute path of file: %q -- %v", fp, err)
		}
		c.InputFiles.Value[i] = abs
	}

	// compile the input file regexeps
	c.compiledInputFileRegexps = make([]*regexp.Regexp, len(c.InputFileRegexps.Value))
	for i, expr := range c.InputFileRegexps.Value {
		rx, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("error compiling regular expression: %q -- %v", expr, err)
		}
		c.compiledInputFileRegexps[i] = rx
	}

	// check that the output filename format contains at most one "%" and
	// that it is followed by an "s" to form the "%s" verb
	if n := strings.Count(c.OutputFileNameFormat.Value, "%"); n == 0 {
		// modify the output filename format
		c.OutputFileNameFormat.Value = "%s" + c.OutputFileNameFormat.Value
	} else if n > 1 || (n == 1 && !strings.Contains(c.OutputFileNameFormat.Value, "%s")) {
		return fmt.Errorf("bad output filename format: %q", c.OutputFileNameFormat.Value)
	}

	// check the field key configuration for errors
	rxFCKTag := regexp.MustCompile(`^(?:[A-Za-z_]\w*)?$`)
	if !rxFCKTag.MatchString(c.FieldKeyTag.Value) {
		return fmt.Errorf("bad field key tag: %q", c.FieldKeyTag.Value)
	}
	if len(c.FieldKeySeparator.Value) > 1 {
		return fmt.Errorf("bad field key separator: %q", c.FieldKeySeparator.Value)
	}

	// check custom rules
	ruleNameMap := make(map[string]struct{}) // to ensure uniqueness
	for _, rc := range c.CustomRules {
		if rc == nil {
			continue
		}

		rc.Name = strings.TrimSpace(rc.Name)
		if len(rc.Name) == 0 {
			return fmt.Errorf("missing custom rule name" /*TODO more helpful message */)
		}
		rc.Func = strings.TrimSpace(rc.Func)
		if len(rc.Func) == 0 {
			return fmt.Errorf("missing custom rule func" /*TODO more helpful message */)
		}

		if _, ok := ruleNameMap[rc.Name]; ok {
			return fmt.Errorf("duplicate custom rule name" /*TODO more helpful message */)
		}
		ruleNameMap[rc.Name] = struct{}{}

		// split function name from package
		if i := strings.LastIndex(rc.Func, "."); i < 0 {
			return fmt.Errorf("bad custom rule func format: %q", rc.Func)
		} else {
			rc.funcPkg, rc.funcName = rc.Func[:i], rc.Func[i+1:]
		}

		// make sure the function is actually exported
		if !token.IsExported(rc.funcName) {
			return fmt.Errorf("only exported functions can be used "+
				"for custom rules, %q is unexported", rc.Func)
		}
	}

	return nil
}

// examineDir reports if the directory at the given path is the root directory
// of a git project and if it is, it will also report the name of the isvalid config
// file (either ".isvalid" or ".isvalid.json") if such a file exists in that root
// directory, otherwise confName will be empty.
func examineDir(path string) (isRoot bool, confName string, err error) {
	d, err := os.Open(path)
	if err != nil {
		return false, "", err
	}
	defer d.Close()

	infoList, err := d.Readdir(-1)
	if err != nil {
		return false, "", err
	}

	for _, info := range infoList {
		name := info.Name()
		if name == ".git" && info.IsDir() {
			isRoot = true
		}
		if (name == ".isvalid" || name == ".isvalid.json") && !info.IsDir() {
			confName = name
		}
	}

	// NOTE(mkopriva): currently we don't care about .isvalid files that live outside
	// of the git project root directory, if, in the future, the rules are expanded
	// then this will need to be either removed or accordingly updated.
	if !isRoot {
		confName = ""
	}
	return isRoot, confName, nil
}

// String implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type String struct {
	Value string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (s String) Get() interface{} {
	return s.Value
}

// String implements the flag.Value interface.
func (s String) String() string {
	return s.Value
}

// Set implements the flag.Value interface.
func (s *String) Set(value string) error {
	s.Value = value
	s.IsSet = true
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *String) UnmarshalJSON(data []byte) error {
	if !s.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		s.Value = value
		s.IsSet = true
	}
	return nil
}

// Bool implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type Bool struct {
	Value bool
	IsSet bool
}

// IsBoolFlag indicates that the Bool type can be used as a boolean flag.
func (b Bool) IsBoolFlag() bool {
	return true
}

// Get implements the flag.Getter interface.
func (b Bool) Get() interface{} {
	return b.String()
}

// String implements the flag.Value interface.
func (b Bool) String() string {
	return strconv.FormatBool(b.Value)
}

// Set implements the flag.Value interface.
func (b *Bool) Set(value string) error {
	if len(value) > 0 {
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		b.Value = v
		b.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (b *Bool) UnmarshalJSON(data []byte) error {
	if !b.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value bool
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		b.Value = value
		b.IsSet = true
	}
	return nil
}

// StringSlice implements both the flag.Value and the json.Unmarshal interfaces
// enforcing priority of flags over json, meaning that json.Unmarshal will
// not override the value if it was previously set by flag.Var.
type StringSlice struct {
	Value []string
	IsSet bool
}

// Get implements the flag.Getter interface.
func (ss StringSlice) Get() interface{} {
	return ss.String()
}

// String implements the flag.Value interface.
func (ss StringSlice) String() string {
	return strings.Join(ss.Value, ",")
}

// Set implements the flag.Value interface.
func (ss *StringSlice) Set(value string) error {
	if len(value) > 0 {
		ss.Value = append(ss.Value, value)
		ss.IsSet = true
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ss *StringSlice) UnmarshalJSON(data []byte) error {
	if !ss.IsSet {
		if len(data) == 0 || string(data) == `null` {
			return nil
		}

		var value []string
		if err := json.Unmarshal(data, &value); err != nil {
			return err
		}
		if len(value) > 0 {
			ss.Value = value
			ss.IsSet = true
		}
	}
	return nil
}
