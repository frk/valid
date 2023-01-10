package generate

import (
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
