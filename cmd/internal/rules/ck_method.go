package rules

import (
	"github.com/frk/valid/cmd/internal/errors"
)

// methodCheck ...
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
			if !m.Type.In[i].Type.Identical(r.Spec.FType.In[i].Type) {
				continue loop
			}
		}
		for i := range m.Type.Out {
			if !m.Type.Out[i].Type.Identical(r.Spec.FType.Out[i].Type) {
				continue loop
			}
		}

		methodFound = true
		break loop
	}
	if !methodFound {
		return errors.TODO("methodCheck: method not found")
	}

	// Check that the arguments specified in the rule tag can be used
	// as the arguments for the corresponding method parameters.
	if err := c.checkRuleArgsAsFuncParams(r); err != nil {
		return err
	}

	return nil
}
