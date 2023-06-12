package rules

import (
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
