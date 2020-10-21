package analysis

import (
	"regexp"
	"strings"
)

type ruleCheckFuncs []func(f *StructField, r *Rule) error

func (cc ruleCheckFuncs) check(f *StructField, r *Rule) error {
	for _, c := range cc {
		if err := c(f, r); err != nil {
			return err
		}
	}
	return nil
}

type ruleTypeMap map[string]map[int]ruleCheckFuncs

func (m ruleTypeMap) find(r *Rule) (ruleCheckFuncs, error) {
	if pm, ok := m[strings.ToLower(r.Name)]; ok {
		var variadic ruleCheckFuncs
		for num, rc := range pm {
			if num == len(r.Args) {
				return rc, nil
			}

			// -1 indicates variadic rule
			if num == -1 {
				variadic = rc
			}
		}
		if variadic != nil {
			return variadic, nil
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
	return c.check(f, r)
}

var ruleTypes = ruleTypeMap{
	"required": {0: {}, 1: {checkRuleRequired}, 2: {checkRuleRequired}},
	"email":    {0: {checkTypeKindString}},
	"url":      {0: {checkTypeKindString}},
	"uri":      {0: {checkTypeKindString}},
	"pan":      {0: {checkTypeKindString}},
	"cvv":      {0: {checkTypeKindString}},
	"ssn":      {0: {checkTypeKindString}},
	"ein":      {0: {checkTypeKindString}},
	"numeric":  {0: {checkTypeKindString}},
	"hex":      {0: {checkTypeKindString}},
	"hexcolor": {0: {checkTypeKindString}},
	"alphanum": {0: {checkTypeKindString}},
	"cidr":     {0: {checkTypeKindString}},
	"phone":    {-1: {checkTypeKindString, checkArgCountryCode}},
	"zip":      {-1: {checkTypeKindString, checkArgCountryCode}},
	"uuid":     {-1: {checkTypeKindString, checkArgUUID}},
	"ip":       {0: {checkTypeKindString}, 1: {checkTypeKindString, checkArgIP}},
	"mac":      {0: {checkTypeKindString}, 1: {checkTypeKindString, checkArgMAC}},
	"iso":      {1: {checkTypeKindString, checkArgISO}},
	"rfc":      {1: {checkTypeKindString, checkArgRFC}},

	"re":       {1: {checkTypeKindString, checkArgRegexp}},
	"contains": {1: {checkTypeKindString, checkArgValue}},
	"prefix":   {1: {checkTypeKindString, checkArgValue}},
	"suffix":   {1: {checkTypeKindString, checkArgValue}},
	"eq":       {1: {checkArgCanCompare}},
	"ne":       {1: {checkArgCanCompare}},
	"gt":       {1: {checkTypeNumeric, checkArgCanCompare}},
	"lt":       {1: {checkTypeNumeric, checkArgCanCompare}},
	"gte":      {1: {checkTypeNumeric, checkArgCanCompare}},
	"lte":      {1: {checkTypeNumeric, checkArgCanCompare}},
	"min":      {1: {checkTypeNumeric, checkArgCanCompare}},
	"max":      {1: {checkTypeNumeric, checkArgCanCompare}},
	"rng":      {2: {checkTypeNumeric, checkArgCanCompare}},
	"len":      {1: {checkTypeHasLength, checkArgLen}, 2: {checkTypeHasLength, checkArgLen}},
}

// preform check specific to the "required" rule.
func checkRuleRequired(f *StructField, r *Rule) error {
	// checks that each arg kind is either context or group key.
	for _, a := range r.Args {
		if a.Kind != ArgKindContext && a.Kind != ArgKindGroupKey {
			return &anError{Code: errRuleArgKind, RuleArg: a}
		}
	}

	if len(r.Args) == 2 {
		// if two args were provided ensure that they're not of the same kind.
		if r.Args[0].Kind == r.Args[1].Kind {
			return &anError{Code: errRuleArgKindConflict, RuleArg: r.Args[1]}
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

// checks that the rule's arg value is a valid country code.
func checkArgCountryCode(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral {
			if !(len(p.Value) == 2 && rxCountryCode2.MatchString(p.Value)) &&
				!(len(p.Value) == 3 && rxCountryCode3.MatchString(p.Value)) {
				return &anError{Code: errRuleArgValueCountryCode, RuleArg: p}
			}
		}
	}
	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// checks that the rule's args values are valid UUID versions.
func checkArgUUID(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && !rxUUIDVer.MatchString(p.Value) {
			return &anError{Code: errRuleArgValueUUIDVer, RuleArg: p}
		}
	}
	return nil
}

var rxIPVer = regexp.MustCompile(`^(?:v?(?:4|6))$`)

// checks that the rule's args values are valid IP versions.
func checkArgIP(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && !rxIPVer.MatchString(p.Value) {
			return &anError{Code: errRuleArgValueIPVer, RuleArg: p}
		}
	}
	return nil
}

var rxMACVer = regexp.MustCompile(`^(?:v?(?:6|8))$`)

// checks that the rule's args values are valid MAC versions.
func checkArgMAC(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && !rxMACVer.MatchString(p.Value) {
			return &anError{Code: errRuleArgValueMACVer, RuleArg: p}
		}
	}
	return nil
}

var rxISOStd = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported ISO standard identifier.
func checkArgISO(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && !rxISOStd.MatchString(p.Value) {
			return &anError{Code: errRuleArgValueISOStd, RuleArg: p}
		}
	}
	return nil
}

var rxRFCStd = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported RFC standard identifier.
func checkArgRFC(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && !rxRFCStd.MatchString(p.Value) {
			return &anError{Code: errRuleArgValueRFCStd, RuleArg: p}
		}
	}
	return nil
}

// checks that the rule's arg values are either of the literal or reference kind.
func checkArgValue(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}
	}
	return nil
}

// checks that the rule's args are strings containing compilable regular expressions.
func checkArgRegexp(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral {
			if p.Type != ArgTypeString {
				return &anError{Code: errRuleArgTypeString, RuleArg: p}
			}
			if _, err := regexp.Compile(p.Value); err != nil {
				return &anError{Code: errRuleArgValueRegexp, RuleArg: p, Err: err}
			}
		}
	}
	return nil
}

// checks that the rule's args are of the uint type.
func checkArgUint(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && p.Type != ArgTypeUint {
			return &anError{Code: errRuleArgTypeUint, RuleArg: p}
		}
	}
	return nil
}

// checks that the rule's args are of the uint type, if present.
func checkArgLen(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral && p.Type != ArgTypeUint && len(p.Value) > 0 {
			return &anError{Code: errRuleArgTypeUint, RuleArg: p}
		}
	}

	if len(r.Args) == 2 {
		if r.Args[0].Value == "" && r.Args[1].Value == "" {
			return &anError{Code: errRuleArgValueLen, RuleArg: r.Args[1]}
		}
	}

	return nil
}

// checks that the rule's args' values can be compared to the field's value.
func checkArgCanCompare(f *StructField, r *Rule) error {
	for _, p := range r.Args {
		if p.Kind != ArgKindReference && p.Kind != ArgKindLiteral {
			return &anError{Code: errRuleArgKind, RuleArg: p}
		}

		if p.Kind == ArgKindLiteral {
			if p.Type == ArgTypeNint && !hasTypeKind(f, TypeKindInt,
				TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
				TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleArgTypeNint, RuleArg: p}
			}
			if p.Type == ArgTypeUint && !hasTypeKind(f, TypeKindInt,
				TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
				TypeKindUint, TypeKindUint8, TypeKindUint16, TypeKindUint32,
				TypeKindUint64, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleArgTypeUint, RuleArg: p}
			}
			if p.Type == ArgTypeFloat && !hasTypeKind(f, TypeKindFloat32, TypeKindFloat64) {
				return &anError{Code: errRuleArgTypeFloat, RuleArg: p}
			}
			if p.Type == ArgTypeString && !hasTypeKind(f, TypeKindString) {
				return &anError{Code: errRuleArgTypeFloat, RuleArg: p}
			}
		}
	}
	return nil
}

// checks that the field's type is of the string kind.
func checkTypeKindString(f *StructField, r *Rule) error {
	if !hasTypeKind(f, TypeKindString) {
		return &anError{Code: errTypeKindString}
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
