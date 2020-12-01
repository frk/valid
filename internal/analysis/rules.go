package analysis

import (
	"go/types"
	//"log"
	"encoding/json"
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

var defaultRuleTypeMap = map[string]RuleType{
	// basic comparison rules
	"eq": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be equal to", OptSep: " or ", WithOpts: true},
		check: isValidRuleValueComparison, optmin: 1, optmax: -1,
	},
	"ne": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must not be equal to", OptSep: " or ", WithOpts: true},
		check: isValidRuleValueComparison, optmin: 1, optmax: -1,
	},
	"gt": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be greater than", WithOpts: true},
		check: isValidRuleNumberComparison, optmin: 1, optmax: 1,
	},
	"lt": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be less than", WithOpts: true},
		check: isValidRuleNumberComparison, optmin: 1, optmax: 1,
	},
	"gte": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be greater than or equal to", WithOpts: true},
		check: isValidRuleNumberComparison, optmin: 1, optmax: 1,
	},
	"lte": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be less than or equal to", WithOpts: true},
		check: isValidRuleNumberComparison, optmin: 1, optmax: 1,
	},
	"min": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be greater than or equal to", WithOpts: true},
		check: isValidRuleNumberComparison, optmin: 1, optmax: 1,
	},
	"max": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be less than or equal to", WithOpts: true},
		check: isValidRuleNumberComparison, optmin: 1, optmax: 1,
	},

	// basic but speciél rules
	"required": RuleTypeBasic{
		Err: ErrMesgConfig{Text: "is required"},
	},
	"notnil": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "cannot be nil"},
		check: isValidRuleNotnil,
	},
	"rng": RuleTypeBasic{
		Err:   ErrMesgConfig{Text: "must be between", OptSep: " and ", WithOpts: true},
		check: isValidRuleRng, optmin: 2, optmax: 2,
	},
	"len":       RuleTypeBasic{check: isValidRuleLen, optmin: 1, optmax: 2},
	"runecount": RuleTypeBasic{check: isValidRuleRuneCount, optmin: 1, optmax: 2},

	// speciél
	"-isvalid": RuleTypeNop{},
	"isvalid":  RuleTypeIsValid{},
	"enum":     RuleTypeEnum{},

	// validators "borrowed" from stdlib
	"prefix": RuleTypeFunc{
		FuncName: "HasPrefix", PkgPath: "strings",
		FieldArgType:   typeString,
		OptionArgTypes: []Type{typeString},
		Err: ErrMesgConfig{
			Text:     "must be prefixed with",
			OptSep:   " or ",
			WithOpts: true,
		},
		LOp: LogicalOr,
	},
	"suffix": RuleTypeFunc{
		FuncName: "HasSuffix", PkgPath: "strings",
		FieldArgType:   typeString,
		OptionArgTypes: []Type{typeString},
		Err: ErrMesgConfig{
			Text:     "must be suffixed with",
			OptSep:   " or ",
			WithOpts: true,
		},
		LOp: LogicalOr,
	},
	"contains": RuleTypeFunc{
		FuncName: "Contains", PkgPath: "strings",
		FieldArgType:   typeString,
		OptionArgTypes: []Type{typeString},
		Err: ErrMesgConfig{
			Text:     "must contain substring",
			OptSep:   " or ",
			WithOpts: true,
		},
		LOp: LogicalOr,
	},
}

type RuleConfig struct {
	Name   string `json:"name"`
	OptMin *int   `json:"opt_min"`
	OptMax *int   `json:"opt_max"`
	OptMap []struct {
		Key   *string `json:"key"`
		Value string  `json:"value"`
	} `json:"opt_map"`
	Err ErrMesgConfig   `json:"err"`
	LOp LogicalOperator `json:"log_op"`
}

type ErrMesgConfig struct {
	Text     string `json:"text"`
	OptSep   string `json:"opt_sep"`
	WithOpts bool   `json:"with_opts"`
}

func (conf RuleConfig) RuleTypeFunc(fn *types.Func, isCustom bool) (RuleTypeFunc, error) {
	sig := fn.Type().(*types.Signature)
	p, r := sig.Params(), sig.Results()
	if p.Len() < 1 || r.Len() != 1 {
		return RuleTypeFunc{}, &anError{Code: errRuleFuncSignature, fn: fn}
	}
	if !isBool(r.At(0).Type()) {
		return RuleTypeFunc{}, &anError{Code: errRuleFuncSignature, fn: fn}
	}

	rt := RuleTypeFunc{}
	rt.FuncName = fn.Name()
	rt.PkgPath = fn.Pkg().Path()
	rt.IsVariadic = sig.Variadic()
	rt.FieldArgType = analyzeType0(p.At(0).Type())
	for i := 1; i < p.Len(); i++ {
		rt.OptionArgTypes = append(rt.OptionArgTypes, analyzeType0(p.At(i).Type()))
	}
	rt.Err = conf.Err
	rt.LOp = conf.LOp
	rt.typ = fn

	if conf.OptMin != nil || conf.OptMax != nil {
		rt.acount = new(ruleOptCount)
		if conf.OptMin != nil {
			rt.acount.min = *conf.OptMin
		}
		if conf.OptMax != nil {
			rt.acount.max = *conf.OptMax
		}
	}

	if len(conf.OptMap) > 0 {
		rt.OptionValues = make(map[interface{}]*RuleOption)
		for _, kv := range conf.OptMap {
			opt := parseRuleTagOption(kv.Value)
			if kv.Key != nil {
				rt.OptionValues[*kv.Key] = opt
			} else {
				rt.OptionValues[nil] = opt
			}
		}
	}

	if !isCustom {
		// some builtin function need additional help for type checking
		switch conf.Name {
		case "alpha":
			rt.check = isValidLanguageTag
		case "alnum":
			rt.check = isValidLanguageTag
		case "phone":
			rt.check = isValidCountryCode
		case "zip":
			rt.check = isValidCountryCode
		case "uuid":
			rt.check = isValidRuleUUID
		case "ip":
			rt.check = isValidRuleIP
		case "mac":
			rt.check = isValidRuleMAC
		case "iso":
			rt.check = isValidRuleISO
		case "rfc":
			rt.check = isValidRuleRFC
		case "re":
			rt.UseRawString = true
			rt.check = isValidRuleRegexp
		}
	}

	return rt, nil
}

// LoadRuleTypeFunc loads info for pre-defined function rule types. LoadRuleTypeFunc
// should be invoked only once and before starting the first analysis.
func LoadRuleTypeFunc(ast search.AST) {
	// load functions from the "github.com/frk/isvalid" package
	search.LoadBuiltinFuncs(ast, func(confjson []byte, fn *types.Func) error {
		conf := RuleConfig{}
		if err := json.Unmarshal(confjson, &conf); err != nil {
			panic("bad json for RuleConfig:" + err.Error() + "\n" + string(confjson))
		}

		rt, err := conf.RuleTypeFunc(fn, false)
		if err != nil {
			panic(err.Error())
		}
		defaultRuleTypeMap[conf.Name] = rt
		return nil
	})
}

func isValidRuleNotnil(a *analysis, r *Rule, t Type, f *StructField) error {
	// field's type must be nilable
	if !hasTypeKind(t, TypeKindPtr, TypeKindSlice, TypeKindMap, TypeKindInterface) {
		return &anError{Code: errRuleFieldNonNilable, a: a, f: f, r: r}
	}
	return nil
}

// check that the rule's options are strings containing compilable regular expressions.
func isValidRuleRegexp(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.Type != OptionTypeField {
			if _, err := regexp.Compile(opt.Value); err != nil {
				return &anError{Code: errRuleOptionValueRegexp, a: a,
					f: f, r: r, opt: opt, err: err}
			}
		}
	}
	return nil
}

// check that the StructField and the RuleOptions represent a valid "value comparison" rule.
func isValidRuleValueComparison(a *analysis, r *Rule, t Type, f *StructField) error {
	// rule option must be comparable to the field's type
	typ := t.PtrBase()
	for _, opt := range r.Options {
		if !canConvertRuleOption(a, typ, opt) {
			return &anError{Code: errRuleBasicOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// check that the StructField and the RuleOptions represent a valid "number comparison" rule.
func isValidRuleNumberComparison(a *analysis, r *Rule, t Type, f *StructField) error {
	// the field's type must be numeric
	if err := typeIsNumeric(a, r, t, f); err != nil {
		return err
	}

	// rule option must be comparable to the field's type
	typ := t.PtrBase()
	for _, opt := range r.Options {
		if !canConvertRuleOption(a, typ, opt) {
			return &anError{Code: errRuleBasicOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// check that the StructField and the RuleOptions represent a valid "len" rule.
func isValidRuleLen(a *analysis, r *Rule, t Type, f *StructField) error {
	// associated field's type must have length
	t = t.PtrBase()
	if !hasTypeKind(t, TypeKindString, TypeKindArray, TypeKindSlice, TypeKindMap) {
		return &anError{Code: errRuleFieldLengthless, a: a, f: f, r: r}
	}

	// if 2, then make sure the values represent valid upper to lower bounds
	if len(r.Options) == 2 && r.Options[0].Type != OptionTypeField && r.Options[1].Type != OptionTypeField {
		var lower, upper *uint64
		if opt := r.Options[0]; len(opt.Value) > 0 {
			u64, err := strconv.ParseUint(opt.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleBasicOptionTypeUint, a: a, f: f, r: r, opt: opt, err: err}
			}
			lower = &u64
		}
		if opt := r.Options[1]; len(opt.Value) > 0 {
			u64, err := strconv.ParseUint(opt.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleBasicOptionTypeUint, a: a, f: f, r: r, opt: opt, err: err}
			}
			upper = &u64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return &anError{Code: errRuleOptionValueBounds, a: a, f: f, r: r}
		}
	}

	// rule options must be comparable to a positive integer (the return value of the len builtin)
	typ := Type{Kind: TypeKindUint}
	for _, opt := range r.Options {
		if !canConvertRuleOption(a, typ, opt) {
			return &anError{Code: errRuleBasicOptionTypeUint, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// check that the StructField and the RuleOptions represent a valid "runecount" rule.
func isValidRuleRuneCount(a *analysis, r *Rule, t Type, f *StructField) error {
	// associated field's type must string kind or byte slice
	t = t.PtrBase()
	if t.Kind != TypeKindString && (t.Kind != TypeKindSlice || !t.Elem.IsByte) {
		return &anError{Code: errRuleFieldRuneless,
			a: a, f: f, r: r}
	}

	// if 2, then make sure the values represent valid upper to lower bounds
	if len(r.Options) == 2 && r.Options[0].Type != OptionTypeField && r.Options[1].Type != OptionTypeField {
		var lower, upper *uint64
		if opt := r.Options[0]; len(opt.Value) > 0 {
			u64, err := strconv.ParseUint(opt.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleBasicOptionTypeUint,
					a: a, f: f, r: r, opt: opt, err: err}
			}
			lower = &u64
		}
		if opt := r.Options[1]; len(opt.Value) > 0 {
			u64, err := strconv.ParseUint(opt.Value, 10, 64)
			if err != nil {
				return &anError{Code: errRuleBasicOptionTypeUint,
					a: a, f: f, r: r, opt: opt, err: err}
			}
			upper = &u64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return &anError{Code: errRuleOptionValueBounds, a: a, f: f, r: r}
		}
	}

	// rule options must be comparable to a positive integer (the return value of the utf8.RuneCount function)
	typ := Type{Kind: TypeKindUint}
	for _, opt := range r.Options {
		if !canConvertRuleOption(a, typ, opt) {
			return &anError{Code: errRuleBasicOptionTypeUint,
				a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// check that the StructField and the RuleOptions represent a valid "rng" rule.
func isValidRuleRng(a *analysis, r *Rule, t Type, f *StructField) error {
	// the field's type must be numeric
	if err := typeIsNumeric(a, r, t, f); err != nil {
		return err
	}

	// make sure the rule's option values represent valid upper to lower bounds
	if r.Options[0].Type != OptionTypeField && r.Options[1].Type != OptionTypeField {
		var lower, upper *float64
		if opt := r.Options[0]; len(opt.Value) > 0 {
			f64, err := strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt, err: err}
			}
			lower = &f64
		}
		if opt := r.Options[1]; len(opt.Value) > 0 {
			f64, err := strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt, err: err}
			}
			upper = &f64
		}

		// make sure at least one bound was specified and if both then
		// ensure that the lower bound is less than the upper bound
		if (lower == nil && upper == nil) || (lower != nil && upper != nil && *lower >= *upper) {
			return &anError{Code: errRuleOptionValueBounds, a: a, f: f, r: r}
		}
	}

	// rule options must be comparable to the field's type
	typ := t.PtrBase()
	for _, opt := range r.Options {
		if !canConvertRuleOption(a, typ, opt) {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}

	return nil
}

var rxUUIDVer = regexp.MustCompile(`^(?:v?[1-5])$`)

// check that the RuleOptions are valid UUID versions.
func isValidRuleUUID(a *analysis, r *Rule, t Type, f *StructField) error {
	versions := map[string]struct{}{} // track encountered versions
	for _, opt := range r.Options {
		if opt.Type == OptionTypeString || opt.IsUInt() {
			if !rxUUIDVer.MatchString(opt.Value) {
				return &anError{Code: errRuleOptionValueUUIDVer, a: a, f: f, r: r, opt: opt}
			}

			if len(opt.Value) > 1 && (opt.Value[0] == 'v' || opt.Value[0] == 'V') {
				opt.Value = opt.Value[1:]
				opt.Type = OptionTypeInt
			}
			if _, exists := versions[opt.Value]; exists {
				return &anError{Code: errRuleOptionValueConflict, a: a, f: f, r: r, opt: opt}
			} else {
				versions[opt.Value] = struct{}{}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// checks that the rule options' values are valid IP versions.
func isValidRuleIP(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.IsUInt() {
			if opt.Value != "0" && opt.Value != "4" && opt.Value != "6" {
				return &anError{Code: errRuleOptionValueIPVer, a: a, f: f, r: r, opt: opt}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// checks that the rule's options values are valid MAC versions.
func isValidRuleMAC(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.IsUInt() {
			if opt.Value != "0" && opt.Value != "6" && opt.Value != "8" {
				return &anError{Code: errRuleOptionValueMACVer, a: a, f: f, r: r, opt: opt}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

var rxISO = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's option value is a supported ISO standard identifier.
func isValidRuleISO(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.IsUInt() {
			// TODO(mkopriva): Remove the regex and instead check
			// against a list of supported builtin ISO validators.
			if !rxISO.MatchString(opt.Value) {
				// TODO(mkopriva): once the above TODO item is
				// taken care of, this needs a test cases added.
				return &anError{Code: errRuleOptionValueISONum, a: a, f: f, r: r, opt: opt}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

var rxRFC = regexp.MustCompile(`^(?:[1-9][0-9]*)$`) // non-zero unsigned int

// checks that the rule's option value is a supported RFC standard identifier.
func isValidRuleRFC(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.IsUInt() {
			// TODO(mkopriva): Remove the regex and instead check
			// against a list of supported builtin RFC validators.
			if !rxRFC.MatchString(opt.Value) {
				// TODO(mkopriva): once the above TODO item is
				// taken care of, this needs a test cases added.
				return &anError{Code: errRuleOptionValueRFCNum, a: a, f: f, r: r, opt: opt}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

var rxCountryCode2 = regexp.MustCompile(`^(?i:a(?:d|e|f|g|i|l|m|o|q|r|s|t|u|w|x|z)|b(?:a|b|d|e|f|g|h|i|j|l|m|n|o|q|r|s|t|v|w|y|z)|c(?:a|c|d|f|g|h|i|k|l|m|n|o|r|u|v|w|x|y|z)|d(?:e|j|k|m|o|z)|e(?:c|e|g|h|r|s|t)|f(?:i|j|k|m|o|r)|g(?:a|b|d|e|f|g|h|i|l|m|n|p|q|r|s|t|u|w|y)|h(?:k|m|n|r|t|u)|i(?:d|e|l|m|n|o|q|r|s|t)|j(?:e|m|o|p)|k(?:e|g|h|i|m|n|p|r|w|y|z)|l(?:a|b|c|i|k|r|s|t|u|v|y)|m(?:a|c|d|e|f|g|h|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)|n(?:a|c|e|f|g|i|l|o|p|r|u|z)|om|p(?:a|e|f|g|h|k|l|m|n|r|s|t|w|y)|qa|r(?:e|o|s|u|w)|s(?:a|b|c|d|e|g|h|i|j|k|l|m|n|o|r|s|t|v|x|y|z)|t(?:c|d|f|g|h|j|k|l|m|n|o|r|t|v|w|z)|u(?:a|g|m|s|y|z)|v(?:a|c|e|g|i|n|u)|w(?:f|s)|y(?:e|t)|z(?:a|m|w))$`)
var rxCountryCode3 = regexp.MustCompile(`^(?i:a(?:bw|fg|go|ia|la|lb|nd|re|rg|rm|sm|ta|tf|tg|us|ut|ze)|b(?:di|el|en|es|fa|gd|gr|hr|hs|ih|lm|lr|lz|mu|ol|ra|rb|rn|tn|vt|wa)|c(?:af|an|ck|he|hl|hn|iv|mr|od|og|ok|ol|om|pv|ri|ub|uw|xr|ym|yp|ze)|d(?:eu|ji|ma|nk|om|za)|e(?:cu|gy|ri|sh|sp|st|th)|f(?:in|ji|lk|ra|ro|sm)|g(?:ab|br|eo|gy|ha|ib|in|lp|mb|nb|nq|rc|rd|rl|tm|uf|um|uy)|h(?:kg|md|nd|rv|ti|un)|i(?:dn|mn|nd|ot|rl|rn|rq|sl|sr|ta)|j(?:am|ey|or|pn)|k(?:az|en|gz|hm|ir|na|or|wt)|l(?:ao|bn|br|by|ca|ie|ka|so|tu|ux|va)|m(?:ac|af|ar|co|da|dg|dv|ex|hl|kd|li|lt|mr|ne|ng|np|oz|rt|sr|tq|us|wi|ys|yt)|n(?:am|cl|er|fk|ga|ic|iu|ld|or|pl|ru|zl)|omn|p(?:ak|an|cn|er|hl|lw|ng|ol|ri|rk|rt|ry|se|yf)|qat|r(?:eu|ou|us|wa)|s(?:au|dn|en|gp|gs|hn|jm|lb|le|lv|mr|om|pm|rb|sd|tp|ur|vk|vn|we|wz|xm|yc|yr)|t(?:ca|cd|go|ha|jk|kl|km|ls|on|to|un|ur|uv|wn|za)|u(?:ga|kr|mi|ry|sa|zb)|v(?:at|ct|en|gb|ir|nm|ut)|w(?:lf|sm)|yem|z(?:af|mb|we)|)$`)

// check that the rule's option value is a valid country code.
func isValidCountryCode(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.Type == OptionTypeString {
			if !(len(opt.Value) == 2 && rxCountryCode2.MatchString(opt.Value)) &&
				!(len(opt.Value) == 3 && rxCountryCode3.MatchString(opt.Value)) {
				return &anError{Code: errRuleOptionValueCountryCode, a: a, f: f, r: r, opt: opt}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
		}
	}
	return nil
}

// supported language tags
// TODO(mkopriva): expand!
var rxLanguageTag = regexp.MustCompile(`^(?i:be|bg|cnr|cs|en|mk|pl|ru|sh|sk|sl|sr|uk|wen)$`)

// check that the rule's option value is one of the supported language tags.
func isValidLanguageTag(a *analysis, r *Rule, t Type, f *StructField) error {
	for _, opt := range r.Options {
		if opt.Type == OptionTypeString {
			if !rxLanguageTag.MatchString(opt.Value) {
				return &anError{Code: errRuleOptionValueLanguageTag, a: a, f: f, r: r, opt: opt}
			}
		} else if opt.Type != OptionTypeField {
			return &anError{Code: errRuleFuncOptionType, a: a, f: f, r: r, opt: opt}
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
		return &anError{Code: errRuleFieldNonNumeric, a: a, f: f, r: r}
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
