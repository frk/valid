package rules

// Rule represents the rule as parsed from a struct tag.
type Rule struct {
	// The name of the rule.
	Name string
	// The arguments of the rule.
	Args []*Arg
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

// // IsBasic reports whether or not the rule is a basic rule.
// // What "basic" means at this point is that the rule is NOT a
// // function that returns an error as its second return value.
// func (r Rule) IsBasic() bool {
// 	return r.Spec.Kind != FUNCTION ||
// 		!r.Spec.FType.CanError()
// }

// // CanAssignTo reports whether or not the Arg can be assigned to the Go type
// // represented by t. The keyMap, if provided, is used to resolve ARG_FIELD args.
// func (a *Arg) CanAssignTo(t *xtypes.Type, keyMap map[string]*FieldNode) bool {
// 	if a.Type == ARG_FIELD_ABS || a.Type == ARG_FIELD_REL {
// 		if f, ok := keyMap[a.Value]; ok && f != nil {
// 			// can use the addr, accept
// 			if t.PtrOf(f.Type.Type) {
// 				return true
// 			}
// 			return t.CanAssign(f.Type.Type) != xtypes.ASSIGNMENT_INVALID
// 		}
// 		return false
// 	}
//
// 	// t is interface{} or string, accept
// 	if t.IsEmptyInterface() || t.Kind == xtypes.K_STRING {
// 		return true
// 	}
//
// 	// arg is unknown, accept
// 	if a.Type == ARG_UNKNOWN {
// 		return true
// 	}
//
// 	// both are booleans, accept
// 	if t.Kind == xtypes.K_BOOL && a.Type == ARG_BOOL {
// 		return true
// 	}
//
// 	// t is float and option is numeric, accept
// 	if t.Kind.IsFloat() && (a.Type == ARG_INT || a.Type == ARG_FLOAT) {
// 		return true
// 	}
//
// 	// both are integers, accept
// 	if t.Kind.IsInteger() && a.Type == ARG_INT {
// 		return true
// 	}
//
// 	// t is unsigned and option is not negative, accept
// 	if t.Kind.IsUnsigned() && a.Type == ARG_INT && a.Value[0] != '-' {
// 		return true
// 	}
//
// 	// arg is string & string can be converted to t, accept
// 	if a.Type == ARG_STRING && (t.Kind == xtypes.K_STRING || (t.Kind == xtypes.K_SLICE &&
// 		t.Elem.Type.Name == "" && (t.Elem.Type.Kind == xtypes.K_UINT8 || t.Elem.Type.Kind == xtypes.K_INT32))) {
// 		return true
// 	}
//
// 	return false
// }

////////////////////////////////////////////////////////////////////////////////
// helper type

type Set struct {
	Is  List
	Pre List
}

type List struct {
	Rules []*Rule
}

// Contains reports whether or not the List
// contains a rule with the given name.
func (ls List) Contains(name string) bool {
	for _, r := range ls.Rules {
		if r.Name == name {
			return true
		}
	}
	return false
}

// Add adds the given rule to the List.
func (ls *List) Add(r *Rule) {
	for _, e := range ls.Rules {
		if e.Name == r.Name {
			return
		}
	}
	ls.Rules = append(ls.Rules, r)
}

// Remove removes a rule from the List by name.
func (ls *List) Remove(name string) {
	for i, r := range ls.Rules {
		if r.Name == name {
			ls.Rules = append(append(
				[]*Rule{},
				ls.Rules[:i]...),
				ls.Rules[i+1:]...)
			return
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// the following are convenience methods intended primarily for the generator

// Empty reports whether or not the slice is empty.
func (ls List) Empty() bool {
	return len(ls.Rules) == 0
}

// One reports whether or not there is exactly one rule in the slice.
func (ls List) One() bool {
	return len(ls.Rules) == 1
}

// Many reports whether or not there are multiple rules in the slice.
func (ls List) Many() bool {
	return len(ls.Rules) > 1
}
