package analysis

import (
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"

	"github.com/frk/tagutil"
)

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
	// FieldMap maintains a map of StructField pointers to the fields'
	// related go/types specific information. Intended for error reporting.
	FieldMap map[*StructField]FieldVar
	// Maps RuleArgs of kind ArgTypeReference to information
	// that's used by the analysis and by the generator.
	ArgReferenceMap map[*RuleArg]*ArgReferenceInfo
}

// analysis holds the state of the analyzer.
type analysis struct {
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
}

// used for error reporting only
type needsContext struct {
	field *StructField
	rule  *Rule
}

type Config struct {
	FieldKeyTag       string
	FieldKeySeparator string
	FieldKeyBase      bool
}

func (c Config) Analyze(fset *token.FileSet, named *types.Named, pos token.Pos, info *Info) (*ValidatorStruct, error) {
	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		panic(named.Obj().Name() + " must be a struct type.") // this shouldn't happen
	}

	a := new(analysis)
	a.fset = fset
	a.named = named
	a.pkgPath = named.Obj().Pkg().Path()
	a.keys = make(map[string]uint)

	a.info = info
	a.info.FileSet = fset
	a.info.PkgPath = a.pkgPath
	a.info.TypeName = named.Obj().Name()
	a.info.TypeNamePos = pos
	a.info.FieldMap = make(map[*StructField]FieldVar)
	a.info.ArgReferenceMap = make(map[*RuleArg]*ArgReferenceInfo)

	a.fieldKey = makeFieldKeyFunc(c)
	return analyzeValidatorStruct(a, structType)
}

func (a *analysis) anError(e interface{}, f *StructField, r *Rule) error {
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

	if fv, ok := a.info.FieldMap[f]; ok {
		pos := a.fset.Position(fv.Var.Pos())
		err.FieldFileName = pos.Filename
		err.FieldFileLine = pos.Line
		if f.Type.Kind == TypeKindInvalid {
			err.FieldType = fv.Var.Type().String()
		}
	}

	obj := a.named.Obj()
	pos := a.fset.Position(obj.Pos())
	err.VtorName = obj.Name()
	err.VtorFileName = pos.Filename
	err.VtorFileLine = pos.Line
	return err
}

// analyzeValidatorStruct runs the analysis of a ValidatorStruct.
func analyzeValidatorStruct(a *analysis, structType *types.Struct) (*ValidatorStruct, error) {
	a.validator = new(ValidatorStruct)
	a.validator.TypeName = a.named.Obj().Name()

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

	// 2. resolve any arg references
	for p, info := range a.info.ArgReferenceMap {
		info.Selector = findSelectorForKey(p.Value, fields)
		if len(info.Selector) == 0 {
			return nil, a.anError(errFieldKeyUnknown, info.StructField, info.Rule)
		}
	}

	// 3. type-check all of the fields' rules
	if err := ruleCheckStructFields(a, fields); err != nil {
		return nil, err
	}

	// 4. ensure that if a rule with context exists, that then also a ContextOptionField exists
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
		a.info.FieldMap[f] = FieldVar{Var: fvar, Tag: ftag}

		// resolve field key for selector & make sure that it is unique
		fsel := append(selector, f)
		f.Key = makeFieldKey(a, fsel)
		if _, ok := a.keys[f.Key]; ok {
			return nil, a.anError(errFieldKeyConflict, f, nil)
		} else {
			a.keys[f.Key] = 1
		}

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

func analyzeErrorHandlerField(a *analysis, f *StructField, isAggregator bool) error {
	if a.validator.ErrorHandler != nil {
		return a.anError(errErrorHandlerFieldConflict, f, nil)
	}

	a.validator.ErrorHandler = new(ErrorHandlerField)
	a.validator.ErrorHandler.Name = f.Name
	a.validator.ErrorHandler.IsAggregator = isAggregator
	return nil
}

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

var rxNint = regexp.MustCompile(`^(?:-[1-9][0-9]*)$`) // negative integer
var rxUint = regexp.MustCompile(`^(?:0|[1-9][0-9]*)$`)
var rxFloat = regexp.MustCompile(`^(?:(?:-?0|[1-9][0-9]*)?\.[0-9]+)$`)
var rxBool = regexp.MustCompile(`^(?:false|true)$`)

func analyzeRules(a *analysis, f *StructField) error {
	for _, s := range f.Tag["is"] {
		r := parseRule(s)
		for _, arg := range r.Args {
			if arg.Type == ArgTypeReference {
				a.info.ArgReferenceMap[arg] = &ArgReferenceInfo{
					Rule:        r,
					StructField: f,
				}
			}
		}
		if len(r.Context) > 0 && a.needsContext == nil {
			a.needsContext = &needsContext{f, r}
		}

		// check that rule type exists
		if _, err := ruleTypes.find(r); err != nil {
			return a.anError(err, f, r)
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

// expected format
// rule_name{ :rule_opt }
func parseRule(str string) *Rule {
	str = strings.TrimSpace(str)
	name := str
	opts := ""

	if i := strings.IndexByte(str, ':'); i > -1 {
		name = str[:i]
		opts = str[i+1:]
	}

	// if the opts string ends with ':' (e.g. `len:4:`) then append
	// an empty RuleArg to the end of the Rule.Args slice.
	var appendEmpty bool

	r := &Rule{Name: name}
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
				arg = &RuleArg{Type: ArgTypeReference}
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
					case rxNint.MatchString(arg.Value):
						arg.Type = ArgTypeNint
					case rxUint.MatchString(arg.Value):
						arg.Type = ArgTypeUint
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
			arg = &RuleArg{Type: ArgTypeString} // empty string
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

func ruleCheckStructFields(a *analysis, fields []*StructField) error {
	for _, f := range fields {
		for _, r := range f.Rules {
			if err := ruleTypes.check(a, f, r); err != nil {
				return a.anError(err, f, r)
			}
		}

		if subfields := f.SubFields(); subfields != nil {
			if err := ruleCheckStructFields(a, subfields); err != nil {
				return err
			}
		}

		//if f.Type.Kind == TypeKindStruct {
		//	if err := ruleCheckStructFields(a, f.Type.Fields); err != nil {
		//		return err
		//	}
		//}
		//if f.Type.Kind == TypeKindPtr && f.Type.Elem.Kind == TypeKindStruct {
		//	if err := ruleCheckStructFields(a, f.Type.Elem.Fields); err != nil {
		//		return err
		//	}
		//}
	}
	return nil
}

func findSelectorForKey(key string, fields []*StructField) []*StructField {
	for _, f := range fields {
		if f.Key == key {
			return []*StructField{f}
		}
		if f.Type.Kind == TypeKindStruct {
			if s := findSelectorForKey(key, f.Type.Fields); len(s) > 0 {
				return append([]*StructField{f}, s...)
			}
		}
		if f.Type.Kind == TypeKindPtr && f.Type.Elem.Kind == TypeKindStruct {
			if s := findSelectorForKey(key, f.Type.Elem.Fields); len(s) > 0 {
				return append([]*StructField{f}, s...)
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
