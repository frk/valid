package spec

import (
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

type Spec struct {
	// The unique name of the rule.
	Name string
	// The kind of the rule.
	Kind Kind
	// Kind=FUNCTION only, the function's identifier.
	FName string
	// Kind=FUNCTION only, the function's type.
	FType *types.Type
	// ArgMin and ArgMax define bounds of allowed
	// number of arguments for the rule.
	ArgMin, ArgMax int
	// Kind=FUNCTION only, the rule's pre-declared argument options.
	ArgOpts []map[string]rules.Arg
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

func Get(name string) *Spec {
	if spec, ok := _custom[name]; ok {
		return spec
	}
	if spec, ok := _included[name]; ok {
		return spec
	}
	if spec, ok := _builtin[name]; ok {
		return spec
	}
	return nil
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

// ErrSpec returns the Spec's ErrSpec for the provided rule.
func (s *Spec) ErrSpec(r *rules.Rule) ErrSpec {
	errSpec := s.Err
	if len(s.ErrOpts) > 0 && len(r.Args) > 0 {
		var key string
		for _, a := range r.Args {
			key += ":"
			if len(a.Value) > 0 {
				key += "x"
			}
		}

		key = key[1:]
		if specOpt, ok := s.ErrOpts[key]; ok {
			errSpec = specOpt
		}
	}
	return errSpec
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

func (s *Spec) GetFuncParamByArgIndex(i int) (p *types.Var, pi int) {
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
		return &types.Var{Name: p.Name, Type: p.Type.Elem.Type}, last

	case i > last && s.JoinOp > 0:
		if len(s.ArgOpts)-i >= len(params) {
			pi = i % len(params)
			return params[pi], pi
		}
	}

	return nil, -1
}
