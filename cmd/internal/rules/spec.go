package rules

import (
	"go/types"
	"strings"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/search"

	"gopkg.in/yaml.v3"
)

type SpecKind uint

func (k SpecKind) String() string {
	if int(k) < len(_speckindstring) {
		return _speckindstring[k]
	}
	return "<invalid>"
}

const (
	_ SpecKind = iota

	REQUIRED   // required, notnil
	COMPARABLE // =, !=
	ORDERED    // >, >=, <, <=, min, max
	LENGTH     // len, runecount
	RANGE      // rng
	ENUM       // enum
	FUNCTION   // <custom/builtin/included func rules>
	METHOD     // isvalid (implicit), ...

	// "modifiers"
	OPTIONAL // omitnil [is the default rule for pointers] (ptr only), optional (ptr & base)
	NOGUARD  // nonilguard
	REMOVE   // -isvalid

	// "preprocessors"
	PREPROC // <custom/builtin/included func rules>
)

var _speckindstring = [...]string{
	REQUIRED:   "REQUIRED",
	COMPARABLE: "COMPARABLE",
	ORDERED:    "ORDERED",
	LENGTH:     "LENGTH",
	RANGE:      "RANGE",
	ENUM:       "ENUM",
	FUNCTION:   "FUNCTION",
	METHOD:     "METHOD",
	OPTIONAL:   "OPTIONAL",
	NOGUARD:    "NOGUARD",
	REMOVE:     "REMOVE",
	PREPROC:    "PREPROC",
}

type Spec struct {
	// The unique name of the rule.
	Name string
	// The kind of the rule.
	Kind SpecKind
	// Kind=FUNCTION only, the function's identifier.
	FName string
	// Kind=FUNCTION only, the function's type.
	FType *gotype.Type
	// ArgMin and ArgMax define bounds of allowed
	// number of arguments for the rule.
	ArgMin, ArgMax int
	// Kind=FUNCTION only, the rule's pre-declared argument options.
	ArgOpts []map[string]Arg
	// The join operator that should be used for joining
	// multiple instances of the rule into a single one.
	JoinOp JoinOp
	// The spec for the error message that the
	// generator should produce for the rule.
	Err ErrSpec
	// The error options for specific argument combinations
	ErrOpts map[string]ErrSpec
	// Indicates that the generated code should use raw
	// strings for any string arguments of the rule.
	UseRawString bool
}

type ErrSpec struct {
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

var _joinOps = [...]JoinOp{
	config.JOIN_NOT: JOIN_NOT,
	config.JOIN_AND: JOIN_AND,
	config.JOIN_OR:  JOIN_OR,
}

// InitSpecs loads the rule specs for function rules implemented by
// the github.com/frk/valid package and then initializes custom rules
// from the given config.
//
// InitSpecs should be invoked only once
// and before starting the first rule-check.
func InitSpecs(cfg config.Config, a *search.AST) error {
	if err := loadIncludedSpecs(a); err != nil {
		return err
	}
	return initCustomSpecs(cfg, a)
}

// loadIncludedSpecs loads the rule specs for function
// rules implemented by the github.com/frk/valid package.
func loadIncludedSpecs(a *search.AST) error {
	// loads functions from the "github.com/frk/valid" package
	return search.FindIncludedFuncs(a, func(ftyp *types.Func, rawCfg []byte) error {
		cfg := new(config.RuleSpec)
		if err := yaml.Unmarshal(rawCfg, cfg); err != nil {
			return &Error{C: ERR_CONFIG_INVALID, a: a, ft: ftyp, err: err}
		}
		spec, err := specFromFunc(a, ftyp, cfg)
		if err != nil {
			return err
		}
		_included[spec.Name] = spec
		return nil
	})
}

// initCustomSpecs initializes custom rules from the given config.
func initCustomSpecs(cfg config.Config, a *search.AST) error {
	for _, rc := range cfg.Rules {
		ftyp, rawCfg, err := search.FindFunc(rc.Func.Pkg, rc.Func.Name, a)
		if err != nil {
			return &Error{C: ERR_CONFIG_FUNCSEARCH, a: a, c: &cfg, err: err}
		}

		// Use the config from the function's documentation
		// if none was provided in the config file itself.
		if rc.Rule == nil && len(rawCfg) > 0 {
			rc.Rule = new(config.RuleSpec)
			if err := yaml.Unmarshal(rawCfg, rc.Rule); err != nil {
				return &Error{C: ERR_CONFIG_INVALID, a: a, c: &cfg, ft: ftyp, err: err}
			}
		}
		if rc.Rule == nil {
			return &Error{C: ERR_CONFIG_MISSING, a: a, c: &cfg, ft: ftyp, err: err}
		}

		spec, err := specFromFunc(a, ftyp, rc.Rule)
		if err != nil {
			return extendError(err, func(e *Error) { e.c = &cfg })
		}
		_custom[spec.Name] = spec
	}
	return nil
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

func specFromFunc(a *search.AST, f *types.Func, rs *config.RuleSpec) (*Spec, error) {
	if rs.Name == "" {
		return nil, &Error{C: ERR_CONFIG_NONAME, a: a, ft: f, rs: rs}
	}
	if _reserved[rs.Name] {
		return nil, &Error{C: ERR_CONFIG_RESERVED, a: a, ft: f, rs: rs}
	}

	an := gotype.NewAnalyzer(f.Pkg())
	ty := an.Object(f)
	jop := _joinOps[rs.JoinOp]

	specKind := FUNCTION
	if strings.HasPrefix(rs.Name, "pre:") {
		specKind = PREPROC
	}

	switch specKind {
	case FUNCTION:
		// Make sure the function's signature is alright.
		if len(ty.In) < 1 || (len(ty.Out) != 1 && len(ty.Out) != 2) {
			return nil, &Error{C: ERR_CONFIG_FUNCTYPE, a: a, ft: f, rs: rs}
		}
		if ty.Out[0].Type.Kind != gotype.K_BOOL || (len(ty.Out) > 1 && !ty.Out[1].Type.IsGoError()) {
			return nil, &Error{C: ERR_CONFIG_FUNCTYPE, a: a, ft: f, rs: rs}
		}
	case PREPROC:
		// Make sure the function's signature is alright.
		if len(ty.In) < 1 || (len(ty.Out) != 1 || !ty.In[0].Type.IsIdentical(ty.Out[0].Type)) {
			return nil, &Error{C: ERR_CONFIG_PREFUNCTYPE, a: a, ft: f, rs: rs}
		}

		// Joins & Errors are NOT supported for preprocs.
		// These two could probably be just warnings.
		if jop > 0 {
			return nil, &Error{C: ERR_CONFIG_PREPROCJOIN, a: a, ft: f, rs: rs}
		}
		if rs.Error != (config.RuleErrorConfig{}) {
			return nil, &Error{C: ERR_CONFIG_PREPROCERROR, a: a, ft: f, rs: rs}
		}
	}

	// If args were specified in the configuration then make sure
	// that their number is compatible with the function's signature.
	if nargs := len(rs.Args); nargs > 0 {
		if !isValidNumberOfArgs(nargs, ty, jop) {
			return nil, &Error{C: ERR_CONFIG_ARGNUM, a: a, ft: f, rs: rs}
		}
	}
	// If arg bounds were specified in the configuration then make
	// sure that they are compatible with the function's signature.
	if min := rs.ArgMin; min != nil {
		if !isValidNumberOfArgs(int(*min), ty, jop) {
			return nil, &Error{C: ERR_CONFIG_ARGBOUNDS, a: a, ft: f, rs: rs}
		}
	}
	if max := rs.ArgMax; max != nil {
		if !isValidNumberOfArgs(*max, ty, jop) {
			return nil, &Error{C: ERR_CONFIG_ARGBOUNDS, a: a, ft: f, rs: rs}
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
		argOpts := make(map[string]Arg)
		if def := a.Default; def != nil {
			argOpts[""] = Arg{
				Type:  _scalarargs[def.Type],
				Value: def.Value,
			}
		}
		for _, opt := range a.Options {
			argOpts[opt.Value.Value] = Arg{
				Type:  _scalarargs[opt.Value.Type],
				Value: opt.Value.Value,
			}
			if len(opt.Alias) > 0 {
				argOpts[opt.Alias] = Arg{
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

func isValidNumberOfArgs(nargs int, ft *gotype.Type, joinOp JoinOp) bool {
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
		var p *gotype.Var

		switch last := len(params) - 1; {
		case i < last || (i == last && !s.FType.IsVariadic):
			j = i
			p = params[j]

		case i >= last && s.FType.IsVariadic:
			j = last
			p = &gotype.Var{
				Name: params[j].Name,
				Type: params[j].Type.Elem,
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
				return &Error{C: ERR_CONFIG_ARGTYPE,
					rca: &a, rcai: &i, rcak: &k, fp: p, fpi: &j,
				}
			}
		}
	}
	return nil
}

func (s *Spec) getFuncParamByArgIndex(i int) (p *gotype.Var, pi int) {
	params := s.FType.In
	firstParamIsField := (s.Kind != METHOD)

	if firstParamIsField && (len(s.FType.In) > 1 || !s.FType.IsVariadic) {
		// Drop the first parameter (the field) if there is more
		// than one parameter, or if the function is non-variadic.
		params = s.FType.In[1:]
	}

	switch last := len(params) - 1; {
	case i < last || (i == last && !s.FType.IsVariadic):
		return params[i], i

	case i >= last && s.FType.IsVariadic:
		p := params[last]
		return &gotype.Var{Name: p.Name, Type: p.Type.Elem}, last

	case i > last && s.JoinOp > 0:
		if len(s.ArgOpts)-i >= len(params) {
			pi = i % len(params)
			return params[pi], pi
		}
	}

	return nil, -1
}
