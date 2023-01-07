package specs

import (
	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/search"
	"github.com/frk/valid/cmd/internal/types"
)

// Load loads the rule specs for included rules implemented in
// the github.com/frk/valid package and then initializes custom
// rules from the given config.
//
// Load should be invoked only once
// and before starting the first rule-check.
func Load(cfg config.Config, a *search.AST) error {
	if err := loadBuiltins(a); err != nil {
		return err
	}
	if err := loadIncludedSpecs(a); err != nil {
		return err
	}
	return LoadCustomSpecs(cfg, a)
}

// Get returns the rules.Spec with the given name.
func Get(name string) *rules.Spec {
	if s, ok := custom[name]; ok {
		return s
	}
	if s, ok := included[name]; ok {
		return s
	}
	if s, ok := builtin[name]; ok {
		return s
	}
	return nil
}

// GetFunc returns the types.Func associated with the given spec.
func GetFunc(s *rules.Spec) *types.Func {
	return funcMap[s]
}

// AddCustomSpec is intended to be used by tests in other packages
// that don't normally have write access to the _custom map.
func AddCustomSpec(key string, s *rules.Spec, fn *types.Func) {
	custom[key] = s
	if fn != nil {
		funcMap[s] = fn
	}
}

func GetFuncParamByArgIndex(s *rules.Spec, i int) (p *types.Var, pi int) {
	fn := GetFunc(s)
	params := fn.Type.In
	firstParamIsField := (s.Kind != rules.METHOD)

	if firstParamIsField && (len(fn.Type.In) > 1 || !fn.Type.IsVariadic) {
		// Drop the first parameter (the field) if there is more
		// than one parameter, or if the function is non-variadic.
		params = fn.Type.In[1:]
	}

	switch last := len(params) - 1; {
	case i < last || (i == last && !fn.Type.IsVariadic):
		return params[i], i

	case i >= last && fn.Type.IsVariadic:
		p := params[last]
		return &types.Var{Name: p.Name, Type: p.Type.Elem.Type}, last

	case i > last && s.JoinOp > 0:
		if len(s.ArgOpts)-i >= len(params) {
			pi = i % len(params)
			return params[pi], pi
		}
	}

	return nil, -1
}

// CanReturnError reports whether the function associated with
// the given spec returns an error as its second return value.
func CanReturnError(s *rules.Spec) bool {
	return s.Kind == rules.FUNCTION && GetFunc(s).Type.CanError()
}

// set by loadIncludedSpecs
var _regexpFunc *types.Func

func RegisterRegexpFunc() *types.Func {
	return _regexpFunc
}
