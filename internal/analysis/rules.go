package analysis

import (
	"regexp"
	"strconv"
)

const (
	pkgisvalid = "github.com/frk/isvalid"
)

var (
	typeInt         = Type{Kind: TypeKindInt}
	typeUint        = Type{Kind: TypeKindUint}
	typeString      = Type{Kind: TypeKindString}
	typeStringSlice = Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindString}}
	typeIntSlice    = Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}}
	typeEmptyIface  = Type{Kind: TypeKindInterface, IsEmptyInterface: true}
)

func rulechecks(checks ...func(a *analysis, f *StructField, r *Rule) error) (check func(a *analysis, f *StructField, r *Rule) error) {
	return func(a *analysis, f *StructField, r *Rule) error {
		for _, chk := range checks {
			if err := chk(a, f, r); err != nil {
				return err
			}
		}
		return nil
	}
}

var defaultRuleSpecMap = map[string]RuleSpec{
	// speciél
	"isvalid": RuleIsValid{},

	// basic rules
	"required": RuleBasic{check: isValidRuleRequired},
	"notnil":   RuleBasic{check: isValidRuleNotnil},
	"rng":      RuleBasic{check: isValidRuleRng},
	"len":      RuleBasic{check: isValidRuleLen},
	"eq":       RuleBasic{check: isValidRuleValueComparison},
	"ne":       RuleBasic{check: isValidRuleValueComparison},
	"gt":       RuleBasic{check: isValidRuleNumberComparison},
	"lt":       RuleBasic{check: isValidRuleNumberComparison},
	"gte":      RuleBasic{check: isValidRuleNumberComparison},
	"lte":      RuleBasic{check: isValidRuleNumberComparison},
	"min":      RuleBasic{check: isValidRuleNumberComparison},
	"max":      RuleBasic{check: isValidRuleNumberComparison},

	// predefined functions
	"email":    RuleFunc{FuncName: "Email", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"url":      RuleFunc{FuncName: "URL", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"uri":      RuleFunc{FuncName: "URI", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"pan":      RuleFunc{FuncName: "PAN", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"cvv":      RuleFunc{FuncName: "CVV", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"ssn":      RuleFunc{FuncName: "SSN", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"ein":      RuleFunc{FuncName: "EIN", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"numeric":  RuleFunc{FuncName: "Numeric", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"hex":      RuleFunc{FuncName: "Hex", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"hexcolor": RuleFunc{FuncName: "HexColor", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"alphanum": RuleFunc{FuncName: "Alphanum", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},
	"cidr":     RuleFunc{FuncName: "CIDR", PkgPath: pkgisvalid, ArgTypes: []Type{typeString}},

	"phone": RuleFunc{
		FuncName: "Phone", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeStringSlice},
		check:    isValidCountryCode},
	"zip": RuleFunc{
		FuncName: "Zip", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeStringSlice},
		check:    isValidCountryCode},
	"uuid": RuleFunc{
		FuncName: "UUID", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeIntSlice},
		check:    isValidRuleUUID},
	"ip": RuleFunc{
		FuncName: "IP", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeIntSlice},
		check:    isValidRuleIP},
	"mac": RuleFunc{
		FuncName: "MAC", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeIntSlice},
		check:    isValidRuleMAC},
	"iso": RuleFunc{
		FuncName: "ISO", PkgPath: pkgisvalid,
		ArgTypes: []Type{typeString, typeInt},
		check:    isValidRuleISO},
	"rfc": RuleFunc{
		FuncName: "RFC", PkgPath: pkgisvalid,
		ArgTypes: []Type{typeString, typeInt},
		check:    isValidRuleRFC},
	"re": RuleFunc{
		FuncName: "Match", PkgPath: pkgisvalid,
		ArgTypes: []Type{typeString, typeString}, UseRawStrings: true,
		check: isValidRuleRegexp},
	"prefix": RuleFunc{
		FuncName: "HasPrefix", PkgPath: "strings",
		ArgTypes: []Type{typeString, typeString},
		BoolConn: RuleFuncBoolOr},
	"suffix": RuleFunc{
		FuncName: "HasSuffix", PkgPath: "strings",
		ArgTypes: []Type{typeString, typeString},
		BoolConn: RuleFuncBoolOr},
	"contains": RuleFunc{
		FuncName: "Contains", PkgPath: "strings",
		ArgTypes: []Type{typeString, typeString},
		BoolConn: RuleFuncBoolOr},
	"notcontains": RuleFunc{
		FuncName: "Contains", PkgPath: "strings",
		ArgTypes: []Type{typeString, typeString},
		BoolConn: RuleFuncBoolAnd},
}

func isValidRuleRequired(a *analysis, f *StructField, r *Rule) error {
	// rule must have 0 args
	if len(r.Args) != 0 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}
	return nil
}

func isValidRuleNotnil(a *analysis, f *StructField, r *Rule) error {
	// rule must have 0 args
	if len(r.Args) != 0 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	// field's type must be nilable
	if !hasTypeKind(f, TypeKindPtr, TypeKindSlice, TypeKindMap, TypeKindInterface) {
		return &anError{Code: errTypeNil}
	}
	return nil
}

// check that the rule's args are strings containing compilable regular expressions.
func isValidRuleRegexp(a *analysis, f *StructField, r *Rule) error {
	for _, a := range r.Args {
		if a.Type != ArgTypeField {
			if _, err := regexp.Compile(a.Value); err != nil {
				return &anError{Code: errRuleArgValueRegexp, RuleArg: a, Err: err}
			}
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "value comparison" rule.
func isValidRuleValueComparison(a *analysis, f *StructField, r *Rule) error {
	// rule must have at least 1 arg, no less
	if len(r.Args) < 1 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	// rule arg must be comparable to the field's type
	typ := f.Type
	for typ.Kind == TypeKindPtr {
		typ = *typ.Elem
	}
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return &anError{Code: errRuleFuncRuleArgType, RuleArg: ra}
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "number comparison" rule.
func isValidRuleNumberComparison(a *analysis, f *StructField, r *Rule) error {
	// rule can have exactly 1 arg, no more no less
	if len(r.Args) != 1 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	// the field's type must be numeric
	if err := typeIsNumeric(a, f, r); err != nil {
		return err
	}

	// rule arg must be comparable to the field's type
	typ := f.Type
	for typ.Kind == TypeKindPtr {
		typ = *typ.Elem
	}
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return &anError{Code: errRuleFuncRuleArgType, RuleArg: ra}
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "len" rule.
func isValidRuleLen(a *analysis, f *StructField, r *Rule) error {
	// associated field's type must have length
	if !hasTypeKind(f, TypeKindString, TypeKindArray, TypeKindSlice, TypeKindMap) {
		return &anError{Code: errTypeLength}
	}

	// rule can have 1 or 2 args, no more no less
	if len(r.Args) < 1 || len(r.Args) > 2 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	// if 2, then make sure the values represent valid upper to lower bounds
	if len(r.Args) == 2 && r.Args[0].Type != ArgTypeField && r.Args[1].Type != ArgTypeField {
		var lower, upper *uint64
		if a := r.Args[0]; len(a.Value) > 0 {
			u64, err := strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleFuncRuleArgType, RuleArg: a, Err: err}
			}
			lower = &u64
		}
		if a := r.Args[1]; len(a.Value) > 0 {
			u64, err := strconv.ParseUint(a.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleFuncRuleArgType, RuleArg: a, Err: err}
			}
			upper = &u64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return &anError{Code: errRuleArgValueBounds, RuleArg: r.Args[1]}
		}
	}

	// rule args must be comparable to a positive integer (the return value of the len builtin)
	typ := Type{Kind: TypeKindUint}
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return &anError{Code: errRuleFuncRuleArgType, RuleArg: ra}
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "rng" rule.
func isValidRuleRng(a *analysis, f *StructField, r *Rule) error {
	// rule must have exactly 2 args, no more no less
	if len(r.Args) != 2 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	// the field's type must be numeric
	if err := typeIsNumeric(a, f, r); err != nil {
		return err
	}

	// make sure the rule's arg values represent valid upper to lower bounds
	if r.Args[0].Type != ArgTypeField && r.Args[1].Type != ArgTypeField {
		var lower, upper *float64
		if a := r.Args[0]; len(a.Value) > 0 {
			f64, err := strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return &anError{Code: errRuleFuncRuleArgType, RuleArg: a, Err: err}
			}
			lower = &f64
		}
		if a := r.Args[1]; len(a.Value) > 0 {
			f64, err := strconv.ParseFloat(a.Value, 64)
			if err != nil {
				return &anError{Code: errRuleFuncRuleArgType, RuleArg: a, Err: err}
			}
			upper = &f64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return &anError{Code: errRuleArgValueBounds, RuleArg: r.Args[1]}
		}
	}

	// rule args must be comparable to the field's type
	typ := f.Type
	for typ.Kind == TypeKindPtr {
		typ = *typ.Elem
	}
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return &anError{Code: errRuleFuncRuleArgType, RuleArg: ra}
		}
	}

	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// check that the RuleArgs represent a valid "uuid" rule.
func isValidRuleUUID(a *analysis, f *StructField, r *Rule) error {
	// rule can have at most 5 args, no more
	if len(r.Args) > 5 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	versions := map[string]struct{}{} // track encountered versions
	for _, a := range r.Args {
		if a.Type == ArgTypeString || a.IsUInt() {
			if !rxUUIDVer.MatchString(a.Value) {
				return &anError{Code: errRuleArgValueUUIDVer, RuleArg: a}
			}

			if len(a.Value) > 1 && (a.Value[0] == 'v' || a.Value[0] == 'V') {
				a.Value = a.Value[1:]
				a.Type = ArgTypeInt
			}
			if _, exists := versions[a.Value]; exists {
				return &anError{Code: errRuleArgValueConflict, RuleArg: a}
			} else {
				versions[a.Value] = struct{}{}
			}
		} else if a.Type != ArgTypeField {
			return &anError{Code: errRuleArgType, RuleArg: a}
		}
	}
	return nil
}

var rxIPVer = regexp.MustCompile(`^(?:v?(?:4|6))$`)

// checks that the rule args' values are valid IP versions.
func isValidRuleIP(a *analysis, f *StructField, r *Rule) error {
	// rule can have at most 2 args, no more
	if len(r.Args) > 2 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	var version string // the first version
	for _, a := range r.Args {
		if a.Type == ArgTypeString || a.IsUInt() {
			if !rxIPVer.MatchString(a.Value) {
				return &anError{Code: errRuleArgValueIPVer, RuleArg: a}
			}

			if len(a.Value) > 1 && (a.Value[0] == 'v' || a.Value[0] == 'V') {
				a.Value = a.Value[1:]
				a.Type = ArgTypeInt
			}
			if len(version) > 0 && version == a.Value {
				return &anError{Code: errRuleArgValueConflict, RuleArg: a}
			} else {
				version = a.Value
			}
		} else if a.Type != ArgTypeField {
			return &anError{Code: errRuleArgType, RuleArg: a}
		}
	}
	return nil
}

var rxMACVer = regexp.MustCompile(`^(?:v?(?:6|8))$`)

// checks that the rule's args values are valid MAC versions.
func isValidRuleMAC(a *analysis, f *StructField, r *Rule) error {
	// rule can have at most 2 args, no more
	if len(r.Args) > 2 {
		return &anError{Code: errRuleFuncRuleArgCount}
	}

	var version string // the first version
	for _, a := range r.Args {
		if a.Type == ArgTypeString || a.IsUInt() {
			if !rxMACVer.MatchString(a.Value) {
				return &anError{Code: errRuleArgValueMACVer, RuleArg: a}
			}

			if len(a.Value) > 1 && (a.Value[0] == 'v' || a.Value[0] == 'V') {
				a.Value = a.Value[1:]
				a.Type = ArgTypeInt
			}
			if len(version) > 0 && version == a.Value {
				return &anError{Code: errRuleArgValueConflict, RuleArg: a}
			} else {
				version = a.Value
			}
		} else if a.Type != ArgTypeField {
			return &anError{Code: errRuleArgType, RuleArg: a}
		}
	}
	return nil
}

var rxISO = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported ISO standard identifier.
func isValidRuleISO(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		if arg.IsUInt() {
			// TODO remove the regex and instead check against
			// a list of supported ISO validators
			if !rxISO.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueISOStd, RuleArg: arg}
			}
		} else if arg.Type != ArgTypeField {
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxRFC = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported RFC standard identifier.
func isValidRuleRFC(a *analysis, f *StructField, r *Rule) error {
	for _, arg := range r.Args {
		if arg.IsUInt() {
			// TODO remove the regex and instead check against
			// a list of supported RFC validators
			if !rxRFC.MatchString(arg.Value) {
				return &anError{Code: errRuleArgValueRFCStd, RuleArg: arg}
			}
		} else if arg.Type != ArgTypeField {
			return &anError{Code: errRuleArgType, RuleArg: arg}
		}
	}
	return nil
}

var rxCountryCode2 = regexp.MustCompile(`^(?i:a(?:d|e|f|g|i|l|m|o|q|r|s|t|u|w|x|z)|b(?:a|b|d|e|f|g|h|i|j|l|m|n|o|q|r|s|t|v|w|y|z)|c(?:a|c|d|f|g|h|i|k|l|m|n|o|r|u|v|w|x|y|z)|d(?:e|j|k|m|o|z)|e(?:c|e|g|h|r|s|t)|f(?:i|j|k|m|o|r)|g(?:a|b|d|e|f|g|h|i|l|m|n|p|q|r|s|t|u|w|y)|h(?:k|m|n|r|t|u)|i(?:d|e|l|m|n|o|q|r|s|t)|j(?:e|m|o|p)|k(?:e|g|h|i|m|n|p|r|w|y|z)|l(?:a|b|c|i|k|r|s|t|u|v|y)|m(?:a|c|d|e|f|g|h|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)|n(?:a|c|e|f|g|i|l|o|p|r|u|z)|om|p(?:a|e|f|g|h|k|l|m|n|r|s|t|w|y)|qa|r(?:e|o|s|u|w)|s(?:a|b|c|d|e|g|h|i|j|k|l|m|n|o|r|s|t|v|x|y|z)|t(?:c|d|f|g|h|j|k|l|m|n|o|r|t|v|w|z)|u(?:a|g|m|s|y|z)|v(?:a|c|e|g|i|n|u)|w(?:f|s)|y(?:e|t)|z(?:a|m|w))$`)
var rxCountryCode3 = regexp.MustCompile(`^(?i:a(?:bw|fg|go|ia|la|lb|nd|re|rg|rm|sm|ta|tf|tg|us|ut|ze)|b(?:di|el|en|es|fa|gd|gr|hr|hs|ih|lm|lr|lz|mu|ol|ra|rb|rn|tn|vt|wa)|c(?:af|an|ck|he|hl|hn|iv|mr|od|og|ok|ol|om|pv|ri|ub|uw|xr|ym|yp|ze)|d(?:eu|ji|ma|nk|om|za)|e(?:cu|gy|ri|sh|sp|st|th)|f(?:in|ji|lk|ra|ro|sm)|g(?:ab|br|eo|gy|ha|ib|in|lp|mb|nb|nq|rc|rd|rl|tm|uf|um|uy)|h(?:kg|md|nd|rv|ti|un)|i(?:dn|mn|nd|ot|rl|rn|rq|sl|sr|ta)|j(?:am|ey|or|pn)|k(?:az|en|gz|hm|ir|na|or|wt)|l(?:ao|bn|br|by|ca|ie|ka|so|tu|ux|va)|m(?:ac|af|ar|co|da|dg|dv|ex|hl|kd|li|lt|mr|ne|ng|np|oz|rt|sr|tq|us|wi|ys|yt)|n(?:am|cl|er|fk|ga|ic|iu|ld|or|pl|ru|zl)|omn|p(?:ak|an|cn|er|hl|lw|ng|ol|ri|rk|rt|ry|se|yf)|qat|r(?:eu|ou|us|wa)|s(?:au|dn|en|gp|gs|hn|jm|lb|le|lv|mr|om|pm|rb|sd|tp|ur|vk|vn|we|wz|xm|yc|yr)|t(?:ca|cd|go|ha|jk|kl|km|ls|on|to|un|ur|uv|wn|za)|u(?:ga|kr|mi|ry|sa|zb)|v(?:at|ct|en|gb|ir|nm|ut)|w(?:lf|sm)|yem|z(?:af|mb|we)|)$`)

// check that the rule's arg value is a valid country code.
func isValidCountryCode(a *analysis, f *StructField, r *Rule) error {
	for _, a := range r.Args {
		if a.Type == ArgTypeString {
			if !(len(a.Value) == 2 && rxCountryCode2.MatchString(a.Value)) &&
				!(len(a.Value) == 3 && rxCountryCode3.MatchString(a.Value)) {
				return &anError{Code: errRuleArgValueCountryCode, RuleArg: a}
			}
		} else if a.Type != ArgTypeField {
			return &anError{Code: errRuleArgType, RuleArg: a}
		}
	}
	return nil
}

// checks that the StructField's type is one of the int/uint/float types.
func typeIsNumeric(a *analysis, f *StructField, r *Rule) error {
	if !hasTypeKind(f, TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
		TypeKindUint, TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
		TypeKindFloat32, TypeKindFloat64) {
		return &anError{Code: errTypeNumeric}
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
