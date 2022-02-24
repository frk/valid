package rules

// comparableCheck checks that the rule's arguments can be converted
// to a Go value of a type that is *comparable* with the node's type.
func (c *Checker) comparableCheck(n *Node, r *Rule) error {
	for _, arg := range r.Args {
		if !c.canConvertRuleArg(n.Type, arg) {
			return &Error{C: ERR_ARG_BADCMP, ty: n.Type, r: r, ra: arg}
		}
	}
	return nil
}
