package rules

import (
	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/gotype"
)

// Rule represents the rule as parsed from a struct tag.
type Rule struct {
	// The name of the rule.
	Name string
	// The arguments of the rule.
	Args []*Arg
	// The spec associated with the rule. Note that this is
	// not populated by the parser but instead by the Checker.
	Spec *Spec
	// The context property of the rule.
	// NOTE currently not used by the generator.
	Context string
}

func (r Rule) String() (out string) {
	out = r.Name
	for i := range r.Args {
		out += ":"
		if r.Args[i].Type == ARG_FIELD {
			out += "&"
		}
		out += r.Args[i].Value
	}
	return out
}

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

// CanAssignTo reports whether or not the Arg can be assigned to the Go type
// represented by t. The keyMap, if provided, is used to resolve ARG_FIELD args.
func (a *Arg) CanAssignTo(t *gotype.Type, keyMap map[string]*FieldNode) bool {
	if a.Type == ARG_FIELD {
		if f, ok := keyMap[a.Value]; ok && f != nil {
			// can use the addr, accept
			if t.PtrOf(f.Type.Type) {
				return true
			}
			return t.CanAssign(f.Type.Type) != gotype.ASSIGNMENT_INVALID
		}
		return false
	}

	// t is interface{} or string, accept
	if t.IsEmptyInterface() || t.Kind == gotype.K_STRING {
		return true
	}

	// arg is unknown, accept
	if a.Type == ARG_UNKNOWN {
		return true
	}

	// both are booleans, accept
	if t.Kind == gotype.K_BOOL && a.Type == ARG_BOOL {
		return true
	}

	// t is float and option is numeric, accept
	if t.Kind.IsFloat() && (a.Type == ARG_INT || a.Type == ARG_FLOAT) {
		return true
	}

	// both are integers, accept
	if t.Kind.IsInteger() && a.Type == ARG_INT {
		return true
	}

	// t is unsigned and option is not negative, accept
	if t.Kind.IsUnsigned() && a.Type == ARG_INT && a.Value[0] != '-' {
		return true
	}

	// arg is string & string can be converted to t, accept
	if a.Type == ARG_STRING && (t.Kind == gotype.K_STRING || (t.Kind == gotype.K_SLICE &&
		t.Elem.Name == "" && (t.Elem.Kind == gotype.K_UINT8 || t.Elem.Kind == gotype.K_INT32))) {
		return true
	}

	return false
}

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
	ARG_FIELD
)

var _argtypestring = [...]string{
	ARG_UNKNOWN: "<unknown>",
	ARG_BOOL:    "bool",
	ARG_INT:     "int",
	ARG_FLOAT:   "float",
	ARG_STRING:  "string",
	ARG_FIELD:   "<field>",
}

var _scalarargs = [...]ArgType{
	config.NIL:    ARG_UNKNOWN,
	config.BOOL:   ARG_BOOL,
	config.INT:    ARG_INT,
	config.FLOAT:  ARG_FLOAT,
	config.STRING: ARG_STRING,
}

////////////////////////////////////////////////////////////////////////////////
// helper type

type RuleList []*Rule

// Contains reports whether or not the RuleList
// contains a rule with the given name.
func (ls RuleList) Contains(name string) bool {
	for _, r := range ls {
		if r.Name == name {
			return true
		}
	}
	return false
}

// Add adds the given rule to the RuleList.
func (ls *RuleList) Add(r *Rule) {
	for _, e := range *ls {
		if e.Name == r.Name {
			return
		}
	}
	*ls = append(*ls, r)
}

// Remove removes a rule from the RuleList by name.
func (ls *RuleList) Remove(name string) {
	for i, r := range *ls {
		if r.Name == name {
			*ls = append(append(
				[]*Rule{},
				(*ls)[:i]...),
				(*ls)[i+1:]...)
			return
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// the following are convenience methods intended primarily for the generator

// Empty reports whether or not the slice is empty.
func (ls RuleList) Empty() bool {
	return len(ls) == 0
}

// One reports whether or not there is exactly one rule in the slice.
func (ls RuleList) One() bool {
	return len(ls) == 1
}

// Many reports whether or not there are multiple rules in the slice.
func (ls RuleList) Many() bool {
	return len(ls) > 1
}
