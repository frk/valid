package rules

import (
	"reflect"
	"strings"

	"github.com/frk/valid"
)

// A Tag is a tree-like representation of a parsed "rule" struct tag.
type Tag struct {
	// The list of rules present at this level of the tree.
	Rules []*Rule
	// Key and Elem are child nodes of this Tag.
	Key, Elem *Tag
	// The struct tag key used for extracting the struct
	// tag value, intended for error reporting only.
	stkey string `cmp:"-"`
}

func (t *Tag) String() (out string) {
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

// HasElem reports whether or not the Tag has an Key child.
func (t *Tag) HasKey() bool {
	return t != nil && t.Key != nil
}

// HasElem reports whether or not the Tag has an Elem child.
func (t *Tag) HasElem() bool {
	return t != nil && t.Elem != nil
}

func (t *Tag) GetKey() *Tag {
	if t != nil {
		return t.Key
	}
	return nil
}

func (t *Tag) GetElem() *Tag {
	if t != nil {
		return t.Elem
	}
	return nil
}

func (t *Tag) GetRules() []*Rule {
	if t != nil {
		return t.Rules
	}
	return nil
}

// HasRule reports whether or not the Tag
// contains a rule with the given name.
func (t *Tag) HasRule(name string) bool {
	if t != nil {
		for _, r := range t.Rules {
			if r.Name == name {
				return true
			}
		}
	}
	return false
}

// GetRule looks up a rule in the Tag by
// the given name and, if present, returns it.
func (t *Tag) GetRule(name string) (r *Rule, ok bool) {
	if t != nil {
		for _, r := range t.Rules {
			if r.Name == name {
				return r, true
			}
		}
	}
	return nil, false
}

// AddRule adds the given rule to the Tag.
func (t *Tag) AddRule(r *Rule) {
	if t != nil {
		for i := range t.Rules {
			if t.Rules[i].Name == r.Name {
				return
			}
		}
		t.Rules = append(t.Rules, r)
	}
}

// RmRule removes a rule from the Tag by name.
func (t *Tag) RmRule(name string) {
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

// IsEmpty reports whether or not the Tag,
// including its sub-nodes, are empty.
func (t *Tag) IsEmpty() bool {
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

// parseTag parses the given tag and returns the resulting AST.
// The string "is" is used as the key in the struct tag, e.g. `is:"rule"`.
//
// If the "is" key is not present in the tag, or the associated value
// is empty or equal to "-", then the returned AST will be nil.
//
// Following is a description of the rule syntax using EBNF:
//
//      node      = rule | [ "[" [ node ] "]" ] [ ( node | rule "," node ) ] .
//      rule      = rule_name [ { ":" rule_arg } ] { "," rule } .
//      rule_name = identifier .
//      rule_arg  = | boolean_lit | integer_lit | float_lit | string_lit | quoted_string_lit | field_reference .
//
//      boolean_lit       = "true" | "false" .
//      integer_lit       = "0" | [ "-" ] "1"…"9" { "0"…"9" } .
//      float_lit         = [ "-" ] ( "0" | "1"…"9" { "0"…"9" } ) "." "0"…"9" { "0"…"9" } .
//      string_lit        = .
//      quoted_string_lit = `"` `"` .
//
//      field_reference = field_ref_abs | field_ref_rel
//      field_ref_abs   = "&" field_key .
//      field_ref_rel   = "." field_key .
//      field_key       = identifier { field_key_separator identifier } .
//      field_key_sep   = "." | (* optionally specified by the user *)
//
//      identifier = letter { letter } .
//      letter     = "A"…"Z" | "a"…"z" | "_" .
//
func parseTag(tag, key string) *Tag {
	str, ok := reflect.StructTag(tag).Lookup(key)
	if !ok || str == "-" || len(str) == 0 {
		return nil
	}
	return parseRule(str, key)
}

// `is:".field.bar"`

// parseRule parses the given rule string and returns the AST.
func parseRule(str, stkey string) *Tag {
	tag := &Tag{stkey: stkey}
	for str != "" {
		// skip leading space
		i := 0
		for i < len(str) && str[i] == ' ' {
			i++
		}
		str = str[i:]
		if str == "" {
			break
		}

		// parse bracketed nodes
		if str[0] == '[' {

			// scan up to the *matching* closing bracket
			i, n := 1, 0
			for i < len(str) && (str[i] != ']' || n > 0) {
				// adjust nesting level
				if str[i] == '[' {
					n++
				} else if str[i] == ']' {
					n--
				}
				i++

				// scan quoted string, ignoring brackets inside quotes
				if str[i-1] == '"' {
					for i < len(str) && str[i] != '"' {
						if str[i] == '\\' {
							i++
						}
						i++
					}

					// keep the closing double quote, or
					// else the subsequent parser calls
					// will be confused without it
					if i < len(str) {
						i++
					}
				}
			}

			// recursively invoke parser for key
			if ktag := str[1:i]; len(ktag) > 0 {
				tag.Key = parseRule(ktag, stkey)
			}
			// recursively invoke parser for elem
			if etag := str[i:]; len(etag) > 1 {
				etag = etag[1:] // drop the leading ']'
				tag.Elem = parseRule(etag, stkey)
			}

			// done; exit
			return tag
		}

		// scan to the end of a rule's name
		i = 0
		for i < len(str) && str[i] != ',' && str[i] != ':' {
			i++
		}

		// empty name's no good; next
		if str[:i] == "" {
			str = str[1:]
			continue
		}

		rule := &Rule{Name: str[:i]}
		tag.AddRule(rule)

		// this rule's done; next or exit
		if str = str[i:]; str == "" {
			break
		} else if str[0] == ',' {
			str = str[1:]
			continue
		}

		// scan the rule's arguments
		for str != "" {
			str = str[1:] // drop the leading ':'

			// quoted argument value; scan to the end quote
			if len(str) > 0 && str[0] == '"' {
				i := 1
				for i < len(str) && str[i] != '"' {
					if str[i] == '\\' {
						i++
					}
					i++
				}

				a := &Arg{}
				a.Type = ARG_STRING
				a.Value = str[1:i]
				rule.Args = append(rule.Args, a)

				str = str[i:]

				// drop the closing quote
				if len(str) > 0 && str[0] == '"' {
					str = str[1:]
				}

				// next argument?
				if len(str) > 0 && str[0] == ':' {
					continue
				}

				// drop rule separator
				if len(str) > 0 && str[0] == ',' {
					str = str[1:]
				}

				// this rule's done; exit
				break
			}

			// scan to the end of a rule's argument
			i := 0
			for i < len(str) && str[i] != ':' && str[i] != ',' {
				i++
			}

			astr := str[:i]
			arg := parseArg(astr)
			rule.Args = append(rule.Args, arg)

			str = str[i:]
			if str == "" {
				break
			} else if str[0] == ',' {
				str = str[1:]
				break
			}
		}
	}
	return tag
}

// parseArg parses the given string as an Arg and returns the result.
func parseArg(str string) (a *Arg) {
	a = &Arg{}
	if len(str) > 0 {
		if str[0] == '&' {
			a.Type = ARG_FIELD_ABS
			a.Value = str[1:]
		} else if str[0] == '.' {
			a.Type = ARG_FIELD_REL
			a.Value = str[1:]
		} else {
			switch {
			case valid.Int(str):
				a.Type = ARG_INT
			case valid.Float(str):
				a.Type = ARG_FLOAT
			case str == "true" || str == "false":
				a.Type = ARG_BOOL
			case str != `nil`:
				a.Type = ARG_STRING
			}
			a.Value = str
		}
	}
	return a
}
