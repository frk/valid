package rules

type Spec struct {
	// The unique name of the rule.
	Name string
	// The kind of the rule.
	Kind Kind
	// Func is set to the function's identifier when the spec's
	// Kind is FUNCTION/METHOD/PREPROC, otherwise it will be nil.
	Func *FuncIdent
	// ArgMin and ArgMax define bounds of allowed
	// number of arguments for the rule.
	ArgMin, ArgMax int
	// The rule's pre-declared argument options.
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

type FuncIdent struct {
	// The function package's import path. Will be
	// empty when the FuncIdent identifies a method.
	Pkg string
	// The name of the function/method.
	Name string
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

// ErrSpec returns the rule's ErrSpec.
func (r *Rule) ErrSpec() ErrSpec {
	errSpec := r.Spec.Err
	if len(r.Spec.ErrOpts) > 0 && len(r.Args) > 0 {
		var key string
		for _, a := range r.Args {
			key += ":"
			if len(a.Value) > 0 {
				key += "x"
			}
		}

		key = key[1:]
		if specOpt, ok := r.Spec.ErrOpts[key]; ok {
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
