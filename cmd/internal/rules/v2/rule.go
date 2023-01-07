package rules

// Rule represents the rule as parsed from a struct tag.
type Rule struct {
	// The name of the rule.
	Name string
	// The arguments of the rule.
	Args []*Arg
	// The spec associated with the rule. Note that this is not
	// populated by the parser but instead by the rules/checker.
	Spec *Spec
}

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

func (r *Rule) Is(kk ...Kind) bool {
	for _, k := range kk {
		if r.Spec.Kind == k {
			return true
		}
	}
	return false
}

type List []*Rule

func (l List) Has(kk ...Kind) bool {
	for _, r := range l {
		if r.Is(kk...) {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////

// Arg represents an argument for a rule.
type Arg struct {
	// The type of the argument.
	Type ArgType
	// The literal string representation of the value.
	Value string
}

// IsUInt reports whether or not Arg value is a valid uint "candidate".
func (a *Arg) IsUInt() bool {
	return a.Type == ARG_INT && a.Value[0] != '-'
}

// IsEmpty reports whether or not Arg value is empty.
func (a *Arg) IsEmpty() bool {
	return a.Value == ""
}

// IsNumeric reports whether or not Arg value is numeric.
func (a *Arg) IsNumeric() bool {
	return a.Type == ARG_INT || a.Type == ARG_FLOAT
}

// IsFieldRef reports whether or not Arg value is a field reference.
func (a *Arg) IsFieldRef() bool {
	return a.Type == ARG_FIELD_ABS || a.Type == ARG_FIELD_REL
}

////////////////////////////////////////////////////////////////////////////////

// ArgType indicates the type of a rule argument value.
type ArgType uint

func (t ArgType) String() string {
	if int(t) < len(_argtypestring) {
		return _argtypestring[t]
	}
	return "<invalid>"
}

const (
	ARG_UNKNOWN ArgType = iota
	ARG_BOOL
	ARG_INT
	ARG_FLOAT
	ARG_STRING
	ARG_FIELD_ABS
	ARG_FIELD_REL
)

var _argtypestring = [...]string{
	ARG_UNKNOWN:   "<unknown>",
	ARG_BOOL:      "bool",
	ARG_INT:       "int",
	ARG_FLOAT:     "float",
	ARG_STRING:    "string",
	ARG_FIELD_ABS: "<field_abs>",
	ARG_FIELD_REL: "<field_rel>",
}
