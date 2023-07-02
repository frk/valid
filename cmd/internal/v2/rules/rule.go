package rules

import (
	"strings"

	"github.com/frk/valid/cmd/internal/v2/config"
	"github.com/frk/valid/cmd/internal/v2/source"
	"github.com/frk/valid/cmd/internal/v2/types"
)

// JoinOp represents the boolean operator that can be used
// to join multiple instances of a rule into a single one.
//
// NOTE(mkopriva): Because the generated code will be looking
// for **invalid values, as opposed to valid ones**, the actual
// expressions generated based on these operators will be the
// inverse of what their names indicate, see the comments next
// to the operators for an example.
type JoinOp uint

const (
	_        JoinOp = iota
	JOIN_NOT        // x || x || x....
	JOIN_AND        // !x || !x || !x....
	JOIN_OR         // !x && !x && !x....
)

var _configJoinOpMap = [...]JoinOp{
	config.JOIN_NOT: JOIN_NOT,
	config.JOIN_AND: JOIN_AND,
	config.JOIN_OR:  JOIN_OR,
}

// ErrMesg defines how a rule's error message should be generated.
type ErrMesg struct {
	// The text of the error message.
	Text string
	// If true the generated error message
	// will include the rule's arguments.
	WithArgs bool
	// The separator used to join the rule's
	// arguments for the error message.
	ArgSep string
	// The text to be appended after the list of arguments.
	ArgSuffix string
}

////////////////////////////////////////////////////////////////////////////////

// Rule represents a validation rule.
type Rule struct {
	// The name of the rule.
	Name string
	// The kind of the rule.
	Kind Kind
	// Func is set to the rule's function type if the rule's kind
	// is FUNCTION, METHOD, or PREPROC, otherwise it will be nil.
	Func *types.Func
	// ArgMin and ArgMax define bounds of allowed
	// number of arguments for the rule.
	ArgMin, ArgMax int
	// The rule's pre-declared argument options.
	ArgOpts []map[string]Arg
	// The join operator that should be used for joining
	// multiple instances of the rule into a single one.
	JoinOp JoinOp
	// The error message that should be generated for the rule.
	Err ErrMesg
	// The error options for specific argument combinations
	ErrOpts map[string]ErrMesg
	// Indicates that the generated code should use raw
	// strings for any string arguments of the rule.
	UseRawString bool
	// The actual arguments of the rule.
	Args []*Arg
}

// String implements the fmt.Stringer interface.
func (r Rule) String() (out string) {
	out = r.Name
	for i := range r.Args {
		out += ":"
		switch r.Args[i].Type {
		case ARG_FIELD_ABS:
			out += "&"
		case ARG_FIELD_REL:
			out += "."
		}
		out += r.Args[i].Value
	}
	return out
}

// Is reports if the rule is one of the provided kinds.
func (r Rule) Is(kinds ...Kind) bool {
	for _, k := range kinds {
		if r.Kind == k {
			return true
		}
	}
	return false
}

type List []*Rule

// Has reports if the list contains a rule of one of the provided kinds.
func (l List) Has(kinds ...Kind) bool {
	for _, r := range l {
		if r.Is(kinds...) {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////

// ArgType indicates the type of a rule argument value.
type ArgType uint

// String implements the fmt.Stringer interface.
func (t ArgType) String() string {
	if int(t) < len(_argtypestring) {
		return _argtypestring[t]
	}
	return "<invalid>"
}

const (
	ARG_UNKNOWN ArgType = iota
	ARG_FIELD_ABS
	ARG_FIELD_REL
	ARG_BOOL
	ARG_INT
	ARG_FLOAT
	ARG_STRING
)

var _argtypestring = [...]string{
	ARG_UNKNOWN:   "<unknown>",
	ARG_FIELD_ABS: "<field_abs>",
	ARG_FIELD_REL: "<field_rel>",
	ARG_BOOL:      "bool",
	ARG_INT:       "int",
	ARG_FLOAT:     "float",
	ARG_STRING:    "string",
}

var _configArgTypeMap = [...]ArgType{
	config.NIL:    ARG_UNKNOWN,
	config.BOOL:   ARG_BOOL,
	config.INT:    ARG_INT,
	config.FLOAT:  ARG_FLOAT,
	config.STRING: ARG_STRING,
}

// Arg represents a parsed argument of a rule.
type Arg struct {
	// The type of the argument.
	Type ArgType
	// The literal string representation of the value.
	Value string
}

// IsUnknown reports whether or not the rule argument is unknown.
func (a *Arg) IsUnknown() bool {
	return a.Type == ARG_UNKNOWN
}

// IsEmpty reports whether or not the rule argument is empty.
func (a *Arg) IsEmpty() bool {
	return a.Value == ""
}

// IsField reports whether or not the rule argument is a field reference.
func (a *Arg) IsField() bool {
	return a.Type == ARG_FIELD_ABS ||
		a.Type == ARG_FIELD_REL
}

// IsLiteral reports whether or not the rule argument is a literal value.
func (a *Arg) IsLiteral() bool {
	return a.Type == ARG_BOOL ||
		a.Type == ARG_INT ||
		a.Type == ARG_FLOAT ||
		a.Type == ARG_STRING
}

// IsNumeric reports whether or not the rule argument is a numeric literal.
func (a *Arg) IsNumeric() bool {
	return a.Type == ARG_INT || a.Type == ARG_FLOAT
}

// IsInt reports whether or not the rule argument is an int literal.
func (a *Arg) IsInt() bool {
	return a.Type == ARG_INT
}

// IsFloat reports whether or not the rule argument is a float literal.
func (a *Arg) IsFloat() bool {
	return a.Type == ARG_FLOAT
}

// IsUInt reports whether or not the rule argument is a uint literal.
func (a *Arg) IsUInt() bool {
	return a.Type == ARG_INT && a.Value[0] != '-'
}

////////////////////////////////////////////////////////////////////////////////

// reserved rule names
var _reserved = map[string]bool{
	"omitkey":  true, //non-rule
	"optional": true,
	"omitnil":  true,
	"required": true,
	"notnil":   true,
	"noguard":  true,
	"isvalid":  true,
	"-isvalid": true,
	"enum":     true,
}

func newRuleFromFunc(fn *source.Func, src *source.Source) (*Rule, error) {
	cfg := new(config.RuleConfig)
	if err := fn.DecodeConfig(cfg); err != nil {
		return nil, &Error{C: E_CONFIG_INVALID, s: src, fn: fn, err: err}
	}
	if cfg.Name == "" {
		return nil, &Error{C: E_CONFIG_NONAME, s: src, fn: fn, rc: cfg}
	}
	if _reserved[cfg.Name] {
		return nil, &Error{C: E_CONFIG_RESERVED, s: src, fn: fn, rc: cfg}
	}

	fun := types.AnalyzeFunc(fn.Func, src)
	join := _configJoinOpMap[cfg.JoinOp]
	kind := FUNCTION
	if strings.HasPrefix(cfg.Name, "pre:") {
		kind = PREPROC
	}

	switch kind {
	case FUNCTION:
		// Make sure the function's signature is alright.
		if len(fun.Type.In) < 1 || (len(fun.Type.Out) != 1 && len(fun.Type.Out) != 2) {
			return nil, &Error{C: E_CONFIG_FUNCTYPE, s: src, fn: fn, rc: cfg}
		}
		if fun.Type.Out[0].Type.Kind != types.BOOL || (len(fun.Type.Out) > 1 && !fun.Type.Out[1].Type.IsGoError()) {
			return nil, &Error{C: E_CONFIG_FUNCTYPE, s: src, fn: fn, rc: cfg}
		}
	case PREPROC:
		// Make sure the function's signature is alright.
		if len(fun.Type.In) < 1 || (len(fun.Type.Out) != 1 || !fun.Type.In[0].Type.IsIdenticalTo(fun.Type.Out[0].Type)) {
			return nil, &Error{C: E_CONFIG_PREFUNCTYPE, s: src, fn: fn, rc: cfg}
		}

		// Joins & Errors are NOT supported for preprocs.
		// These two could probably be just warnings.
		if join > 0 {
			return nil, &Error{C: E_CONFIG_PREPROCJOIN, s: src, fn: fn, rc: cfg}
		}
		if cfg.Error != (config.RuleErrMesg{}) {
			return nil, &Error{C: E_CONFIG_PREPROCERROR, s: src, fn: fn, rc: cfg}
		}
	}

	// If args were specified in the configuration then make sure
	// that their number is compatible with the function's signature.
	if nargs := len(cfg.Args); nargs > 0 {
		if !isValidNumberOfArgs(nargs, fun, join) {
			return nil, &Error{C: E_CONFIG_ARGNUM, s: src, fn: fn, rc: cfg}
		}
	}
	// If arg bounds were specified in the configuration then make
	// sure that they are compatible with the function's signature.
	if min := cfg.ArgMin; min != nil {
		if !isValidNumberOfArgs(int(*min), fun, join) {
			return nil, &Error{C: E_CONFIG_ARGBOUNDS, s: src, fn: fn, rc: cfg}
		}
	}
	if max := cfg.ArgMax; max != nil {
		if !isValidNumberOfArgs(*max, fun, join) {
			return nil, &Error{C: E_CONFIG_ARGBOUNDS, s: src, fn: fn, rc: cfg}
		}
	}

	// Build the rule from the func's type & the config.
	rule := new(Rule)
	rule.Name = cfg.Name
	rule.Kind = kind
	rule.Func = fun
	rule.JoinOp = join
	rule.Err = ErrMesg(cfg.Error)

	// the "re" (regexp) rule should use raw strings for arguments
	if rule.Name == "re" {
		rule.UseRawString = true
	}

	// resolve the arg bounds
	rule.ArgMin = len(fun.Type.In) - 1 // -1 for the field argument
	rule.ArgMax = len(fun.Type.In) - 1 // -1 for the field argument
	if cfg.ArgMin != nil {
		rule.ArgMin = int(*cfg.ArgMin)
	} else if fun.Type.IsVariadic {
		rule.ArgMin -= 1
	}
	if cfg.ArgMax != nil {
		rule.ArgMax = *cfg.ArgMax
	} else if fun.Type.IsVariadic || join > 0 {
		rule.ArgMax = -1
	}

	// make the arg options more convenient to use
	for _, a := range cfg.Args {
		argOpts := make(map[string]Arg)
		if def := a.Default; def != nil {
			argOpts[""] = Arg{
				Type:  _configArgTypeMap[def.Type],
				Value: def.Value,
			}
		}
		for _, opt := range a.Options {
			argOpts[opt.Value.Value] = Arg{
				Type:  _configArgTypeMap[opt.Value.Type],
				Value: opt.Value.Value,
			}
			if len(opt.Alias) > 0 {
				argOpts[opt.Alias] = Arg{
					Type:  _configArgTypeMap[opt.Value.Type],
					Value: opt.Value.Value,
				}
			}
		}
		rule.ArgOpts = append(rule.ArgOpts, argOpts)
	}

	// If args were provided make sure that their types are
	// assignable to their corresponding parameters' types.
	if len(rule.ArgOpts) > 0 {
		if err := validateArgOptsAsFuncParams(rule); err != nil {
			e := err.(*Error)
			e.s, e.fn, e.rc = src, fn, cfg
			return nil, e
		}
	}

	return rule, nil
}

func isValidNumberOfArgs(nargs int, fun *types.Func, join JoinOp) bool {
	nparams := len(fun.Type.In) - 1 // -1 for the field argument

	// If not variadic & not joinable, require exact match.
	if !fun.Type.IsVariadic && join == 0 && nargs != nparams {
		return false
	}
	// If not variadic but joinable, require that the args
	// can be spread evenly across multiple joined instances.
	if !fun.Type.IsVariadic && join > 1 && (nargs%nparams) != 0 {
		return false
	}
	// If variadic, require only that enough args were specified.
	if fun.Type.IsVariadic && nargs < nparams-1 {
		return false
	}

	return true
}

func validateArgOptsAsFuncParams(r *Rule) error {
	params := r.Func.Type.In
	firstParamIsField := (r.Kind != METHOD)

	if firstParamIsField && (len(r.Func.Type.In) > 1 || !r.Func.Type.IsVariadic) {
		// Drop the first parameter (the field) if there is more
		// than one parameter, or if the function is non-variadic.
		params = r.Func.Type.In[1:]
	}

	for i, argOpts := range r.ArgOpts {
		var j int
		var p *types.Var

		switch last := len(params) - 1; {
		case i < last || (i == last && !r.Func.Type.IsVariadic):
			j = i
			p = params[j]

		case i >= last && r.Func.Type.IsVariadic:
			j = last
			p = &types.Var{
				Name: params[j].Name,
				Type: params[j].Type.Elem,
			}

		case i > last && r.JoinOp > 0:
			if len(r.ArgOpts)-i >= len(params) {
				j = i % len(params)
				p = params[j]
			}
		}
		if p == nil {
			// NOTE this relies on the fact that ArgOpts is
			// constructured from the rule's Args *after* it
			// has been confirmed, with isValidNumberOfArgs,
			// that the number of Args is ok, given the
			// associated function's type.
			panic("shouldn't reach")
		}

		for k, a := range argOpts {
			if !canAssignArgTo(&a, p.Type) {
				a, i, k := a, i, k // copy
				return &Error{C: E_CONFIG_ARGTYPE,
					rca: &a, rcai: &i, rcak: &k,
					fp: p, fpi: &j}
			}
		}
	}
	return nil
}

// canAssignArgTo reports whether or not the Arg can
// be assigned to the Go type represented by t.
func canAssignArgTo(a *Arg, t *types.Type) bool {
	// NOTE(self): remember this assumes a.Type is not a field reference
	// since field reference arg-options don't make much sense in specs.

	// arg is unknown, accept
	if a.Type == ARG_UNKNOWN {
		return true
	}

	// t is interface{} or string, accept
	if t.IsEmptyIface() || t.Kind == types.STRING {
		return true
	}

	// both are booleans, accept
	if a.Type == ARG_BOOL && t.Kind == types.BOOL {
		return true
	}

	// t is float and arg is numeric, accept
	if a.IsNumeric() && t.Kind.IsFloat() {
		return true
	}

	// both are integers, accept
	if a.Type == ARG_INT && t.Kind.IsInteger() {
		return true
	}

	// t is unsigned and arg is not negative, accept
	if a.IsUInt() && t.Kind.IsUnsigned() {
		return true
	}

	// arg is string & a Go string literal can be converted to t, accept
	if a.Type == ARG_STRING {
		tt := &types.Type{Kind: types.STRING}
		if tt.IsConvertibleTo(t) {
			return true
		}
	}
	return false
}
