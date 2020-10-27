package analysis

import (
	"go/types"
	"strconv"

	"github.com/frk/tagutil"
)

type (
	// ValidatorStruct represents the result of the analysis of a target struct type.
	ValidatorStruct struct {
		// Name of the validator struct type.
		TypeName string
		// The primary fields of the validator struct.
		Fields []*StructField
		// Info on the isvalid.ErrorConstructor or gosql.ErrorAggregator
		// field of the validator struct type, or nil.
		ErrorHandler *ErrorHandlerField
	}

	// StructField
	StructField struct {
		// Name of the field.
		Name string
		// The key of the StructField (used for errors, reference args, etc.),
		// the value of this is determined by the "field key" setting, if not
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
		// The list of validation rules, as parsed from the struct's tag,
		// that need to be applied to the field.
		Rules []*Rule
	}

	// Type
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

	// ErrorHandlerField is the result of analyzing a validator struct's field whose
	// type implements the isvalid.ErrorConstructor or isvalid.ErrorAggregator interface.
	ErrorHandlerField struct {
		// Name of the field (case preserved).
		Name string
		// Indicates whether or not the field's type implements
		// the isvalid.ErrorAggregator interface.
		IsAggregator bool
	}

	// Rule
	Rule struct {
		// Name of the rule
		Name string
		// The args of the rule
		Args []*RuleArg
		// The context in which the rule should be applied.
		Context string
		SetKey  string
	}

	// RuleArg
	RuleArg struct {
		// The type of the arg value.
		Type ArgType
		// The arg value, may be empty string.
		Value string
	}

	// ArgReferenceInfo holds information on a RuleArg of kind ArgTypeReference.
	ArgReferenceInfo struct {
		// The Rule to which the reference RuleArg belongs.
		Rule *Rule
		// The StructField to which the reference RuleArg belongs.
		StructField *StructField
		// Selector of the StructField referenced by the RuleArg.
		Selector []*StructField
	}

	// FieldVar holds the types.Var represenation and the raw tag of a struct field.
	FieldVar struct {
		// types.Var representation of the struct field.
		Var *types.Var
		// The raw string value of the field's tag.
		Tag string
	}
)

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

func (f *StructField) RulesCopy() []*Rule {
	if f.Rules == nil {
		return nil
	}

	rules := make([]*Rule, len(f.Rules))
	copy(rules, f.Rules)
	return rules
}

func (f *StructField) SubFields() []*StructField {
	typ := f.Type
	for typ.Kind == TypeKindPtr {
		typ = *typ.Elem // deref pointer
	}
	if typ.Kind == TypeKindStruct {
		return typ.Fields
	}
	return nil
}

func (f *StructField) HasRuleRequired() bool {
	for _, r := range f.Rules {
		if r.Name == "required" {
			return true
		}
	}
	return false
}

func (f *StructField) HasRuleNotnil() bool {
	for _, r := range f.Rules {
		if r.Name == "notnil" {
			return true
		}
	}
	return false
}

func (i ArgReferenceInfo) SelectorLast() *StructField {
	return i.Selector[len(i.Selector)-1]
}

// ArgType indicates the type of a rule arg value.
type ArgType uint

const (
	ArgTypeString ArgType = iota // default is string, i.e. r.Value == "" (empty string)
	ArgTypeNint                  // negative integer
	ArgTypeUint                  // unsigned integer
	ArgTypeFloat
	ArgTypeBool
	ArgTypeReference
)

// TypeKind indicates the specific kind of a Go type.
type TypeKind uint

const (
	// basic
	TypeKindInvalid TypeKind = iota

	_basic_kind_start
	TypeKindBool
	_numeric_kind_start
	TypeKindInt
	TypeKindInt8
	TypeKindInt16
	TypeKindInt32
	TypeKindInt64
	TypeKindUint
	TypeKindUint8
	TypeKindUint16
	TypeKindUint32
	TypeKindUint64
	TypeKindUintptr
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

func (k TypeKind) IsNumeric() bool { return _numeric_kind_start < k && k < _numeric_kind_end }

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
