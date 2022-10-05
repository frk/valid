package rules

import (
	"github.com/frk/valid/cmd/internal/gotype"
)

func GetSpec(name string) *Spec {
	if spec, ok := _custom[name]; ok {
		return spec
	}
	if spec, ok := _included[name]; ok {
		return spec
	}
	if spec, ok := _builtin_and_stdlib[name]; ok {
		return spec
	}
	return nil
}

// A set of rule specs populated by initCustomSpecs.
var _custom = map[string]*Spec{}

// A set of rule specs populated by loadIncludedSpecs.
var _included = map[string]*Spec{}

// A set of specs for builtin rules and rules implemented
// with functions from the standard library.
var _builtin_and_stdlib = map[string]*Spec{}

func init() {
	specs := joinSpecLists(
		_builtin_specs,
		_special_specs,
		_stdlib_specs,
		_stdlib_pre_specs,
	)
	for _, s := range specs {
		key := s.Name
		if s.Kind == PREPROC {
			key = "pre:" + key
		}
		_builtin_and_stdlib[key] = s
	}
}

// A list of specs that utilize the Go language's primitive
// operators and builtin functions to do the validation.
var _builtin_specs = []*Spec{{
	Name:   "required",
	Kind:   REQUIRED,
	ArgMin: 0,
	ArgMax: 0,
	Err:    ErrSpec{Text: "is required"},
}, {
	Name:   "notnil",
	Kind:   REQUIRED,
	ArgMin: 0,
	ArgMax: 0,
	Err:    ErrSpec{Text: "cannot be nil"},
}, {
	Name:   "optional",
	Kind:   OPTIONAL,
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:   "omitnil",
	Kind:   OPTIONAL,
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:   "eq",
	Kind:   COMPARABLE,
	ArgMin: 1,
	ArgMax: -1,
	JoinOp: JOIN_OR,
	Err: ErrSpec{
		Text:     "must be equal to",
		WithArgs: true,
		ArgSep:   " or ",
	},
}, {
	Name:   "ne",
	Kind:   COMPARABLE,
	ArgMin: 1,
	ArgMax: -1,
	JoinOp: JOIN_OR,
	Err: ErrSpec{
		Text:     "must not be equal to",
		WithArgs: true,
		ArgSep:   " or ",
	},
}, {
	Name:   "gt",
	Kind:   ORDERED,
	ArgMin: 1,
	ArgMax: 1,
	Err: ErrSpec{
		Text:     "must be greater than",
		WithArgs: true,
	},
}, {
	Name:   "lt",
	Kind:   ORDERED,
	ArgMin: 1,
	ArgMax: 1,
	Err: ErrSpec{
		Text:     "must be less than",
		WithArgs: true,
	},
}, {
	Name:   "gte",
	Kind:   ORDERED,
	ArgMin: 1,
	ArgMax: 1,
	Err: ErrSpec{
		Text:     "must be greater than or equal to",
		WithArgs: true,
	},
}, {
	Name:   "lte",
	Kind:   ORDERED,
	ArgMin: 1,
	ArgMax: 1,
	Err: ErrSpec{
		Text:     "must be less than or equal to",
		WithArgs: true,
	},
}, { // alias for gte
	Name:   "min",
	Kind:   ORDERED,
	ArgMin: 1,
	ArgMax: 1,
	Err: ErrSpec{
		Text:     "must be greater than or equal to",
		WithArgs: true,
	},
}, { // alias for lte
	Name:   "max",
	Kind:   ORDERED,
	ArgMin: 1,
	ArgMax: 1,
	Err: ErrSpec{
		Text: "must be less than or equal to", WithArgs: true,
	},
}, {
	Name:   "rng",
	Kind:   RANGE,
	ArgMin: 2,
	ArgMax: 2,
	Err: ErrSpec{
		Text: "must be between", WithArgs: true,
		ArgSep: " and ",
	},
}, {
	Name:   "enum",
	Kind:   ENUM,
	ArgMin: 0,
	ArgMax: 0,
	Err: ErrSpec{
		// NOTE since "enum" takes no arguments it is the
		// generator's responsibility to use the enum's
		// constants as the arguments to the error message.
		Text:     "must be one of",
		WithArgs: true,
		ArgSep:   " or ",
	},
}, {
	Name:   "len",
	Kind:   LENGTH,
	ArgMin: 1,
	ArgMax: 2,
	ErrOpts: map[string]ErrSpec{
		"x":  {Text: "must be of length", WithArgs: true},
		"x:": {Text: "must be of length at least", WithArgs: true},
		":x": {Text: "must be of length at most", WithArgs: true},
		"x:x": {Text: "must be of length between",
			WithArgs:  true,
			ArgSep:    " and ",
			ArgSuffix: "(inclusive)",
		},
	},
}}

// A list of specs for "special" rules.
var _special_specs = []*Spec{{
	Name:   "noguard",
	Kind:   NOGUARD,
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "isvalid",
	Kind:  METHOD,
	FName: "IsValid",
	FType: &gotype.Type{
		Pkg:  gotype.Pkg{},
		In:   []*gotype.Var{},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_BOOL}}},
		Kind: gotype.K_FUNC,
	},
	Err: ErrSpec{Text: "is not valid"},
}, {
	Name:   "-isvalid",
	Kind:   REMOVE,
	ArgMin: 0,
	ArgMax: 0,
}}

// A list of specs for validation rules that are implmeneted
// with functions from the Go standard library.
var _stdlib_specs = []*Spec{{
	Name:   "runecount",
	Kind:   LENGTH,
	ArgMin: 1,
	ArgMax: 2,
	ErrOpts: map[string]ErrSpec{
		"x":  {Text: "must have rune count", WithArgs: true},
		"x:": {Text: "must have rune count at least", WithArgs: true},
		":x": {Text: "must have rune count at most", WithArgs: true},
		"x:x": {Text: "must have rune count between",
			WithArgs:  true,
			ArgSep:    " and ",
			ArgSuffix: "(inclusive)",
		},
	},
}, {
	Name:  "prefix",
	Kind:  FUNCTION,
	FName: "HasPrefix",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_BOOL}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: -1,
	JoinOp: JOIN_OR,
	Err: ErrSpec{
		Text:     "must be prefixed with",
		WithArgs: true,
		ArgSep:   " or ",
	},
}, {
	Name:  "suffix",
	Kind:  FUNCTION,
	FName: "HasSuffix",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_BOOL}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: -1,
	JoinOp: JOIN_OR,
	Err: ErrSpec{
		Text:     "must be suffixed with",
		WithArgs: true,
		ArgSep:   " or ",
	},
}, {
	// NOTE this, together with the Type.JoinOp set to JOIN_OR
	// is effectively "contains one of" changing JoinOp to JOIN_AND
	// would turn it to "contains all", which may be handy
	// as well and might get added later.
	Name:  "contains",
	Kind:  FUNCTION,
	FName: "Contains",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_BOOL}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: -1,
	JoinOp: JOIN_OR,
	Err: ErrSpec{
		Text:     "must contain substring",
		WithArgs: true,
		ArgSep:   " or ",
	},
}}

// A list of specs for preprocessor rules that are implmeneted
// with functions from the Go standard library.
var _stdlib_pre_specs = []*Spec{{
	Name:  "repeat",
	Kind:  PREPROC,
	FName: "Repeat",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_INT}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: 1,
}, {
	Name:  "replace",
	Kind:  PREPROC,
	FName: "Replace",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_INT}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin:  2,
	ArgMax:  3,
	ArgOpts: []map[string]Arg{{}, {}, {"": {ARG_INT, "-1"}}},
}, {
	Name:  "lower",
	Kind:  PREPROC,
	FName: "ToLower",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "upper",
	Kind:  PREPROC,
	FName: "ToUpper",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "title",
	Kind:  PREPROC,
	FName: "ToTitle",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "validutf8",
	Kind:  PREPROC,
	FName: "ToValidUTF8",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin:  0,
	ArgMax:  1,
	ArgOpts: []map[string]Arg{{"": {ARG_STRING, ""}}},
}, {
	Name:  "trim",
	Kind:  PREPROC,
	FName: "TrimSpace",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "ltrim",
	Kind:  PREPROC,
	FName: "TrimLeft",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: 1,
}, {
	Name:  "rtrim",
	Kind:  PREPROC,
	FName: "TrimRight",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: 1,
}, {
	Name:  "trimprefix",
	Kind:  PREPROC,
	FName: "TrimPrefix",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: 1,
}, {
	Name:  "trimsuffix",
	Kind:  PREPROC,
	FName: "TrimSuffix",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strings",
			Path: "strings",
		},
		In: []*gotype.Var{
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
			{Type: &gotype.Type{Kind: gotype.K_STRING}},
		},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 1,
	ArgMax: 1,
}, {
	Name:  "quote",
	Kind:  PREPROC,
	FName: "Quote",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strconv",
			Path: "strconv",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "quoteascii",
	Kind:  PREPROC,
	FName: "QuoteToASCII",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strconv",
			Path: "strconv",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "quotegraphic",
	Kind:  PREPROC,
	FName: "QuoteToGraphic",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "strconv",
			Path: "strconv",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "urlqueryesc",
	Kind:  PREPROC,
	FName: "QueryEscape",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "url",
			Path: "net/url",
		},
		// XXX
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "urlpathesc",
	Kind:  PREPROC,
	FName: "PathEscape",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "url",
			Path: "net/url",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "htmlesc",
	Kind:  PREPROC,
	FName: "EscapeString",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "html",
			Path: "html",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "htmlunesc",
	Kind:  PREPROC,
	FName: "UnescapeString",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "html",
			Path: "html",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_STRING}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}, {
	Name:  "round",
	Kind:  PREPROC,
	FName: "Round",
	FType: &gotype.Type{
		Pkg: gotype.Pkg{
			Name: "math",
			Path: "math",
		},
		In:   []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_FLOAT64}}},
		Out:  []*gotype.Var{{Type: &gotype.Type{Kind: gotype.K_FLOAT64}}},
		Kind: gotype.K_FUNC,
	},
	ArgMin: 0,
	ArgMax: 0,
}}

// helper
func joinSpecLists(ss ...[]*Spec) (out []*Spec) {
	for i := range ss {
		out = append(out, ss[i]...)
	}
	return out
}
