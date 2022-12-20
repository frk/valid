package spec

import (
	stdtypes "go/types"
	"strings"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"

	"gopkg.in/yaml.v3"
)

// Load loads the rule specs for included rules implemented in
// the github.com/frk/valid package and then initializes custom
// rules from the given config.
//
// Load should be invoked only once
// and before starting the first rule-check.
func Load(cfg config.Config, a *search.AST) error {
	if err := LoadIncludedSpecs(a); err != nil {
		return err
	}
	return LoadCustomSpecs(cfg, a)
}

// LoadIncludedSpecs loads the rule specs for function
// rules implemented by the github.com/frk/valid package.
func LoadIncludedSpecs(a *search.AST) error {
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
		_included[s.Name] = s
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
		_custom[s.Name] = s
	}
	return nil
}

func specFromFunc(a *search.AST, f *stdtypes.Func, rs *config.RuleSpec) (*Spec, error) {
	if rs.Name == "" {
		return nil, &Error{C: E_CONFIG_NONAME, a: a, ft: f, rs: rs}
	}
	if _reserved[rs.Name] {
		return nil, &Error{C: E_CONFIG_RESERVED, a: a, ft: f, rs: rs}
	}

	fn := types.AnalyzeFunc(f, a)
	ty := fn.Type
	jop := _joinOps[rs.JoinOp]

	specKind := FUNCTION
	if strings.HasPrefix(rs.Name, "pre:") {
		specKind = PREPROC
	}

	switch specKind {
	case FUNCTION:
		// Make sure the function's signature is alright.
		if len(ty.In) < 1 || (len(ty.Out) != 1 && len(ty.Out) != 2) {
			return nil, &Error{C: E_CONFIG_FUNCTYPE, a: a, ft: f, rs: rs}
		}
		if ty.Out[0].Type.Kind != types.BOOL || (len(ty.Out) > 1 && !ty.Out[1].Type.IsGoError()) {
			return nil, &Error{C: E_CONFIG_FUNCTYPE, a: a, ft: f, rs: rs}
		}
	case PREPROC:
		// Make sure the function's signature is alright.
		if len(ty.In) < 1 || (len(ty.Out) != 1 || !ty.In[0].Type.IsIdenticalTo(ty.Out[0].Type)) {
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
		if !isValidNumberOfArgs(nargs, ty, jop) {
			return nil, &Error{C: E_CONFIG_ARGNUM, a: a, ft: f, rs: rs}
		}
	}
	// If arg bounds were specified in the configuration then make
	// sure that they are compatible with the function's signature.
	if min := rs.ArgMin; min != nil {
		if !isValidNumberOfArgs(int(*min), ty, jop) {
			return nil, &Error{C: E_CONFIG_ARGBOUNDS, a: a, ft: f, rs: rs}
		}
	}
	if max := rs.ArgMax; max != nil {
		if !isValidNumberOfArgs(*max, ty, jop) {
			return nil, &Error{C: E_CONFIG_ARGBOUNDS, a: a, ft: f, rs: rs}
		}
	}

	// Build the rule spec from the func's type & the config.
	spec := new(Spec)
	spec.Name = rs.Name
	spec.Kind = specKind
	spec.FName = f.Name()
	spec.FType = ty
	spec.JoinOp = jop
	spec.Err = ErrSpec(rs.Error)

	// the "re" (regexp) rule should use raw strings for arguments
	if spec.Name == "re" {
		spec.UseRawString = true
	}

	// resolve the arg bounds
	spec.ArgMin = len(ty.In) - 1 // -1 for the field argument
	spec.ArgMax = len(ty.In) - 1 // -1 for the field argument
	if rs.ArgMin != nil {
		spec.ArgMin = int(*rs.ArgMin)
	} else if ty.IsVariadic {
		spec.ArgMin -= 1
	}
	if rs.ArgMax != nil {
		spec.ArgMax = *rs.ArgMax
	} else if ty.IsVariadic || jop > 0 {
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
		if err := validateArgOptsAsFuncParams(spec); err != nil {
			e := err.(*Error)
			e.a, e.ft, e.rs = a, f, rs
			return nil, e
		}
	}

	return spec, nil
}

func isValidNumberOfArgs(nargs int, ft *types.Type, joinOp JoinOp) bool {
	nparams := len(ft.In) - 1 // -1 for the field argument

	// If not variadic & not joinable, require exact match.
	if !ft.IsVariadic && joinOp == 0 && nargs != nparams {
		return false
	}
	// If not variadic but joinable, require that the args
	// can be spread evenly across multiple joined instances.
	if !ft.IsVariadic && joinOp > 1 && (nargs%nparams) != 0 {
		return false
	}
	// If variadic, require only that enough args were specified.
	if ft.IsVariadic && nargs < nparams-1 {
		return false
	}

	return true
}

// NOTE(mkopriva): while not identical to, this does share some logic
// with checkRuleArgsAsFuncParams. If something in that logic changes
// here, make sure to apply that change to the "cousin" as well.
func validateArgOptsAsFuncParams(s *Spec) error {
	params := s.FType.In
	firstParamIsField := (s.Kind != METHOD)

	if firstParamIsField && (len(s.FType.In) > 1 || !s.FType.IsVariadic) {
		// Drop the first parameter (the field) if there is more
		// than one parameter, or if the function is non-variadic.
		params = s.FType.In[1:]
	}

	for i, argOpts := range s.ArgOpts {
		var j int
		var p *types.Var

		switch last := len(params) - 1; {
		case i < last || (i == last && !s.FType.IsVariadic):
			j = i
			p = params[j]

		case i >= last && s.FType.IsVariadic:
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
			if !a.CanAssignTo(p.Type, nil) {
				a, i, k := a, i, k // copy
				return &Error{C: E_CONFIG_ARGTYPE,
					rca: &a, rcai: &i, rcak: &k, fp: p, fpi: &j,
				}
			}
		}
	}
	return nil
}

var _joinOps = [...]JoinOp{
	config.JOIN_NOT: JOIN_NOT,
	config.JOIN_AND: JOIN_AND,
	config.JOIN_OR:  JOIN_OR,
}

var _scalarargs = [...]rules.ArgType{
	config.NIL:    rules.ARG_UNKNOWN,
	config.BOOL:   rules.ARG_BOOL,
	config.INT:    rules.ARG_INT,
	config.FLOAT:  rules.ARG_FLOAT,
	config.STRING: rules.ARG_STRING,
}
