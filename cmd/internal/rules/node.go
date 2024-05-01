package rules

import (
	"fmt"
	"strings"

	"github.com/frk/valid/cmd/internal/gotype"
)

type Node struct {
	// The type associated with the node.
	Type *gotype.Type
	// List of preprocessor rules associated with the node.
	PreRules RuleList
	// List of validation rules associated with the node.
	IsRules RuleList
	// If the Node represents a map, then Key
	// will hold information about the map's key.
	Key *Node
	// If the Node represents a map, array, or slice, then
	// Elem will hold information about the element's type.
	//
	// If the Node represents a pointer, then Elem will
	// hold information about the pointer's base type.
	Elem *Node
	// If the Node represents a struct, then Fields
	// will hold information about the struct's fields.
	Fields []*FieldNode
	// Ptr holds the Node whose type is the pointer to
	// this Node's type. If this Node's type is not pointed
	// to by any pointer then this field will be nil.
	Ptr *Node
}

type FieldNode struct {
	Field *gotype.StructField
	// The Node representation of the field's type.
	Type *Node
	// The unique key of the field.
	Key string
	// The field's selector.
	Selector gotype.FieldSelector
}

// makeNode creates a Node representation of the given
// Go type and its associated "is" & "pre" rule tags.
func (c *Checker) makeNode(t *gotype.Type, is, pre *Tag, fs gotype.FieldSelector) (_ *Node, err error) {
	var root, base *Node
	if t.Kind == gotype.K_PTR {
		root, base, err = c.makeNodeFromPtr(t, is, pre, fs)
		if err != nil {
			return nil, c.err(err, errOpts{sf: fs.Last()})
		}
	} else {
		isList, preList, err := c.makeRuleLists(is, pre, fs)
		if err != nil {
			return nil, c.err(err, errOpts{sf: fs.Last(), ty: t})
		}

		root = &Node{Type: t}
		root.PreRules = preList
		root.IsRules = isList
		base = root
	}

	// use the base type for further checks
	t = base.Type

	if false {
		fmt.Printf("type=%s is=%s pre=%s\n", t.TypeString(nil), is.String(), pre.String())
	}

	// Make sure that the rule structure in both tags
	// is compatible with the associated type.
	if is.HasKey() && !t.Is(gotype.K_MAP) {
		return nil, &Error{C: ERR_RULE_KEY, sf: fs.Last(), ty: t, tag: is}
	}
	if is.HasElem() && !t.Is(gotype.K_ARRAY, gotype.K_SLICE, gotype.K_MAP) {
		return nil, &Error{C: ERR_RULE_ELEM, sf: fs.Last(), ty: t, tag: is}
	}
	if pre.HasKey() && !t.Is(gotype.K_MAP) {
		return nil, &Error{C: ERR_RULE_KEY, sf: fs.Last(), ty: t, tag: pre}
	}
	if pre.HasElem() && !t.Is(gotype.K_ARRAY, gotype.K_SLICE, gotype.K_MAP) {
		return nil, &Error{C: ERR_RULE_ELEM, sf: fs.Last(), ty: t, tag: pre}
	}

	// descend the type hierarchy
	switch t.Kind {
	case gotype.K_ARRAY:
		if base.Elem, err = c.makeNode(t.Elem, is.GetElem(), pre.GetElem(), fs); err != nil {
			return nil, err
		}
	case gotype.K_SLICE:
		if base.Elem, err = c.makeNode(t.Elem, is.GetElem(), pre.GetElem(), fs); err != nil {
			return nil, err
		}
	case gotype.K_MAP:
		if base.Key, err = c.makeNode(t.Key, is.GetKey(), pre.GetKey(), fs); err != nil {
			return nil, err
		}
		if base.Elem, err = c.makeNode(t.Elem, is.GetElem(), pre.GetElem(), fs); err != nil {
			return nil, err
		}
	case gotype.K_STRUCT:
		for _, f := range t.Fields {
			if !f.CanAccess(c.pkg) {
				continue
			}

			node, err := c.makeFieldNode(f, fs)
			if err != nil {
				return nil, err
			}
			base.Fields = append(base.Fields, node)
		}
	}

	return root, nil
}

func (c *Checker) makeNodeFromPtr(t *gotype.Type, is, pre *Tag, fs gotype.FieldSelector) (root, base *Node, err error) {
	root = &Node{Type: t}
	base = root

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
		required *Rule // should apply to pointers and base (with easy to generate zero value)
		notnil   *Rule // should apply to pointers and nilable base
		optional *Rule // should apply to pointers and base
		omitnil  *Rule // should apply to pointers only
		noguard  *Rule // should apply to pointers only

		// rest of the rules should apply to base only
		rr []*Rule
	)

	for _, r := range is.GetRules() {
		if r.Name == "omitkey" { // non-rule
			continue
		}

		if r.Spec = GetSpec(r.Name); r.Spec == nil {
			return nil, nil, &Error{C: ERR_RULE_UNDEFINED, ty: t, tag: is, r: r}
		}

		switch {
		case r.Spec.Kind == REQUIRED && r.Spec.Name == "required":
			required = r
		case r.Spec.Kind == REQUIRED && r.Spec.Name == "notnil":
			notnil = r
		case r.Spec.Kind == OPTIONAL && r.Spec.Name == "optional":
			optional = r
		case r.Spec.Kind == OPTIONAL && r.Spec.Name == "omitnil":
			omitnil = r
		case r.Spec.Kind == NOGUARD:
			noguard = r
		default:
			rr = append(rr, r)
		}

		for _, a := range r.Args {
			if a.Type == ARG_FIELD_REL {
				c.normalizeRelFieldValue(a, fs)
			}
		}
	}

	// for pointer types the "omitnil" is applied by default
	// if no other pointer-specific rule was explicitly provided
	if required == nil && notnil == nil && optional == nil && noguard == nil {
		r := &Rule{Name: "omitnil"}
		if r.Spec = GetSpec(r.Name); r.Spec == nil {
			panic("shoudn't happen")
		}
		omitnil = r
	}

	// apply the pointer specific rules to
	// every pointer in the pointer-chain
	for t.Kind == gotype.K_PTR {
		if required != nil {
			base.IsRules = append(base.IsRules, required)
		}
		if notnil != nil {
			base.IsRules = append(base.IsRules, notnil)
		}
		if optional != nil {
			base.IsRules = append(base.IsRules, optional)
		}
		if omitnil != nil {
			base.IsRules = append(base.IsRules, omitnil)
		}
		if noguard != nil {
			base.IsRules = append(base.IsRules, noguard)
		}

		// resolve base
		ptr := *base
		base.Elem = &Node{Type: t.Elem, Ptr: &ptr}
		base = base.Elem
		t = t.Elem
	}

	// apply base specific rules
	if required != nil && (t.Kind.IsBasic() || t.IsNilable() || (len(t.Name) > 0 && t.IsComparable())) {
		base.IsRules = append(base.IsRules, required)
	}
	if notnil != nil && t.IsNilable() {
		base.IsRules = append(base.IsRules, notnil)
	}
	if optional != nil {
		base.IsRules = append(base.IsRules, optional)
	}
	if len(rr) > 0 {
		base.IsRules = append(base.IsRules, rr...)
	}

	// Load the specs for preprocessor rules. Note that the
	// preprocessor rules currently apply to base only.
	for _, r := range pre.GetRules() {
		if r.Spec = GetSpec("pre:" + r.Name); r.Spec == nil {
			return nil, nil, &Error{C: ERR_RULE_UNDEFINED, ty: t, tag: pre, r: r}
		}
		base.PreRules = append(base.PreRules, r)

		for _, a := range r.Args {
			if a.Type == ARG_FIELD_REL {
				c.normalizeRelFieldValue(a, fs)
			}
		}
	}

	return root, base, nil
}

// makeRuleLists creates RuleLists from the given Tags.
func (c *Checker) makeRuleLists(is, pre *Tag, fs gotype.FieldSelector) (isList, preList RuleList, err error) {
	// Load the specs for validation the rules
	// and sort the rules according to "priority".
	var (
		r0 []*Rule // REQUIRED/OPTIONAL/NOGUAR take precedence
		rr []*Rule // rest of the rules
	)
	for _, r := range is.GetRules() {
		if r.Name == "omitkey" { // non-rule
			continue
		}

		if r.Spec = GetSpec(r.Name); r.Spec == nil {
			return nil, nil, &Error{C: ERR_RULE_UNDEFINED, tag: is, r: r}
		}
		switch r.Spec.Kind {
		case REQUIRED, OPTIONAL, NOGUARD:
			r0 = append(r0, r)
		default:
			rr = append(rr, r)
		}

		for _, a := range r.Args {
			if a.Type == ARG_FIELD_REL {
				c.normalizeRelFieldValue(a, fs)
			}
		}
	}
	isList = append(r0, rr...)

	// Load the specs for preprocessor rules.
	for _, r := range pre.GetRules() {
		if r.Spec = GetSpec("pre:" + r.Name); r.Spec == nil {
			return nil, nil, &Error{C: ERR_RULE_UNDEFINED, tag: pre, r: r}
		}
		preList = append(preList, r)

		for _, a := range r.Args {
			if a.Type == ARG_FIELD_REL {
				c.normalizeRelFieldValue(a, fs)
			}
		}
	}

	return isList, preList, nil
}

// makeFieldNode creates a FieldNode for the given struct field.
func (c *Checker) makeFieldNode(f *gotype.StructField, fs gotype.FieldSelector) (n *FieldNode, err error) {
	n = &FieldNode{Field: f}
	n.Selector = fs.CopyWith(f)
	n.Key = c.fieldKey(n.Selector, false)

	is, pre := parseTag(f.Tag, "is"), parseTag(f.Tag, "pre")
	if n.Type, err = c.makeNode(f.Type, is, pre, n.Selector); err != nil {
		return nil, err
	}

	c.Info.KeyMap[n.Key] = n
	return n, nil
}

////////////////////////////////////////////////////////////////////////////////
// the following are convenience methods intended primarily for the generator
//

// IsPtr reports whether or not n's type is a pointer type.
func (n *Node) IsPtr() bool {
	return n.Type.Kind == gotype.K_PTR
}

// IsStruct reports whether or not n's type is a struct type.
func (n *Node) IsStruct() bool {
	return n.Type.Kind == gotype.K_STRUCT
}

// Base returns the pointer base Node of n. If n is not
// a pointer then the returned Node will be n itself.
func (n *Node) Base() (m *Node) {
	for m = n; m.IsPtr(); {
		m = m.Elem
	}
	return m
}

// Root returns the pointer root Node of n. If n is not
// a pointer base then the returned Node will be n itself.
func (n *Node) Root() (m *Node) {
	for m = n; m.Ptr != nil; {
		m = m.Ptr
	}
	return m
}

// PtrDepth reports the "pointer depth" of Node n,
// i.e. the number of pointer Nodes from root to n.
func (n *Node) PtrDepth() int {
	if n.Ptr == nil {
		return 0
	}
	return n.Ptr.PtrDepth() + 1
}

// IsShallow reports whether the n's pointer depth is "shallow" or not.
// A "shallow" pointer depth is 0 or 1, i.e. a depth that allows the base
// type's fields and/or methods to be accessed without explicit indirection.
func (n *Node) IsShallow() bool {
	return n.PtrDepth() < 2
}

// IsDeep reports whether the n's pointer depth is "deep" or not.
// A "deep" pointer depth is 2 and above, i.e. a depth that requires
// explicit indirection to access the base type's fields and/or methods.
func (n *Node) IsDeep() bool {
	return n.PtrDepth() > 1
}

// NumIsRules returns the length of the IsRules slice.
func (n *Node) NumIsRules() int {
	return len(n.IsRules)
}

// HasRules reports whether or not n, or any of its child
// Nodes, have rules other than OPTIONAL and "noguard".
func (n *Node) HasRules() bool {
	if n == nil {
		return false
	}

	if len(n.PreRules) > 0 {
		return true
	}
	for _, r := range n.IsRules {
		if r.Spec.Kind != OPTIONAL && r.Spec.Kind != NOGUARD {
			return true
		}
	}

	if n.Key.HasRules() {
		return true
	}
	if n.Elem.HasRules() {
		return true
	}
	for i := range n.Fields {
		if n.Fields[i].Type.HasRules() {
			return true
		}
	}

	return false
}

// HasSubRules reports whether or not any of n's child
// Nodes have rules other than OPTIONAL and "noguard".
func (n *Node) HasSubRules() bool {
	if n == nil {
		return false
	}

	if n.Key.HasRules() {
		return true
	}
	if n.Elem.HasRules() {
		return true
	}
	for i := range n.Fields {
		if n.Fields[i].Type.HasRules() {
			return true
		}
	}

	return false
}

// IsOptional reports whether or not n has an OPTIONAL is-rule.
func (n *Node) IsOptional() bool {
	return len(n.IsRules) > 0 && n.IsRules[0].Spec.Kind == OPTIONAL
}

// IsRequired reports whether or not n has a REQUIRED is-rule.
func (n *Node) IsRequired() bool {
	return len(n.IsRules) > 0 && n.IsRules[0].Spec.Kind == REQUIRED
}

// IsNoGuard reports whether or not n has a NOGUARD is-rule.
func (n *Node) IsNoGuard() bool {
	return len(n.IsRules) > 0 && n.IsRules[0].Spec.Kind == NOGUARD
}

func (n *Node) NeedsTempVar() bool {
	return n.IsDeep() && (n.IsRules.Many() || n.IsStruct())
}

func (c *Checker) normalizeRelFieldValue(a *Arg, fs gotype.FieldSelector) {
	sep := "."
	if len(fs) > 1 {
		a.Value = c.fieldKey(fs[:len(fs)-1], true) + sep + a.Value
	}
	a.Value = strings.TrimPrefix(a.Value, sep)
	a.Value = strings.TrimSuffix(a.Value, sep)
}
