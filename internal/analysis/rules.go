package analysis

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/internal/search"
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

var defaultRuleSpecMap = map[string]RuleSpec{
	// speciél
	"-isvalid": RuleNop{},
	"isvalid":  RuleIsValid{},
	"enum":     RuleEnum{},

	// basic rules
	"required": RuleBasic{},
	"notnil":   RuleBasic{check: isValidRuleNotnil},
	"rng":      RuleBasic{check: isValidRuleRng, amin: 2, amax: 2},
	"len":      RuleBasic{check: isValidRuleLen, amin: 1, amax: 2},

	"eq":  RuleBasic{check: isValidRuleValueComparison, amin: 1, amax: -1},
	"ne":  RuleBasic{check: isValidRuleValueComparison, amin: 1, amax: -1},
	"gt":  RuleBasic{check: isValidRuleNumberComparison, amin: 1, amax: 1},
	"lt":  RuleBasic{check: isValidRuleNumberComparison, amin: 1, amax: 1},
	"gte": RuleBasic{check: isValidRuleNumberComparison, amin: 1, amax: 1},
	"lte": RuleBasic{check: isValidRuleNumberComparison, amin: 1, amax: 1},
	"min": RuleBasic{check: isValidRuleNumberComparison, amin: 1, amax: 1},
	"max": RuleBasic{check: isValidRuleNumberComparison, amin: 1, amax: 1},

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
		check:    isValidRuleUUID, acount: &ruleArgCount{min: 0, max: 5}},
	"ip": RuleFunc{
		FuncName: "IP", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeIntSlice},
		check:    isValidRuleIP, acount: &ruleArgCount{min: 0, max: 2}},
	"mac": RuleFunc{
		FuncName: "MAC", PkgPath: pkgisvalid, IsVariadic: true,
		ArgTypes: []Type{typeString, typeIntSlice},
		check:    isValidRuleMAC, acount: &ruleArgCount{min: 0, max: 2}},
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
}

// LoadBuiltinFuncTypes loads the type information for builtin rule functions.
// The function should be invoked only once and before executing the first analysis.
func LoadBuiltinFuncTypes(ast search.AST) {
	for rule, spec := range defaultRuleSpecMap {
		rf, ok := spec.(RuleFunc)
		if !ok {
			continue
		}

		typ, err := search.FindFunc(rf.PkgPath, rf.FuncName, ast)
		if err != nil {
			// It is possible that the user of the cmd/isvalid tool does not
			// have github.com/frk/isvalid source on the user's machine, which
			// is ok because the source would be downloaded automatically as
			// soon as the user attempts to run the generated code, or maybe
			// the user does not intend to use the builtin rules, or perhaps
			// the user has supplied a set of custom rules that override
			// the builtin ones anyway.
			//
			// In case the error is genuine the code should keep working without
			// issues, it's just that the reporting of user errors will be poorer.
			continue
		}

		rf.typ = typ
		defaultRuleSpecMap[rule] = rf
	}
}

func isValidRuleNotnil(a *analysis, r *Rule, t Type, f *StructField) error {
	// field's type must be nilable
	if !hasTypeKind(t, TypeKindPtr, TypeKindSlice, TypeKindMap, TypeKindInterface) {
		return a.anError(errRuleFieldNonNilable, f, r)
	}
	return nil
}

// check that the rule's args are strings containing compilable regular expressions.
func isValidRuleRegexp(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, ra := range r.Args {
		if ra.Type != ArgTypeField {
			if _, err := regexp.Compile(ra.Value); err != nil {
				return a.anError(&anError{Code: errRuleArgValueRegexp, ra: ra, err: err}, f, r)
			}
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "value comparison" rule.
func isValidRuleValueComparison(a *analysis, r *Rule, t Type, f *StructField) error {
	// rule arg must be comparable to the field's type
	typ := t.PtrBase()
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return a.anError(&anError{Code: errRuleBasicArgType, ra: ra}, f, r)
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "number comparison" rule.
func isValidRuleNumberComparison(a *analysis, r *Rule, t Type, f *StructField) error {
	// the field's type must be numeric
	if err := typeIsNumeric(a, r, t, f); err != nil {
		return err
	}

	// rule arg must be comparable to the field's type
	typ := t.PtrBase()
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return a.anError(&anError{Code: errRuleBasicArgType, ra: ra}, f, r)
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "len" rule.
func isValidRuleLen(a *analysis, r *Rule, t Type, f *StructField) error {
	// associated field's type must have length
	t = t.PtrBase()
	if !hasTypeKind(t, TypeKindString, TypeKindArray, TypeKindSlice, TypeKindMap) {
		return a.anError(errRuleFieldLengthless, f, r)
	}

	// if 2, then make sure the values represent valid upper to lower bounds
	if len(r.Args) == 2 && r.Args[0].Type != ArgTypeField && r.Args[1].Type != ArgTypeField {
		var lower, upper *uint64
		if ra := r.Args[0]; len(ra.Value) > 0 {
			u64, err := strconv.ParseUint(ra.Value, 10, 64)
			if err != nil {
				return a.anError(&anError{Code: errRuleBasicArgTypeUint, ra: ra, err: err}, f, r)
			}
			lower = &u64
		}
		if ra := r.Args[1]; len(ra.Value) > 0 {
			u64, err := strconv.ParseUint(ra.Value, 10, 64)
			if err != nil {
				return a.anError(&anError{Code: errRuleBasicArgTypeUint, ra: ra, err: err}, f, r)
			}
			upper = &u64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return a.anError(&anError{Code: errRuleArgValueBounds}, f, r)
		}
	}

	// rule args must be comparable to a positive integer (the return value of the len builtin)
	typ := Type{Kind: TypeKindUint}
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return a.anError(&anError{Code: errRuleBasicArgTypeUint, ra: ra}, f, r)
		}
	}
	return nil
}

// check that the StructField and the RuleArgs represent a valid "rng" rule.
func isValidRuleRng(a *analysis, r *Rule, t Type, f *StructField) error {
	// the field's type must be numeric
	if err := typeIsNumeric(a, r, t, f); err != nil {
		return err
	}

	// make sure the rule's arg values represent valid upper to lower bounds
	if r.Args[0].Type != ArgTypeField && r.Args[1].Type != ArgTypeField {
		var lower, upper *float64
		if ra := r.Args[0]; len(ra.Value) > 0 {
			f64, err := strconv.ParseFloat(ra.Value, 64)
			if err != nil {
				return a.anError(&anError{Code: errRuleFuncArgType, ra: ra, err: err}, f, r)
			}
			lower = &f64
		}
		if ra := r.Args[1]; len(ra.Value) > 0 {
			f64, err := strconv.ParseFloat(ra.Value, 64)
			if err != nil {
				return a.anError(&anError{Code: errRuleFuncArgType, ra: ra, err: err}, f, r)
			}
			upper = &f64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return a.anError(&anError{Code: errRuleArgValueBounds}, f, r)
		}
	}

	// rule args must be comparable to the field's type
	typ := t.PtrBase()
	for _, ra := range r.Args {
		if !canConvertRuleArg(a, typ, ra) {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}

	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// check that the RuleArgs are valid UUID versions.
func isValidRuleUUID(a *analysis, r *Rule, t Type, f *StructField) error {
	versions := map[string]struct{}{} // track encountered versions
	for _, ra := range r.Args {
		if ra.Type == ArgTypeString || ra.IsUInt() {
			if !rxUUIDVer.MatchString(ra.Value) {
				return a.anError(&anError{Code: errRuleArgValueUUIDVer, ra: ra}, f, r)
			}

			if len(ra.Value) > 1 && (ra.Value[0] == 'v' || ra.Value[0] == 'V') {
				ra.Value = ra.Value[1:]
				ra.Type = ArgTypeInt
			}
			if _, exists := versions[ra.Value]; exists {
				return a.anError(&anError{Code: errRuleArgValueConflict, ra: ra}, f, r)
			} else {
				versions[ra.Value] = struct{}{}
			}
		} else if ra.Type != ArgTypeField {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	return nil
}

var rxIPVer = regexp.MustCompile(`^(?:v?(?:4|6))$`)

// checks that the rule args' values are valid IP versions.
func isValidRuleIP(a *analysis, r *Rule, t Type, f *StructField) error {
	var version string // the first version
	for _, ra := range r.Args {
		if ra.Type == ArgTypeString || ra.IsUInt() {
			if !rxIPVer.MatchString(ra.Value) {
				return a.anError(&anError{Code: errRuleArgValueIPVer, ra: ra}, f, r)
			}

			if len(ra.Value) > 1 && (ra.Value[0] == 'v' || ra.Value[0] == 'V') {
				ra.Value = ra.Value[1:]
				ra.Type = ArgTypeInt
			}
			if len(version) > 0 && version == ra.Value {
				return a.anError(&anError{Code: errRuleArgValueConflict, ra: ra}, f, r)
			} else {
				version = ra.Value
			}
		} else if ra.Type != ArgTypeField {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	return nil
}

var rxMACVer = regexp.MustCompile(`^(?:v?(?:6|8))$`)

// checks that the rule's args values are valid MAC versions.
func isValidRuleMAC(a *analysis, r *Rule, t Type, f *StructField) error {
	var version string // the first version
	for _, ra := range r.Args {
		if ra.Type == ArgTypeString || ra.IsUInt() {
			if !rxMACVer.MatchString(ra.Value) {
				return a.anError(&anError{Code: errRuleArgValueMACVer, ra: ra}, f, r)
			}

			if len(ra.Value) > 1 && (ra.Value[0] == 'v' || ra.Value[0] == 'V') {
				ra.Value = ra.Value[1:]
				ra.Type = ArgTypeInt
			}
			if len(version) > 0 && version == ra.Value {
				return a.anError(&anError{Code: errRuleArgValueConflict, ra: ra}, f, r)
			} else {
				version = ra.Value
			}
		} else if ra.Type != ArgTypeField {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	return nil
}

var rxISO = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported ISO standard identifier.
func isValidRuleISO(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, ra := range r.Args {
		if ra.IsUInt() {
			// TODO(mkopriva): Remove the regex and instead check
			// against a list of supported builtin ISO validators.
			if !rxISO.MatchString(ra.Value) {
				// TODO(mkopriva): once the above TODO item is
				// taken care of, this needs a test cases added.
				return a.anError(&anError{Code: errRuleArgValueISONum, ra: ra}, f, r)
			}
		} else if ra.Type != ArgTypeField {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	return nil
}

var rxRFC = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's arg value is a supported RFC standard identifier.
func isValidRuleRFC(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, ra := range r.Args {
		if ra.IsUInt() {
			// TODO(mkopriva): Remove the regex and instead check
			// against a list of supported builtin RFC validators.
			if !rxRFC.MatchString(ra.Value) {
				// TODO(mkopriva): once the above TODO item is
				// taken care of, this needs a test cases added.
				return a.anError(&anError{Code: errRuleArgValueRFCNum, ra: ra}, f, r)
			}
		} else if ra.Type != ArgTypeField {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	return nil
}

var rxCountryCode2 = regexp.MustCompile(`^(?i:a(?:d|e|f|g|i|l|m|o|q|r|s|t|u|w|x|z)|b(?:a|b|d|e|f|g|h|i|j|l|m|n|o|q|r|s|t|v|w|y|z)|c(?:a|c|d|f|g|h|i|k|l|m|n|o|r|u|v|w|x|y|z)|d(?:e|j|k|m|o|z)|e(?:c|e|g|h|r|s|t)|f(?:i|j|k|m|o|r)|g(?:a|b|d|e|f|g|h|i|l|m|n|p|q|r|s|t|u|w|y)|h(?:k|m|n|r|t|u)|i(?:d|e|l|m|n|o|q|r|s|t)|j(?:e|m|o|p)|k(?:e|g|h|i|m|n|p|r|w|y|z)|l(?:a|b|c|i|k|r|s|t|u|v|y)|m(?:a|c|d|e|f|g|h|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)|n(?:a|c|e|f|g|i|l|o|p|r|u|z)|om|p(?:a|e|f|g|h|k|l|m|n|r|s|t|w|y)|qa|r(?:e|o|s|u|w)|s(?:a|b|c|d|e|g|h|i|j|k|l|m|n|o|r|s|t|v|x|y|z)|t(?:c|d|f|g|h|j|k|l|m|n|o|r|t|v|w|z)|u(?:a|g|m|s|y|z)|v(?:a|c|e|g|i|n|u)|w(?:f|s)|y(?:e|t)|z(?:a|m|w))$`)
var rxCountryCode3 = regexp.MustCompile(`^(?i:a(?:bw|fg|go|ia|la|lb|nd|re|rg|rm|sm|ta|tf|tg|us|ut|ze)|b(?:di|el|en|es|fa|gd|gr|hr|hs|ih|lm|lr|lz|mu|ol|ra|rb|rn|tn|vt|wa)|c(?:af|an|ck|he|hl|hn|iv|mr|od|og|ok|ol|om|pv|ri|ub|uw|xr|ym|yp|ze)|d(?:eu|ji|ma|nk|om|za)|e(?:cu|gy|ri|sh|sp|st|th)|f(?:in|ji|lk|ra|ro|sm)|g(?:ab|br|eo|gy|ha|ib|in|lp|mb|nb|nq|rc|rd|rl|tm|uf|um|uy)|h(?:kg|md|nd|rv|ti|un)|i(?:dn|mn|nd|ot|rl|rn|rq|sl|sr|ta)|j(?:am|ey|or|pn)|k(?:az|en|gz|hm|ir|na|or|wt)|l(?:ao|bn|br|by|ca|ie|ka|so|tu|ux|va)|m(?:ac|af|ar|co|da|dg|dv|ex|hl|kd|li|lt|mr|ne|ng|np|oz|rt|sr|tq|us|wi|ys|yt)|n(?:am|cl|er|fk|ga|ic|iu|ld|or|pl|ru|zl)|omn|p(?:ak|an|cn|er|hl|lw|ng|ol|ri|rk|rt|ry|se|yf)|qat|r(?:eu|ou|us|wa)|s(?:au|dn|en|gp|gs|hn|jm|lb|le|lv|mr|om|pm|rb|sd|tp|ur|vk|vn|we|wz|xm|yc|yr)|t(?:ca|cd|go|ha|jk|kl|km|ls|on|to|un|ur|uv|wn|za)|u(?:ga|kr|mi|ry|sa|zb)|v(?:at|ct|en|gb|ir|nm|ut)|w(?:lf|sm)|yem|z(?:af|mb|we)|)$`)

// check that the rule's arg value is a valid country code.
func isValidCountryCode(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, ra := range r.Args {
		if ra.Type == ArgTypeString {
			if !(len(ra.Value) == 2 && rxCountryCode2.MatchString(ra.Value)) &&
				!(len(ra.Value) == 3 && rxCountryCode3.MatchString(ra.Value)) {
				return a.anError(&anError{Code: errRuleArgValueCountryCode, ra: ra}, f, r)
			}
		} else if ra.Type != ArgTypeField {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	return nil
}

// checks that the StructField's type is one of the int/uint/float types.
func typeIsNumeric(a *analysis, r *Rule, t Type, f *StructField) error {
	t = t.PtrBase()
	if !hasTypeKind(t, TypeKindInt, TypeKindInt8, TypeKindInt16, TypeKindInt32, TypeKindInt64,
		TypeKindUint, TypeKindUint8, TypeKindUint16, TypeKindUint32, TypeKindUint64,
		TypeKindFloat32, TypeKindFloat64) {
		return a.anError(errRuleFieldNonNumeric, f, r)
	}
	return nil
}

func hasTypeKind(t Type, kinds ...TypeKind) bool {
	for _, k := range kinds {
		if k == t.Kind {
			return true
		}
	}
	return false
}
