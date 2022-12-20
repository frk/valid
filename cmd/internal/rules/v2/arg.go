package rules

import (
	"github.com/frk/valid/cmd/internal/types"
)

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

// CanAssignTo reports whether or not the Arg can be assigned to the Go type
// represented by t. The frefs, if provided, is used to resolve ARG_FIELD args.
func (a *Arg) CanAssignTo(t *types.Type, frefs map[*Arg]*types.StructField) bool {
	if a.IsFieldRef() {
		if f, ok := frefs[a]; ok && f != nil {
			// If the type is a pointer and the referenced field's type
			// is identical to the pointer's base, then accept because the
			// generator can use deref expression.
			if t.Kind == types.PTR {
				if t.Elem.Type.IsIdenticalTo(f.Obj.Type) {
					return true
				}
			}

			return f.Obj.Type.IsConvertibleTo(t)
		}
		return false
	}

	// t is interface{} or string, accept
	if t.IsEmptyIface() || t.Kind == types.STRING {
		return true
	}

	// arg is unknown, accept
	if a.Type == ARG_UNKNOWN {
		return true
	}

	// both are booleans, accept
	if t.Kind == types.BOOL && a.Type == ARG_BOOL {
		return true
	}

	// t is float and arg is numeric, accept
	if t.Kind.IsFloat() && a.IsNumeric() {
		return true
	}

	// both are integers, accept
	if t.Kind.IsInteger() && a.Type == ARG_INT {
		return true
	}

	// t is unsigned and arg is not negative, accept
	if t.Kind.IsUnsigned() && a.IsUInt() {
		return true
	}

	// arg is string & a Go string literal can be converted to t, accept
	if a.Type == ARG_STRING {
		tt := &types.Type{Kind: types.STRING}
		if tt.IsConvertibleTo(t) {
			return true
		}
	}
	return false
}
