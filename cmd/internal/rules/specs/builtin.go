package specs

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
)

var builtin_list = []*rules.Spec{
	////////////////////////////////////////////////////////////////////////
	// A list of specs that utilize the Go language's primitive
	// operators and builtin functions to do the validation.
	////////////////////////////////////////////////////////////////////////
	{
		Name: "required",
		Kind: rules.REQUIRED,
		Err:  rules.ErrSpec{Text: "is required"},
	}, {
		Name: "notnil",
		Kind: rules.REQUIRED,
		Err:  rules.ErrSpec{Text: "cannot be nil"},
	}, {
		Name: "optional",
		Kind: rules.OPTIONAL,
	}, {
		Name: "omitnil",
		Kind: rules.OPTIONAL,
	}, {
		Name:   "eq",
		Kind:   rules.COMPARABLE,
		ArgMin: 1,
		ArgMax: -1,
		JoinOp: rules.JOIN_OR,
		Err: rules.ErrSpec{
			Text:     "must be equal to",
			WithArgs: true,
			ArgSep:   " or ",
		},
	}, {
		Name:   "ne",
		Kind:   rules.COMPARABLE,
		ArgMin: 1,
		ArgMax: -1,
		JoinOp: rules.JOIN_OR,
		Err: rules.ErrSpec{
			Text:     "must not be equal to",
			WithArgs: true,
			ArgSep:   " or ",
		},
	}, {
		Name:   "gt",
		Kind:   rules.ORDERED,
		ArgMin: 1,
		ArgMax: 1,
		Err: rules.ErrSpec{
			Text:     "must be greater than",
			WithArgs: true,
		},
	}, {
		Name:   "lt",
		Kind:   rules.ORDERED,
		ArgMin: 1,
		ArgMax: 1,
		Err: rules.ErrSpec{
			Text:     "must be less than",
			WithArgs: true,
		},
	}, {
		Name:   "gte",
		Kind:   rules.ORDERED,
		ArgMin: 1,
		ArgMax: 1,
		Err: rules.ErrSpec{
			Text:     "must be greater than or equal to",
			WithArgs: true,
		},
	}, {
		Name:   "lte",
		Kind:   rules.ORDERED,
		ArgMin: 1,
		ArgMax: 1,
		Err: rules.ErrSpec{
			Text:     "must be less than or equal to",
			WithArgs: true,
		},
	}, { // alias for gte
		Name:   "min",
		Kind:   rules.ORDERED,
		ArgMin: 1,
		ArgMax: 1,
		Err: rules.ErrSpec{
			Text:     "must be greater than or equal to",
			WithArgs: true,
		},
	}, { // alias for lte
		Name:   "max",
		Kind:   rules.ORDERED,
		ArgMin: 1,
		ArgMax: 1,
		Err: rules.ErrSpec{
			Text:     "must be less than or equal to",
			WithArgs: true,
		},
	}, {
		Name:   "rng",
		Kind:   rules.RANGE,
		ArgMin: 2,
		ArgMax: 2,
		Err: rules.ErrSpec{
			Text:     "must be between",
			WithArgs: true,
			ArgSep:   " and ",
		},
	}, { // alias for rng
		Name:   "between",
		Kind:   rules.RANGE,
		ArgMin: 2,
		ArgMax: 2,
		Err: rules.ErrSpec{
			Text:     "must be between",
			WithArgs: true,
			ArgSep:   " and ",
		},
	}, {
		Name: "enum",
		Kind: rules.ENUM,
		Err: rules.ErrSpec{
			// NOTE since "enum" takes no arguments it is the
			// generator's responsibility to use the enum's
			// constants as the arguments to the error message.
			Text:     "must be one of",
			WithArgs: true,
			ArgSep:   " or ",
		},
	}, {
		Name:   "len",
		Kind:   rules.LENGTH,
		ArgMin: 1,
		ArgMax: 2,
		ErrOpts: map[string]rules.ErrSpec{
			"x":  {Text: "must be of length", WithArgs: true},
			"x:": {Text: "must be of length at least", WithArgs: true},
			":x": {Text: "must be of length at most", WithArgs: true},
			"x:x": {
				Text:      "must be of length between",
				WithArgs:  true,
				ArgSep:    " and ",
				ArgSuffix: "(inclusive)",
			},
		},
	}, {
		Name:   "runecount",
		Kind:   rules.LENGTH,
		ArgMin: 1,
		ArgMax: 2,
		ErrOpts: map[string]rules.ErrSpec{
			"x":  {Text: "must have rune count", WithArgs: true},
			"x:": {Text: "must have rune count at least", WithArgs: true},
			":x": {Text: "must have rune count at most", WithArgs: true},
			"x:x": {
				Text:      "must have rune count between",
				WithArgs:  true,
				ArgSep:    " and ",
				ArgSuffix: "(inclusive)",
			},
		},
	},

	////////////////////////////////////////////////////////////////////////
	// A list of specs for "special" rules.
	////////////////////////////////////////////////////////////////////////
	{

		Name: "noguard",
		Kind: rules.NOGUARD,
	}, {
		Name: "-isvalid",
		Kind: rules.REMOVE,
	}, {
		Name: "isvalid",
		Kind: rules.METHOD,
		Func: &rules.FuncIdent{Name: "IsValid"},
		Err:  rules.ErrSpec{Text: "is not valid"},
	},

	////////////////////////////////////////////////////////////////////////
	// A list of specs for validation and preprocessor rules that are
	// implmeneted with functions from the Go standard library.
	////////////////////////////////////////////////////////////////////////
	{
		Name:   "prefix",
		Kind:   rules.FUNCTION,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "HasPrefix"},
		ArgMin: 1,
		ArgMax: -1,
		JoinOp: rules.JOIN_OR,
		Err: rules.ErrSpec{
			Text:     "must be prefixed with",
			WithArgs: true,
			ArgSep:   " or ",
		},
	}, {
		Name:   "suffix",
		Kind:   rules.FUNCTION,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "HasSuffix"},
		ArgMin: 1,
		ArgMax: -1,
		JoinOp: rules.JOIN_OR,
		Err: rules.ErrSpec{
			Text:     "must be suffixed with",
			WithArgs: true,
			ArgSep:   " or ",
		},
	}, {
		// NOTE this, together with the Type.JoinOp set to JOIN_OR
		// is effectively "contains one of" changing JoinOp to JOIN_AND
		// would turn it to "contains all", which may be handy
		// as well and might get added later.
		Name:   "contains",
		Kind:   rules.FUNCTION,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "Contains"},
		ArgMin: 1,
		ArgMax: -1,
		JoinOp: rules.JOIN_OR,
		Err: rules.ErrSpec{
			Text:     "must contain substring",
			WithArgs: true,
			ArgSep:   " or ",
		},
	}, {
		Name:   "repeat",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "Repeat"},
		ArgMin: 1,
		ArgMax: 1,
	}, {
		Name:   "replace",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "Replace"},
		ArgMin: 2,
		ArgMax: 3,
		ArgOpts: []map[string]rules.Arg{
			{},
			{},
			{"": {rules.ARG_INT, "-1"}},
		},
	}, {
		Name: "lower",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strings", Name: "ToLower"},
	}, {
		Name: "upper",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strings", Name: "ToUpper"},
	}, {
		Name: "title",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strings", Name: "ToTitle"},
	}, {
		Name:   "validutf8",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "ToValidUTF8"},
		ArgMin: 0,
		ArgMax: 1,
		ArgOpts: []map[string]rules.Arg{
			{"": {rules.ARG_STRING, ""}},
		},
	}, {
		Name: "trim",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strings", Name: "TrimSpace"},
	}, {
		Name:   "ltrim",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "TrimLeft"},
		ArgMin: 1,
		ArgMax: 1,
	}, {
		Name:   "rtrim",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "TrimRight"},
		ArgMin: 1,
		ArgMax: 1,
	}, {
		Name:   "trimprefix",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "TrimPrefix"},
		ArgMin: 1,
		ArgMax: 1,
	}, {
		Name:   "trimsuffix",
		Kind:   rules.PREPROC,
		Func:   &rules.FuncIdent{Pkg: "strings", Name: "TrimSuffix"},
		ArgMin: 1,
		ArgMax: 1,
	}, {
		Name: "quote",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strconv", Name: "Quote"},
	}, {
		Name: "quoteascii",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strconv", Name: "QuoteToASCII"},
	}, {
		Name: "quotegraphic",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "strconv", Name: "QuoteToGraphic"},
	}, {
		Name: "urlqueryesc",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "net/url", Name: "QueryEscape"},
	}, {
		Name: "urlpathesc",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "net/url", Name: "PathEscape"},
	}, {
		Name: "htmlesc",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "html", Name: "EscapeString"},
	}, {
		Name: "htmlunesc",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "html", Name: "UnescapeString"},
	}, {
		Name: "round",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "math", Name: "Round"},
	}, {
		Name: "ceil",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "math", Name: "Ceil"},
	}, {
		Name: "floor",
		Kind: rules.PREPROC,
		Func: &rules.FuncIdent{Pkg: "math", Name: "Floor"},
	},
}
