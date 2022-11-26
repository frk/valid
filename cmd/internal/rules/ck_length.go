package rules

import (
	"strconv"

	"github.com/frk/valid/cmd/internal/gotype"
)

// checkLength checks that the Node and Rule
// can be used to generate a "length" rule.
func (c *Checker) lengthCheck(n *Node, r *Rule) error {
	switch r.Name {
	case "len":
		// tn must have length
		if !n.Type.HasLength() {
			return &Error{C: ERR_LENGTH_NOLEN, ty: n.Type, r: r}
		}
	case "runecount":
		// tn must be string kind or byte slice
		if n.Type.Kind != gotype.K_STRING && (n.Type.Kind != gotype.K_SLICE || !n.Type.Elem.Type.IsByte) {
			return &Error{C: ERR_LENGTH_NORUNE, ty: n.Type, r: r}
		}
	}

	// Make sure that at least one arg was provided.
	if (len(r.Args) == 2 && r.Args[0].IsEmpty() && r.Args[1].IsEmpty()) ||
		(len(r.Args) == 1 && r.Args[0].IsEmpty()) {
		return &Error{C: ERR_LENGTH_NOARG, ty: n.Type, r: r}
	}

	// If two args were provided and both are uint values, then make
	// sure that those values represent valid lower to upper bounds.
	if len(r.Args) == 2 && r.Args[0].IsUInt() && r.Args[1].IsUInt() {
		bounds := [2]uint64{}
		for i, a := range r.Args {
			u64, err := strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				panic(err) // this shouldn't happen
			}
			bounds[i] = u64
		}

		if bounds[0] >= bounds[1] {
			return &Error{C: ERR_LENGTH_BOUNDS, ty: n.Type, r: r}
		}
		return nil
	}

	for _, a := range r.Args {
		switch {
		case a.Type == ARG_UNKNOWN:
			continue

		// If the argument is a field then make sure that that field's
		// type can be converted to an int, which is the type of the
		// return value of the `len` and `ut8.RuneCount` functions.
		case a.Type == ARG_FIELD_ABS || a.Type == ARG_FIELD_REL:
			tt := &gotype.Type{Kind: gotype.K_INT}
			ft := c.KeyMap[a.Value].Type.Type // field's type
			if !tt.CanConvert(ft) {
				return &Error{C: ERR_LENGTH_ARGTYPE, r: r, ra: a}
			}
			continue

		case !a.IsUInt():
			return &Error{C: ERR_LENGTH_ARGTYPE, r: r, ra: a}
		}
	}

	// TODO-maybe: based on arg combination link the ErrOpts[<args>] to the rule.
	return nil
}
