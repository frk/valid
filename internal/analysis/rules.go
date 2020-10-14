package analysis

import (
	"regexp"
)

type ruleCheck func(f *StructField, r *Rule) error

type ruleType struct {
	// name of the rule
	name string
	// the number of parameters that this specific rule type can take
	numParams int
}

var ruleTypes = map[ruleType]ruleCheck{
	{name: "required", numParams: 0}: nil,
	{name: "required", numParams: 1}: checkRuleRequired,
	{name: "required", numParams: 2}: checkRuleRequired,
	{name: "email", numParams: 0}:    checkTypeString,
	{name: "url", numParams: 0}:      checkTypeString,
	{name: "uri", numParams: 0}:      checkTypeString,
	{name: "pan", numParams: 0}:      checkTypeString,
	{name: "cvv", numParams: 0}:      checkTypeString,
	{name: "ssn", numParams: 0}:      checkTypeString,
	{name: "ein", numParams: 0}:      checkTypeString,
	{name: "numeric", numParams: 0}:  checkTypeString,
	{name: "hex", numParams: 0}:      checkTypeString,
	{name: "hexcolor", numParams: 0}: checkTypeString,
	{name: "alphanum", numParams: 0}: checkTypeString,
	{name: "cidr", numParams: 0}:     checkTypeString,
	{name: "phone", numParams: 0}:    checkTypeString,
	{name: "phone", numParams: 1}:    checkMulti(checkTypeString, checkParamCountryCode),
	{name: "zip", numParams: 0}:      checkTypeString,
	{name: "zip", numParams: 1}:      checkMulti(checkTypeString, checkParamCountryCode),
	{name: "uuid", numParams: 0}:     checkTypeString,
	{name: "uuid", numParams: 1}:     checkMulti(checkTypeString, checkParamUUID),
	{name: "ip", numParams: 0}:       checkTypeString,
	{name: "ip", numParams: 1}:       checkMulti(checkTypeString, checkParamIP),
	{name: "mac", numParams: 0}:      checkTypeString,
	{name: "mac", numParams: 1}:      checkMulti(checkTypeString, checkParamMAC),
	{name: "iso", numParams: 1}:      checkMulti(checkTypeString, checkParamISO),
	{name: "rfc", numParams: 1}:      checkMulti(checkTypeString, checkParamRFC),
	{name: "re", numParams: 1}:       checkMulti(checkTypeString, checkParamRegexp),
	{name: "contains", numParams: 1}: checkMulti(checkTypeString, checkParamValue),
	{name: "prefix", numParams: 1}:   checkMulti(checkTypeString, checkParamValue),
	{name: "suffix", numParams: 1}:   checkMulti(checkTypeString, checkParamValue),
	{name: "eq", numParams: 1}:       checkParamCanCompare,
	{name: "ne", numParams: 1}:       checkParamCanCompare,
	{name: "gt", numParams: 1}:       checkMulti(checkTypeNumeric, checkParamCanCompare),
	{name: "lt", numParams: 1}:       checkMulti(checkTypeNumeric, checkParamCanCompare),
	{name: "gte", numParams: 1}:      checkMulti(checkTypeNumeric, checkParamCanCompare),
	{name: "lte", numParams: 1}:      checkMulti(checkTypeNumeric, checkParamCanCompare),
	{name: "min", numParams: 1}:      checkMulti(checkTypeNumeric, checkParamCanCompare),
	{name: "max", numParams: 1}:      checkMulti(checkTypeNumeric, checkParamCanCompare),
	{name: "rng", numParams: 2}:      checkMulti(checkTypeNumeric, checkParamCanCompare),

	{name: "len", numParams: 1}: checkMulti(checkTypeHasLength, checkParamLen),
	{name: "len", numParams: 2}: checkMulti(checkTypeHasLength, checkParamLen),
}

func checkMulti(cc ...ruleCheck) ruleCheck {
	return func(f *StructField, r *Rule) error {
		for _, c := range cc {
			if err := c(f, r); err != nil {
				return err
			}
		}
		return nil
	}
}

// preform check specific to the "required" rule.
func checkRuleRequired(f *StructField, r *Rule) error {
	// checks that each param kind is either context or group key.
	for _, p := range r.Params {
		if p.Kind != ParamKindContext && p.Kind != ParamKindGroupKey {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}
	}

	if len(r.Params) == 2 {
		// if two params were provided ensure that they're not of the same kind.
		if r.Params[0].Kind == r.Params[1].Kind {
			return &anError{Code: errRuleParamKindConflict, RuleParam: r.Params[1]}
		}
	}
	return nil
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

// checks that the rule's param value is a valid country code.
func checkParamCountryCode(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral {
			if !(len(p.Value) == 2 && rxCountryCode2.MatchString(p.Value)) &&
				!(len(p.Value) == 3 && rxCountryCode3.MatchString(p.Value)) {
				return &anError{Code: errRuleParamValueCountryCode, RuleParam: p}
			}
		}
	}
	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// checks that the rule's params values are valid UUID versions.
func checkParamUUID(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral && !rxUUIDVer.MatchString(p.Value) {
			return &anError{Code: errRuleParamValueUUID, RuleParam: p}
		}
	}
	return nil
}

var rxIPVer = regexp.MustCompile(`^(?:v?(?:4|6))$`)

// checks that the rule's params values are valid IP versions.
func checkParamIP(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral && !rxIPVer.MatchString(p.Value) {
			return &anError{Code: errRuleParamValueIP, RuleParam: p}
		}
	}
	return nil
}

var rxMACVer = regexp.MustCompile(`^(?:v?(?:6|8))$`)

// checks that the rule's params values are valid MAC versions.
func checkParamMAC(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral && !rxMACVer.MatchString(p.Value) {
			return &anError{Code: errRuleParamValueMAC, RuleParam: p}
		}
	}
	return nil
}

// checks that the rule's param value is a supported ISO standard identifier.
func checkParamISO(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		// TODO if literal match value against a list of supported ISO standards
	}
	return nil
}

// checks that the rule's param value is a supported RFC standard identifier.
func checkParamRFC(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		// TODO if literal match value against a list of supported RFC standards
	}
	return nil
}

// checks that the rule's param values are either of the literal or reference kind.
func checkParamValue(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}
	}
	return nil
}

// checks that the rule's params are strings containing compilable regular expressions.
func checkParamRegexp(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral {
			if p.Type != ParamTypeString {
				return &anError{Code: errRuleParamTypeString, RuleParam: p}
			}
			if _, err := regexp.Compile(p.Value); err != nil {
				return &anError{Code: errRuleParamValueRegexp, RuleParam: p, Err: err}
			}
		}
	}
	return nil
}

// checks that the rule's params are of the uint type.
func checkParamUint(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral && p.Type != ParamTypeUint {
			return &anError{Code: errRuleParamTypeUint, RuleParam: p}
		}
	}
	return nil
}

// checks that the rule's params are of the uint type, if present.
func checkParamLen(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral && p.Type != ParamTypeUint && len(p.Value) > 0 {
			return &anError{Code: errRuleParamTypeUint, RuleParam: p}
		}
	}

	if len(r.Params) == 2 {
		if r.Params[0].Value == "" && r.Params[1].Value == "" {
			return &anError{Code: errRuleParamValueLen, RuleParam: r.Params[1]}
		}
	}

	return nil
}

// checks that the rule's params' values can be compared to the field's value.
func checkParamCanCompare(f *StructField, r *Rule) error {
	for _, p := range r.Params {
		if p.Kind != ParamKindReference && p.Kind != ParamKindLiteral {
			return &anError{Code: errRuleParamKind, RuleParam: p}
		}

		if p.Kind == ParamKindLiteral {
			if p.Type == ParamTypeNint && !hasTypeKind(f, TypeKindInt,
				TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
				TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleParamTypeNint, RuleParam: p}
			}
			if p.Type == ParamTypeUint && !hasTypeKind(f, TypeKindInt,
				TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
				TypeKindUint, TypeKindUint8, TypeKindUint16, TypeKindUint32,
				TypeKindUint64, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleParamTypeUint, RuleParam: p}
			}
			if p.Type == ParamTypeFloat && !hasTypeKind(f, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleParamTypeFloat, RuleParam: p}
			}
			if p.Type == ParamTypeString && !hasTypeKind(f, TypeKindString) {
				return &anError{Code: errRuleParamTypeFloat, RuleParam: p}
			}
		}
	}
	return nil
}

// checks that the field's type is of the string kind.
func checkTypeString(f *StructField, r *Rule) error {
	if !hasTypeKind(f, TypeKindString) {
		return &anError{Code: errTypeString}
	}
	return nil
}

// checks that the field's type is one of the int/uint/float types.
func checkTypeNumeric(f *StructField, r *Rule) error {
	ok := hasTypeKind(f, TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
		TypeKindUint, TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
		TypeKindFloat32, TypeKindFloat64)

	if !ok {
		return &anError{Code: errTypeNumeric}
	}
	return nil
}

// checks that the field's type can be passed to the builtin len func.
func checkTypeHasLength(f *StructField, r *Rule) error {
	ok := hasTypeKind(f, TypeKindString, TypeKindArray, TypeKindSlice, TypeKindMap)
	if !ok {
		return &anError{Code: errTypeLength}
	}
	return nil
}

func hasTypeKind(f *StructField, kinds ...TypeKind) bool {
	kind := f.Type.Kind
	for kind == TypeKindPtr {
		kind = f.Type.Elem.Kind
	}

	for _, k := range kinds {
		if kind == k {
			return true
		}
	}
	return false
}
