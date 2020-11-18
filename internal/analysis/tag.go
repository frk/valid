package analysis

import (
	"reflect"
	"regexp"
)

type RuleTag struct {
	Rules     []*Rule
	Key, Elem *RuleTag
}

func (rt *RuleTag) HasRuleRequired() bool {
	if rt != nil {
		for _, r := range rt.Rules {
			if r.Name == "required" {
				return true
			}
		}
	}
	return false
}

func (rt *RuleTag) HasRuleNotnil() bool {
	if rt != nil {
		for _, r := range rt.Rules {
			if r.Name == "notnil" {
				return true
			}
		}
	}
	return false
}

// ContainsRules reports whether or not the RuleTag rt, or any of
// the RuleTags in the key-elem hierarchy of rt, contain validation rules.
func (rt *RuleTag) ContainsRules() bool {
	if rt != nil {
		if len(rt.Rules) > 0 {
			return true
		}
		if rt.Key.ContainsRules() {
			return true
		}
		if rt.Elem.ContainsRules() {
			return true
		}
	}
	return false
}

var rxInt = regexp.MustCompile(`^(?:0|-?[1-9][0-9]*)$`)
var rxFloat = regexp.MustCompile(`^(?:(?:-?0|[1-9][0-9]*)?\.[0-9]+)$`)
var rxBool = regexp.MustCompile(`^(?:false|true)$`)

func parseRuleTag(tag string) (*RuleTag, error) {
	val, ok := reflect.StructTag(tag).Lookup("is")
	if !ok || val == "-" || len(val) == 0 {
		return &RuleTag{}, nil
	}

	// parser is invoked recursively to parse tags enclosed in square brackets.
	var parser func(tag string) (*RuleTag, error)
	parser = func(tag string) (*RuleTag, error) {
		rt := &RuleTag{}
		for tag != "" {
			// skip leading space
			i := 0
			for i < len(tag) && tag[i] == ' ' {
				i++
			}
			tag = tag[i:]
			if tag == "" {
				break
			}

			// parse bracketed rules
			if tag[0] == '[' {

				// scan up to the *matching* closing bracket
				i, n := 1, 0
				for i < len(tag) && (tag[i] != ']' || n > 0) {
					// adjust nesting level
					if tag[i] == '[' {
						n++
					} else if tag[i] == ']' {
						n--
					}
					i++

					// scan quoted string, ignoring brackets inside quotes
					if tag[i-1] == '"' {
						for i < len(tag) && tag[i] != '"' {
							if tag[i] == '\\' {
								i++
							}
							i++
						}

						// keep the closing double quote, or
						// else the subsequent parser calls
						// will be confused without it
						if i < len(tag) {
							i++
						}
					}
				}

				// recursively invoke parser for key
				if ktag := tag[1:i]; len(ktag) > 0 {
					key, err := parser(ktag)
					if err != nil {
						return nil, err
					}
					rt.Key = key
				}
				// recursively invoke parser for elem
				if etag := tag[i:]; len(etag) > 1 {
					etag = etag[1:] // drop the leading ']'
					elem, err := parser(etag)
					if err != nil {
						return nil, err
					}
					rt.Elem = elem
				}

				// done; exit
				return rt, nil
			}

			// scan to the end of a rule's name
			i = 0
			for i < len(tag) && tag[i] != ',' && tag[i] != ':' {
				i++
			}

			// empty name's no good; next
			if tag[:i] == "" {
				tag = tag[1:]
				continue
			}

			r := &Rule{Name: tag[:i]}
			rt.Rules = append(rt.Rules, r)

			// this rule's done; next or exit
			if tag = tag[i:]; tag == "" {
				break
			} else if tag[0] == ',' {
				tag = tag[1:]
				continue
			}

			// scan the rule's arguments
			for tag != "" {
				tag = tag[1:] // drop the leading ':'

				// quoted arg value; scan to the end quote
				if len(tag) > 0 && tag[0] == '"' {
					i := 1
					for i < len(tag) && tag[i] != '"' {
						if tag[i] == '\\' {
							i++
						}
						i++
					}

					ra := &RuleArg{}
					ra.Value = tag[1:i]
					ra.Type = ArgTypeString
					r.Args = append(r.Args, ra)

					tag = tag[i:]

					// drop the closing quote
					if len(tag) > 0 && tag[0] == '"' {
						tag = tag[1:]
					}

					// next arg?
					if len(tag) > 0 && tag[0] == ':' {
						continue
					}

					// drop rule separator
					if len(tag) > 0 && tag[0] == ',' {
						tag = tag[1:]
					}

					// this rule's done; exit
					break
				}

				// scan to the end of a rule's argument
				i := 0
				for i < len(tag) && tag[i] != ':' && tag[i] != ',' {
					i++
				}

				arg := tag[:i]
				tag = tag[i:]
				ra := &RuleArg{}

				// resolve the type of the rule's argument
				if len(arg) > 0 {
					switch arg[0] {
					case '@':
						ra = nil // don't append non-arg
						r.Context = arg[1:]
					case '&':
						ra.Value = arg[1:]
						ra.Type = ArgTypeField
					default:
						ra.Value = arg
						switch {
						case rxInt.MatchString(arg):
							ra.Type = ArgTypeInt
						case rxFloat.MatchString(arg):
							ra.Type = ArgTypeFloat
						case rxBool.MatchString(arg):
							ra.Type = ArgTypeBool
						default:
							ra.Type = ArgTypeString
						}
					}
				}

				if ra != nil {
					r.Args = append(r.Args, ra)
				}

				if tag == "" {
					break
				} else if tag[0] == ',' {
					tag = tag[1:]
					break
				}
			}
		}
		return rt, nil
	}

	return parser(val)
}
