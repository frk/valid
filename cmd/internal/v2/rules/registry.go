package rules

import (
	"github.com/frk/valid/cmd/internal/v2/config"
	"github.com/frk/valid/cmd/internal/v2/source"
)

type Registry interface {
	Lookup(name string) (r Rule, ok bool)
}

type registry struct {
	builtins map[string]Rule
	included map[string]Rule
	custom   map[string]Rule
}

func (rr registry) Lookup(name string) (r Rule, ok bool) {
	if r, ok := rr.custom[name]; ok {
		return r, ok
	}
	if r, ok := rr.included[name]; ok {
		return r, ok
	}
	if r, ok := rr.builtins[name]; ok {
		return r, ok
	}
	return Rule{}, false
}

////////////////////////////////////////////////////////////////////////////////

var rreg registry

func Load(cfg config.Config, src *source.Source) error {
	if err := rreg.load_builtins(src); err != nil {
		return err
	}
	if err := rreg.load_included(src); err != nil {
		return err
	}
	return rreg.load_custom(cfg, src)
}

func (rr *registry) load_custom(cfg config.Config, src *source.Source) error {
	rr.custom = make(map[string]Rule)
	return nil
}

func (rr *registry) load_included(src *source.Source) error {
	rr.included = make(map[string]Rule)

	// loads functions from the "github.com/frk/valid" package
	funcs, err := src.GetIncludedRuleFuncs()
	if err != nil {
		return err
	}
	for _, fn := range funcs {
		r, err := newRuleFromFunc(fn, src)
		if err != nil {
			return err
		}
		rr.included[r.Name] = *r
		return nil
	}
	return nil
}

func (rr *registry) load_builtins(src *source.Source) error {
	rr.builtins = make(map[string]Rule)
	return nil
}

//var builtins = []Rule{
//	////////////////////////////////////////////////////////////////////////
//	// A list of rules that utilize the Go language's primitive
//	// operators and builtin functions to do the validation.
//	////////////////////////////////////////////////////////////////////////
//	{
//		Name: "required",
//		Kind: REQUIRED,
//		Err:  ErrMesg{Text: "is required"},
//	}, {
//		Name: "notnil",
//		Kind: REQUIRED,
//		Err:  ErrMesg{Text: "cannot be nil"},
//	}, {
//		Name: "optional",
//		Kind: OPTIONAL,
//	}, {
//		Name: "omitnil",
//		Kind: OPTIONAL,
//	}, {
//		Name:   "eq",
//		Kind:   COMPARABLE,
//		ArgMin: 1,
//		ArgMax: -1,
//		JoinOp: JOIN_OR,
//		Err: ErrMesg{
//			Text:     "must be equal to",
//			WithArgs: true,
//			ArgSep:   " or ",
//		},
//	}, {
//		Name:   "ne",
//		Kind:   COMPARABLE,
//		ArgMin: 1,
//		ArgMax: -1,
//		JoinOp: JOIN_OR,
//		Err: ErrMesg{
//			Text:     "must not be equal to",
//			WithArgs: true,
//			ArgSep:   " or ",
//		},
//	}, {
//		Name:   "gt",
//		Kind:   ORDERED,
//		ArgMin: 1,
//		ArgMax: 1,
//		Err: ErrMesg{
//			Text:     "must be greater than",
//			WithArgs: true,
//		},
//	}, {
//		Name:   "lt",
//		Kind:   ORDERED,
//		ArgMin: 1,
//		ArgMax: 1,
//		Err: ErrMesg{
//			Text:     "must be less than",
//			WithArgs: true,
//		},
//	}, {
//		Name:   "gte",
//		Kind:   ORDERED,
//		ArgMin: 1,
//		ArgMax: 1,
//		Err: ErrMesg{
//			Text:     "must be greater than or equal to",
//			WithArgs: true,
//		},
//	}, {
//		Name:   "lte",
//		Kind:   ORDERED,
//		ArgMin: 1,
//		ArgMax: 1,
//		Err: ErrMesg{
//			Text:     "must be less than or equal to",
//			WithArgs: true,
//		},
//	}, { // alias for gte
//		Name:   "min",
//		Kind:   ORDERED,
//		ArgMin: 1,
//		ArgMax: 1,
//		Err: ErrMesg{
//			Text:     "must be greater than or equal to",
//			WithArgs: true,
//		},
//	}, { // alias for lte
//		Name:   "max",
//		Kind:   ORDERED,
//		ArgMin: 1,
//		ArgMax: 1,
//		Err: ErrMesg{
//			Text:     "must be less than or equal to",
//			WithArgs: true,
//		},
//	}, {
//		Name:   "rng",
//		Kind:   RANGE,
//		ArgMin: 2,
//		ArgMax: 2,
//		Err: ErrMesg{
//			Text:     "must be between",
//			WithArgs: true,
//			ArgSep:   " and ",
//		},
//	}, { // alias for rng
//		Name:   "between",
//		Kind:   RANGE,
//		ArgMin: 2,
//		ArgMax: 2,
//		Err: ErrMesg{
//			Text:     "must be between",
//			WithArgs: true,
//			ArgSep:   " and ",
//		},
//	}, {
//		Name: "enum",
//		Kind: ENUM,
//		Err: ErrMesg{
//			// NOTE since "enum" takes no arguments it is the
//			// generator's responsibility to use the enum's
//			// constants as the arguments to the error message.
//			Text:     "must be one of",
//			WithArgs: true,
//			ArgSep:   " or ",
//		},
//	}, {
//		Name:   "len",
//		Kind:   LENGTH,
//		ArgMin: 1,
//		ArgMax: 2,
//		ErrOpts: map[string]ErrMesg{
//			"x":  {Text: "must be of length", WithArgs: true},
//			"x:": {Text: "must be of length at least", WithArgs: true},
//			":x": {Text: "must be of length at most", WithArgs: true},
//			"x:x": {
//				Text:      "must be of length between",
//				WithArgs:  true,
//				ArgSep:    " and ",
//				ArgSuffix: "(inclusive)",
//			},
//		},
//	}, {
//		Name:   "runecount",
//		Kind:   LENGTH,
//		ArgMin: 1,
//		ArgMax: 2,
//		ErrOpts: map[string]ErrMesg{
//			"x":  {Text: "must have rune count", WithArgs: true},
//			"x:": {Text: "must have rune count at least", WithArgs: true},
//			":x": {Text: "must have rune count at most", WithArgs: true},
//			"x:x": {
//				Text:      "must have rune count between",
//				WithArgs:  true,
//				ArgSep:    " and ",
//				ArgSuffix: "(inclusive)",
//			},
//		},
//	},
//
//	////////////////////////////////////////////////////////////////////////
//	// A list of specs for "special"
//	////////////////////////////////////////////////////////////////////////
//	{
//		Name: "noguard",
//		Kind: NOGUARD,
//	}, {
//		Name: "-isvalid",
//		Kind: REMOVE,
//	}, {
//		Name: "isvalid",
//		Kind: METHOD,
//		Func: &FuncIdent{Name: "IsValid"},
//		Err:  ErrMesg{Text: "is not valid"},
//	},
//
//	////////////////////////////////////////////////////////////////////////
//	// A list of specs for validation and preprocessor rules that are
//	// implmeneted with functions from the Go standard library.
//	////////////////////////////////////////////////////////////////////////
//	{
//		Name:   "prefix",
//		Kind:   FUNCTION,
//		Func:   &FuncIdent{Pkg: "strings", Name: "HasPrefix"},
//		ArgMin: 1,
//		ArgMax: -1,
//		JoinOp: JOIN_OR,
//		Err: ErrMesg{
//			Text:     "must be prefixed with",
//			WithArgs: true,
//			ArgSep:   " or ",
//		},
//	}, {
//		Name:   "suffix",
//		Kind:   FUNCTION,
//		Func:   &FuncIdent{Pkg: "strings", Name: "HasSuffix"},
//		ArgMin: 1,
//		ArgMax: -1,
//		JoinOp: JOIN_OR,
//		Err: ErrMesg{
//			Text:     "must be suffixed with",
//			WithArgs: true,
//			ArgSep:   " or ",
//		},
//	}, {
//		// NOTE this, together with the Type.JoinOp set to JOIN_OR
//		// is effectively "contains one of" changing JoinOp to JOIN_AND
//		// would turn it to "contains all", which may be handy
//		// as well and might get added later.
//		Name:   "contains",
//		Kind:   FUNCTION,
//		Func:   &FuncIdent{Pkg: "strings", Name: "Contains"},
//		ArgMin: 1,
//		ArgMax: -1,
//		JoinOp: JOIN_OR,
//		Err: ErrMesg{
//			Text:     "must contain substring",
//			WithArgs: true,
//			ArgSep:   " or ",
//		},
//	}, {
//		Name:   "repeat",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "Repeat"},
//		ArgMin: 1,
//		ArgMax: 1,
//	}, {
//		Name:   "replace",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "Replace"},
//		ArgMin: 2,
//		ArgMax: 3,
//		ArgOpts: []map[string]Arg{
//			{},
//			{},
//			{"": {ARG_INT, "-1"}},
//		},
//	}, {
//		Name: "lower",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strings", Name: "ToLower"},
//	}, {
//		Name: "upper",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strings", Name: "ToUpper"},
//	}, {
//		Name: "title",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strings", Name: "ToTitle"},
//	}, {
//		Name:   "validutf8",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "ToValidUTF8"},
//		ArgMin: 0,
//		ArgMax: 1,
//		ArgOpts: []map[string]Arg{
//			{"": {ARG_STRING, ""}},
//		},
//	}, {
//		Name: "trim",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strings", Name: "TrimSpace"},
//	}, {
//		Name:   "ltrim",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "TrimLeft"},
//		ArgMin: 1,
//		ArgMax: 1,
//	}, {
//		Name:   "rtrim",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "TrimRight"},
//		ArgMin: 1,
//		ArgMax: 1,
//	}, {
//		Name:   "trimprefix",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "TrimPrefix"},
//		ArgMin: 1,
//		ArgMax: 1,
//	}, {
//		Name:   "trimsuffix",
//		Kind:   PREPROC,
//		Func:   &FuncIdent{Pkg: "strings", Name: "TrimSuffix"},
//		ArgMin: 1,
//		ArgMax: 1,
//	}, {
//		Name: "quote",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strconv", Name: "Quote"},
//	}, {
//		Name: "quoteascii",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strconv", Name: "QuoteToASCII"},
//	}, {
//		Name: "quotegraphic",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "strconv", Name: "QuoteToGraphic"},
//	}, {
//		Name: "urlqueryesc",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "net/url", Name: "QueryEscape"},
//	}, {
//		Name: "urlpathesc",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "net/url", Name: "PathEscape"},
//	}, {
//		Name: "htmlesc",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "html", Name: "EscapeString"},
//	}, {
//		Name: "htmlunesc",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "html", Name: "UnescapeString"},
//	}, {
//		Name: "round",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "math", Name: "Round"},
//	}, {
//		Name: "ceil",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "math", Name: "Ceil"},
//	}, {
//		Name: "floor",
//		Kind: PREPROC,
//		Func: &FuncIdent{Pkg: "math", Name: "Floor"},
//	},
//}
