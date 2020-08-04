package analysis

import (
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"github.com/frk/tagutil"
)

// analysis holds the state of the analyzer.
type analysis struct {
	fset *token.FileSet
	// The named type under analysis.
	named *types.Named
	// The package path of the type under analysis.
	pkgPath string
	// This field will hold the result of the analysis.
	validator *ValidatorStruct
	keys      map[string]struct{}
	//
	info *Info
}

func (a *analysis) anError(e interface{}, f *StructField, r *Rule) error {
	var err *anError

	switch v := e.(type) {
	case errorCode:
		err = &anError{Code: v}
	case *anError:
		err = v
	}

	// ...

	return err
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
	// FieldMap maintains a map of StructField pointers to the fields'
	// related go/types specific information. Intended for error reporting.
	FieldMap map[*StructField]FieldVar
	// ...
	ParamReferenceMap map[*RuleParam]*ParamReferenceInfo
}

func Run(fset *token.FileSet, named *types.Named, pos token.Pos, info *Info) (*ValidatorStruct, error) {
	structType, ok := named.Underlying().(*types.Struct)
	if !ok {
		panic(named.Obj().Name() + " must be a struct type.") // this shouldn't happen
	}

	a := new(analysis)
	a.fset = fset
	a.named = named
	a.pkgPath = named.Obj().Pkg().Path()
	a.keys = make(map[string]struct{})

	a.info = info
	a.info.FileSet = fset
	a.info.PkgPath = a.pkgPath
	a.info.TypeName = named.Obj().Name()
	a.info.TypeNamePos = pos
	a.info.FieldMap = make(map[*StructField]FieldVar)
	a.info.ParamReferenceMap = make(map[*RuleParam]*ParamReferenceInfo)

	return analyzeValidatorStruct(a, structType)
}

// analyzeValidatorStruct runs the analysis of a ValidatorStruct.
func analyzeValidatorStruct(a *analysis, structType *types.Struct) (*ValidatorStruct, error) {
	a.validator = new(ValidatorStruct)
	a.validator.TypeName = a.named.Obj().Name()

	typName := strings.ToLower(a.validator.TypeName)
	if !strings.HasSuffix(typName, "validator") {
		panic(a.validator.TypeName + " struct type has unsupported name suffix.") // this shouldn't happen
	}

	fields, err := analyzeStructFields(a, structType, true, true)
	if err != nil {
		return nil, err
	}

	for p, info := range a.info.ParamReferenceMap {
		info.Selector = findSelectorForKey(p.Value, fields)
		if len(info.Selector) == 0 {
			return nil, a.anError(errFieldKeyUnknown, info.StructField, info.Rule)
		}
	}

	a.validator.Fields = fields
	return a.validator, nil
}

func analyzeStructFields(a *analysis, structType *types.Struct, root bool, local bool) (fields []*StructField, err error) {
	for i := 0; i < structType.NumFields(); i++ {
		fvar := structType.Field(i)
		if !local && !fvar.Exported() {
			continue
		}

		ftag := structType.Tag(i)
		tag := tagutil.New(ftag)
		istag := tag.First("is")
		if istag == "-" || (!root && len(tag["is"]) == 0) {
			continue
		}

		f := new(StructField)
		f.Tag = tag
		f.Name = fvar.Name()
		f.Key = resolveFieldKey(a, fvar, ftag)
		f.IsEmbedded = fvar.Embedded()
		f.IsExported = fvar.Exported()
		a.info.FieldMap[f] = FieldVar{Var: fvar, Tag: ftag}

		// make sure the field key is unique
		if _, exists := a.keys[f.Key]; exists {
			return nil, a.anError(errFieldKeyConflict, f, nil)
		} else {
			a.keys[f.Key] = struct{}{}
		}

		typ, err := analyzeType(a, fvar.Type())
		if err != nil {
			return nil, err
		}
		f.Type = typ

		if len(istag) > 0 {
			rules, err := analyzeRules(a, f)
			if err != nil {
				return nil, err
			}
			f.Rules = rules
		}

		fields = append(fields, f)
	}
	return fields, nil
}

func analyzeType(a *analysis, t types.Type) (typ Type, err error) {
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
		elem, err := analyzeType(a, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
	case *types.Array:
		elem, err := analyzeType(a, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
		typ.ArrayLen = T.Len()
	case *types.Map:
		key, err := analyzeType(a, T.Key())
		if err != nil {
			return Type{}, err
		}
		elem, err := analyzeType(a, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Key = &key
		typ.Elem = &elem
	case *types.Pointer:
		elem, err := analyzeType(a, T.Elem())
		if err != nil {
			return Type{}, err
		}
		typ.Elem = &elem
	case *types.Interface:
		typ.IsEmptyInterface = T.NumMethods() == 0
	case *types.Struct:
		fields, err := analyzeStructFields(a, T, false, !typ.IsImported)
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

func analyzeRules(a *analysis, f *StructField) (rules []*Rule, err error) {
	for _, s := range f.Tag["is"] {
		r := parseRule(s)
		for _, p := range r.Params {
			if p.Kind == ParamKindReference {
				a.info.ParamReferenceMap[p] = &ParamReferenceInfo{
					Rule:        r,
					StructField: f,
				}
			}
		}

		// rule type check
		var rt ruleType
		rt.name = strings.ToLower(r.Name)
		rt.numParams = len(r.Params)
		if check, ok := ruleTypes[rt]; !ok {
			return nil, a.anError(errRuleUnknown, f, r)
		} else if check != nil {
			if err := check(f, r); err != nil {
				return nil, a.anError(err, f, r)
			}
		}

		rules = append(rules, r)
	}
	return rules, nil
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

func resolveFieldKey(a *analysis, fvar *types.Var, ftag string) string {
	// TODO
	// The key of the StructField (used for errors, reference args, etc.),
	// the value of this is determined by the "f key" setting, if not
	// specified by the user it will default to the value of the f's name.
	return fvar.Name() // default
}

// isImportedType reports whether or not the given type is imported based on
// on the package in which the target of the analysis is declared.
func isImportedType(a *analysis, named *types.Named) bool {
	return named != nil && named.Obj().Pkg().Path() != a.pkgPath
}

// expected format
// rule_name{ :rule_param }
func parseRule(str string) *Rule {
	str = strings.TrimSpace(str)
	name := str
	params := ""

	if i := strings.IndexByte(str, ':'); i > -1 {
		name = str[:i]
		params = str[i+1:]
	}

	r := &Rule{Name: name}

	for len(params) > 0 {
		p := &RuleParam{}
		if i := strings.IndexByte(params, ':'); i > -1 {
			p.Value = params[:i]
			params = params[i+1:]
		} else {
			p.Value = params
			params = ""
		}

		if len(p.Value) > 0 {
			switch p.Value[0] {
			case '&':
				p.Kind = ParamKindReference
				p.Value = p.Value[1:]
			case '#':
				p.Kind = ParamKindGroupKey
				p.Value = p.Value[1:]
			case '@':
				p.Kind = ParamKindContext
				p.Value = p.Value[1:]
			default:
				p.Kind = ParamKindLiteral

				// if value is surrounded by double quotes, remove both of them
				if n := len(p.Value); n > 1 && p.Value[0] == '"' && p.Value[n-1] == '"' {
					p.Value = p.Value[1 : n-1]
					p.Type = ParamTypeString
				}
			}

			if p.Kind == ParamKindLiteral && p.Type == 0 {
				switch {
				case rxNint.MatchString(p.Value):
					p.Type = ParamTypeNint
				case rxUint.MatchString(p.Value):
					p.Type = ParamTypeUint
				case rxFloat.MatchString(p.Value):
					p.Type = ParamTypeFloat
				default:
					p.Type = ParamTypeString
				}
			}
		}

		r.Params = append(r.Params, p)
	}
	return r
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
