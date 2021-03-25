package analysis

import (
	"go/types"
	"strconv"
	"strings"

	"github.com/frk/isvalid/internal/search"
	"github.com/frk/tagutil"
)

type (
	// ValidatorStruct represents the result of the analysis of a validator struct type.
	ValidatorStruct struct {
		// Name of the validator struct type.
		TypeName string
		// The primary fields of the validator struct.
		Fields []*StructField
		// Info on the isvalid.ErrorConstructor or isvalid.ErrorAggregator
		// field of the validator struct type, or nil.
		ErrorHandler *ErrorHandlerField
		// Info on the validator type's field named "context" (case insensitive), or nil.
		ContextOption *ContextOptionField
		// Info on the validator type's method named "beforevalidate" (case insensitive), or nil.
		BeforeValidate *MethodInfo
		// Info on the validator type's method named "aftervalidate" (case insensitive), or nil.
		AfterValidate *MethodInfo
	}

	// StructField describes a single struct field in a ValidatorStruct or
	// in any of a ValidatorStruct's members that themselves are structs.
	StructField struct {
		// Name of the field.
		Name string
		// The unique key of the StructField (used for errors, field args, etc.),
		// the value of this is determined by the "field key" settings, if not
		// specified by the user it will default to the value of the field's name.
		Key string
		// The field's type.
		Type Type
		// The field's parsed tag.
		Tag tagutil.Tag
		// Indicates whether or not the field is embedded.
		IsEmbedded bool
		// Indicates whether or not the field is exported.
		IsExported bool
		// Indicates that the field, if nilable, is guaranteed to *not* be nil.
		OmitNilGuard bool
		// The field's analyzed "rule" struct tag.
		RuleTag *TagNode
	}

	// StructFieldSelector is a list of fields that represents a chain of
	// selectors where, the 0th field is the "root" field and the len-1
	// field is the "leaf" field.
	StructFieldSelector []*StructField

	// Type is the representation of a Go type.
	Type struct {
		// The name of a named type or empty string for unnamed types
		Name string
		// The kind of the go type.
		Kind TypeKind
		// The package import path.
		PkgPath string
		// The package's name.
		PkgName string
		// The local package name (including ".").
		PkgLocal string
		// Indicates whether or not the package is imported.
		IsImported bool
		// Indicates whether or not the field is exported.
		IsExported bool
		// Indicates that the type satisfies the IsValider interface.
		CanIsValid bool
		// If the base type's an array type, this field will hold the array's length.
		ArrayLen int64
		// Indicates whether or not the type is an empty interface type.
		IsEmptyInterface bool
		// Indicates whether or not the type is the "byte" alias type.
		IsByte bool
		// Indicates whether or not the type is the "rune" alias type.
		IsRune bool
		// If kind is map, key will hold the info on the map's key type.
		Key *Type
		// If kind is map, elem will hold the info on the map's value type.
		// If kind is ptr, elem will hold the info on pointed-to type.
		// If kind is slice/array, elem will hold the info on slice/array element type.
		Elem *Type
		// If kind is struct, Fields will hold the list of the struct's fields.
		Fields []*StructField
	}

	// Const represents the identifier of a declared constant.
	Const struct {
		// Name of the constant.
		Name string
		// The import path of the package to which the constant belongs.
		PkgPath string
	}

	// ErrorHandlerField is the result of analyzing a validator struct's field whose
	// type implements the isvalid.ErrorConstructor or isvalid.ErrorAggregator interface.
	ErrorHandlerField struct {
		// Name of the field (case preserved).
		Name string
		// Indicates whether or not the field's type implements
		// the isvalid.ErrorAggregator interface.
		IsAggregator bool
	}

	// ContextOptionField is the result of analyzing a validator struct's
	// field whose name is equal to "context" (case insensitive).
	ContextOptionField struct {
		// Name of the field (case preserved).
		Name string
	}

	// MethodInfo represents the result of analysing a type's method.
	MethodInfo struct {
		// The name of the method (case preserved).
		Name string
	}

	// Rule holds the basic rule information as parsed from a "rule" tag.
	Rule struct {
		// The name of the rule.
		Name string
		// The options of the rule.
		Options []*RuleOption
		// The context property of the rule.
		Context string
	}

	// RuleOption represents a rule option as parsed from a "rule" tag.
	RuleOption struct {
		// The option value, or empty string.
		Value string
		// The type of the option value.
		Type OptionType
	}

	// RuleType implementations indicate to the generator what code it should
	// produce for a Rule. In addition to that a RuleType implementation also
	// describes how a Rule, its Options, and its related StructField should be
	// checked for correctness before code generation.
	//
	// NOTE(mkopriva): RuleType, instead of being a direct member of
	// the Rule struct, is mapped to each Rule by the Rule's Name.
	RuleType interface {
		ErrConf() ErrMesgConfig
		// Should return the expected option count.
		optCount() ruleOptCount
		// Should check the given Rule for correctness.
		checkRule(a *analysis, r *Rule, t Type, f *StructField) error
	}

	// RuleTypeNop is mapped to Rules that should produce NO code.
	RuleTypeNop struct{}

	// RuleTypeIsValid is mapped to Rules that should produce code that
	// validates a value by invoking the "IsValid()" method on that value.
	RuleTypeIsValid struct{}

	// RuleTypeEnum is mapped to Rules that should produce code that validates
	// a value against a set of constants declared with that value's type.
	RuleTypeEnum struct{}

	// RuleTypeBasic is mapped to Rules that should produce code that
	// validates a value using basic expressions with comparison operators.
	RuleTypeBasic struct {
		Err ErrMesgConfig
		// check is a plugin used by the checkRule method.
		check func(a *analysis, r *Rule, t Type, f *StructField) error
		// option count requirements, used by the optCount method.
		optmin, optmax int
	}

	// RuleTypeFunc is mapped to Rules that should produce code that
	// validates a value by invoking a function.
	RuleTypeFunc struct {
		// The name of the function.
		FuncName string
		// The function's package import path.
		PkgPath string
		// The types of the function's 1st argument which will always
		// be the associate field value or a field's element value.
		FieldArgType Type
		// The types of the options to the function. Will always be of
		// length at least 1 where the 0th option represents the field
		// to be validated.
		OptionArgTypes []Type

		OptionValues []map[interface{}]*RuleOption
		// Indicates whether or not the function's signature is variadic.
		IsVariadic bool
		// NOTE(mkopriva): Although currently not enforced, this field is
		// intended to be used only with binary functions, i.e. functions
		// that take exactly two arguments, no more, no less.
		//
		// If set, the generator will produce one call expression for
		// each of the associated Rule's options and then join those
		// call expressions into a boolean expression using the *inverse*
		// of the field's logical operator value.
		LOp LogicalOperator
		// Indicates that the generated code should use raw strings
		// for any string options passed to the function.
		UseRawString bool
		// If set, will be used by the generator for producing error messages.
		Err ErrMesgConfig
		// check is a plugin used by the checkRule method.
		check func(a *analysis, r *Rule, t Type, f *StructField) error
		// If set, it will be returned by the optCount method. If left
		// unset, the optCount method will return a value inferred from
		// the LOp and OptionArgTypes fields.
		acount *ruleOptCount
		// Used for error reporting.
		typ *types.Func `cmp:"+"`
	}
)

// Accepts no options.
func (RuleTypeNop) optCount() ruleOptCount {
	return ruleOptCount{}
}

func (RuleTypeNop) ErrConf() ErrMesgConfig {
	return ErrMesgConfig{}
}

func (rt RuleTypeNop) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.optCount().check(len(r.Options)); !ok {
		return &anError{Code: errRuleOptionCount, a: a, f: f, r: r}
	}
	return nil
}

// Accepts no options.
func (RuleTypeIsValid) optCount() ruleOptCount {
	return ruleOptCount{}
}

func (RuleTypeIsValid) ErrConf() ErrMesgConfig { return ErrMesgConfig{Text: "is not valid"} }

func (rt RuleTypeIsValid) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.optCount().check(len(r.Options)); !ok {
		return &anError{Code: errRuleOptionCount, a: a, f: f, r: r}
	}
	return nil
}

// Accepts no options.
func (RuleTypeEnum) optCount() ruleOptCount {
	return ruleOptCount{}
}

func (RuleTypeEnum) ErrConf() ErrMesgConfig { return ErrMesgConfig{Text: "is not valid"} }

// checkRule checks whether the given Type is compatible with a RuleTypeEnum,
// the *Rule and *StructField arguments are used for error reporting.
func (rt RuleTypeEnum) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.optCount().check(len(r.Options)); !ok {
		return &anError{Code: errRuleOptionCount, a: a, f: f, r: r}
	}

	typ := t.PtrBase()
	if len(typ.Name) == 0 {
		return &anError{Code: errRuleEnumTypeUnnamed, a: a, f: f, r: r}
	}
	if !typ.Kind.IsBasic() {
		return &anError{Code: errRuleEnumType, a: a, f: f, r: r}
	}

	ident := typ.PkgPath + "." + typ.Name
	if _, ok := a.info.EnumMap[ident]; ok { // already done?
		return nil
	}

	enums := []Const{}
	consts := search.FindConstantsByType(typ.PkgPath, typ.Name, a.ast)
	for _, c := range consts {
		name := c.Name()
		pkgpath := c.Pkg().Path()
		// blank, skip
		if name == "_" {
			continue
		}
		// imported but not exported, skip
		if a.pkgPath != pkgpath && !c.Exported() {
			continue
		}
		enums = append(enums, Const{Name: name, PkgPath: pkgpath})
	}
	if len(enums) == 0 {
		return &anError{Code: errRuleEnumTypeNoConst, a: a, f: f, r: r}
	}

	a.info.EnumMap[ident] = enums
	return nil
}

// Returns a count based on RuleTypeBasic's optmin & optmax values.
func (rt RuleTypeBasic) optCount() ruleOptCount {
	return ruleOptCount{min: rt.optmin, max: rt.optmax}
}

func (rt RuleTypeBasic) ErrConf() ErrMesgConfig { return rt.Err }

// checkRule invokes the function of RuleTypeBasic's check field, if set.
func (rt RuleTypeBasic) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.optCount().check(len(r.Options)); !ok {
		return &anError{Code: errRuleOptionCount, a: a, f: f, r: r}
	}
	if rt.check != nil {
		return rt.check(a, r, t, f)
	}
	return nil
}

// Returns a count based on RuleTypeFunc's properties.
func (rt RuleTypeFunc) optCount() ruleOptCount {
	if rt.acount != nil {
		return *rt.acount
	} else if rt.LOp > 0 {
		return ruleOptCount{1, -1}
	}

	expected := len(rt.OptionArgTypes)
	if rt.IsVariadic {
		return ruleOptCount{expected - 1, -1}
	}
	return ruleOptCount{expected, expected}
}

func (rt RuleTypeFunc) ErrConf() ErrMesgConfig { return rt.Err }

// checkRule checks whether or not the Rule and its associated field Type can be
// used together with the RuleTypeFunc to produce code that compiles without errors.
func (rt RuleTypeFunc) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	rt.adjustRule(r)
	if ok := rt.optCount().check(len(r.Options)); !ok {
		return &anError{Code: errRuleOptionCount, a: a, f: f, r: r}
	}

	// field type cannot be converted to func arg type, fail
	fldType, argType := t.PtrBase(), rt.FieldArgType
	if rt.IsVariadic && len(rt.OptionArgTypes) == 0 {
		argType = *argType.Elem
	}
	if !canConvert(argType, fldType) {
		return &anError{Code: errRuleFuncFieldType, a: a, f: f, r: r}
	}

	// optional check returns error, fail
	if rt.check != nil {
		if err := rt.check(a, r, t, f); err != nil {
			return err
		}
	}

	// rule option cannot be converted to func arg, fail
	if len(rt.OptionArgTypes) > 0 {
		optypes := rt.OptionArgTypes[:]
		for i, optype := range optypes {
			if rt.IsVariadic && i == (len(optypes)-1) {
				// don't do the last one if it's variadic
				// will be done outside the loop
				break
			}

			opt := r.Options[i]
			if !canConvertRuleOption(a, optype, opt) {
				return &anError{Code: errRuleFuncOptionType,
					a: a, f: f, r: r, opt: opt}
			}
		}

		if rt.IsVariadic {
			optype := *(optypes[len(optypes)-1]).Elem
			for _, opt := range r.Options[len(optypes)-1:] {
				if !canConvertRuleOption(a, optype, opt) {
					return &anError{Code: errRuleFuncOptionType,
						a: a, f: f, r: r, opt: opt}
				}
			}
		} else if rt.LOp > 0 {
			optype := rt.OptionArgTypes[0]
			for _, opt := range r.Options {
				if !canConvertRuleOption(a, optype, opt) {
					return &anError{Code: errRuleFuncOptionType,
						a: a, f: f, r: r, opt: opt}
				}
			}
		}
	}
	return nil
}

// Adjusts the Rule according to the OptMap.
func (rt RuleTypeFunc) adjustRule(r *Rule) {
	for i, optmap := range rt.OptionValues {
		if len(r.Options) <= i {
			// If no option was provided for the ith argument
			// then initialize it to an "unknown" and see if
			// the map contains a default (key=nil) entry.
			opt := &RuleOption{Type: OptionTypeUnknown}
			if val, ok := optmap[nil]; ok {
				*opt = *val
			}
			r.Options = append(r.Options, opt)
			continue
		}

		opt := r.Options[i]

		// If the map contains a default (key=nil) entry and the *Rule's
		// option is "unknown", then update it with the default.
		if opt.Value == "" && opt.Type == OptionTypeUnknown {
			if val, ok := optmap[nil]; ok {
				*opt = *val
			}
			continue
		}

		// If a RuleOption's Value matches an entry in the
		// options map, then update the RuleOption.
		if val, ok := optmap[opt.Value]; ok {
			*opt = *val
		}
	}
}

// PkgName returns the name of the package to which the function belongs.
func (rt *RuleTypeFunc) PkgName() string {
	if len(rt.PkgPath) > 0 {
		if i := strings.LastIndexByte(rt.PkgPath, '/'); i > -1 {
			return rt.PkgPath[i+1:]
		}
		return rt.PkgPath
	}
	return ""
}

// TypesForOptions returns an adjusted version of the RuleTypeFunc's OptionArgTypes slice.
// The returned Type slice will match in length the given slice of RuleOptions.
func (rt *RuleTypeFunc) TypesForOptions(opts []*RuleOption) (types []Type) {
	types = append(types, rt.OptionArgTypes...)
	if rt.IsVariadic {
		last := rt.lastArgType().Elem
		if len(types) > 0 {
			types[len(types)-1] = *last
		} else {
			types = []Type{*last}
		}

		diff := len(opts) - len(types)
		for i := 0; i < diff; i++ {
			types = append(types, *last)
		}
		return types
	}

	last := rt.lastArgType()
	diff := len(opts) - len(types)
	for i := 0; i < diff; i++ {
		types = append(types, last)
	}
	return types
}

// lastArgType returns the type of the function's last argument.
func (rt *RuleTypeFunc) lastArgType() Type {
	if len(rt.OptionArgTypes) > 0 {
		return rt.OptionArgTypes[len(rt.OptionArgTypes)-1]
	}
	return rt.FieldArgType
}

// ruleOptCount is a helper type that represents the number of options
// a rule can take. It is used for type checking and error reporting.
type ruleOptCount struct {
	min, max int
}

func (c ruleOptCount) check(num int) bool {
	if num < c.min || (num > c.max && c.max != -1) {
		return false
	}
	return true
}

// ContainsRules reports whether or not the StructField f, or any of
// the StructFields in the type hierarchy of f, contain validation rules.
func (f *StructField) ContainsRules() bool {
	if f.RuleTag.ContainsRules() {
		return true
	}

	// walk recursively traverses the hierarchy of the given type and
	// invokes ContainsRules on any struct fields it encounters.
	var walk func(Type) bool
	walk = func(typ Type) bool {
		typ = typ.PtrBase()
		switch typ.Kind {
		case TypeKindStruct:
			for _, f := range typ.Fields {
				if f.ContainsRules() {
					return true
				}
			}
			return false
		case TypeKindArray, TypeKindSlice:
			return walk(*typ.Elem)
		case TypeKindMap:
			if walk(*typ.Key) {
				return true
			}
			return walk(*typ.Elem)
		}
		return false
	}
	return walk(f.Type)
}

// PtrBase returns the pointer base type of t.
func (t Type) PtrBase() Type {
	for t.Kind == TypeKindPtr {
		t = *t.Elem
	}
	return t
}

// String retruns a string representation of the t Type.
func (t Type) String() string {
	if len(t.Name) > 0 {
		if t.IsImported {
			return t.PkgName + "." + t.Name
		}
		return t.Name
	}

	if t.IsByte {
		return "byte"
	} else if t.IsRune {
		return "rune"
	} else if t.Kind.IsBasic() {
		return typeKinds[t.Kind]
	}

	switch t.Kind {
	case TypeKindArray:
		return "[" + strconv.FormatInt(t.ArrayLen, 10) + "]" + t.Elem.String()
	case TypeKindInterface:
		if !t.IsEmptyInterface {
			return "interface{ ... }"
		}
		return "interface{}"
	case TypeKindMap:
		return "map[" + t.Key.String() + "]" + t.Elem.String()
	case TypeKindPtr:
		return "*" + t.Elem.String()
	case TypeKindSlice:
		return "[]" + t.Elem.String()
	case TypeKindStruct:
		if len(t.Fields) > 0 {
			return "struct{ ... }"
		}
		return "struct{}"
	case TypeKindChan:
		return "<chan>"
	case TypeKindFunc:
		return "<func>"
	}
	return "<unknown>"
}

// Reports whether the types represented by t and u are equal. Note that this
// does not handle unnamed struct, interface (non-empty), func, and channel types.
func (t Type) Equals(u Type) bool {
	if t.Kind != u.Kind {
		return false
	}

	if len(t.Name) > 0 || len(u.Name) > 0 {
		return t.Name == u.Name && t.PkgPath == u.PkgPath
	}
	if t.Kind.IsBasic() {
		return t.Kind == u.Kind
	}

	switch t.Kind {
	case TypeKindArray:
		return t.ArrayLen == u.ArrayLen && t.Elem.Equals(*u.Elem)
	case TypeKindMap:
		return t.Key.Equals(*u.Key) && t.Elem.Equals(*u.Elem)
	case TypeKindSlice, TypeKindPtr:
		return t.Elem.Equals(*u.Elem)
	case TypeKindInterface:
		return t.IsEmptyInterface && u.IsEmptyInterface
	}
	return false
}

// Reports whether or not a value of type t needs to be converted before
// it can be assigned to a variable of type u.
func (t Type) NeedsConversion(u Type) bool {
	if u.Equals(t) {
		return false
	}
	if u.IsEmptyInterface {
		return false
	}
	return true
}

// Reports whether or not the type t represents a pointer type of u.
func (t Type) PtrOf(u Type) bool {
	return t.Kind == TypeKindPtr && t.Elem.Equals(u)
}

// Last returns the last element of s, if s has 0 elements Last will panic.
func (s StructFieldSelector) Last() *StructField {
	return s[len(s)-1]
}

// IsUInt reports whether or not RuleOption value is a valid uint candidate.
func (a *RuleOption) IsUInt() bool {
	return a.Type == OptionTypeInt && a.Value[0] != '-'
}

// OptionType indicates the type of a rule option value.
type OptionType uint

const (
	OptionTypeUnknown OptionType = iota
	OptionTypeBool
	OptionTypeInt
	OptionTypeFloat
	OptionTypeString
	OptionTypeField
)

var optTypes = [...]string{
	OptionTypeUnknown: "<unknown>",
	OptionTypeBool:    "bool",
	OptionTypeInt:     "int",
	OptionTypeFloat:   "float",
	OptionTypeString:  "string",
	OptionTypeField:   "<field>",
}

func (t OptionType) String() string {
	if int(t) < len(optTypes) {
		return optTypes[t]
	}
	return "<invalid>"
}

// LogicalOperator represents the logical operator that, when used between
// multiple calls of a RuleType function would produce the boolean value true.
//
// NOTE(mkopriva): Because the generated code will be looking for invalid values,
// as opposed to valid ones, the actual expressions generated based on these operators
// will be the inverse of what they represent, see the comments next to the operators
// for an example.
type LogicalOperator uint

const (
	_          LogicalOperator = iota
	LogicalNot                 // x || x || x....
	LogicalAnd                 // !x || !x || !x....
	LogicalOr                  // !x && !x && !x....
)

// TypeKind indicates the specific kind of a Go type.
type TypeKind uint

const (
	// basic
	TypeKindInvalid TypeKind = iota

	_basic_kind_start
	TypeKindBool
	_numeric_kind_start // int/uint/float
	_integer_kind_start // int
	TypeKindInt
	TypeKindInt8
	TypeKindInt16
	TypeKindInt32
	TypeKindInt64
	_integer_kind_end
	_unsigned_kind_start // uint
	TypeKindUint
	TypeKindUint8
	TypeKindUint16
	TypeKindUint32
	TypeKindUint64
	TypeKindUintptr
	_unsigned_kind_end
	TypeKindFloat32
	TypeKindFloat64
	_numeric_kind_end
	TypeKindComplex64
	TypeKindComplex128
	TypeKindString
	TypeKindUnsafePointer
	_basic_kind_end

	// non-basic
	TypeKindArray     // try to validate individual elements
	TypeKindInterface // try to validate ... ???
	TypeKindMap       // try to validate individual elements
	TypeKindPtr       // try to validate the element
	TypeKindSlice     // try to validate the individual elements
	TypeKindStruct    // try to validate the individual fields
	TypeKindChan      // don't validate
	TypeKindFunc      // don't validate

	// alisases (basic)
	TypeKindByte = TypeKindUint8
	TypeKindRune = TypeKindInt32
)

// Reports whether or not k is of a basic kind.
func (k TypeKind) IsBasic() bool { return _basic_kind_start < k && k < _basic_kind_end }

// Reports whether or not k is of the numeric kind, note that this
// does not include the complex64 and complex128 kinds.
func (k TypeKind) IsNumeric() bool { return _numeric_kind_start < k && k < _numeric_kind_end }

// Reports whether or not k is one of the int / uint types.
func (k TypeKind) IsInteger() bool { return _integer_kind_start < k && k < _integer_kind_end }

// Reports whether or not k is one of the uint types.
func (k TypeKind) IsUnsigned() bool { return _unsigned_kind_start < k && k < _unsigned_kind_end }

// Reports whether or not k is one of the float types.
func (k TypeKind) IsFloat() bool { return TypeKindFloat32 == k || k == TypeKindFloat64 }

// BasicString returns a string representation of k.
func (k TypeKind) BasicString() string {
	if k.IsBasic() {
		return typeKinds[k]
	}
	return "<unknown>"
}

func (k TypeKind) String() string {
	if int(k) < len(typeKinds) {
		return typeKinds[k]
	}
	return "<unknown>"
}

// Type kind string represenations indexed by typeKind.
var typeKinds = [...]string{
	TypeKindInvalid:    "<invalid>",
	TypeKindBool:       "bool",
	TypeKindInt:        "int",
	TypeKindInt8:       "int8",
	TypeKindInt16:      "int16",
	TypeKindInt32:      "int32",
	TypeKindInt64:      "int64",
	TypeKindUint:       "uint",
	TypeKindUint8:      "uint8",
	TypeKindUint16:     "uint16",
	TypeKindUint32:     "uint32",
	TypeKindUint64:     "uint64",
	TypeKindUintptr:    "uintptr",
	TypeKindFloat32:    "float32",
	TypeKindFloat64:    "float64",
	TypeKindComplex64:  "complex64",
	TypeKindComplex128: "complex128",
	TypeKindString:     "string",

	// ...
	TypeKindArray:     "array",
	TypeKindInterface: "interface",
	TypeKindMap:       "map",
	TypeKindPtr:       "ptr",
	TypeKindSlice:     "slice",
	TypeKindStruct:    "struct",
	TypeKindChan:      "chan",
	TypeKindFunc:      "func",
}

// typeKinds indexed by types.BasicKind.
var typesBasicKindToTypeKind = [...]TypeKind{
	types.Invalid:       TypeKindInvalid,
	types.Bool:          TypeKindBool,
	types.Int:           TypeKindInt,
	types.Int8:          TypeKindInt8,
	types.Int16:         TypeKindInt16,
	types.Int32:         TypeKindInt32,
	types.Int64:         TypeKindInt64,
	types.Uint:          TypeKindUint,
	types.Uint8:         TypeKindUint8,
	types.Uint16:        TypeKindUint16,
	types.Uint32:        TypeKindUint32,
	types.Uint64:        TypeKindUint64,
	types.Uintptr:       TypeKindUintptr,
	types.Float32:       TypeKindFloat32,
	types.Float64:       TypeKindFloat64,
	types.Complex64:     TypeKindComplex64,
	types.Complex128:    TypeKindComplex128,
	types.String:        TypeKindString,
	types.UnsafePointer: TypeKindUnsafePointer,
}
