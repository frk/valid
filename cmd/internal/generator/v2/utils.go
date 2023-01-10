package generate

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

// p_end returns the pointer's "end" base object.
func p_end(o *types.Obj) *types.Obj {
	for o.Type.Kind == types.PTR {
		o = o.Type.Elem
	}
	return o
}

// p_rules returns the "pointer rules" from the given object
// and its base objects, all the way to its "end" base.
func p_rules(o *types.Obj) (rr []*rules.Rule) {
	for o.Type.Kind == types.PTR {
		if o.Has(rules.OPTIONAL, rules.REQUIRED) {
			rr = append(rr, o.IsRules[0])
			o.IsRules = o.IsRules[1:]
		}
		o = o.Type.Elem
	}
	if o.Has(rules.OPTIONAL, rules.REQUIRED) {
		rr = append(rr, o.IsRules[0])
		o.IsRules = o.IsRules[1:]
	}
	return rr
}

// can_join_opt reports whether or not the given object can be
// joined in an if-statement directly to a set of optional rules.
func can_join_opt(o *types.Obj) bool {
	if len(o.PreRules) == 0 && len(o.IsRules) == 1 {
		if is_boolean_rule(o.IsRules[0]) {
			return true
		}
	}
	return false
}

// is_boolean_rule reports whether or not the given
// rule's expression evaluates to a single bool value.
func is_boolean_rule(r *rules.Rule) bool {
	if r.Spec.Kind != rules.FUNCTION {
		return true
	}

	out := specs.GetFunc(r.Spec).Type.Out
	if len(out) == 1 && out[0].Type.Kind == types.BOOL {
		return true
	}

	return false
}

// is_multi_expr reports whether or not the given
// rule produces more than one simple expression.
func (g *generator) is_multi_expr(r *rules.Rule) bool {
	switch r.Spec.Kind {
	case rules.COMPARABLE:
		return len(r.Args) > 1
	case rules.LENGTH:
		return len(r.Args) == 2 &&
			r.Args[0].Value != "" &&
			r.Args[1].Value != ""
	case rules.RANGE:
		return true
	case rules.ENUM:
		o := g.info.RuleObjMap[r]
		return len(g.info.EnumMap[o.Type]) > 1
	case rules.FUNCTION:
		return r.Spec.JoinOp > 0 && len(r.Args) > 1
	}
	return false
}
