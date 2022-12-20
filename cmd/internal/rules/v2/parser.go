package rules

import (
	"reflect"

	"github.com/frk/valid"
)

// ParseTag parses the given tag and returns the resulting TagNode.
// The string "is" is used as the key in the struct tag, e.g. `is:"rule"`.
//
// If the "is" key is not present in the tag, or the associated value
// is empty or equal to "-", then the returned TagNode will be nil.
//
// Following is a description of the rule syntax using EBNF:
//
//      node      = rule | [ "[" [ node ] "]" ] [ ( node | rule "," node ) ] .
//      rule      = rule_name [ { ":" rule_arg } ] { "," rule } .
//      rule_name = identifier .
//      rule_arg  = | boolean_lit | integer_lit | float_lit | string_lit | quoted_string_lit | field_ref .
//
//      boolean_lit       = "true" | "false" .
//      integer_lit       = "0" | [ "-" ] "1"…"9" { "0"…"9" } .
//      float_lit         = [ "-" ] ( "0" | "1"…"9" { "0"…"9" } ) "." "0"…"9" { "0"…"9" } .
//      string_lit        = .
//      quoted_string_lit = `"` `"` .
//
//      field_ref = field_ref_abs | field_ref_rel
//      field_ref_abs   = "&" field_key .
//      field_ref_rel   = "." field_key .
//      field_key       = identifier { field_key_sep identifier } .
//      field_key_sep   = "."
//
//      identifier = letter [ { letter | digit } ] .
//      letter     = "A"…"Z" | "a"…"z" | "_" .
//      digit      = "0"…"9"
//
func ParseTag(tag string, key string) *TagNode {
	str, ok := reflect.StructTag(tag).Lookup(key)
	if !ok || str == "-" || len(str) == 0 {
		return nil
	}
	return parseTag(str, key)
}

// parseTag parses the given rule string and returns the AST.
func parseTag(str string, stkey string) *TagNode {
	tag := &TagNode{stkey: stkey}
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
				tag.Key = parseTag(ktag, stkey)
			}
			// recursively invoke parser for elem
			if etag := str[i:]; len(etag) > 1 {
				etag = etag[1:] // drop the leading ']'
				tag.Elem = parseTag(etag, stkey)
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
