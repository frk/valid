package rules

import (
	"strconv"
)

// rangeCheck checks that the Node and Rule
// can be used to generate a "range" rule.
func (c *Checker) rangeCheck(n *Node, r *Rule) error {
	// n must be numeric
	if !n.Type.Kind.IsNumeric() {
		return &Error{C: ERR_RANGE_TYPE, ty: n.Type, r: r}
	}

	// Make sure that both args were provided.
	if r.Args[0].IsEmpty() || r.Args[1].IsEmpty() {
		return &Error{C: ERR_RANGE_NOARG, ty: n.Type, r: r}
	}

	// The rule's args must be comparable to n.
	for _, a := range r.Args {
		if !c.canConvertRuleArg(n.Type, a) {
			return &Error{C: ERR_RANGE_ARGTYPE, ty: n.Type, r: r, ra: a}
		}
	}

	// If both args are numeric constants, then make sure
	// that the values represent valid upper to lower bounds.
	if r.Args[0].IsNumeric() && r.Args[1].IsNumeric() {
		bounds := [2]float64{}
		for i := range r.Args {
			f64, err := strconv.ParseFloat(r.Args[i].Value, 64)
			if err != nil {
				panic(err) // shouldn't happen
			}
			bounds[i] = f64
		}
		if bounds[0] >= bounds[1] {
			return &Error{C: ERR_RANGE_BOUNDS, ty: n.Type, r: r}
		}
	}

	return nil
}
