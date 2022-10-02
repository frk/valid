package rules

// methodCheck checks that the node can be used to
// generate a call to the rule's method, e.g. "IsValid".
func (c *Checker) methodCheck(n *Node, r *Rule) error {
	// Check that n implements a method that matches the rule type.
	var methodFound bool
loop:
	for _, m := range n.Type.Methods {
		if m.Name != r.Spec.FName {
			continue loop
		}

		if len(m.Type.In) != len(r.Spec.FType.In) {
			continue loop
		}
		if len(m.Type.Out) != len(r.Spec.FType.Out) {
			continue loop
		}

		for i := range m.Type.In {
			if !m.Type.In[i].Type.IsIdentical(r.Spec.FType.In[i].Type) {
				continue loop
			}
		}
		for i := range m.Type.Out {
			if !m.Type.Out[i].Type.IsIdentical(r.Spec.FType.Out[i].Type) {
				continue loop
			}
		}

		methodFound = true
		break loop
	}
	if !methodFound {
		return &Error{C: ERR_METHOD_TYPE, ty: n.Type, r: r}
	}

	// Check that the arguments specified in the rule tag can be used
	// as the arguments for the corresponding method parameters.
	if err := c.checkRuleArgsAsFuncParams(r); err != nil {
		// TODO would be nice to have a test case for this
		// but for that I'd need to add other METHOD rules
		// since currently there's only one and that one
		// takes no input arguments.
		return err
	}

	return nil
}
