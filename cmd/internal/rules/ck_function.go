package rules

import (
	"regexp"

	"github.com/frk/valid"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/internal/cldr"
	"github.com/frk/valid/internal/tables"
)

// functionCheck ...
func (c *Checker) functionCheck(n *Node, r *Rule) error {
	// Check if an instance of n can be passed
	// to the function as its first argument.
	p := r.Spec.FType.In[0].Type
	if r.Spec.FType.IsVariadic && len(r.Spec.FType.In) == 1 {
		p = p.Elem.Type
	}
	if p.CanAssign(n.Type) == gotype.ASSIGNMENT_INVALID {
		return &Error{C: ERR_FUNCTION_INTYPE, ty: n.Type, r: r}
	}

	// Check that the arguments specified in the rule tag can be used
	// as the arguments for the corresponding function parameters.
	if err := c.checkRuleArgsAsFuncParams(r); err != nil {
		return c.err(err, errOpts{C: ERR_FUNCTION_ARGTYPE, ty: n.Type})
	}

	// Some included functions accept arguments of a known set of valid
	// values, check that the rule argument's values belong to that set.
	if r.Spec.FType.IsIncluded() {
		if err := c.checkIncludedRuleArgValues(r); err != nil {
			return c.err(err, errOpts{C: ERR_FUNCTION_ARGVALUE, ty: n.Type})
		}
	}

	return nil
}

// NOTE(mkopriva): while not identical to, this does share some logic
// with validateArgOptsAsFuncParams. If something in that logic changes
// here, make sure to apply that change to the "cousin" as well.
func (c *Checker) checkRuleArgsAsFuncParams(r *Rule) error {
	params := r.Spec.FType.In
	firstParamIsField := (r.Spec.Kind != METHOD)

	if firstParamIsField && (len(r.Spec.FType.In) > 1 || !r.Spec.FType.IsVariadic) {
		// Drop the first parameter (the field) if there is more
		// than one parameter, or if the function is non-variadic.
		params = r.Spec.FType.In[1:]
	}

	for i := range r.Args {
		var j int
		var p *gotype.Var

		switch last := len(params) - 1; {
		case i < last || (i == last && !r.Spec.FType.IsVariadic):
			j = i
			p = params[j]

		case i >= last && r.Spec.FType.IsVariadic:
			j = last
			p = &gotype.Var{
				Name: params[j].Name,
				Type: params[j].Type.Elem.Type,
			}

		case i > last && r.Spec.JoinOp > 0:
			if len(r.Args)-i >= len(params) {
				j = i % len(params)
				p = params[j]
			}
		}
		if p == nil {
			// NOTE this relies on the fact that ArgOpts is
			// constructured from the rule's Args *after* it
			// has been confirmed, with isValidNumberOfArgs,
			// that the number of Args is ok, given the
			// associated function's type.
			panic("shouldn't reach")
		}

		if !c.canConvertRuleArg(p.Type, r.Args[i]) {
			return &Error{r: r, ra: r.Args[i], fp: p, fpi: &j}
		}
	}
	return nil
}

// checkIncludedRuleArgValues is used to validate the literal arguments
// provided to the rules of the github.com/frk/valid package.
func (c *Checker) checkIncludedRuleArgValues(r *Rule) error {
	if len(r.Args) == 0 { // no args to check?
		return nil
	}

	// The majority of the included rules need only their
	// first argument checked, the rest will be handled on
	// a per-rule basis.
	var a0 *Arg
	if len(r.Args) > 0 && (r.Args[0].Type != ARG_FIELD_ABS && r.Args[0].Type != ARG_FIELD_REL) {
		a0 = r.Args[0]
	}

	switch r.Spec.Name {

	// both alpha & alnum expect an ISO-639 language code as argument
	case "alpha", "alnum":
		if a0 != nil && !valid.ISO639(a0.Value, 0) {
			p, pi := r.Spec.getFuncParamByArgIndex(0)
			return &Error{r: r, ra: a0, fp: p, fpi: &pi}
		}

	// ccy expects a valid ISO-4217 currency code as argument
	case "ccy":
		if a0 != nil && !valid.ISO4217(a0.Value) {
			p, pi := r.Spec.getFuncParamByArgIndex(0)
			return &Error{r: r, ra: a0, fp: p, fpi: &pi}
		}

	// decimal expects a cldr-supported locale as argument
	case "decimal":
		if a0 != nil {
			if _, ok := cldr.Locale(a0.Value); !ok {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// hash expects an argument present in the HashAlgoLen table
	case "hash":
		if a0 != nil {
			if _, ok := tables.HashAlgoLen[a0.Value]; !ok {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// ip expects an integer specifying a valid ip version as
	// argument, additionally the value 0 is also accepted which
	// allows the validation to validate against all versions
	case "ip":
		if a0 != nil {
			if a0.Value != "4" && a0.Value != "6" && a0.Value != "0" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// isbn expects an integer specifying a valid isbn version as
	// argument, additionally the value 0 is also accepted which
	// allows the validation to validate against all versions
	case "isbn":
		if a0 != nil {
			if a0.Value != "10" && a0.Value != "13" && a0.Value != "0" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// iso639 expects an integer specifying a valid iso639 version as
	// argument, additionally the value 0 is also accepted which allows
	// the validation to validate against all versions
	case "iso639":
		if a0 != nil {
			if a0.Value != "1" && a0.Value != "2" && a0.Value != "0" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// iso31661a expects an integer specifying a valid iso31661a version as
	// argument, additionally the value 0 is also accepted which allows the
	// validation to validate against all versions
	case "iso31661a":
		if a0 != nil {
			if a0.Value != "2" && a0.Value != "3" && a0.Value != "0" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// mac expects an integer specifying a valid mac version as argument,
	// additionally the value 0 is also accepted which allows the validation
	// to validate against all versions
	case "mac":
		if a0 != nil {
			if a0.Value != "6" && a0.Value != "8" && a0.Value != "0" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// re expects a valid regular expression as argument
	case "re":
		if a0 != nil {
			if _, err := regexp.Compile(a0.Value); err != nil {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi, err: err}
			}
		}

	// uuid expects an integer specifying a supported uuid version
	case "uuid":
		if a0 != nil {
			if a0.Value != "3" && a0.Value != "4" && a0.Value != "5" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a0, fp: p, fpi: &pi}
			}
		}

	// phone, var, and zip all expect a valid ISO-3166-1A country code
	case "phone", "vat", "zip":
		if a0 != nil && !valid.ISO31661A(a0.Value, 0) {
			p, pi := r.Spec.getFuncParamByArgIndex(0)
			return &Error{r: r, ra: a0, fp: p, fpi: &pi}
		}
	}
	return nil
}
