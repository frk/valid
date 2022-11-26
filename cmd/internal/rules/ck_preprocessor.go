package rules

import (
	"github.com/frk/valid/cmd/internal/gotype"
)

func (c *Checker) preprocessorCheck(n *Node, r *Rule) error {
	// Make sure that an instance of n can be passed
	// to the function as its first argument.
	paramType := r.Spec.FType.In[0].Type
	if r.Spec.FType.IsVariadic && len(r.Spec.FType.In) == 1 {
		paramType = paramType.Elem.Type
	}
	if paramType.CanAssign(n.Type) == gotype.ASSIGNMENT_INVALID {
		return &Error{C: ERR_PREPROC_INTYPE, ty: n.Type, r: r}
	}

	// Make sure that an instance of the function
	// output type can be assign to n.
	outputType := r.Spec.FType.Out[0].Type
	if n.Type.CanAssign(outputType) == gotype.ASSIGNMENT_INVALID {
		return &Error{C: ERR_PREPROC_OUTTYPE, ty: n.Type, r: r}
	}

	// Check that the arguments specified in the rule tag can be used
	// as the arguments for the corresponding function's parameters.
	if err := c.checkRuleArgsAsFuncParams(r); err != nil {
		return c.err(err, errOpts{C: ERR_PREPROC_ARGTYPE, ty: n.Type})
	}
	return nil
}
