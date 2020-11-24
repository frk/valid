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
		// The arguments of the rule.
		Args []*RuleArg
		// The context property of the rule.
		Context string
	}

	// RuleArg represents a rule argument as parsed from a "rule" tag.
	RuleArg struct {
		// The argument value, or empty string.
		Value string
		// The type of the argument value.
		Type ArgType
	}

	// RuleType implementations indicate to the generator what code it should
	// produce for a Rule. In addition to that a RuleType implementation also
	// describes how a Rule, its Args, and its related StructField should be
	// checked for correctness before code generation.
	//
	// NOTE(mkopriva): RuleType, instead of being a direct member of
	// the Rule struct, is mapped to each Rule by the Rule's Name.
	RuleType interface {
		// Should return the expected argument count.
		argCount() ruleArgCount
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
		// check is a plugin used by the checkRule method.
		check func(a *analysis, r *Rule, t Type, f *StructField) error
		// arg count requirements, used by the argCount method.
		amin, amax int
	}

	// RuleTypeFunc is mapped to Rules that should produce code that
	// validates a value by invoking a function.
	RuleTypeFunc struct {
		// The name of the function.
		FuncName string
		// The function's package import path.
		PkgPath string
		// The types of the arguments to the function. Will always be of
		// length at least 1 where the 0th argument represents the field
		// to be validated.
		ArgTypes []Type
		// Indicates whether or not the function's signature is variadic.
		IsVariadic bool
		// NOTE(mkopriva): Although currently not enforced, this field is
		// intended to be used only with binary functions, i.e. functions
		// that take exactly two arguments.
		//
		// If set, the generator will produce one call expression for
		// each of the associated Rule's arguments and then join those
		// call expressions into a boolean expression using the *inverse*
		// of the field's logical operator value.
		LOp LogicalOperator
		// Indicates that this is a custom user provided RuleTypeFunc.
		IsCustom bool
		// check is a plugin used by the checkRule method.
		check func(a *analysis, r *Rule, t Type, f *StructField) error
		// If set, it will be returned by the argCount method. If left
		// unset, the argCount method will return a value inferred from
		// the LOp and ArgTypes fields.
		acount *ruleArgCount
		// Used for error reporting.
		typ *types.Func `cmp:"+"`
	}
)

// accepts no args
func (RuleTypeNop) argCount() ruleArgCount {
	return ruleArgCount{}
}

func (rt RuleTypeNop) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.argCount().check(len(r.Args)); !ok {
		return a.anError(&anError{Code: errRuleArgCount}, f, r)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// accepts no args
func (RuleTypeIsValid) argCount() ruleArgCount {
	return ruleArgCount{}
}

func (rt RuleTypeIsValid) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.argCount().check(len(r.Args)); !ok {
		return a.anError(&anError{Code: errRuleArgCount}, f, r)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// accepts no args
func (RuleTypeEnum) argCount() ruleArgCount {
	return ruleArgCount{}
}

// checkRule checks whether the given Type is compatible with a RuleTypeEnum,
// the *Rule and *StructField arguments are used for error reporting
func (rt RuleTypeEnum) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.argCount().check(len(r.Args)); !ok {
		return a.anError(&anError{Code: errRuleArgCount}, f, r)
	}

	typ := t.PtrBase()
	if len(typ.Name) == 0 {
		return a.anError(errRuleEnumTypeUnnamed, f, r)
	}
	if !typ.Kind.IsBasic() {
		return a.anError(errRuleEnumType, f, r)
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
		return a.anError(errRuleEnumTypeNoConst, f, r)
	}

	a.info.EnumMap[ident] = enums
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// returns a count based on RuleTypeBasic's amin & amax values
func (rt RuleTypeBasic) argCount() ruleArgCount {
	return ruleArgCount{min: rt.amin, max: rt.amax}
}

// checkRule invokes the function of RuleTypeBasic's check field, if set
func (rt RuleTypeBasic) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.argCount().check(len(r.Args)); !ok {
		return a.anError(&anError{Code: errRuleArgCount}, f, r)
	}
	if rt.check != nil {
		return rt.check(a, r, t, f)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// returns a count based on RuleTypeFunc's properties
func (rt RuleTypeFunc) argCount() ruleArgCount {
	if rt.acount != nil {
		return *rt.acount
	} else if rt.LOp > 0 {
		return ruleArgCount{1, -1}
	}

	// 1st is the field arg, which is implicit, no need to count it
	expected := len(rt.ArgTypes[1:])
	if rt.IsVariadic {
		return ruleArgCount{expected - 1, -1}
	}
	return ruleArgCount{expected, expected}
}

func (rt RuleTypeFunc) checkRule(a *analysis, r *Rule, t Type, f *StructField) error {
	if ok := rt.argCount().check(len(r.Args)); !ok {
		return a.anError(&anError{Code: errRuleArgCount}, f, r)
	}

	// field type cannot be converted to func 0th arg type, fail
	fldType, argType := t.PtrBase(), rt.ArgTypes[0]
	if rt.IsVariadic && len(rt.ArgTypes) == 1 {
		argType = *argType.Elem
	}
	if !canConvert(argType, fldType) {
		return a.anError(&anError{Code: errRuleFuncFieldType}, f, r)
	}

	// optional check returns error, fail
	if rt.check != nil {
		if err := rt.check(a, r, t, f); err != nil {
			return err
		}
	}

	// rule arg cannot be converted to func arg, fail
	fatypes := rt.ArgTypes[1:]
	if rt.IsVariadic && len(fatypes) > 0 {
		fatypes = fatypes[:len(fatypes)-1]
	}
	for i, fatyp := range fatypes {
		ra := r.Args[i]
		if !canConvertRuleArg(a, fatyp, ra) {
			return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
		}
	}
	if rt.IsVariadic {
		fatyp := rt.ArgTypes[len(rt.ArgTypes)-1]
		fatyp = *fatyp.Elem
		for _, ra := range r.Args[len(fatypes):] {
			if !canConvertRuleArg(a, fatyp, ra) {
				return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
			}
		}
	} else if rt.LOp > 0 {
		fatyp := rt.ArgTypes[1]
		for _, ra := range r.Args {
			if !canConvertRuleArg(a, fatyp, ra) {
				return a.anError(&anError{Code: errRuleFuncArgType, ra: ra}, f, r)
			}
		}
	}
	return nil
}

func (rt *RuleTypeFunc) PkgName() string {
	if len(rt.PkgPath) > 0 {
		if i := strings.LastIndexByte(rt.PkgPath, '/'); i > -1 {
			return rt.PkgPath[i+1:]
		}
		return rt.PkgPath
	}
	return ""
}

// TypesForArgs returns an adjusted version of the RuleTypeFunc's ArgTypes slice.
// The returned Type slice will match in length the given slice of RuleArgs.
func (rt *RuleTypeFunc) TypesForArgs(args []*RuleArg) (types []Type) {
	types = append(types, rt.ArgTypes[1:]...)
	if rt.IsVariadic {
		last := rt.ArgTypes[len(rt.ArgTypes)-1].Elem
		if len(types) > 0 {
			types[len(types)-1] = *last
		} else {
			types = []Type{*last}
		}

		diff := len(args) - len(types)
		for i := 0; i < diff; i++ {
			types = append(types, *last)
		}
		return types
	}

	last := rt.ArgTypes[len(rt.ArgTypes)-1]
	diff := len(args) - len(types)
	for i := 0; i < diff; i++ {
		types = append(types, last)
	}
	return types
}

////////////////////////////////////////////////////////////////////////////////

// ruleArgCount is a helper type that represents the number of arguments
// a rule can take. It is used for type checking and error reporting.
type ruleArgCount struct {
	min, max int
}

func (c ruleArgCount) check(num int) bool {
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

// Reports whether or not t represents a type that can be indexed (array/slice/map).
func (t Type) CanIndex() bool {
	return t.Kind == TypeKindArray || t.Kind == TypeKindSlice || t.Kind == TypeKindMap ||
		(t.Kind == TypeKindPtr && t.Elem.Kind == TypeKindArray)
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

func (f *StructField) SubFields() []*StructField {
	typ := f.Type.PtrBase()
	if typ.Kind == TypeKindStruct {
		return typ.Fields
	}
	return nil
}

func (s StructFieldSelector) Last() *StructField {
	return s[len(s)-1]
}

func (a *RuleArg) IsUInt() bool {
	return a.Type == ArgTypeInt && a.Value[0] != '-'
}

// ArgType indicates the type of a rule arg value.
type ArgType uint

const (
	ArgTypeUnknown ArgType = iota
	ArgTypeBool
	ArgTypeInt
	ArgTypeFloat
	ArgTypeString
	ArgTypeField
)

var argTypes = [...]string{
	ArgTypeUnknown: "<unknown>",
	ArgTypeBool:    "bool",
	ArgTypeInt:     "int",
	ArgTypeFloat:   "float",
	ArgTypeString:  "string",
	ArgTypeField:   "<field>",
}

func (t ArgType) String() string {
	if int(t) < len(argTypes) {
		return argTypes[t]
	}
	return "<invalid>"
}

// LogicalOperator represents the logical operator that, when used between
// multiple calls of a RuleType function would produce the boolean value true.
//
// NOTE(mkopriva): Because the generated code will be looking for invalid values, as
// opposed to valid ones, the actual expressions generated based on these operators
// will be the inverse of what they represent, see the comments next to the operators.
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
