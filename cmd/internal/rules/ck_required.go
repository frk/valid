package rules

// requiredCheck checks that the REQUIRED rule can be applied to the node.
func (c *Checker) requiredCheck(n *Node, r *Rule) error {
	switch r.Name {
	case "required":
		// nothing to do
	case "notnil":
		// type must be nilable
		if !n.Type.IsNilable() {
			return &Error{C: ERR_NOTNIL_TYPE, ty: n.Type, r: r}
		}
	}
	return nil
}
