package rules

// enumCheck checks that the "enum" rule can be applied to the node's type.
func (c *Checker) enumCheck(n *Node, r *Rule) error {
	if len(n.Type.Name) == 0 {
		return &Error{C: ERR_ENUM_NONAME, ty: n.Type, r: r}
	}
	if !n.Type.Kind.IsBasic() {
		return &Error{C: ERR_ENUM_KIND, ty: n.Type, r: r}
	}

	// Make sure that the type actually has some accessible constants declared.
	consts := c.an.Consts(n.Type, c.ast)
	if len(consts) == 0 {
		return &Error{C: ERR_ENUM_NOCONST, ty: n.Type, r: r}
	}

	c.Info.EnumMap[n.Type] = consts
	return nil
}
