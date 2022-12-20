package rules

import (
	"strings"
)

// A TagNode is a tree-like representation of a parsed "rule" struct tag.
type TagNode struct {
	// The list of rules present at this level of the tree.
	Rules []*Rule
	// Key and Elem are child nodes of this TagNode.
	Key, Elem *TagNode

	// The struct tag key used for extracting the struct
	// tag value, intended for error reporting only.
	stkey string `cmp:"-"`
}

func (t *TagNode) String() (out string) {
	if t == nil {
		return out
	}

	var ss []string
	for i := range t.Rules {
		ss = append(ss, t.Rules[i].String())
	}

	var key, elem string
	if t.Key != nil {
		key = t.Key.String()
	}
	if t.Elem != nil {
		elem = t.Elem.String()
	}
	if len(key) > 0 || len(elem) > 0 {
		ss = append(ss, "["+key+"]"+elem)
	}

	return strings.Join(ss, ",")
}

func (t *TagNode) STKey() string {
	if t != nil {
		return t.stkey
	}
	return ""
}

// HasElem reports whether or not the TagNode has an Key child.
func (t *TagNode) HasKey() bool {
	return t != nil && t.Key != nil
}

// HasElem reports whether or not the TagNode has an Elem child.
func (t *TagNode) HasElem() bool {
	return t != nil && t.Elem != nil
}

func (t *TagNode) GetKey() *TagNode {
	if t != nil {
		return t.Key
	}
	return nil
}

func (t *TagNode) GetElem() *TagNode {
	if t != nil {
		return t.Elem
	}
	return nil
}

func (t *TagNode) GetRules() []*Rule {
	if t != nil {
		return t.Rules
	}
	return nil
}

// HasRule reports whether or not the TagNode
// contains a rule with the given name.
func (t *TagNode) HasRule(name string) bool {
	if t != nil {
		for _, r := range t.Rules {
			if r.Name == name {
				return true
			}
		}
	}
	return false
}

// GetRule looks up a rule in the TagNode by
// the given name and, if present, returns it.
func (t *TagNode) GetRule(name string) (r *Rule, ok bool) {
	if t != nil {
		for _, r := range t.Rules {
			if r.Name == name {
				return r, true
			}
		}
	}
	return nil, false
}

// AddRule adds the given rule to the TagNode.
func (t *TagNode) AddRule(r *Rule) {
	if t != nil {
		for i := range t.Rules {
			if t.Rules[i].Name == r.Name {
				return
			}
		}
		t.Rules = append(t.Rules, r)
	}
}

// RmRule removes a rule from the TagNode by name.
func (t *TagNode) RmRule(name string) {
	if t != nil {
		for i, r := range t.Rules {
			if r.Name == name {
				t.Rules = append(
					append([]*Rule{}, t.Rules[:i]...),
					t.Rules[i+1:]...,
				)
				break
			}
		}
	}
}

// IsEmpty reports whether or not the TagNode,
// including its sub-nodes, are empty.
func (t *TagNode) IsEmpty() bool {
	if t != nil {
		if len(t.Rules) > 0 {
			return false
		}
		if !t.Key.IsEmpty() {
			return false
		}
		if !t.Elem.IsEmpty() {
			return false
		}
	}
	return true
}
