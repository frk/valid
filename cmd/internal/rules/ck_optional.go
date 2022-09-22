package rules

// optionalCheck checks that the use of the optional rule is sound.
func (c *Checker) optionalCheck(n *Node, r *Rule) error {
	if r.Name == "optional" {
		for _, r2 := range n.IsRules {
			if r2.Name == "required" || r2.Name == "notnil" {
				return &Error{C: ERR_OPTIONAL_CONFLICT, ty: n.Type, r: r, r2: r2}
			}
		}
	}
	if r.Name == "omitnil" {
		for _, r2 := range n.IsRules {
			if r2.Name == "notnil" {
				return &Error{C: ERR_OPTIONAL_CONFLICT, ty: n.Type, r: r, r2: r2}
			}
		}
	}
	return nil
}
