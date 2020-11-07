package analysis

import (
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"

	"github.com/frk/tagutil"
)

type Config struct {
	FieldKeyTag       string
	FieldKeySeparator string
	FieldKeyBase      bool

	// map of custom RuleSpecs
	ruleSpecMap map[string]RuleSpec
}

// AddRuleFunc is used to register a custom RuleFuncs with the Config. The
// custom function MUST have at least one parameter item and, it MUST have
// exactly one result item which MUST be of type bool.
func (c *Config) AddRuleFunc(ruleName string, ruleFunc *types.Func) error {
	sig := ruleFunc.Type().(*types.Signature)
	p, r := sig.Params(), sig.Results()
	if p.Len() < 1 || r.Len() != 1 {
		return &anError{Code: errRuleFuncParamCount}
	}
	if !isBool(r.At(0).Type()) {
		return &anError{Code: errRuleFuncResultType}
	}

	rf := RuleFunc{iscustom: true}
	rf.FuncName = ruleFunc.Name()
	rf.PkgPath = ruleFunc.Pkg().Path()
	rf.IsVariadic = sig.Variadic()
	for i := 0; i < p.Len(); i++ {
		rf.ArgTypes = append(rf.ArgTypes, analyzeType0(p.At(i).Type()))
	}

	if c.ruleSpecMap == nil {
		c.ruleSpecMap = make(map[string]RuleSpec)
	}
	c.ruleSpecMap[ruleName] = rf
	return nil
}

func (c Config) Analyze(fset *token.FileSet, named *types.Named, pos token.Pos, info *Info) (*ValidatorStruct, error) {
	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		panic(named.Obj().Name() + " must be a struct type.") // this shouldn't happen
	}

	a := new(analysis)
	a.conf = c
	a.fset = fset
	a.named = named
	a.pkgPath = named.Obj().Pkg().Path()
	a.keys = make(map[string]uint)
	a.fieldVarMap = make(map[*StructField]fieldVar)

	a.info = info
	a.info.FileSet = fset
	a.info.PkgPath = a.pkgPath
	a.info.TypeName = named.Obj().Name()
	a.info.TypeNamePos = pos
	a.info.SelectorMap = make(map[string]StructFieldSelector)

	a.fieldKey = makeFieldKeyFunc(c)
	vs, err := analyzeValidatorStruct(a, structType)
	if err != nil {
		return nil, err
	}

	// merge the rule func maps into one for the generator to use
	a.info.RuleSpecMap = make(map[string]RuleSpec)
	for k, v := range defaultRuleSpecMap {
		a.info.RuleSpecMap[k] = v
	}
	for k, v := range c.ruleSpecMap {
		a.info.RuleSpecMap[k] = v
	}

	return vs, nil
}

// Info holds information related to an analyzed ValidatorStruct. If the analysis
// returns an error, the collected information will be incomplete.
type Info struct {
	// The FileSet associated with the analyzed ValidatorStruct.
	FileSet *token.FileSet
	// The package path of the analyzed ValidatorStruct.
	PkgPath string
	// The type name of the analyzed ValidatorStruct.
	TypeName string
	// The soruce position of the ValidatorStruct's type name.
	TypeNamePos token.Pos
	// RuleSpecMap will be populated by all the registered RuleSpecs.
	RuleSpecMap map[string]RuleSpec
	// SelectorMap maps field keys to their respective field selectors.
	SelectorMap map[string]StructFieldSelector
}

// analysis holds the state of the analyzer.
type analysis struct {
	conf Config
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
	rule  *RuleTag
}

// used for error reporting only
type fieldVar struct {
	v   *types.Var
	tag string
}

func (a *analysis) anError(e interface{}, f *StructField, r *RuleTag) error {
	var err *anError

	switch v := e.(type) {
	case errorCode:
		err = &anError{Code: v}
	case *anError:
		err = v
	}

	if f != nil {
		err.FieldName = f.Name
		err.FieldTag = f.Tag
		err.FieldType = f.Type.String()
	}
	if r != nil {
		err.RuleName = r.Name
	}

	if fv, ok := a.fieldVarMap[f]; ok {
		pos := a.fset.Position(fv.v.Pos())
		err.FieldFileName = pos.Filename
		err.FieldFileLine = pos.Line
		if f.Type.Kind == TypeKindInvalid {
			err.FieldType = fv.v.Type().String()
		}
	}

	obj := a.named.Obj()
	pos := a.fset.Position(obj.Pos())
	err.VtorName = obj.Name()
	err.VtorFileName = pos.Filename
	err.VtorFileLine = pos.Line
	return err
}

// analyzeValidatorStruct is the entry point of a ValidatorStruct analysis.
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
		return nil, a.anError(errEmptyValidator, nil, nil)
	}

	// 2. type-check all of the fields' rules
	if err := typeCheckRules(a, fields); err != nil {
		return nil, err
	}

	// 3. ensure that if a rule with context exists, that also a ContextOptionField exists
	if a.needsContext != nil && a.validator.ContextOption == nil {
		return nil, a.anError(errContextOptionFieldRequired, a.needsContext.field, a.needsContext.rule)
	}

	a.validator.Fields = fields
	return a.validator, nil
}

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
		// a separate field's rule with a reference arg.
		if istag == "-" {
			continue
		}

		f := new(StructField)
		f.Tag = tag
		f.Name = fvar.Name()
		f.IsEmbedded = fvar.Embedded()
		f.IsExported = fvar.Exported()

		// map field to fvar for error reporting
		a.fieldVarMap[f] = fieldVar{v: fvar, tag: ftag}

		// resolve field key for selector & make sure that it is unique
		fsel := append(selector, f)
		f.Key = makeFieldKey(a, fsel)
		if _, ok := a.keys[f.Key]; ok {
			return nil, a.anError(errFieldKeyConflict, f, nil)
		} else {
			a.keys[f.Key] = 1
		}

		// map keys to selectors
		a.info.SelectorMap[f.Key] = append(StructFieldSelector{}, fsel...)

		typ, err := analyzeType(a, fsel, fvar.Type())
		if err != nil {
			return nil, err
		}
		f.Type = typ

		if len(istag) > 0 {
			if err := analyzeRules(a, f); err != nil {
				return nil, err
			}
		} else if len(selector) == 0 { // root?
			// Check for untagged, "special" root fields.
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
		return a.anError(errErrorHandlerFieldConflict, f, nil)
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
		return a.anError(errContextOptionFieldConflict, f, nil)
	}
	if f.Type.Kind != TypeKindString {
		return a.anError(errContextOptionFieldType, f, nil)
	}

	a.validator.ContextOption = new(ContextOptionField)
	a.validator.ContextOption.Name = f.Name
	return nil
}

func analyzeType(a *analysis, selector []*StructField, t types.Type) (typ Type, err error) {
	if named, ok := t.(*types.Named); ok {
		pkg := named.Obj().Pkg()
		typ.Name = named.Obj().Name()
		typ.PkgPath = pkg.Path()
		typ.PkgName = pkg.Name()
		typ.PkgLocal = pkg.Name()
		typ.IsImported = isImportedType(a, named)
		typ.IsExported = named.Obj().Exported()
		t = named.Underlying()
	}

	typ.Kind = analyzeTypeKind(t)

	switch T := t.(type) {
	case *types.Basic:
		typ.IsRune = T.Name() == "rune"
		typ.IsByte = T.Name() == "byte"
	case *types.Slice:
		elem, err := analyzeType(a, selector, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
	case *types.Array:
		elem, err := analyzeType(a, selector, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
		typ.ArrayLen = T.Len()
	case *types.Map:
		key, err := analyzeType(a, selector, T.Key())
		if err != nil {
			return Type{}, err
		}
		elem, err := analyzeType(a, selector, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Key = &key
		typ.Elem = &elem
	case *types.Pointer:
		elem, err := analyzeType(a, selector, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
	case *types.Interface:
		typ.IsEmptyInterface = T.NumMethods() == 0
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

var rxInt = regexp.MustCompile(`^(?:0|-?[1-9][0-9]*)$`)
var rxFloat = regexp.MustCompile(`^(?:(?:-?0|[1-9][0-9]*)?\.[0-9]+)$`)
var rxBool = regexp.MustCompile(`^(?:false|true)$`)

func analyzeRules(a *analysis, f *StructField) error {
	for _, s := range f.Tag["is"] {
		r := parseRuleTag(s)
		if len(r.Context) > 0 && a.needsContext == nil {
			a.needsContext = &needsContext{f, r}
		}

		// make sure a spec is registered for the rule
		if _, ok := a.conf.ruleSpecMap[r.Name]; !ok {
			if _, ok := defaultRuleSpecMap[r.Name]; !ok {
				return a.anError(errRuleUnknown, f, r)
			}
		}

		f.Rules = append(f.Rules, r)
	}
	return nil
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

func makeFieldKey(a *analysis, selector []*StructField) (key string) {
	key = a.fieldKey(selector)
	if num, ok := a.keys[key]; ok {
		a.keys[key] = num + 1
		key += "-" + strconv.FormatUint(uint64(num), 10)
	}
	return key
}

// isImportedType reports whether or not the given type is imported based on
// on the package in which the target of the analysis is declared.
func isImportedType(a *analysis, named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkgPath
}

// parses the given string as the contents of a rule tag.
func parseRuleTag(str string) *RuleTag {
	str = strings.TrimSpace(str)
	name := str
	opts := ""

	if i := strings.IndexByte(str, ':'); i > -1 {
		name = str[:i]
		opts = str[i+1:]
	}

	// if the opts string ends with ':' (e.g. `len:4:`) then append
	// an empty RuleArg to the end of the RuleTag.Args slice.
	var appendEmpty bool

	r := &RuleTag{Name: name}
	for len(opts) > 0 {

		var opt string
		if i := strings.IndexByte(opts, ':'); i > -1 {
			appendEmpty = (i == len(opts)-1) // is ':' the last char?
			opt = opts[:i]
			opts = opts[i+1:]
		} else {
			opt = opts
			opts = ""
		}

		var arg *RuleArg
		if len(opt) > 0 {
			switch opt[0] {
			case '#':
				r.SetKey = opt[1:]
			case '@':
				r.Context = opt[1:]
			case '&':
				arg = &RuleArg{Type: ArgTypeField}
				arg.Value = opt[1:]
			default:
				arg = &RuleArg{}
				arg.Value = opt

				// if value is surrounded by double quotes, remove both of them
				if n := len(arg.Value); n > 1 && arg.Value[0] == '"' && arg.Value[n-1] == '"' {
					arg.Value = arg.Value[1 : n-1]
					arg.Type = ArgTypeString
				} else {
					switch {
					case rxInt.MatchString(arg.Value):
						arg.Type = ArgTypeInt
					case rxFloat.MatchString(arg.Value):
						arg.Type = ArgTypeFloat
					case rxBool.MatchString(arg.Value):
						arg.Type = ArgTypeBool
					default:
						arg.Type = ArgTypeString
					}
				}
			}
		} else {
			arg = &RuleArg{}
		}

		if arg != nil {
			r.Args = append(r.Args, arg)
		}
	}

	if appendEmpty {
		r.Args = append(r.Args, &RuleArg{})
	}
	return r
}

// Checks all fields and their rules, and whether each rule can be applied
// to its respective field without causing a compiler error.
func typeCheckRules(a *analysis, fields []*StructField) error {
	for _, f := range fields {
		for _, r := range f.Rules {
			// Ensure that the Value of a RuleArg of type ArgTypeField
			// references a valid field key which will be indicated by
			// a presence of a selector in the SelectorMap.
			for _, arg := range r.Args {
				if arg.Type == ArgTypeField {
					if _, ok := a.info.SelectorMap[arg.Value]; !ok {
						// TODO test
						return a.anError(errFieldKeyUnknown, f, r)
					}
				}
			}

			// Ensure a spec for the specified rule exists.
			spec, ok := a.conf.ruleSpecMap[r.Name]
			if !ok {
				spec, ok = defaultRuleSpecMap[r.Name]
				if !ok {
					return a.anError(errRuleUnknown, f, r)
				}
			}

			switch s := spec.(type) {
			case RuleBasic:
				if err := s.check(a, f, r); err != nil {
					return a.anError(err, f, r)
				}
			case RuleFunc:
				if err := typeCheckRuleFunc(a, f, r, s); err != nil {
					return a.anError(err, f, r)
				}
			}
		}

		if subfields := f.SubFields(); subfields != nil {
			if err := typeCheckRules(a, subfields); err != nil {
				return err
			}
		}
	}
	return nil
}

// checks whether the rule func can be applied to its respective field.
func typeCheckRuleFunc(a *analysis, f *StructField, r *RuleTag, rf RuleFunc) error {
	if rf.BoolConn > RuleFuncBoolNone {
		// func with bool connective but rule with no args, fail
		if len(r.Args) < 1 {
			return a.anError(errRuleFuncRuleArgCount, f, r)
		}
	} else {
		// rule arg count and func arg count are not compatibale, fail
		numreq := len(rf.ArgTypes[1:])
		if rf.IsVariadic {
			numreq -= 1
		}
		if numarg := len(r.Args); numreq > numarg || (numreq < numarg && !rf.IsVariadic) {
			return a.anError(errRuleFuncRuleArgCount, f, r)
		}
	}

	// field type cannot be converted to func 0th arg type, fail
	fldType, argType := f.Type.PtrBase(), rf.ArgTypes[0]
	if rf.IsVariadic && len(rf.ArgTypes) == 1 {
		argType = *argType.Elem
	}
	if !canConvert(argType, fldType) {
		return a.anError(errRuleFuncFieldArgType, f, r)
	}

	// optional check returns error, fail
	if rf.check != nil {
		if err := rf.check(a, f, r); err != nil {
			return err
		}
	}

	// rule arg cannot be converted to func arg, fail
	fatypes := rf.ArgTypes[1:]
	if rf.IsVariadic && len(fatypes) > 0 {
		fatypes = fatypes[:len(fatypes)-1]
	}
	for i, fatyp := range fatypes {
		ra := r.Args[i]
		if !canConvertRuleArg(a, fatyp, ra) {
			return a.anError(&anError{Code: errRuleFuncRuleArgType, RuleArg: ra}, f, r)
		}
	}
	if rf.IsVariadic {
		fatyp := rf.ArgTypes[len(rf.ArgTypes)-1]
		fatyp = *fatyp.Elem
		for _, ra := range r.Args[len(fatypes):] {
			if !canConvertRuleArg(a, fatyp, ra) {
				return a.anError(&anError{Code: errRuleFuncRuleArgType, RuleArg: ra}, f, r)
			}
		}
	} else if rf.BoolConn > RuleFuncBoolNone {
		fatyp := rf.ArgTypes[1]
		for _, ra := range r.Args {
			if !canConvertRuleArg(a, fatyp, ra) {
				return a.anError(&anError{Code: errRuleFuncRuleArgType, RuleArg: ra}, f, r)
			}
		}
	}
	return nil
}

// Reports whether src type can be converted to dst type. Note that this does
// not handle unnamed struct, interface, func, and channel types.
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

func canConvertRuleArg(a *analysis, dst Type, src *RuleArg) bool {
	if src.Type == ArgTypeField {
		field := a.info.SelectorMap[src.Value].Last()
		return canConvert(dst, field.Type)
	}

	// dst is interface{} or string, accept
	if dst.IsEmptyInterface || dst.Kind == TypeKindString {
		return true
	}

	// src is unknown, accept
	if src.Type == ArgTypeUnknown {
		return true
	}

	// both are booleans, accept
	if dst.Kind == TypeKindBool && src.Type == ArgTypeBool {
		return true
	}

	// dst is float and arg is numeric, accept
	if dst.Kind.IsFloat() && (src.Type == ArgTypeInt || src.Type == ArgTypeFloat) {
		return true
	}

	// both are integers, accept
	if dst.Kind.IsInteger() && src.Type == ArgTypeInt {
		return true
	}

	// dst is unsigned and arg is not negative, accept
	if dst.Kind.IsUnsigned() && src.Type == ArgTypeInt && src.Value[0] != '-' {
		return true
	}

	// src is string & dst is convertable from string, accept
	if src.Type == ArgTypeString && (dst.Kind == TypeKindString || (dst.Kind == TypeKindSlice &&
		dst.Elem.Name == "" && (dst.Elem.Kind == TypeKindUint8 || dst.Elem.Kind == TypeKindInt32))) {
		return true
	}

	return false
}

func findSelectorForKey(key string, fields []*StructField) StructFieldSelector {
	for _, f := range fields {
		if f.Key == key {
			return StructFieldSelector{f}
		}
		if f.Type.Kind == TypeKindStruct {
			if s := findSelectorForKey(key, f.Type.Fields); len(s) > 0 {
				return append(StructFieldSelector{f}, s...)
			}
		}
		if f.Type.Kind == TypeKindPtr && f.Type.Elem.Kind == TypeKindStruct {
			if s := findSelectorForKey(key, f.Type.Elem.Fields); len(s) > 0 {
				return append(StructFieldSelector{f}, s...)
			}
		}
	}
	return nil
}

func makeFieldKeyFunc(conf Config) (fn func([]*StructField) string) {
	if len(conf.FieldKeyTag) > 0 {
		if !conf.FieldKeyBase {
			return func(sel []*StructField) string {
				return fieldKeyFromTag(sel, conf.FieldKeyTag, conf.FieldKeySeparator)
			}
		}
		return func(sel []*StructField) string {
			return fieldKeyFromTagBase(sel, conf.FieldKeyTag)
		}
	}

	if !conf.FieldKeyBase {
		return func(sel []*StructField) string {
			return fieldKeyFromName(sel, conf.FieldKeySeparator)
		}
	}
	return fieldKeyFromNameBase
}

func fieldKeyFromNameBase(selector []*StructField) (key string) {
	f := selector[len(selector)-1]
	return f.Name
}

func fieldKeyFromName(selector []*StructField, sep string) (key string) {
	for _, f := range selector {
		if f.IsEmbedded {
			continue
		}
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

func fieldKeyFromTagBase(selector []*StructField, tag string) (key string) {
	f := selector[len(selector)-1]
	key = f.Tag.First(tag)
	if len(key) == 0 {
		key = f.Name
	}
	return key
}

func fieldKeyFromTag(selector []*StructField, tag, sep string) (key string) {
	for _, f := range selector {
		if f.IsEmbedded {
			continue
		}
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
