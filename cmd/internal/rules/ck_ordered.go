package rules

import (
	"github.com/frk/valid/cmd/internal/gotype"
)

// orderedCheck checks that the Node's type is *ordered* and
// that the Rule's arguments can be converted to a Go value of
// a type that is *comparable* with the Node's type.
func (c *Checker) orderedCheck(n *Node, r *Rule) error {
	// the type must be numeric or string
	if !n.Type.Kind.IsNumeric() && n.Type.Kind != gotype.K_STRING {
		return &Error{C: ERR_ORDERED_TYPE, ty: n.Type, r: r}
	}

	// rule args must be comparable with n.Type
	for _, a := range r.Args {
		if !c.canConvertRuleArg(n.Type, a) {
			return &Error{C: ERR_ORDERED_ARGTYPE, ty: n.Type, r: r, ra: a}
		}
	}
	return nil
}
