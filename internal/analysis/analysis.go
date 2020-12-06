package analysis

import (
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/frk/isvalid/internal/search"
	"github.com/frk/tagutil"
)

// A Config specifies the configuration for the analysis.
type Config struct {
	// If set, specifies the struct tag to be used to produce a struct
	// field's key. If not set, the field's name will be used by default.
	FieldKeyTag string
	// If set to true, a nested struct field's key will be produced by
	// joining it together with all of its parent fields. If set to false,
	// such a field's key will be produced only from that field's name/tag.
	FieldKeyJoin bool
	// If set, specifies the separator that will be used when producing a
	// joined key for a nested struct field. If not set, the field's key
	// will be joined without a separator.
	//
	// This field is only used if FieldKeyJoin is set to true.
	FieldKeySeparator string

	// map of custom RuleTypes
	customTypeMap map[string]RuleType
}

// AddRuleFunc is used to register a custom RuleFunc with the Config. The
// custom function MUST have at least one parameter item and, it MUST have
// exactly one result item which MUST be of type bool.
func (c *Config) AddRuleFunc(ruleName string, typ *types.Func) error {
	if name := strings.ToLower(ruleName); name == "isvalid" || name == "-isvalid" || name == "enum" {
		return &anError{Code: errRuleNameReserved, r: &Rule{Name: ruleName}}
	}

	conf := RuleConfig{Name: ruleName}
	rt, err := conf.RuleTypeFunc(typ, true)
	if err != nil {
		return err
	}

	if c.customTypeMap == nil {
		c.customTypeMap = make(map[string]RuleType)
	}
	c.customTypeMap[ruleName] = rt
	return nil
}

// Analyze runs the analysis of the validator struct represented by the given *search.Match.
// If successful, the returned *ValidatorStruct value is ready to be fed to the generator.
func (c Config) Analyze(ast search.AST, match *search.Match, info *Info) (*ValidatorStruct, error) {
	structType, ok := match.Named.Underlying().(*types.Struct)
	if !ok {
		panic(match.Named.Obj().Name() + " must be a struct type.") // this shouldn't happen
	}

	a := new(analysis)
	a.conf = c
	a.ast = ast
	a.fset = match.Fset
	a.named = match.Named
	a.pkgPath = match.Named.Obj().Pkg().Path()
	a.keys = make(map[string]uint)
	a.fieldVarMap = make(map[*StructField]fieldVar)

	a.info = info
	a.info.FileSet = match.Fset
	a.info.PkgPath = a.pkgPath
	a.info.TypeName = match.Named.Obj().Name()
	a.info.TypeNamePos = match.Pos
	a.info.SelectorMap = make(map[string]StructFieldSelector)
	a.info.EnumMap = make(map[string][]Const)

	a.fieldKey = fieldKeyFunc(c)
	vs, err := analyzeValidatorStruct(a, structType)
	if err != nil {
		return nil, err
	}

	// merge the rule func maps into one for the generator to use
	a.info.RuleTypeMap = make(map[string]RuleType)
	for k, v := range defaultRuleTypeMap {
		a.info.RuleTypeMap[k] = v
	}
	for k, v := range c.customTypeMap {
		a.info.RuleTypeMap[k] = v
	}

	return vs, nil
}

// Info holds information related to the analysis and inteded to be used by the generator.
// If the analysis returns an error, the collected information will be incomplete.
type Info struct {
	// The FileSet associated with the analyzed ValidatorStruct.
	FileSet *token.FileSet
	// The package path of the analyzed ValidatorStruct.
	PkgPath string
	// The type name of the analyzed ValidatorStruct.
	TypeName string
	// The soruce position of the ValidatorStruct's type name.
	TypeNamePos token.Pos
	// RuleTypeMap will be populated by all the registered RuleTypes.
	RuleTypeMap map[string]RuleType
	// SelectorMap maps field keys to their related field selectors.
	SelectorMap map[string]StructFieldSelector
	// EnumMap maps package-path qualified type names to a slice of
	// constants declared with that type.
	EnumMap map[string][]Const
}

// analysis holds the state of the analyzer.
type analysis struct {
	conf Config
	// The AST as populated by search.Search.
	ast search.AST
	// The FileSet associated with the type under analysis,
	// used primarily for error reporting.
	fset *token.FileSet
	// The named type under analysis.
	named *types.Named
	// The package path of the type under analysis.
	pkgPath string
	// This field will hold the result of the analysis.
	validator *ValidatorStruct
	// Tracks already created field keys to ensure uniqueness.
	keys map[string]uint
	// Holds useful information aggregated during analysis.
	info *Info
	// Constructs a field key for the given selector, initialized from Config.
	fieldKey func([]*StructField) (key string)
	// For error reporting. If not nil it will hold the last encountered
	// rule & field that need the ValidatorStruct to have a "context" field.
	needsContext *needsContext
	// fieldVarMap maintains a map of StructField pointers to the fields'
	// related go/types specific information. Intended for error reporting.
	fieldVarMap map[*StructField]fieldVar
}

// used for error reporting only
type needsContext struct {
	field *StructField
	rule  *Rule
}

// used for error reporting only
type fieldVar struct {
	v   *types.Var
	tag string
}

// used by tests to trim away developer specific file system location of
// the project from testdata files' filepaths.
var filenamehook = func(name string) string { return name }

// analyzeValidatorStruct is the entry point of the analysis of a ValidatorStruct type.
func analyzeValidatorStruct(a *analysis, structType *types.Struct) (*ValidatorStruct, error) {
	a.validator = new(ValidatorStruct)
	a.validator.TypeName = a.named.Obj().Name()
	if name := lookupBeforeValidate(a.named); len(name) > 0 {
		a.validator.BeforeValidate = &MethodInfo{Name: name}
	}
	if name := lookupAfterValidate(a.named); len(name) > 0 {
		a.validator.AfterValidate = &MethodInfo{Name: name}
	}

	typName := strings.ToLower(a.validator.TypeName)
	if !strings.HasSuffix(typName, "validator") {
		panic(a.validator.TypeName + " struct type has unsupported name suffix.") // this shouldn't happen
	}

	// 1. analyze all fields
	fields, err := analyzeStructFields(a, structType, nil, true)
	if err != nil {
		return nil, err
	} else if len(fields) == 0 {
		return nil, &anError{Code: errValidatorNoField, a: a}
	}

	// 2. type-check all of the fields' rules
	if err := typeCheckRules(a, fields); err != nil {
		return nil, err
	}

	// 3. ensure that if a rule with context exists, that also a ContextOptionField exists
	if a.needsContext != nil && a.validator.ContextOption == nil {
		return nil, &anError{Code: errContextOptionFieldRequired, a: a,
			f: a.needsContext.field, r: a.needsContext.rule}
	}

	a.validator.Fields = fields
	return a.validator, nil
}

// analyzeStructFields analyzes the given *types.Struct's fields.
func analyzeStructFields(a *analysis, structType *types.Struct, selector []*StructField, local bool) (fields []*StructField, err error) {
	for i := 0; i < structType.NumFields(); i++ {
		fvar := structType.Field(i)
		ftag := structType.Tag(i)
		tag := tagutil.New(ftag)
		istag := tag.First("is")

		// Skip imported but unexported fields.
		if !local && !fvar.Exported() {
			continue
		}

		// Skip fields with blank name.
		if fvar.Name() == "_" {
			continue
		}

		// Skip fields that were explicitly flagged; fields with
		// no `is` tag may still be useful if they are referenced by
		// a separate field's rule with a field option.
		if istag == "-" {
			continue
		}

		f := new(StructField)
		f.Tag = tag
		f.Name = fvar.Name()
		f.IsEmbedded = fvar.Embedded()
		f.IsExported = fvar.Exported()
		f.RuleTag, _ = parseRuleTag(ftag)

		// map field to fvar for error reporting
		a.fieldVarMap[f] = fieldVar{v: fvar, tag: ftag}

		// resolve field key for selector & make sure that it is unique
		fsel := append(selector, f)
		f.Key = makeFieldKey(a, fsel)
		if _, ok := a.keys[f.Key]; ok {
			// NOTE(mkopriva): this shouldn't happen given that
			// makeFieldKey already checks for duplicates and if
			// one is found then it modifies the key so as to make
			// it unique, regardless this stays here just in case
			// an update to makeFieldKey, or something else, breaks
			// that expected behaviour.
			panic("shouldn't reach")
			return nil, nil
		} else {
			a.keys[f.Key] = 1
		}

		// map keys to selectors
		a.info.SelectorMap[f.Key] = append(StructFieldSelector{}, fsel...)

		typ, err := analyzeType(a, fvar.Type(), fsel)
		if err != nil {
			return nil, err
		}
		f.Type = typ

		// Check for untagged, "special" root fields.
		if len(istag) == 0 && len(selector) == 0 {
			if isErrorConstructor(fvar.Type()) {
				if err := analyzeErrorHandlerField(a, f, false); err != nil {
					return nil, err
				}
				continue
			} else if isErrorAggregator(fvar.Type()) {
				if err := analyzeErrorHandlerField(a, f, true); err != nil {
					return nil, err
				}
				continue
			} else if strings.ToLower(fvar.Name()) == "context" {
				if err := analyzeContextOptionField(a, f); err != nil {
					return nil, err
				}
				continue
			}
		}

		fields = append(fields, f)
	}
	return fields, nil
}

// analyzeErrorHandlerField analyzes the given field as a ErrorHandlerField.
// The field's type is known to implement either the ErrorConstructor or the
// ErrorAggregator interface.
func analyzeErrorHandlerField(a *analysis, f *StructField, isAggregator bool) error {
	if a.validator.ErrorHandler != nil {
		return &anError{Code: errErrorHandlerFieldConflict, a: a, f: f}
	}

	a.validator.ErrorHandler = new(ErrorHandlerField)
	a.validator.ErrorHandler.Name = f.Name
	a.validator.ErrorHandler.IsAggregator = isAggregator
	return nil
}

// analyzeContextOptionField analyzes the given field as a ContextOptionField.
// The field's name is known to be "context" (case insensitive).
func analyzeContextOptionField(a *analysis, f *StructField) error {
	if a.validator.ContextOption != nil {
		return &anError{Code: errContextOptionFieldConflict, a: a, f: f}
	}
	if f.Type.Kind != TypeKindString {
		return &anError{Code: errContextOptionFieldType, a: a, f: f}
	}

	a.validator.ContextOption = new(ContextOptionField)
	a.validator.ContextOption.Name = f.Name
	return nil
}

// analyzeType analyzes the given types.Type.
func analyzeType(a *analysis, t types.Type, selector []*StructField) (typ Type, err error) {
	if named, ok := t.(*types.Named); ok {
		pkg := named.Obj().Pkg()
		typ.Name = named.Obj().Name()
		typ.PkgPath = pkg.Path()
		typ.PkgName = pkg.Name()
		typ.PkgLocal = pkg.Name()
		typ.IsImported = isImportedType(a, named)
		typ.IsExported = named.Obj().Exported()
		typ.CanIsValid = canIsValid(t)
		t = named.Underlying()
	}

	typ.Kind = analyzeTypeKind(t)

	switch T := t.(type) {
	case *types.Basic:
		typ.IsRune = T.Name() == "rune"
		typ.IsByte = T.Name() == "byte"
	case *types.Slice:
		elem, err := analyzeType(a, T.Elem(), selector)
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
	case *types.Array:
		elem, err := analyzeType(a, T.Elem(), selector)
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
		typ.ArrayLen = T.Len()
	case *types.Map:
		key, err := analyzeType(a, T.Key(), selector)
		if err != nil {
			return Type{}, err
		}
		elem, err := analyzeType(a, T.Elem(), selector)
		if err != nil {
			return Type{}, err
		}
		typ.Key = &key
		typ.Elem = &elem
	case *types.Pointer:
		elem, err := analyzeType(a, T.Elem(), selector)
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
	case *types.Interface:
		typ.IsEmptyInterface = T.NumMethods() == 0
		typ.CanIsValid = canIsValid(t)
	case *types.Struct:
		fields, err := analyzeStructFields(a, T, selector, !typ.IsImported)
		if err != nil {
			return Type{}, err
		}
		typ.Fields = fields
	}

	return typ, nil
}

// a simplified version of the above
func analyzeType0(t types.Type) (typ Type) {
	if named, ok := t.(*types.Named); ok {
		pkg := named.Obj().Pkg()
		typ.Name = named.Obj().Name()
		typ.PkgPath = pkg.Path()
		typ.PkgName = pkg.Name()
		typ.PkgLocal = pkg.Name()
		typ.IsExported = named.Obj().Exported()
		t = named.Underlying()
	}

	typ.Kind = analyzeTypeKind(t)

	switch T := t.(type) {
	case *types.Basic:
		typ.IsRune = T.Name() == "rune"
		typ.IsByte = T.Name() == "byte"
	case *types.Slice:
		elem := analyzeType0(T.Elem())
		typ.Elem = &elem
	case *types.Array:
		elem := analyzeType0(T.Elem())
		typ.Elem = &elem
		typ.ArrayLen = T.Len()
	case *types.Map:
		key := analyzeType0(T.Key())
		elem := analyzeType0(T.Elem())
		typ.Key = &key
		typ.Elem = &elem
	case *types.Pointer:
		elem := analyzeType0(T.Elem())
		typ.Elem = &elem
	case *types.Interface:
		typ.IsEmptyInterface = T.NumMethods() == 0
	case *types.Struct, *types.Chan:
		// TODO probably return an error
	}

	return typ
}

// analyzeTypeKind returns the TypeKind for the given types.Type.
func analyzeTypeKind(typ types.Type) TypeKind {
	switch x := typ.(type) {
	case *types.Basic:
		return typesBasicKindToTypeKind[x.Kind()]
	case *types.Array:
		return TypeKindArray
	case *types.Chan:
		return TypeKindChan
	case *types.Signature:
		return TypeKindFunc
	case *types.Interface:
		return TypeKindInterface
	case *types.Map:
		return TypeKindMap
	case *types.Pointer:
		return TypeKindPtr
	case *types.Slice:
		return TypeKindSlice
	case *types.Struct:
		return TypeKindStruct
	case *types.Named:
		return analyzeTypeKind(x.Underlying())
	}
	return 0 // unsupported / unknown
}

// Checks all fields and their rules, and whether each rule can be applied
// to its related field without causing a compiler error.
func typeCheckRules(a *analysis, fields []*StructField) error {

	// typwalk recursively traverses the hierarchy of the given type and
	// invokes typeCheckRules for all nested struct fields it encounters.
	var typwalk func(a *analysis, typ Type) error
	typwalk = func(a *analysis, typ Type) error {
		typ = typ.PtrBase()
		switch typ.Kind {
		case TypeKindStruct:
			return typeCheckRules(a, typ.Fields)
		case TypeKindArray, TypeKindSlice:
			return typwalk(a, *typ.Elem)
		case TypeKindMap:
			if err := typwalk(a, *typ.Key); err != nil {
				return err
			}
			return typwalk(a, *typ.Elem)
		}
		return nil
	}

	// tagcheck checks the given tag's Rules and, if the tag has a Key or Elem then
	// tagcheck will recursively invoke itself with those Key/Elem instances of *TagNode.
	var tagcheck func(a *analysis, tag *TagNode, typ Type, f *StructField) error
	tagcheck = func(a *analysis, tag *TagNode, typ Type, f *StructField) error {

		// First handle the "isvalid" rule. The rule does not have to be specified explicitly,
		// instead it will be applied automatically if a type implements the "IsValid() bool" method.
		// The explicit "-isvalid" rule can be used to disable the automatic "isvalid" rule.
		canisvalid := typ.PtrBase().CanIsValid
		hasisvalid := false
		omitisvalid := false
		for _, r := range tag.Rules {
			if r.Name == "isvalid" {
				hasisvalid = true
			} else if r.Name == "-isvalid" {
				omitisvalid = true
			}
		}
		if !canisvalid && hasisvalid {
			// can't invoke IsValid() method; TODO should return error
		} else if canisvalid && !hasisvalid && !omitisvalid {
			tag.Rules = append(tag.Rules, &Rule{Name: "isvalid"})
		}

		// handle the rest
		for _, r := range tag.Rules {
			// Ensure that the Value of a RuleOption of type OptionTypeField
			// references a valid field key which will be indicated by
			// a presence of a selector in the SelectorMap.
			for _, opt := range r.Options {
				if opt.Type == OptionTypeField {
					if _, ok := a.info.SelectorMap[opt.Value]; !ok {
						return &anError{Code: errRuleOptionFieldUnknown,
							a: a, f: f, r: r, opt: opt}
					}
				}
			}

			if len(r.Context) > 0 && a.needsContext == nil {
				a.needsContext = &needsContext{f, r}
			}

			// Ensure a RuleType for the specified rule exists.
			rt, ok := a.conf.customTypeMap[r.Name]
			if !ok {
				rt, ok = defaultRuleTypeMap[r.Name]
				if !ok {
					return &anError{Code: errRuleUnknown, a: a, f: f, r: r}
				}
			}

			if err := rt.checkRule(a, r, typ, f); err != nil {
				return err
			}
		}

		// descend if key/elem are present
		if tag.Key != nil {
			typ = typ.PtrBase()
			if typ.Kind != TypeKindMap {
				return &anError{Code: errRuleKey, a: a, f: f}
			}
			if err := tagcheck(a, tag.Key, *typ.Key, f); err != nil {
				return err
			}
		}
		if tag.Elem != nil {
			typ = typ.PtrBase()
			if typ.Kind != TypeKindArray && typ.Kind != TypeKindSlice && typ.Kind != TypeKindMap {
				return &anError{Code: errRuleElem, a: a, f: f}
			}
			if err := tagcheck(a, tag.Elem, *typ.Elem, f); err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range fields {
		if f.RuleTag != nil {
			if err := tagcheck(a, f.RuleTag, f.Type, f); err != nil {
				return err
			}
		}
		if err := typwalk(a, f.Type); err != nil {
			return err
		}
	}
	return nil
}

// canConvert reports whether src type can be converted to dst type. Note that
// this does not handle unnamed struct, interface, func, and channel types.
func canConvert(dst, src Type) bool {
	// if same, accept
	if src.Equals(dst) {
		return true
	}

	// if dst is interface{}, accept
	if dst.IsEmptyInterface {
		return true
	}

	// same basic kind, accept
	if dst.Kind == src.Kind && dst.Kind.IsBasic() {
		return true
	}

	// both numeric, accept
	if dst.Kind.IsNumeric() && src.Kind.IsNumeric() {
		return true
	}

	// string from []byte, []rune, []uint8, and []int32, accept
	if dst.Kind == TypeKindString && src.Kind == TypeKindSlice && src.Elem.Name == "" &&
		(src.Elem.Kind == TypeKindUint8 || src.Elem.Kind == TypeKindInt32) {
		return true
	}
	// string to []byte, []rune, []uint8, and []int32, accept
	if src.Kind == TypeKindString && dst.Kind == TypeKindSlice && dst.Elem.Name == "" &&
		(dst.Elem.Kind == TypeKindUint8 || dst.Elem.Kind == TypeKindInt32) {
		return true
	}

	// element types (and key & len) of non-basic are equal, accept
	if dst.Kind == src.Kind && !dst.Kind.IsBasic() {
		switch dst.Kind {
		case TypeKindArray:
			return dst.ArrayLen == src.ArrayLen && dst.Elem.Equals(*src.Elem)
		case TypeKindMap:
			return dst.Key.Equals(*src.Key) && dst.Elem.Equals(*src.Elem)
		case TypeKindSlice, TypeKindPtr:
			return dst.Elem.Equals(*src.Elem)
		}
	}
	return false
}

// canConvertRuleOption reports whether or not the src RuleOption's literal
// value can be converted to the type represented by dst.
func canConvertRuleOption(a *analysis, dst Type, src *RuleOption) bool {
	if src.Type == OptionTypeField {
		field := a.info.SelectorMap[src.Value].Last()
		if canConvert(dst, field.Type) {
			return true
		}

		// can use the addr, accept
		if dst.Kind == TypeKindPtr && dst.Elem.Equals(field.Type) {
			return true // TODO add test
		}

		return false
	}

	// dst is interface{} or string, accept
	if dst.IsEmptyInterface || dst.Kind == TypeKindString {
		return true
	}

	// src is unknown, accept
	if src.Type == OptionTypeUnknown {
		return true
	}

	// both are booleans, accept
	if dst.Kind == TypeKindBool && src.Type == OptionTypeBool {
		return true
	}

	// dst is float and option is numeric, accept
	if dst.Kind.IsFloat() && (src.Type == OptionTypeInt || src.Type == OptionTypeFloat) {
		return true
	}

	// both are integers, accept
	if dst.Kind.IsInteger() && src.Type == OptionTypeInt {
		return true
	}

	// dst is unsigned and option is not negative, accept
	if dst.Kind.IsUnsigned() && src.Type == OptionTypeInt && src.Value[0] != '-' {
		return true
	}

	// src is string & dst is convertable from string, accept
	if src.Type == OptionTypeString && (dst.Kind == TypeKindString || (dst.Kind == TypeKindSlice &&
		dst.Elem.Name == "" && (dst.Elem.Kind == TypeKindUint8 || dst.Elem.Kind == TypeKindInt32))) {
		return true
	}

	return false
}

// isImportedType reports whether or not the given type is imported based on
// on the package in which the target of the analysis is declared.
func isImportedType(a *analysis, named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkgPath
}

// makeFieldKey constructs a unique field key for the given selector.
func makeFieldKey(a *analysis, selector []*StructField) (key string) {
	key = a.fieldKey(selector)
	if num, ok := a.keys[key]; ok {
		a.keys[key] = num + 1
		key += "-" + strconv.FormatUint(uint64(num), 10)
	}
	return key
}

// fieldKeyFunc returns a function that, based on the given config, produces
// field keys from a list of struct fields.
func fieldKeyFunc(conf Config) (fn func([]*StructField) string) {
	if len(conf.FieldKeyTag) > 0 {
		if conf.FieldKeyJoin {
			tag := conf.FieldKeyTag
			sep := conf.FieldKeySeparator
			// Returns the joined tag values of the fields in the given slice.
			// If one of the fields does not have a tag value set, their name
			// will be used in the join as default.
			return func(sel []*StructField) (key string) {
				for _, f := range sel {
					if f.Tag.Contains("isvalid", "omitkey") {
						continue
					}

					v := f.Tag.First(tag)
					if len(v) == 0 {
						v = f.Name
					}
					key += v + sep
				}
				if len(sep) > 0 && len(key) > len(sep) {
					return key[:len(key)-len(sep)]
				}
				return key
			}
		}

		// Returns the tag value of the last field, if no value was
		// set the field's name will be returned instead.
		return func(sel []*StructField) string {
			if key := sel[len(sel)-1].Tag.First(conf.FieldKeyTag); len(key) > 0 {
				return key
			}
			return sel[len(sel)-1].Name
		}
	}

	if conf.FieldKeyJoin {
		sep := conf.FieldKeySeparator
		// Returns the joined names of the fields in the given slice.
		return func(sel []*StructField) (key string) {
			for _, f := range sel {
				if f.Tag.Contains("isvalid", "omitkey") {
					continue
				}
				key += f.Name + sep
			}
			if len(sep) > 0 && len(key) > len(sep) {
				return key[:len(key)-len(sep)]
			}
			return key
		}
	}

	// Returns the name of the last field.
	return func(sel []*StructField) string {
		return sel[len(sel)-1].Name
	}
}
