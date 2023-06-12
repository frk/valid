// Following is a description of the rule syntax using EBNF:
//
//	node      = rule | [ "[" [ node ] "]" ] [ ( node | rule "," node ) ] .
//	rule      = rule_name [ { ":" rule_arg } ] { "," rule } .
//	rule_name = identifier .
//	rule_arg  = | boolean_lit | integer_lit | float_lit | string_lit | quoted_string_lit | field_ref .
//
//	boolean_lit       = "true" | "false" .
//	integer_lit       = "0" | [ "-" ] "1"…"9" { "0"…"9" } .
//	float_lit         = [ "-" ] ( "0" | "1"…"9" { "0"…"9" } ) "." "0"…"9" { "0"…"9" } .
//	string_lit        = .
//	quoted_string_lit = `"` `"` .
//
//	field_ref = field_ref_abs | field_ref_rel
//	field_ref_abs   = "&" field_key .
//	field_ref_rel   = "." field_key .
//	field_key       = identifier { field_key_sep identifier } .
//	field_key_sep   = "."
//
//	identifier = letter [ { letter | digit } ] .
//	letter     = "A"…"Z" | "a"…"z" | "_" .
//	digit      = "0"…"9"
package rules

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/frk/valid"
)

type UndefinedRuleError struct {
	RuleName string
	TagKey   string
	TagValue string
}

func (e *UndefinedRuleError) Error() string {
	return fmt.Sprintf("undefined rule: %q", e.RuleName)
}

// ParseTag parses the given struct tag and returns the resulting *Obj.
//
// The key argument is used to lookup the specific struct tag value that
// should be parsed. If the no value for the key is present in the tag, or
// the value is empty or equal to "-", then the returned *Obj will be nil.
func ParseTag(tag string, key string, rr Registry) (*Obj, error) {
	str, ok := reflect.StructTag(tag).Lookup(key)
	if !ok || str == "-" || len(str) == 0 {
		return nil, nil
	}
	obj, err := Parse(str, rr)
	if err != nil {
		if e, ok := err.(*UndefinedRuleError); ok && e != nil {
			e.TagKey = key
			e.TagValue = tag
		}
		return nil, err
	}
	return obj, nil
}

// Obj is the result of parsing a rule string. It contains a set
// of rules that should be applied to an instance of types.Type.
type Obj struct {
	// The list of rules that should be applied to the target type proper.
	Rules []*Rule
	// Key contains rules that should be applied to the Key of the target type.
	Key *Obj
	// Key contains rules that should be applied to the Elem of the target type.
	Elem *Obj
}

// Parse parses the given string and returns the resulting *Obj.
func Parse(str string, rr Registry) (*Obj, error) {
	o := new(Obj)
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
			if key := str[1:i]; len(key) > 0 {
				obj, err := Parse(key, rr)
				if err != nil {
					return nil, err
				}
				o.Key = obj
			}
			// recursively invoke parser for elem
			if elem := str[i:]; len(elem) > 1 {
				elem = elem[1:] // drop the leading ']'
				obj, err := Parse(elem, rr)
				if err != nil {
					return nil, err
				}
				o.Elem = obj
			}

			// done; exit
			return o, nil
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

		name := str[:i]
		rule, ok := rr.Lookup(name)
		if !ok {
			return nil, &UndefinedRuleError{RuleName: name}
		}
		o.AddRule(&rule)

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

				a := new(Arg)
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

			arg := parseArg(str[:i])
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
	return o, nil
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

// String implements the fmt.Stringer interface.
func (o *Obj) String() (out string) {
	if o == nil {
		return out
	}

	var ss []string
	for i := range o.Rules {
		ss = append(ss, o.Rules[i].String())
	}

	var key, elem string
	if o.Key != nil {
		key = o.Key.String()
	}
	if o.Elem != nil {
		elem = o.Elem.String()
	}
	if len(key) > 0 || len(elem) > 0 {
		ss = append(ss, "["+key+"]"+elem)
	}

	return strings.Join(ss, ",")
}

// AddRule adds the given rule to the Obj iff it isn't already present.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) AddRule(r *Rule) {
	if o == nil {
		return
	}

	for i := range o.Rules {
		if o.Rules[i].Name == r.Name {
			return
		}
	}
	o.Rules = append(o.Rules, r)
}

// RemoveRule removes a rule from the Obj by name.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) RemoveRule(name string) {
	if o == nil {
		return
	}

	for i, r := range o.Rules {
		if r.Name == name {
			head, tail := o.Rules[:i], o.Rules[i+1:]
			o.Rules = append(append([]*Rule{}, head...), tail...)
			break
		}
	}
}

// IsEmpty reports whether or not the Obj, including its Key and Elem sub-objects,
// are empty. The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) IsEmpty() bool {
	if o != nil {
		if len(o.Rules) > 0 {
			return false
		}
		if !o.Key.IsEmpty() {
			return false
		}
		if !o.Elem.IsEmpty() {
			return false
		}
	}
	return true
}

// HasKey reports whether or not the Obj has a non-nil Key.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) HasKey() bool {
	return o != nil && o.Key != nil
}

// HasElem reports whether or not the Obj has a non-nil Elem.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) HasElem() bool {
	return o != nil && o.Elem != nil
}

// HasRule reports whether or not the Obj contains a rule with the given name.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) HasRule(name string) bool {
	if o != nil {
		for _, r := range o.Rules {
			if r.Name == name {
				return true
			}
		}
	}
	return false
}

// GetRules is a convenience method that returns the Obj's Rules.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) GetRules() []*Rule {
	if o != nil {
		return o.Rules
	}
	return nil
}

// GetRule looks up a rule in the Obj by the given name and, if present, returns it.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) GetRule(name string) (r *Rule, ok bool) {
	if o == nil {
		return nil, false
	}

	for _, r := range o.Rules {
		if r.Name == name {
			return r, true
		}
	}
	return nil, false
}

// GetKey is a convenience method that returns the Obj's Key.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) GetKey() *Obj {
	if o != nil {
		return o.Key
	}
	return nil
}

// GetElem is a convenience method that returns the Obj's Elem.
// The method is safe to be invoked on a nil instance of *Obj.
func (o *Obj) GetElem() *Obj {
	if o != nil {
		return o.Elem
	}
	return nil
}
