package analysis

import (
	"regexp"
	"strconv"
	"strings"
)

type ruleCheckFuncs []func(a *analysis, f *StructField, r *Rule) error

func (cc ruleCheckFuncs) check(a *analysis, f *StructField, r *Rule) error {
	for _, c := range cc {
		if err := c(a, f, r); err != nil {
			return err
		}
	}
	return nil
}

type ruleTypeMap map[string]map[struct {
	// concrete number of allowed arguments, if set to -1 a variadic number is allowed
	// if num is set, mix & max are ignored
	num int
	// min & max can be used to specify a range of the allowed number of arguments
	// if max is set to -1 it will be interpreted as a range with no upper limit
	min, max int
}]ruleCheckFuncs

func (m ruleTypeMap) find(r *Rule) (ruleCheckFuncs, error) {
	if pm, ok := m[strings.ToLower(r.Name)]; ok {
		for arg, rc := range pm {
			if arg.num != 0 {
				if arg.num == len(r.Args) {
					return rc, nil
				}
				if arg.num == -1 {
					return rc, nil
				}
				continue
			}

			if arg.max != 0 {
				if arg.max >= len(r.Args) && arg.min <= len(r.Args) {
					return rc, nil
				}
				if arg.max == -1 && arg.min <= len(r.Args) {
					return rc, nil
				}
				continue
			}

			if arg.num == len(r.Args) {
				return rc, nil
			}
		}

		return nil, &anError{Code: errRuleArgNum}
	}
	return nil, &anError{Code: errRuleUnknown}
}

func (m ruleTypeMap) check(a *analysis, f *StructField, r *Rule) error {
	c, err := m.find(r)
	if err != nil {
		return err
	}
	return c.check(a, f, r)
}

var ruleTypes = ruleTypeMap{
	"required": {{num: 0}: {}},
	"notnil":   {{num: 0}: {checkTypeIsNilable}},

	"email":    {{num: 0}: {checkTypeKindString}},
	"url":      {{num: 0}: {checkTypeKindString}},
	"uri":      {{num: 0}: {checkTypeKindString}},
	"pan":      {{num: 0}: {checkTypeKindString}},
	"cvv":      {{num: 0}: {checkTypeKindString}},
	"ssn":      {{num: 0}: {checkTypeKindString}},
	"ein":      {{num: 0}: {checkTypeKindString}},
	"numeric":  {{num: 0}: {checkTypeKindString}},
	"hex":      {{num: 0}: {checkTypeKindString}},
	"hexcolor": {{num: 0}: {checkTypeKindString}},
	"alphanum": {{num: 0}: {checkTypeKindString}},
	"cidr":     {{num: 0}: {checkTypeKindString}},

	"phone":    {{num: -1}: {checkTypeKindString, checkArgCountryCode}},
	"zip":      {{num: -1}: {checkTypeKindString, checkArgCountryCode}},
	"uuid":     {{max: 5}: {checkTypeKindString, checkArgUUID}},
	"ip":       {{max: 2}: {checkTypeKindString, checkArgIP}},
	"mac":      {{max: 2}: {checkTypeKindString, checkArgMAC}},
	"iso":      {{num: 1}: {checkTypeKindString, checkArgISO}},
	"rfc":      {{num: 1}: {checkTypeKindString, checkArgRFC}},
	"re":       {{num: 1}: {checkTypeKindString, checkArgRegexp}},
	"prefix":   {{min: 1, max: -1}: {checkTypeKindString, checkArgValueString}},
	"suffix":   {{min: 1, max: -1}: {checkTypeKindString, checkArgValueString}},
	"contains": {{min: 1, max: -1}: {checkTypeKindString, checkArgValueString}},
	"eq":       {{min: 1, max: -1}: {checkArgCanCompare}},
	"ne":       {{min: 1, max: -1}: {checkArgCanCompare}},
	"gt":       {{num: 1}: {checkTypeNumeric, checkArgCanCompare}},
	"lt":       {{num: 1}: {checkTypeNumeric, checkArgCanCompare}},
	"gte":      {{num: 1}: {checkTypeNumeric, checkArgCanCompare}},
	"lte":      {{num: 1}: {checkTypeNumeric, checkArgCanCompare}},
	"min":      {{num: 1}: {checkTypeNumeric, checkArgCanCompare}},
	"max":      {{num: 1}: {checkTypeNumeric, checkArgCanCompare}},
	"rng":      {{num: 2}: {checkTypeNumeric, checkArgCanCompare, checkArgRange}},
	"len":      {{min: 1, max: 2}: {checkTypeHasLength, checkArgLen}},
}

var rxCountryCode2 = regexp.MustCompile(`^(?i:` +
	`a(?:d|e|f|g|i|l|m|o|q|r|s|t|u|w|x|z)|` +
	`b(?:a|b|d|e|f|g|h|i|j|l|m|n|o|q|r|s|t|v|w|y|z)|` +
	`c(?:a|c|d|f|g|h|i|k|l|m|n|o|r|u|v|w|x|y|z)|` +
	`d(?:e|j|k|m|o|z)|` +
	`e(?:c|e|g|h|r|s|t)|` +
	`f(?:i|j|k|m|o|r)|` +
	`g(?:a|b|d|e|f|g|h|i|l|m|n|p|q|r|s|t|u|w|y)|` +
	`h(?:k|m|n|r|t|u)|` +
	`i(?:d|e|l|m|n|o|q|r|s|t)|` +
	`j(?:e|m|o|p)|` +
	`k(?:e|g|h|i|m|n|p|r|w|y|z)|` +
	`l(?:a|b|c|i|k|r|s|t|u|v|y)|` +
	`m(?:a|c|d|e|f|g|h|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)|` +
	`n(?:a|c|e|f|g|i|l|o|p|r|u|z)|` +
	`om|` +
	`p(?:a|e|f|g|h|k|l|m|n|r|s|t|w|y)|` +
	`qa|` +
	`r(?:e|o|s|u|w)|` +
	`s(?:a|b|c|d|e|g|h|i|j|k|l|m|n|o|r|s|t|v|x|y|z)|` +
	`t(?:c|d|f|g|h|j|k|l|m|n|o|r|t|v|w|z)|` +
	`u(?:a|g|m|s|y|z)|` +
	`v(?:a|c|e|g|i|n|u)|` +
	`w(?:f|s)|` +
	`y(?:e|t)|` +
	`z(?:a|m|w)` +
	`)$`)

var rxCountryCode3 = regexp.MustCompile(`^(?i:` +
	`a(?:bw|fg|go|ia|la|lb|nd|re|rg|rm|sm|ta|tf|tg|us|ut|ze)|` +
	`b(?:di|el|en|es|fa|gd|gr|hr|hs|ih|lm|lr|lz|mu|ol|ra|rb|rn|tn|vt|wa)|` +
	`c(?:af|an|ck|he|hl|hn|iv|mr|od|og|ok|ol|om|pv|ri|ub|uw|xr|ym|yp|ze)|` +
	`d(?:eu|ji|ma|nk|om|za)|` +
	`e(?:cu|gy|ri|sh|sp|st|th)|` +
	`f(?:in|ji|lk|ra|ro|sm)|` +
	`g(?:ab|br|eo|gy|ha|ib|in|lp|mb|nb|nq|rc|rd|rl|tm|uf|um|uy)|` +
	`h(?:kg|md|nd|rv|ti|un)|` +
	`i(?:dn|mn|nd|ot|rl|rn|rq|sl|sr|ta)|` +
	`j(?:am|ey|or|pn)|` +
	`k(?:az|en|gz|hm|ir|na|or|wt)|` +
	`l(?:ao|bn|br|by|ca|ie|ka|so|tu|ux|va)|` +
	`m(?:ac|af|ar|co|da|dg|dv|ex|hl|kd|li|lt|mr|ne|ng|np|oz|rt|sr|tq|us|wi|ys|yt)|` +
	`n(?:am|cl|er|fk|ga|ic|iu|ld|or|pl|ru|zl)|` +
	`omn|` +
	`p(?:ak|an|cn|er|hl|lw|ng|ol|ri|rk|rt|ry|se|yf)|` +
	`qat|` +
	`r(?:eu|ou|us|wa)|` +
	`s(?:au|dn|en|gp|gs|hn|jm|lb|le|lv|mr|om|pm|rb|sd|tp|ur|vk|vn|we|wz|xm|yc|yr)|` +
	`t(?:ca|cd|go|ha|jk|kl|km|ls|on|to|un|ur|uv|wn|za)|` +
	`u(?:ga|kr|mi|ry|sa|zb)|` +
	`v(?:at|ct|en|gb|ir|nm|ut)|` +
	`w(?:lf|sm)|` +
	`yem|` +
	`z(?:af|mb|we)|` +
	`)$`)

// checks that the rule's arg value is a valid country code.
func checkArgCountryCode(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeString:
			if !(len(arg.Value) == 2 && rxCountryCode2.MatchString(arg.Value)) &&
				!(len(arg.Value) == 3 && rxCountryCode3.MatchString(arg.Value)) {
				return &anError{Code: errRuleArgValueCountryCode, RuleArg: arg}
			}
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// checks that the rule's args values are valid UUID versions.
func checkArgUUID(a *analysis, f *StructField, r *Rule) error {
	versions := map[string]struct{}{} // track encountered versions
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeString, ArgTypeUint:
			if !rxUUIDVer.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueUUIDVer, RuleArg: arg}
			}

			v := arg.Value
			if len(v) > 1 {
				v = v[1:]
			}
			if _, exists := versions[v]; exists {
				return &anError{Code: errRuleArgValueConflict, RuleArg: arg}
			} else {
				versions[v] = struct{}{}
			}
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString, TypeKindUint,
				TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
				TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxIPVer = regexp.MustCompile(`^(?:v?(?:4|6))$`)

// checks that the rule's args values are valid IP versions.
func checkArgIP(a *analysis, f *StructField, r *Rule) error {
	var version string // the first version
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeString, ArgTypeUint:
			if !rxIPVer.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueIPVer, RuleArg: arg}
			}

			v := arg.Value
			if len(v) > 1 {
				v = v[1:]
			}
			if len(version) > 0 && version == v {
				return &anError{Code: errRuleArgValueConflict, RuleArg: arg}
			} else {
				version = v
			}
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString, TypeKindUint,
				TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
				TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxMACVer = regexp.MustCompile(`^(?:v?(?:6|8))$`)

// checks that the rule's args values are valid MAC versions.
func checkArgMAC(a *analysis, f *StructField, r *Rule) error {
	var version string // the first version
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeString, ArgTypeUint:
			if !rxMACVer.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueMACVer, RuleArg: arg}
			}

			v := arg.Value
			if len(v) > 1 {
				v = v[1:]
			}
			if len(version) > 0 && version == v {
				return &anError{Code: errRuleArgValueConflict, RuleArg: arg}
			} else {
				version = v
			}
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString, TypeKindUint,
				TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
				TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxISOStd = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported ISO standard identifier.
func checkArgISO(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeUint:
			if !rxISOStd.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueISOStd, RuleArg: arg}
			}
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString, TypeKindUint,
				TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
				TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxRFCStd = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported RFC standard identifier.
func checkArgRFC(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeUint:
			if !rxRFCStd.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueRFCStd, RuleArg: arg}
			}
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString, TypeKindUint,
				TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
				TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

// checks that the rule's args are strings containing compilable regular expressions.
func checkArgRegexp(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			if _, err := regexp.Compile(arg.Value); err != nil {
				return &anError{Code: errRuleArgValueRegexp, RuleArg: arg, Err: err}
			}
		}
	}
	return nil
}

// checks that the rule's arg values are strings.
func checkArgValueString(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindString) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			// any of the non-reference types can be represented as string...
		}
	}
	return nil
}

// checks that the args are valid for use with the "rng" rule.
func checkArgRange(a *analysis, f *StructField, r *Rule) error {
	// r.Args is known to be of length 2
	if r.Args[0].Type == ArgTypeReference || r.Args[1].Type == ArgTypeReference {
		// nothing to do if we don't know both the values beforehand
		return nil
	}

	var lower, upper *float64
	if a := r.Args[0]; len(a.Value) > 0 {
		f64, err := strconv.ParseFloat(a.Value, 64)
		if err != nil {
			return &anError{Code: errRuleArgValueParseFloat, RuleArg: a, Err: err}
		}
		lower = &f64
	}
	if a := r.Args[1]; len(a.Value) > 0 {
		f64, err := strconv.ParseFloat(a.Value, 64)
		if err != nil {
			return &anError{Code: errRuleArgValueParseFloat, RuleArg: a, Err: err}
		}
		upper = &f64
	}

	// make sure at least one bound was specified and if both then
	// ensure that the lower bound is less than the upper bound
	if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
		return &anError{Code: errRuleArgValueBounds, RuleArg: r.Args[1]}
	}
	return nil
}

// checks that the rule's args are of the uint type, if present.
func checkArgLen(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeUint:
			// ok
		case ArgTypeString:
			if len(arg.Value) > 0 {
				return &anError{Code: errRuleArgType, RuleArg: arg}
			}
			// empty string means no bound (ok)
		case ArgTypeReference:
			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, TypeKindUint, TypeKindUint8, TypeKindUint16,
				TypeKindUint32, TypeKindUint64, TypeKindInt, TypeKindInt8,
				TypeKindInt16, TypeKindInt32, TypeKindInt64) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		default:
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}

	if len(r.Args) == 2 {
		if r.Args[0].Type == ArgTypeReference || r.Args[1].Type == ArgTypeReference {
			// nothing to do if we don't know both the values beforehand
			return nil
		}

		var lower, upper *uint64
		if a := r.Args[0]; len(a.Value) > 0 {
			u64, err := strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleArgValueParseUint, RuleArg: a, Err: err}
			}
			lower = &u64
		}
		if a := r.Args[1]; len(a.Value) > 0 {
			u64, err := strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleArgValueParseUint, RuleArg: a, Err: err}
			}
			upper = &u64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return &anError{Code: errRuleArgValueBounds, RuleArg: r.Args[1]}
		}
	}

	return nil
}

// checks that the rule's args' values can be compared to the field's value.
func checkArgCanCompare(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		switch arg.Type {
		case ArgTypeString:
			if !hasTypeKind(f, TypeKindString) {
				if len(arg.Value) > 0 {
					return &anError{Code: errRuleArgTypeString, RuleArg: arg}
				}

				// arg of string type with empty value can be used for 0 int/uint/float
				if !hasTypeKind(f, TypeKindInt, TypeKindInt8, TypeKindInt16,
					TypeKindInt32, TypeKindInt64, TypeKindUint, TypeKindUint8, TypeKindUint16,
					TypeKindUint32, TypeKindUint64, TypeKindFloat32, TypeKindFloat64) {
					return &anError{Code: errRuleArgTypeString, RuleArg: arg}
				}
			}
		case ArgTypeNint:
			if !hasTypeKind(f, TypeKindString, TypeKindInt, TypeKindInt8, TypeKindInt16,
				TypeKindInt32, TypeKindInt64, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleArgTypeNint, RuleArg: arg}
			}
		case ArgTypeUint:
			if !hasTypeKind(f, TypeKindString, TypeKindInt, TypeKindInt8, TypeKindInt16,
				TypeKindInt32, TypeKindInt64, TypeKindUint, TypeKindUint8, TypeKindUint16,
				TypeKindUint32, TypeKindUint64, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleArgTypeUint, RuleArg: arg}
			}
		case ArgTypeFloat:
			if !hasTypeKind(f, TypeKindString, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleArgTypeFloat, RuleArg: arg}
			}
		case ArgTypeReference:
			kind := f.Type.Kind
			for kind == TypeKindPtr {
				kind = f.Type.Elem.Kind
			}

			f2 := a.info.ArgReferenceMap[arg].SelectorLast()
			if !hasTypeKind(f2, kind) {
				return &anError{Code: errRuleArgTypeReferenceKind, RuleArg: arg}
			}
		}
	}
	return nil
}

// checks that the field's type is of the string kind.
func checkTypeKindString(a *analysis, f *StructField, r *Rule) error {
	if !hasTypeKind(f, TypeKindString) {
		return &anError{Code: errTypeKindString}
	}
	return nil
}

// checks that the field's type is one of the int/uint/float types.
func checkTypeNumeric(a *analysis, f *StructField, r *Rule) error {
	ok := hasTypeKind(f, TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
		TypeKindUint, TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
		TypeKindFloat32, TypeKindFloat64)

	if !ok {
		return &anError{Code: errTypeNumeric}
	}
	return nil
}

// checks that the field's type can be passed to the builtin len func.
func checkTypeHasLength(a *analysis, f *StructField, r *Rule) error {
	ok := hasTypeKind(f, TypeKindString, TypeKindArray, TypeKindSlice, TypeKindMap)
	if !ok {
		return &anError{Code: errTypeLength}
	}
	return nil
}

// checks that the field's type can be nil.
func checkTypeIsNilable(a *analysis, f *StructField, r *Rule) error {
	ok := hasTypeKind(f, TypeKindPtr, TypeKindSlice, TypeKindMap, TypeKindInterface)
	if !ok {
		return &anError{Code: errTypeNil}
	}
	return nil
}

func hasTypeKind(f *StructField, kinds ...TypeKind) bool {
	isptr := f.Type.Kind == TypeKindPtr

	typ := f.Type
	for typ.Kind == TypeKindPtr {
		typ = *typ.Elem
	}
	for _, k := range kinds {
		if k == typ.Kind {
			return true
		}
		if k == TypeKindPtr && isptr {
			return true
		}
	}
	return false
}
