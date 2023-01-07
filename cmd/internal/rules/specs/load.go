package specs

import (
	stdtypes "go/types"
	"strings"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"

	"gopkg.in/yaml.v3"
)

var (
	// A set of specs for builtin rules, special rules, and rules
	// implemented with functions from the standard library.
	builtin = map[string]*rules.Spec{}
	// A set of rule specs for functions declared in github.com/frk/valid.
	included = map[string]*rules.Spec{}
	// A set of specs for rules defined by the user in config.Config.Rules.
	custom = map[string]*rules.Spec{}
	// FuncMap maps specs of kind FUNCTION/PREPROC to their associated functions.
	funcMap = map[*rules.Spec]*types.Func{}
)

// loadBuiltins loads the rule specs from the builtin_list.
func loadBuiltins(a *search.AST) error {
	for _, s := range builtin_list {
		if s.Kind == rules.FUNCTION || s.Kind == rules.PREPROC {
			ftyp, _, err := search.FindFunc(s.Func.Pkg, s.Func.Name, a)
			if err != nil {
				return err
			}

			fn := types.AnalyzeFunc(ftyp, a)
			funcMap[s] = fn
		}

		if s.Kind == rules.METHOD && s.Name == "isvalid" {
			funcMap[s] = &types.Func{
				Name: "IsValid",
				Type: &types.Type{
					Pkg:  types.Pkg{},
					In:   []*types.Var{},
					Out:  []*types.Var{{Type: &types.Type{Kind: types.BOOL}}},
					Kind: types.FUNC,
				},
			}
		}

		key := s.Name
		if s.Kind == rules.PREPROC {
			key = "pre:" + key
		}

		builtin[key] = s
	}

	return nil
}

// loadIncludedSpecs loads the rule specs for function
// rules implemented by the github.com/frk/valid package.
func loadIncludedSpecs(a *search.AST) error {
	ft, _, err := search.FindFunc("github.com/frk/valid", "RegisterRegexp", a)
	if err != nil {
		return err // shouldn't happen
	}
	_regexpFunc = types.AnalyzeFunc(ft, a)

	// loads functions from the "github.com/frk/valid" package
	return search.FindIncludedFuncs(a, func(ftyp *stdtypes.Func, rawCfg []byte) error {
		cfg := new(config.RuleSpec)
		if err := yaml.Unmarshal(rawCfg, cfg); err != nil {
			return &Error{C: E_CONFIG_INVALID, a: a, ft: ftyp, err: err}
		}
		s, err := specFromFunc(a, ftyp, cfg)
		if err != nil {
			return err
		}
		included[s.Name] = s
		return nil
	})
}

// LoadCustomSpecs initializes custom rules from the given config.
func LoadCustomSpecs(cfg config.Config, a *search.AST) error {
	for _, rc := range cfg.Rules {
		ftyp, rawCfg, err := search.FindFunc(rc.Func.Pkg, rc.Func.Name, a)
		if err != nil {
			return &Error{C: E_CONFIG_FUNCSEARCH, a: a, c: &cfg, err: err}
		}

		// Use the config from the function's documentation
		// if none was provided in the config file itself.
		if rc.Rule == nil && len(rawCfg) > 0 {
			rc.Rule = new(config.RuleSpec)
			if err := yaml.Unmarshal(rawCfg, rc.Rule); err != nil {
				return &Error{C: E_CONFIG_INVALID, a: a, c: &cfg, ft: ftyp, err: err}
			}
		}
		if rc.Rule == nil {
			return &Error{C: E_CONFIG_MISSING, a: a, c: &cfg, ft: ftyp, err: err}
		}

		s, err := specFromFunc(a, ftyp, rc.Rule)
		if err != nil {
			return extendError(err, func(e *Error) { e.c = &cfg })
		}
		custom[s.Name] = s
	}
	return nil
}

var _joinops = [...]rules.JoinOp{
	config.JOIN_NOT: rules.JOIN_NOT,
	config.JOIN_AND: rules.JOIN_AND,
	config.JOIN_OR:  rules.JOIN_OR,
}

var _scalarargs = [...]rules.ArgType{
	config.NIL:    rules.ARG_UNKNOWN,
	config.BOOL:   rules.ARG_BOOL,
	config.INT:    rules.ARG_INT,
	config.FLOAT:  rules.ARG_FLOAT,
	config.STRING: rules.ARG_STRING,
}

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

func specFromFunc(a *search.AST, f *stdtypes.Func, rs *config.RuleSpec) (*rules.Spec, error) {
	if rs.Name == "" {
		return nil, &Error{C: E_CONFIG_NONAME, a: a, ft: f, rs: rs}
	}
	if _reserved[rs.Name] {
		return nil, &Error{C: E_CONFIG_RESERVED, a: a, ft: f, rs: rs}
	}

	fn := types.AnalyzeFunc(f, a)
	jop := _joinops[rs.JoinOp]

	specKind := rules.FUNCTION
	if strings.HasPrefix(rs.Name, "pre:") {
		specKind = rules.PREPROC
	}

	switch specKind {
	case rules.FUNCTION:
		// Make sure the function's signature is alright.
		if len(fn.Type.In) < 1 || (len(fn.Type.Out) != 1 && len(fn.Type.Out) != 2) {
			return nil, &Error{C: E_CONFIG_FUNCTYPE, a: a, ft: f, rs: rs}
		}
		if fn.Type.Out[0].Type.Kind != types.BOOL || (len(fn.Type.Out) > 1 && !fn.Type.Out[1].Type.IsGoError()) {
			return nil, &Error{C: E_CONFIG_FUNCTYPE, a: a, ft: f, rs: rs}
		}
	case rules.PREPROC:
		// Make sure the function's signature is alright.
		if len(fn.Type.In) < 1 || (len(fn.Type.Out) != 1 || !fn.Type.In[0].Type.IsIdenticalTo(fn.Type.Out[0].Type)) {
			return nil, &Error{C: E_CONFIG_PREFUNCTYPE, a: a, ft: f, rs: rs}
		}

		// Joins & Errors are NOT supported for preprocs.
		// These two could probably be just warnings.
		if jop > 0 {
			return nil, &Error{C: E_CONFIG_PREPROCJOIN, a: a, ft: f, rs: rs}
		}
		if rs.Error != (config.RuleErrorConfig{}) {
			return nil, &Error{C: E_CONFIG_PREPROCERROR, a: a, ft: f, rs: rs}
		}
	}

	// If args were specified in the configuration then make sure
	// that their number is compatible with the function's signature.
	if nargs := len(rs.Args); nargs > 0 {
		if !isValidNumberOfArgs(nargs, fn, jop) {
			return nil, &Error{C: E_CONFIG_ARGNUM, a: a, ft: f, rs: rs}
		}
	}
	// If arg bounds were specified in the configuration then make
	// sure that they are compatible with the function's signature.
	if min := rs.ArgMin; min != nil {
		if !isValidNumberOfArgs(int(*min), fn, jop) {
			return nil, &Error{C: E_CONFIG_ARGBOUNDS, a: a, ft: f, rs: rs}
		}
	}
	if max := rs.ArgMax; max != nil {
		if !isValidNumberOfArgs(*max, fn, jop) {
			return nil, &Error{C: E_CONFIG_ARGBOUNDS, a: a, ft: f, rs: rs}
		}
	}

	// Build the rule spec from the func's type & the config.
	spec := new(rules.Spec)
	spec.Name = rs.Name
	spec.Kind = specKind
	spec.Func = new(rules.FuncIdent)
	spec.Func.Pkg = fn.Type.Pkg.Path
	spec.Func.Name = fn.Name
	spec.JoinOp = jop
	spec.Err = rules.ErrSpec(rs.Error)

	// the "re" (regexp) rule should use raw strings for arguments
	if spec.Name == "re" {
		spec.UseRawString = true
	}

	// resolve the arg bounds
	spec.ArgMin = len(fn.Type.In) - 1 // -1 for the field argument
	spec.ArgMax = len(fn.Type.In) - 1 // -1 for the field argument
	if rs.ArgMin != nil {
		spec.ArgMin = int(*rs.ArgMin)
	} else if fn.Type.IsVariadic {
		spec.ArgMin -= 1
	}
	if rs.ArgMax != nil {
		spec.ArgMax = *rs.ArgMax
	} else if fn.Type.IsVariadic || jop > 0 {
		spec.ArgMax = -1
	}

	// make the arg options more convenient to use
	for _, a := range rs.Args {
		argOpts := make(map[string]rules.Arg)
		if def := a.Default; def != nil {
			argOpts[""] = rules.Arg{
				Type:  _scalarargs[def.Type],
				Value: def.Value,
			}
		}
		for _, opt := range a.Options {
			argOpts[opt.Value.Value] = rules.Arg{
				Type:  _scalarargs[opt.Value.Type],
				Value: opt.Value.Value,
			}
			if len(opt.Alias) > 0 {
				argOpts[opt.Alias] = rules.Arg{
					Type:  _scalarargs[opt.Value.Type],
					Value: opt.Value.Value,
				}
			}
		}
		spec.ArgOpts = append(spec.ArgOpts, argOpts)
	}

	// If args were provided make sure that their types are
	// assignable to their corresponding parameters' types.
	if len(spec.ArgOpts) > 0 {
		if err := validateArgOptsAsFuncParams(spec, fn); err != nil {
			e := err.(*Error)
			e.a, e.ft, e.rs = a, f, rs
			return nil, e
		}
	}

	funcMap[spec] = fn
	return spec, nil
}

func isValidNumberOfArgs(nargs int, fn *types.Func, joinOp rules.JoinOp) bool {
	nparams := len(fn.Type.In) - 1 // -1 for the field argument

	// If not variadic & not joinable, require exact match.
	if !fn.Type.IsVariadic && joinOp == 0 && nargs != nparams {
		return false
	}
	// If not variadic but joinable, require that the args
	// can be spread evenly across multiple joined instances.
	if !fn.Type.IsVariadic && joinOp > 1 && (nargs%nparams) != 0 {
		return false
	}
	// If variadic, require only that enough args were specified.
	if fn.Type.IsVariadic && nargs < nparams-1 {
		return false
	}

	return true
}

// NOTE(mkopriva): while not identical to, this does share some logic
// with checkRuleArgsAsFuncParams. If something in that logic changes
// here, make sure to apply that change to the "cousin" as well.
func validateArgOptsAsFuncParams(s *rules.Spec, fn *types.Func) error {
	params := fn.Type.In
	firstParamIsField := (s.Kind != rules.METHOD)

	if firstParamIsField && (len(fn.Type.In) > 1 || !fn.Type.IsVariadic) {
		// Drop the first parameter (the field) if there is more
		// than one parameter, or if the function is non-variadic.
		params = fn.Type.In[1:]
	}

	for i, argOpts := range s.ArgOpts {
		var j int
		var p *types.Var

		switch last := len(params) - 1; {
		case i < last || (i == last && !fn.Type.IsVariadic):
			j = i
			p = params[j]

		case i >= last && fn.Type.IsVariadic:
			j = last
			p = &types.Var{
				Name: params[j].Name,
				Type: params[j].Type.Elem.Type,
			}

		case i > last && s.JoinOp > 0:
			if len(s.ArgOpts)-i >= len(params) {
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
					rca: &a, rcai: &i, rcak: &k, fp: p, fpi: &j,
				}
			}
		}
	}
	return nil
}

// canAssignArgTo reports whether or not the Arg can
// be assigned to the Go type represented by t.
func canAssignArgTo(a *rules.Arg, t *types.Type) bool {
	// NOTE(self): remember this assumes a.Type is not a field reference
	// since field reference arg-options don't make much sense in specs.

	// arg is unknown, accept
	if a.Type == rules.ARG_UNKNOWN {
		return true
	}

	// t is interface{} or string, accept
	if t.IsEmptyIface() || t.Kind == types.STRING {
		return true
	}

	// both are booleans, accept
	if a.Type == rules.ARG_BOOL && t.Kind == types.BOOL {
		return true
	}

	// t is float and arg is numeric, accept
	if a.IsNumeric() && t.Kind.IsFloat() {
		return true
	}

	// both are integers, accept
	if a.Type == rules.ARG_INT && t.Kind.IsInteger() {
		return true
	}

	// t is unsigned and arg is not negative, accept
	if a.IsUInt() && t.Kind.IsUnsigned() {
		return true
	}

	// arg is string & a Go string literal can be converted to t, accept
	if a.Type == rules.ARG_STRING {
		tt := &types.Type{Kind: types.STRING}
		if tt.IsConvertibleTo(t) {
			return true
		}
	}
	return false
}
