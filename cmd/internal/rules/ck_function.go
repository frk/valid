package rules

import (
	"regexp"

	"github.com/frk/valid"
	"github.com/frk/valid/cmd/internal/errors"
	"github.com/frk/valid/cmd/internal/gotype"
)

// functionCheck ...
func (c *Checker) functionCheck(n *Node, r *Rule) error {
	// Check if an instance of n can be passed
	// to the function as its first argument.
	p := r.Spec.FType.In[0].Type
	if r.Spec.FType.IsVariadic && len(r.Spec.FType.In) == 1 {
		p = p.Elem
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
				Type: params[j].Type.Elem,
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

////////////////////////////////////////////////////////////////////////////////
//
// Check Argument Values of Included Rules:
//
// The methods here are used to normalize and validate the arguments provided
// to the rules of the github.com/frk/valid package.
//
// Some of the rules in that package take arguments of spe
//
// The arguments' types are checked later by a separate method.
// Because of that these checks SHOULD only validate the *value* of the argument
// if the type of the argument is as expected. If the type of the argument is not
// as expected then skip the value validation.
//
////////////////////////////////////////////////////////////////////////////////

func (c *Checker) checkIncludedRuleArgValues(r *Rule) error {
	switch r.Spec.Name {
	case "alpha":
		if err := c.checkArgLangTag(r.Args[0]); err != nil {
			return err
		}
	case "alnum":
		if err := c.checkArgLangTag(r.Args[0]); err != nil {
			return err
		}
	case "ccy":
		// TODO ccy (currency code)
	case "decimal":
		// TODO decimal (locale)
	case "hash":
		// TODO hash (algo)
	case "ip":
		if err := c.checkArgIPVersion(r.Args[0]); err != nil {
			return err
		}
	case "isbn":
		// TODO isbn (ver)
	case "iso639":
		// TODO iso639 (num)
	case "iso31661a":
		// TODO iso31661a (num)
	case "mac":
		if err := c.checkArgMACVersion(r); err != nil {
			return err
		}
	case "phone":
		if err := c.checkArgCountryCode(r.Args[0]); err != nil {
			return err
		}
	case "re":
		if err := c.checkArgRegexp(r); err != nil {
			return err
		}
	case "uuid":
		if a := r.Args[0]; a.Type != ARG_FIELD {
			if a.Value != "3" && a.Value != "4" && a.Value != "5" {
				p, pi := r.Spec.getFuncParamByArgIndex(0)
				return &Error{r: r, ra: a, fp: p, fpi: &pi}
			}
		}
	case "vat":
		// TODO vat (country code)
	case "zip":
		if err := c.checkArgCountryCode(r.Args[0]); err != nil {
			return err
		}
	}
	return nil
}

// check that the rule's arguments are one of the supported language tags.
func (c *Checker) checkArgLangTag(a *Arg) error {
	if a.Type == ARG_STRING && !valid.ISO639(a.Value, 0) {
		return errors.TODO("checkArgLangTag: language tag argument is not supported: %v", a.Value)
	}
	return nil
}

// check that the rule's arguments are valid country codes.
func (c *Checker) checkArgCountryCode(a *Arg) error {
	if a.Type == ARG_STRING && !valid.ISO31661A(a.Value, 0) {
		return errors.TODO("checkArgCountryCode: arg is not valid country code")
	}
	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// check that the rule's arguments are valid UUID versions.
// NOTE this check also modifies the *Arg to normalize its value.
func (c *Checker) checkArgUUIDVersion(a *Arg) error {
	if a.Type == ARG_FIELD {
		return nil
	}

	if a.Type == ARG_STRING || a.IsUInt() {
		if !rxUUIDVer.MatchString(a.Value) {
			return errors.TODO("checkArgUUIDVersion: arg is not valid UUID version")
		}

		if len(a.Value) > 1 && (a.Value[0] == 'v' || a.Value[0] == 'V') {
			a.Value = a.Value[1:]
			a.Type = ARG_INT
		}
	}
	return nil
}

// checks that the rule's arguments are valid IP versions.
func (c *Checker) checkArgIPVersion(a *Arg) error {
	if a.Type == ARG_FIELD {
		return nil
	}

	if !a.IsUInt() || (a.Value != "0" && a.Value != "4" && a.Value != "6") {
		return errors.TODO("checkArgIPVersion: arg is not valid IP version")
	}
	return nil
}

// checks that the rule's arguments are valid MAC versions.
func (c *Checker) checkArgMACVersion(r *Rule) error {
	for _, arg := range r.Args {
		if arg.IsUInt() {
			if arg.Value != "0" && arg.Value != "6" && arg.Value != "8" {
				return errors.TODO("checkArgMACVersion: arg is not valid MAC version")
			}
		}
	}
	return nil
}

// check that the rule's arguments are strings containing compilable regular expressions.
func (c *Checker) checkArgRegexp(r *Rule) error {
	for _, arg := range r.Args {
		if arg.Type != ARG_FIELD {
			if _, err := regexp.Compile(arg.Value); err != nil {
				return errors.TODO("checkArgRegexp: arg is not valid regular expression")
			}
		}
	}
	return nil
}
