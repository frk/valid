package checker

import (
	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

// walks through the type and does basic preparation for the type checking:
// - parses rule tags
// - loads rule specs
// - orders rules by type/priority
// - loads constants for enum rules
// - generates field keys
// - resolves field references
func (c *checker) walkType(t *types.Type) error {
	if err := c.walkObj(&types.Obj{Type: t}, &walk{}); err != nil {
		return err
	}
	return nil
}

type walk struct {
	is  *rules.TagNode
	pre *rules.TagNode
	ff  types.FieldChain
	ptr []*rules.Rule
}

func (c *checker) walkObj(o *types.Obj, w *walk) error {
	c.ty.ObjFieldMap[o] = w.ff.Last()
	if err := c.objRuleSet(o, w); err != nil {
		return err
	}

	switch t := o.Type; t.Kind {
	case types.PTR:
		c.ty.PtrMap[o.Type.Elem] = o
		w := &walk{is: w.is, pre: w.pre, ff: w.ff, ptr: w.ptr}
		if err := c.walkObj(t.Elem, w); err != nil {
			return err
		}
	case types.ARRAY, types.SLICE:
		w := &walk{is: w.is.GetElem(), pre: w.pre.GetElem(), ff: w.ff}
		if err := c.walkObj(t.Elem, w); err != nil {
			return err
		}
	case types.MAP:
		wk := &walk{is: w.is.GetKey(), pre: w.pre.GetKey(), ff: w.ff}
		if err := c.walkObj(t.Key, wk); err != nil {
			return err
		}
		we := &walk{is: w.is.GetElem(), pre: w.pre.GetElem(), ff: w.ff}
		if err := c.walkObj(t.Elem, we); err != nil {
			return err
		}
	case types.STRUCT:
		for _, f := range t.Fields {
			if f.CanSkip(c.v.Type.Pkg) {
				continue
			}

			ff := w.ff.CopyWith(f)
			c.ty.FKeyMap[f] = c.newFKey(ff)
			c.ty.ChainMap[f] = ff

			is := rules.ParseTag(f.Tag, "is")
			pre := rules.ParseTag(f.Tag, "pre")

			w := &walk{is: is, pre: pre, ff: ff}
			if err := c.walkObj(f.Obj, w); err != nil {
				return c.err(err, errOpts{sf: f})
			}
		}
	}
	return nil
}

// ruleSpec retrieves the spec for the given rule and associates
// it with the rule. The key argument can be either "is" or "pre".
func (c *checker) ruleSpec(r *rules.Rule, w *walk, key string) (*rules.Spec, error) {
	var (
		t *rules.TagNode
		s *rules.Spec
	)

	switch key {
	case "is":
		t = w.is
		s = specs.Get(r.Name)
	case "pre":
		t = w.pre
		s = specs.Get("pre:" + r.Name)
	}
	if s == nil {
		return nil, &Error{C: E_RULE_UNDEFINED, tag: t, r: r}
	}

	r.Spec = s
	return s, nil
}

// checkTagNodes checks that the rule structure of
// the tag is compatible with the object's type.
func (c *checker) checkTagNodes(o *types.Obj, w *walk) error {
	base := o.Type
	for base.Kind == types.PTR {
		base = base.Elem.Type
	}

	// check "is" tag
	if w.is.HasKey() && !base.Is(types.MAP) {
		return &Error{C: E_RULE_KEY, ty: base, tag: w.is}
	}
	if w.is.HasElem() && !base.Is(types.ARRAY, types.SLICE, types.MAP) {
		return &Error{C: E_RULE_ELEM, ty: base, tag: w.is}
	}

	// check "pre" tag
	if w.pre.HasKey() && !base.Is(types.MAP) {
		return &Error{C: E_RULE_KEY, ty: base, tag: w.pre}
	}
	if w.pre.HasElem() && !base.Is(types.ARRAY, types.SLICE, types.MAP) {
		return &Error{C: E_RULE_ELEM, ty: base, tag: w.pre}
	}

	return nil
}

// objRuleSet creates a rules.Set for the given object.
func (c *checker) objRuleSet(o *types.Obj, w *walk) error {
	if err := c.checkTagNodes(o, w); err != nil {
		return err
	}

	switch {
	case o.Type.Kind != types.PTR:
		// Load the specs for validation and preprocessor rules
		// and split the validation rules according to "priority".
		var (
			is0 []*rules.Rule // REQUIRED/OPTIONAL/NOGUARD "is" rules
			is1 []*rules.Rule // rest of the "is" rules
			pre []*rules.Rule // the "pre" rules

			// the "isvalid" rule is implicit for types that implement
			// the method, but can also be specified explicitly, and so
			// it requires a bit of special attention.
			isvalid   *rules.Rule
			rmisvalid bool // remove?
		)

		for _, r := range w.is.GetRules() {
			switch r.Name {
			case "omitkey": // non-rule
				continue
			case "isvalid":
				isvalid = r
				continue
			case "-isvalid":
				rmisvalid = true
				continue
			}

			s, err := c.ruleSpec(r, w, "is")
			if err != nil {
				return err
			}

			// split by type/priority
			switch s.Kind {
			case rules.REQUIRED, rules.OPTIONAL, rules.NOGUARD:
				is0 = append(is0, r)
			default:
				is1 = append(is1, r)
			}

			if s.Kind == rules.ENUM && o.Type.Name != "" && o.Type.Kind.IsBasic() {
				c.ty.EnumMap[o.Type] = types.FindConsts(o.Type, c.ast)
			}
			if err := c.addFRefs(r, w.ff); err != nil {
				return err
			}
		}

		// When the type implements the IsValid() method and "-isvalid"
		// was not used, then make sure the "isvalid" rule is included.
		if isvalid == nil && o.Type.HasIsValid() {
			isvalid = &rules.Rule{Name: "isvalid"}
		}
		if !rmisvalid && isvalid != nil {
			if _, err := c.ruleSpec(isvalid, w, "is"); err != nil {
				return err
			}
			is1 = append(is1, isvalid)
		}

		// When ptr.Rules is not empty then is0 is empty because
		// all the REQUIRED/OPTIONAL/NOGUARD rules were split off by
		// previously processing this object's "parent" object (pointer).
		for _, r := range w.ptr {
			r := *r
			switch t, s := o.Type, r.Spec; true {
			case s.Kind == rules.REQUIRED && s.Name == "required":
				if t.Kind.IsBasic() || t.IsNilable() || (len(t.Name) > 0 && t.IsComparable()) {
					is0 = append(is0, &r)
				}
			case s.Kind == rules.REQUIRED && s.Name == "notnil":
				if o.Type.IsNilable() {
					is0 = append(is0, &r)
				}
			case s.Kind == rules.OPTIONAL && s.Name == "optional":
				is0 = append(is0, &r)
			}
		}

		for _, r := range w.pre.GetRules() {
			if _, err := c.ruleSpec(r, w, "pre"); err != nil {
				return err
			}
			pre = append(pre, r)
			if err := c.addFRefs(r, w.ff); err != nil {
				return err
			}
		}

		o.IsRules = append(is0, is1...)
		o.PreRules = pre

	case o.Type.Kind == types.PTR:
		if len(w.ptr) > 0 {
			for _, r := range w.ptr {
				r := *r
				o.IsRules = append(o.IsRules, &r)
			}
			return nil
		}

		// The REQUIRED/OPTIONAL/NOGUARD rules require special
		// attention and are therefore split off from the rest of
		// the rules.
		//
		// Note that currently there's no REQUIRED that applies
		// to "pointers only"... it may be useful to have that.
		//
		// Note also that the NOGUARD rule can be used together
		// with optional, required, notnil-if-base-is-nilable.
		//
		// The NOGUARD rule makes no sense with optional
		// or notnil-if-base-is-NOT-nilable.
		var (
			required *rules.Rule // should apply to pointers and base (with easy to generate zero value)
			notnil   *rules.Rule // should apply to pointers and nilable base
			optional *rules.Rule // should apply to pointers and base
			omitnil  *rules.Rule // should apply to pointers only
			noguard  *rules.Rule // should apply to pointers only

			// rest of the rules should apply to base only
			rr []*rules.Rule
		)
		for _, r := range w.is.GetRules() {
			if r.Name == "omitkey" { // non-rule
				continue
			}

			s, err := c.ruleSpec(r, w, "is")
			if err != nil {
				return err
			}

			switch {
			case s.Kind == rules.REQUIRED && s.Name == "required":
				required = r
			case s.Kind == rules.REQUIRED && s.Name == "notnil":
				notnil = r
			case s.Kind == rules.OPTIONAL && s.Name == "optional":
				optional = r
			case s.Kind == rules.OPTIONAL && s.Name == "omitnil":
				omitnil = r
			case s.Kind == rules.NOGUARD:
				noguard = r
			default:
				rr = append(rr, r)
			}
		}

		// For pointer types the "omitnil" is applied by default
		// if no other pointer-specific rule was explicitly provided.
		if required == nil && notnil == nil && optional == nil && omitnil == nil && noguard == nil {
			r := &rules.Rule{Name: "omitnil"}
			if _, err := c.ruleSpec(r, w, "is"); err != nil {
				return err
			}
			omitnil = r
		}

		if required != nil {
			w.ptr = append(w.ptr, required)
		}
		if notnil != nil {
			w.ptr = append(w.ptr, notnil)
		}
		if optional != nil {
			w.ptr = append(w.ptr, optional)
		}
		if omitnil != nil {
			w.ptr = append(w.ptr, omitnil)
		}
		if noguard != nil {
			w.ptr = append(w.ptr, noguard)
		}

		// Update the "is" rules.TagNode's rules to
		// hold only the non-pointer specific rules.
		if w.is != nil {
			w.is.Rules = rr
		}

		o.IsRules = w.ptr
	}
	return nil
}

func (c *checker) walkVisible(o *types.Obj, ff types.FieldChain) {
	switch t := o.Type; t.Kind {
	case types.PTR:
		c.walkVisible(t.Elem, ff)
	case types.ARRAY, types.SLICE:
		c.walkVisible(t.Elem, ff)
	case types.MAP:
		c.walkVisible(t.Key, ff)
		c.walkVisible(t.Elem, ff)
	case types.STRUCT:
		for _, f := range t.VisibleFields() {
			if f.CanSkip(c.v.Type.Pkg) {
				continue
			}

			ff := ff.CopyWith(f)
			if _, ok := c.ty.VisibleChainMap[f]; !ok {
				c.ty.VisibleChainMap[f] = ff
			}

			c.walkVisible(f.Obj, ff)
		}
	}
}
