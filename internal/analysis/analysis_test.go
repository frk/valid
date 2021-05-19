package analysis

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/frk/compare"
	"github.com/frk/isvalid/internal/search"
	"github.com/frk/tagutil"
)

var testast search.AST
var testpkg search.Package

func TestMain(m *testing.M) {
	pkgs, err := search.Search("../testdata/analysis", false, nil, &testast)
	if err != nil {
		log.Fatal(err)
	}
	testpkg = *(pkgs[0])

	LoadRuleTypeFunc(testast)

	os.Exit(m.Run())
}

func TestAddRuleFunc(t *testing.T) {
	tests := []struct {
		rulename string
		pkgpath  string
		funcname string
		err      error
		want     Config
		printerr bool
	}{{
		rulename: "isvalid",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyRule",
		err: &anError{Code: errRuleNameReserved, r: &Rule{Name: "isvalid"}},
	}, {
		rulename: "-isvalid",
		pkgpath:  "strings", funcname: "Contains",
		err: &anError{Code: errRuleNameReserved, r: &Rule{Name: "-isvalid"}},
	}, {
		rulename: "enum",
		pkgpath:  "strings", funcname: "Contains",
		err: &anError{Code: errRuleNameReserved, r: &Rule{Name: "enum"}},
	}, {
		rulename: "myrule",
		pkgpath:  "strings", funcname: "TrimLeftFunc",
		err: &anError{Code: errRuleFuncSignature, fn: &types.Func{}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyBadRule1",
		err: &anError{Code: errRuleFuncSignature, fn: &types.Func{}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyBadRule2",
		err: &anError{Code: errRuleFuncSignature, fn: &types.Func{}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyBadRule3",
		err: &anError{Code: errRuleFuncSignature, fn: &types.Func{}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyRule",
		want: Config{customTypeMap: map[string]RuleType{
			"myrule": RuleTypeFunc{
				FuncName:     "MyRule",
				PkgPath:      "github.com/frk/isvalid/internal/testdata/mypkg",
				PkgName:      "mypkg",
				FieldArgType: Type{Kind: TypeKindString},
				typ:          &types.Func{},
			},
		}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyRule2",
		want: Config{customTypeMap: map[string]RuleType{
			"myrule": RuleTypeFunc{
				FuncName:     "MyRule2",
				PkgPath:      "github.com/frk/isvalid/internal/testdata/mypkg",
				PkgName:      "mypkg",
				FieldArgType: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindString}},
				IsVariadic:   true,
				typ:          &types.Func{},
			},
		}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyRule3",
		want: Config{customTypeMap: map[string]RuleType{
			"myrule": RuleTypeFunc{
				FuncName:     "MyRule3",
				PkgPath:      "github.com/frk/isvalid/internal/testdata/mypkg",
				PkgName:      "mypkg",
				FieldArgType: Type{Kind: TypeKindInt64},
				OptionArgTypes: []Type{
					{Kind: TypeKindInt},
					{Kind: TypeKindFloat64},
					{Kind: TypeKindString},
					{Kind: TypeKindBool},
				},
				typ: &types.Func{},
			},
		}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}

	for _, tt := range tests {
		t.Run(tt.rulename, func(t *testing.T) {
			fn, err := search.FindFunc(tt.pkgpath, tt.funcname, testast)
			if err != nil {
				t.Fatal(err)
			}

			var conf Config
			err = conf.AddRuleFunc(tt.rulename, fn)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("Error: %v", e)
			}
			if e := compare.Compare(conf, tt.want); e != nil {
				t.Error(e)
			}

			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}

func TestAnalysisRun(t *testing.T) {
	originalhook := filenamehook
	defer func() { filenamehook = originalhook }()

	filenamehook = func(name string) string {
		if i := strings.LastIndex(name, "github.com/frk/"); i > -1 {
			return name[i:]
		}
		return name
	}

	tests := []struct {
		name     string
		want     *ValidatorStruct
		err      error
		printerr bool
	}{{
		name: "AnalysisTestBAD_ValidatorNoFieldValidator",
		err:  &anError{Code: errValidatorNoField, a: &analysis{}},
	}, {
		name: "AnalysisTestBAD_ValidatorNoField2Validator",
		err:  &anError{Code: errValidatorNoField, a: &analysis{}},
	}, {
		name: "AnalysisTestBAD_ValidatorNoField3Validator",
		err:  &anError{Code: errValidatorNoField, a: &analysis{}},
	}, {
		name: "AnalysisTestBAD_ErrorHandlerFieldConflictValidator",
		err:  &anError{Code: errErrorHandlerFieldConflict, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_ErrorHandlerFieldConflict2Validator",
		err:  &anError{Code: errErrorHandlerFieldConflict, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_ContextOptionFieldConflictValidator",
		err:  &anError{Code: errContextOptionFieldConflict, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_ContextOptionFieldTypeValidator",
		err:  &anError{Code: errContextOptionFieldType, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_ContextOptionFieldRequiredValidator",
		err:  &anError{Code: errContextOptionFieldRequired, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionFieldUnknownValidator",
		err: &anError{Code: errRuleOptionFieldUnknown, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleUnknownValidator",
		err:  &anError{Code: errRuleUnknown, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleUnknown2Validator",
		err:  &anError{Code: errRuleUnknown, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumRequiredValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foobar", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumNotNilValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foobar", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNilNotNilValidator",
		err:  &anError{Code: errRuleFieldNonNilable, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumEmailValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringEmailValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumURLValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringURLValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumPANValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringPANValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumCVVValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringCVVValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumSSNValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringSSNValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumEINValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringEINValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumNumericValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringNumericValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumHexValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringHexValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumHexcolorValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringHexcolorValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumAlnumValidator",
		err: &anError{Code: errRuleOptionValueLanguageTag, a: &analysis{}, f: &StructField{},
			r:   &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringAlnumValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumCIDRValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringCIDRValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeKindStringPhoneValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypePhoneValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "321", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2PhoneValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "true", Type: OptionTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3PhoneValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "0.2", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueCountryCodePhoneValidator",
		err: &anError{Code: errRuleOptionValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueCountryCode2PhoneValidator",
		err: &anError{Code: errRuleOptionValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "ab", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindPhoneValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringZipValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeZipValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "321", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2ZipValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "true", Type: OptionTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3ZipValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "0.2", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueCountryCodeZipValidator",
		err: &anError{Code: errRuleOptionValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueCountryCode2ZipValidator",
		err: &anError{Code: errRuleOptionValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "ab", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindZipValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringUUIDValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeUUIDValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-4", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2UUIDValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "true", Type: OptionTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3UUIDValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "0.2", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueUUIDVerUUIDValidator",
		err: &anError{Code: errRuleOptionValueUUIDVer, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueUUIDVer2UUIDValidator",
		err: &anError{Code: errRuleOptionValueUUIDVer, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "v8", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindUUIDValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "z", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumUUIDValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{
				{Value: "1", Type: OptionTypeInt},
				{Value: "2", Type: OptionTypeInt},
				{Value: "3", Type: OptionTypeInt},
				{Value: "4", Type: OptionTypeInt},
				{Value: "5", Type: OptionTypeInt},
				{Value: "6", Type: OptionTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringIPValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeIPValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-4", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2IPValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "true", Type: OptionTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3IPValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "0.2", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueIPVerIPValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "v7", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueIPVer2IPValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindIPValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumIPValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{
				{Value: "v4", Type: OptionTypeString},
				{Value: "v6", Type: OptionTypeString},
				{Value: "v8", Type: OptionTypeString},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringMACValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeMACValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-6", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2MACValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "true", Type: OptionTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3MACValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "0.2", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueMACVerMACValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "v8", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueMACVer2MACValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "vv8", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindMACValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueConflictMACValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumMACValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{
				{Value: "6", Type: OptionTypeInt},
				{Value: "8", Type: OptionTypeInt},
				{Value: "10", Type: OptionTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringRegexpValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueRegexpRegexpValidator",
		err: &anError{Code: errRuleOptionValueRegexp, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "^($", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindRegexpValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumRegexpValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2RegexpValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "foo", Type: OptionTypeString}, {Value: "bar", Type: OptionTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringPrefixValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindPrefixValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "y", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumPrefixValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeKindStringSuffixValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindSuffixValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "y", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumSuffixValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeKindStringContainsValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindContainsValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "y", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumContainsValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumEQValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeStringEQValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintEQValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeUintEQValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatEQValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindEQValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumNEValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeStringNEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintNEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeUintNEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatNEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindNEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumGTValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2GTValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "1", Type: OptionTypeInt}, {Value: "2", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericGTValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintGTValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatGTValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindGTValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumLTValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2LTValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "1", Type: OptionTypeInt}, {Value: "2", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericLTValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintLTValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatLTValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindLTValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumGTEValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2GTEValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "1", Type: OptionTypeInt}, {Value: "2", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericGTEValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintGTEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatGTEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindGTEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumLTEValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2LTEValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "1", Type: OptionTypeInt}, {Value: "2", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericLTEValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintLTEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatLTEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindLTEValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumMinValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2MinValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "1", Type: OptionTypeInt}, {Value: "2", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericMinValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintMinValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatMinValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindMinValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumMaxValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2MaxValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "1", Type: OptionTypeInt}, {Value: "2", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericMaxValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintMaxValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatMaxValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindMaxValidator",
		err: &anError{Code: errRuleBasicOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "x", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumRngValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2RngValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "123", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum3RngValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{
				{Value: "1", Type: OptionTypeInt},
				{Value: "2", Type: OptionTypeInt},
				{Value: "3", Type: OptionTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericRngValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeStringRngValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeString2RngValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "bar", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeNintRngValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFloatRngValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueBoundsRngValidator",
		err: &anError{Code: errRuleOptionValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "2", Type: OptionTypeInt}, {Value: "1.23", Type: OptionTypeFloat}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueBounds2RngValidator",
		err: &anError{Code: errRuleOptionValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "", Type: OptionTypeUnknown}, {Value: "", Type: OptionTypeUnknown}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindRngValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "y", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumLenValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2LenValidator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{
				{Value: "1", Type: OptionTypeInt},
				{Value: "2", Type: OptionTypeInt},
				{Value: "3", Type: OptionTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeLengthLenValidator",
		err:  &anError{Code: errRuleFieldLengthless, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeLenValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2LenValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3LenValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueBoundsLenValidator",
		err: &anError{Code: errRuleOptionValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "20", Type: OptionTypeInt}, {Value: "10", Type: OptionTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueBounds2LenValidator",
		err: &anError{Code: errRuleOptionValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "", Type: OptionTypeUnknown}, {Value: "", Type: OptionTypeUnknown}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeReferenceKindLenValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "y", Type: OptionTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleOptionNumRuneCountValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionNum2RuneCountValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeRunelessRuneCountValidator",
		err:  &anError{Code: errRuleFieldRuneless, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeRuneless2RuneCountValidator",
		err:  &anError{Code: errRuleFieldRuneless, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeRuneless3RuneCountValidator",
		err:  &anError{Code: errRuleFieldRuneless, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeRuneCountValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString}},
	}, {
		name: "AnalysisTestBAD_RuleOptionType2RuneCountValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "-123", Type: OptionTypeInt}},
	}, {
		name: "AnalysisTestBAD_RuleOptionType3RuneCountValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "1.23", Type: OptionTypeFloat}},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueBoundsRuneCountValidator",
		err:  &anError{Code: errRuleOptionValueBounds, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionValueBounds2RuneCountValidator",
		err: &anError{Code: errRuleOptionValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{{Value: "", Type: OptionTypeUnknown}, {Value: "", Type: OptionTypeUnknown}}}},
	}, {
		name: "AnalysisTestBAD_RuleOptionTypeFieldKindRuneCountValidator",
		err: &anError{Code: errRuleBasicOptionTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "y", Type: OptionTypeField}},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleOptionCountValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleOptionCount2Validator",
		err: &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Options: []*RuleOption{
				{Value: "a", Type: OptionTypeString},
				{Value: "b", Type: OptionTypeString},
				{Value: "c", Type: OptionTypeString},
			}},
		},
	}, {
		name: "AnalysisTestBAD_RuleFuncFieldOptionTypeValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleOptionTypeValidator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "foo", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleOptionType2Validator",
		err: &anError{Code: errRuleFuncOptionType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			opt: &RuleOption{Value: "abc", Type: OptionTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleEnumTypeUnnamedValidator",
		err:  &anError{Code: errRuleEnumTypeUnnamed, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleEnumTypeUnnamed2Validator",
		err:  &anError{Code: errRuleEnumTypeUnnamed, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleEnumTypeValidator",
		err:  &anError{Code: errRuleEnumType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleEnumType2Validator",
		err:  &anError{Code: errRuleEnumType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleEnumTypeNoConstValidator",
		err:  &anError{Code: errRuleEnumTypeNoConst, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleKeyValidator",
		err:  &anError{Code: errRuleKey, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_RuleKey2Validator",
		err:  &anError{Code: errRuleKey, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_RuleKey3Validator",
		err:  &anError{Code: errRuleKey, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_RuleElemValidator",
		err:  &anError{Code: errRuleElem, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_RuleElem2Validator",
		err:  &anError{Code: errRuleElem, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_RuleElem3Validator",
		err:  &anError{Code: errRuleElem, a: &analysis{}, f: &StructField{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionCountKeyValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionCountElemValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleOptionCountSubfieldValidator",
		err:  &anError{Code: errRuleOptionCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestOK_ErrorConstructorValidator",
		want: &ValidatorStruct{
			TypeName: "AnalysisTestOK_ErrorConstructorValidator",
			ErrorHandler: &ErrorHandlerField{
				Name: "MyErrorConstructor", IsAggregator: false,
			},
			Fields: []*StructField{{
				Name: "F", Key: "F",
				Tag:  tagutil.Tag{"is": []string{"required"}},
				Type: Type{Kind: TypeKindString}, IsExported: true,
				RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}}},
			}},
		},
	}, {
		name: "AnalysisTestOK_ErrorAggregatorValidator",
		want: &ValidatorStruct{
			TypeName: "AnalysisTestOK_ErrorAggregatorValidator",
			ErrorHandler: &ErrorHandlerField{
				Name: "erragg", IsAggregator: true,
			},
			Fields: []*StructField{{
				Name: "F", Key: "F",
				Tag:  tagutil.Tag{"is": []string{"required"}},
				Type: Type{Kind: TypeKindString}, IsExported: true,
				RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}}},
			}},
		},
	}, {
		name: "AnalysisTestOK_ContextValidator",
		want: &ValidatorStruct{
			TypeName: "AnalysisTestOK_ContextValidator",
			ContextOption: &ContextOptionField{
				Name: "context",
			},
			Fields: []*StructField{{
				Name: "F", Key: "F",
				Tag:  tagutil.Tag{"is": []string{"required"}},
				Type: Type{Kind: TypeKindString}, IsExported: true,
				RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}}},
			}},
		},
	}, {
		name: "AnalysisTestOK_Context2Validator",
		want: &ValidatorStruct{
			TypeName: "AnalysisTestOK_Context2Validator",
			ContextOption: &ContextOptionField{
				Name: "Context",
			},
			Fields: []*StructField{{
				Name: "F", Key: "F",
				Tag:  tagutil.Tag{"is": []string{"required"}},
				Type: Type{Kind: TypeKindString}, IsExported: true,
				RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}}},
			}},
		},
	}, {
		name:     "AnalysisTestOK_Validator",
		printerr: true,
		want: &ValidatorStruct{
			TypeName: "AnalysisTestOK_Validator",
			ContextOption: &ContextOptionField{
				Name: "Context",
			},
			BeforeValidate: &MethodInfo{Name: "beforevalidate"},
			AfterValidate:  &MethodInfo{Name: "AfterValidate"},
			Fields: []*StructField{{
				Name:         "UserInput",
				Tag:          tagutil.Tag{"isvalid": []string{"omitkey", "omitnilguard"}},
				RuleTag:      &TagNode{},
				OmitNilGuard: true,
				Type: Type{
					Kind: TypeKindPtr,
					Elem: &Type{
						Name:       "UserInput",
						Kind:       TypeKindStruct,
						PkgPath:    "github.com/frk/isvalid/internal/testdata",
						PkgName:    "testdata",
						PkgLocal:   "testdata",
						IsImported: true,
						IsExported: true,
						Fields: []*StructField{{
							Name: "CountryCode", Key: "CountryCode",
							Type: Type{Kind: TypeKindString}, IsExported: true,
							RuleTag: &TagNode{},
						}, {
							Name: "SomeVersion", Key: "SomeVersion",
							Type: Type{Kind: TypeKindInt}, IsExported: true,
							RuleTag: &TagNode{},
						}, {
							Name: "SomeValue", Key: "SomeValue",
							Type: Type{Kind: TypeKindString}, IsExported: true,
							RuleTag: &TagNode{},
						}, {
							Name: "F1", Key: "F1",
							Tag:  tagutil.Tag{"is": []string{"required"}},
							Type: Type{Kind: TypeKindString}, IsExported: true,
							RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}}},
						}, {
							Name: "F2", Key: "F2", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"required:@create"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "required", Context: "create"}}},
						}, {
							Name: "F5", Key: "F5", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"email"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "email"}}},
						}, {
							Name: "F6", Key: "F6", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"url"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "url"}}},
						}, {
							Name: "F8", Key: "F8", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"pan"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "pan"}}},
						}, {
							Name: "F9", Key: "F9", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"cvv"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "cvv"}}},
						}, {
							Name: "F10", Key: "F10", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"ssn"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ssn"}}},
						}, {
							Name: "F11", Key: "F11", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"ein"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ein"}}},
						}, {
							Name: "F12", Key: "F12", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"numeric"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "numeric"}}},
						}, {
							Name: "F13", Key: "F13", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"hex"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "hex"}}},
						}, {
							Name: "F14", Key: "F14", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"hexcolor"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "hexcolor"}}},
						}, {
							Name: "F15", Key: "F15", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"alnum"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "alnum", Options: []*RuleOption{
								{Value: "en", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F16", Key: "F16", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"cidr"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "cidr"}}},
						}, {
							Name: "F17", Key: "F17", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"phone"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "phone", Options: []*RuleOption{
								{Value: "us", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F18", Key: "F18", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"phone:us"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "phone", Options: []*RuleOption{
								{Value: "us", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F19", Key: "F19", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"phone:&CountryCode"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "phone", Options: []*RuleOption{
								{Value: "CountryCode", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F20", Key: "F20", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"zip"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "zip", Options: []*RuleOption{
								{Value: "us", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F21", Key: "F21", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"zip:deu"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "zip", Options: []*RuleOption{
								{Value: "deu", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F22", Key: "F22", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"zip:&CountryCode"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "zip", Options: []*RuleOption{
								{Value: "CountryCode", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F23", Key: "F23", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F24", Key: "F24", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid:3"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Options: []*RuleOption{
								{Value: "3", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F25", Key: "F25", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid:v4"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F26", Key: "F26", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid:&SomeVersion"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Options: []*RuleOption{
								{Value: "SomeVersion", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F27", Key: "F27", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Options: []*RuleOption{
								{Value: "0", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F28", Key: "F28", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip:4"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F29", Key: "F29", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip:6"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Options: []*RuleOption{
								{Value: "6", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F30", Key: "F30", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip:&SomeVersion"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Options: []*RuleOption{
								{Value: "SomeVersion", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F31", Key: "F31", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Options: []*RuleOption{
								{Value: "0", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F32", Key: "F32", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac:6"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Options: []*RuleOption{
								{Value: "6", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F33", Key: "F33", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac:8"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Options: []*RuleOption{
								{Value: "8", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F34", Key: "F34", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac:&SomeVersion"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Options: []*RuleOption{
								{Value: "SomeVersion", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F37", Key: "F37", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`re:"^[a-z]+\[[0-9]+\]$"`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "re", Options: []*RuleOption{
								{Value: `^[a-z]+\[[0-9]+\]$`, Type: OptionTypeString},
							}}}},
						}, {
							Name: "F38", Key: "F38", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`re:"\w+"`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "re", Options: []*RuleOption{
								{Value: `\w+`, Type: OptionTypeString},
							}}}},
						}, {
							Name: "F39", Key: "F39", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`contains:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "contains", Options: []*RuleOption{
								{Value: "foo bar", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F40", Key: "F40", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`contains:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "contains", Options: []*RuleOption{
								{Value: "SomeValue", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F41", Key: "F41", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`prefix:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "prefix", Options: []*RuleOption{
								{Value: "foo bar", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F42", Key: "F42", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`prefix:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "prefix", Options: []*RuleOption{
								{Value: "SomeValue", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F43", Key: "F43", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`suffix:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "suffix", Options: []*RuleOption{
								{Value: "foo bar", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F44", Key: "F44", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`suffix:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "suffix", Options: []*RuleOption{
								{Value: "SomeValue", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F45", Key: "F45", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Options: []*RuleOption{
								{Value: "foo bar", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F46", Key: "F46", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:-123`}},
							Type: Type{Kind: TypeKindInt},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Options: []*RuleOption{
								{Value: "-123", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F47", Key: "F47", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:123.987`}},
							Type: Type{Kind: TypeKindFloat64},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Options: []*RuleOption{
								{Value: "123.987", Type: OptionTypeFloat},
							}}}},
						}, {
							Name: "F48", Key: "F48", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Options: []*RuleOption{
								{Value: "SomeValue", Type: OptionTypeField},
							}}}},
						}, {
							Name: "F49", Key: "F49", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Options: []*RuleOption{
								{Value: "foo bar", Type: OptionTypeString},
							}}}},
						}, {
							Name: "F50", Key: "F50", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:-123`}},
							Type: Type{Kind: TypeKindInt},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Options: []*RuleOption{
								{Value: "-123", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F51", Key: "F51", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:123.987`}},
							Type: Type{Kind: TypeKindFloat64},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Options: []*RuleOption{
								{Value: "123.987", Type: OptionTypeFloat},
							}}}},
						}, {
							Name: "F52", Key: "F52", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Options: []*RuleOption{
								{Value: "SomeValue", Type: OptionTypeField}}}}},
						}, {
							Name: "F53", Key: "F53", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gt:24`, `lt:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gt", Options: []*RuleOption{
									{Value: "24", Type: OptionTypeInt}}},
								{Name: "lt", Options: []*RuleOption{
									{Value: "128", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F54", Key: "F54", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gt:-128`, `lt:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gt", Options: []*RuleOption{
									{Value: "-128", Type: OptionTypeInt}}},
								{Name: "lt", Options: []*RuleOption{
									{Value: "-24", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F55", Key: "F55", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gt:0.24`, `lt:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gt", Options: []*RuleOption{
									{Value: "0.24", Type: OptionTypeFloat}}},
								{Name: "lt", Options: []*RuleOption{
									{Value: "1.28", Type: OptionTypeFloat}}},
							}},
						}, {
							Name: "F56", Key: "F56", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gte:24`, `lte:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gte", Options: []*RuleOption{
									{Value: "24", Type: OptionTypeInt}}},
								{Name: "lte", Options: []*RuleOption{
									{Value: "128", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F57", Key: "F57", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gte:-128`, `lte:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gte", Options: []*RuleOption{
									{Value: "-128", Type: OptionTypeInt}}},
								{Name: "lte", Options: []*RuleOption{
									{Value: "-24", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F58", Key: "F58", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gte:0.24`, `lte:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gte", Options: []*RuleOption{
									{Value: "0.24", Type: OptionTypeFloat}}},
								{Name: "lte", Options: []*RuleOption{
									{Value: "1.28", Type: OptionTypeFloat}}},
							}},
						}, {
							Name: "F59", Key: "F59", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`min:24`, `max:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "min", Options: []*RuleOption{
									{Value: "24", Type: OptionTypeInt}}},
								{Name: "max", Options: []*RuleOption{
									{Value: "128", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F60", Key: "F60", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`min:-128`, `max:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "min", Options: []*RuleOption{
									{Value: "-128", Type: OptionTypeInt}}},
								{Name: "max", Options: []*RuleOption{
									{Value: "-24", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F61", Key: "F61", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`min:0.24`, `max:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "min", Options: []*RuleOption{
									{Value: "0.24", Type: OptionTypeFloat}}},
								{Name: "max", Options: []*RuleOption{
									{Value: "1.28", Type: OptionTypeFloat}}},
							}},
						}, {
							Name: "F62", Key: "F62", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`rng:24:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rng", Options: []*RuleOption{
								{Value: "24", Type: OptionTypeInt},
								{Value: "128", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F63", Key: "F63", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`rng:-128:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rng", Options: []*RuleOption{
								{Value: "-128", Type: OptionTypeInt},
								{Value: "-24", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F64", Key: "F64", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`rng:0.24:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rng", Options: []*RuleOption{
								{Value: "0.24", Type: OptionTypeFloat},
								{Value: "1.28", Type: OptionTypeFloat}}},
							}},
						}, {
							Name: "F65", Key: "F65", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F66", Key: "F66", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F67", Key: "F67", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindMap, Key: &Type{Kind: TypeKindString}, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F68", Key: "F68", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F69", Key: "F69", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F70", Key: "F70", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindMap, Key: &Type{Kind: TypeKindString}, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F71", Key: "F71", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len::28`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "", Type: OptionTypeUnknown},
								{Value: "28", Type: OptionTypeInt}}},
							}},
						}, {
							Name: "F72", Key: "F72", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
								{Value: "", Type: OptionTypeUnknown}}},
							}},
						}, {
							Name: "G1", Key: "G1", IsExported: true,
							Type: Type{Kind: TypeKindStruct,
								Fields: []*StructField{{
									Name: "F1", Key: "G1.F1", IsExported: true,
									Tag:     tagutil.Tag{"is": []string{`required`}},
									Type:    Type{Kind: TypeKindString},
									RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}}},
								}},
							},
							RuleTag: &TagNode{},
						}, {
							Name: "F73", Key: "F73", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{`utf8`}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "utf8"}}},
						}, {
							Name: "F74", Key: "F74", IsExported: true,
							Type: Type{
								Kind:       TypeKindStruct,
								Name:       "Time",
								PkgPath:    "time",
								PkgName:    "time",
								PkgLocal:   "time",
								IsImported: true,
								IsExported: true,
							},
							Tag: tagutil.Tag{"is": []string{`timecheck`, `ifacecheck`}},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "timecheck"},
								{Name: "ifacecheck"},
							}},
						}, {
							Name: "F75", Key: "F75", IsExported: true,
							Type: Type{
								Kind:       TypeKindString,
								Name:       "MyString",
								PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
								PkgName:    "mypkg",
								PkgLocal:   "mypkg",
								IsImported: true,
								IsExported: true,
								CanIsValid: true,
							},
							Tag:     tagutil.Tag{"is": []string{`required`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}, {Name: "isvalid"}}},
						}, {
							Name: "F76", Key: "F76", IsExported: true,
							Type: Type{
								Kind:       TypeKindInt,
								Name:       "MyInt",
								PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
								PkgName:    "mypkg",
								PkgLocal:   "mypkg",
								IsImported: true,
								IsExported: true,
								CanIsValid: true,
							},
							Tag:     tagutil.Tag{"is": []string{`required`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}, {Name: "isvalid"}}},
						}, {
							Name: "F77", Key: "F77", IsExported: true,
							Type: Type{
								Kind: TypeKindPtr,
								Elem: &Type{
									Kind:       TypeKindString,
									Name:       "MyString",
									PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
									PkgName:    "mypkg",
									PkgLocal:   "mypkg",
									IsImported: true,
									IsExported: true,
									CanIsValid: true,
								},
							},
							Tag:     tagutil.Tag{"is": []string{`-isvalid`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "-isvalid"}}},
						}, {
							Name: "F78", Key: "F78", IsExported: true,
							Type: Type{
								Kind: TypeKindPtr,
								Elem: &Type{
									Kind: TypeKindPtr,
									Elem: &Type{
										Kind:       TypeKindInt,
										Name:       "MyInt",
										PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
										PkgName:    "mypkg",
										PkgLocal:   "mypkg",
										IsImported: true,
										IsExported: true,
										CanIsValid: true,
									},
								},
							},
							Tag:     tagutil.Tag{"is": []string{`isvalid:@create`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "isvalid", Context: "create"}}},
						}, {
							Name: "F79", Key: "F79", IsExported: true,
							Type: Type{
								Kind: TypeKindPtr,
								Elem: &Type{
									Kind:       TypeKindInterface,
									CanIsValid: true,
								},
							},
							Tag:     nil,
							RuleTag: &TagNode{Rules: []*Rule{{Name: "isvalid"}}},
						}, {
							Name: "F80", Key: "F80", IsExported: true,
							Type: Type{
								Kind:       TypeKindString,
								Name:       "someKind",
								PkgPath:    "github.com/frk/isvalid/internal/testdata",
								PkgName:    "testdata",
								PkgLocal:   "testdata",
								IsImported: true,
							},
							Tag:     tagutil.Tag{"is": []string{`enum`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "enum"}}},
						}, {
							Name: "F81", Key: "F81", IsExported: true,
							Type: Type{
								Kind:       TypeKindUint,
								Name:       "MyKind",
								PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
								PkgName:    "mypkg",
								PkgLocal:   "mypkg",
								IsImported: true,
								IsExported: true,
							},
							Tag:     tagutil.Tag{"is": []string{`enum`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "enum"}}},
						}, {
							Name: "F82", Key: "F82", IsExported: true,
							Type: Type{
								Kind: TypeKindPtr,
								Elem: &Type{
									Kind: TypeKindPtr,
									Elem: &Type{
										Kind:       TypeKindString,
										Name:       "someKind",
										PkgPath:    "github.com/frk/isvalid/internal/testdata",
										PkgName:    "testdata",
										PkgLocal:   "testdata",
										IsImported: true,
									},
								},
							},
							Tag:     tagutil.Tag{"is": []string{`enum`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "enum"}}},
						}, {
							Name: "F83", Key: "F83", IsExported: true,
							Type: Type{
								Kind: TypeKindPtr,
								Elem: &Type{
									Kind:       TypeKindUint,
									Name:       "MyKind",
									PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
									PkgName:    "mypkg",
									PkgLocal:   "mypkg",
									IsImported: true,
									IsExported: true,
								},
							},
							Tag:     tagutil.Tag{"is": []string{`enum:@create`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "enum", Context: "create"}}},
						}, {
							Name: "F84", Key: "F84", IsExported: true,
							Type: Type{
								Kind: TypeKindSlice,
								Elem: &Type{Kind: TypeKindString},
							},
							Tag:     tagutil.Tag{"is": []string{`[]email`}},
							RuleTag: &TagNode{Elem: &TagNode{Rules: []*Rule{{Name: "email"}}}},
						}, {
							Name: "F85", Key: "F85", IsExported: true,
							Type: Type{
								Kind: TypeKindMap,
								Key:  &Type{Kind: TypeKindString},
								Elem: &Type{Kind: TypeKindString},
							},
							Tag:     tagutil.Tag{"is": []string{`[email]`}},
							RuleTag: &TagNode{Key: &TagNode{Rules: []*Rule{{Name: "email"}}}},
						}, {
							Name: "F86", Key: "F86", IsExported: true,
							Type: Type{
								Kind: TypeKindMap,
								Key:  &Type{Kind: TypeKindString},
								Elem: &Type{Kind: TypeKindString},
							},
							Tag: tagutil.Tag{"is": []string{`[phone:us]zip:us`}},
							RuleTag: &TagNode{
								Key: &TagNode{Rules: []*Rule{
									{Name: "phone", Options: []*RuleOption{
										{Value: "us", Type: OptionTypeString},
									}},
								}},
								Elem: &TagNode{Rules: []*Rule{
									{Name: "zip", Options: []*RuleOption{
										{Value: "us", Type: OptionTypeString},
									}},
								}},
							},
						}, {
							Name: "F87", Key: "F87", IsExported: true,
							Type: Type{
								Kind: TypeKindMap,
								Key:  &Type{Kind: TypeKindString},
								Elem: &Type{Kind: TypeKindPtr, Elem: &Type{
									Kind: TypeKindStruct,
									Fields: []*StructField{{
										Name: "F1", Key: "F87.F1", IsExported: true,
										Type: Type{Kind: TypeKindString},
										Tag:  tagutil.Tag{"is": []string{`len:2:32`}},
										RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
											{Value: "2", Type: OptionTypeInt},
											{Value: "32", Type: OptionTypeInt},
										}}}},
									}, {
										Name: "F2", Key: "F87.F2", IsExported: true,
										Type: Type{Kind: TypeKindString},
										Tag:  tagutil.Tag{"is": []string{`len:2:32`}},
										RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Options: []*RuleOption{
											{Value: "2", Type: OptionTypeInt},
											{Value: "32", Type: OptionTypeInt},
										}}}},
									}, {
										Name: "F3", Key: "F87.F3", IsExported: true,
										Type: Type{Kind: TypeKindString},
										Tag:  tagutil.Tag{"is": []string{`phone`}},
										RuleTag: &TagNode{Rules: []*Rule{{Name: "phone", Options: []*RuleOption{
											{Value: "us", Type: OptionTypeString},
										}}}},
									}},
								}},
							},
							Tag:     tagutil.Tag{"is": []string{`[email]`}},
							RuleTag: &TagNode{Key: &TagNode{Rules: []*Rule{{Name: "email"}}}},
						}, {
							Name: "F88", Key: "F88", IsExported: true,
							Type: Type{
								Kind: TypeKindSlice,
								Elem: &Type{Kind: TypeKindString},
							},
							Tag:     tagutil.Tag{"is": []string{`notnil`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "notnil"}}},
						}, {
							Name: "F89", Key: "F89", IsExported: true,
							Type: Type{
								Kind: TypeKindMap,
								Key:  &Type{Kind: TypeKindString},
								Elem: &Type{Kind: TypeKindString},
							},
							Tag:     tagutil.Tag{"is": []string{`notnil`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "notnil"}}},
						}, {
							Name: "F90", Key: "F90", IsExported: true,
							Type: Type{
								Kind:             TypeKindInterface,
								IsEmptyInterface: true,
							},
							Tag:     tagutil.Tag{"is": []string{`notnil`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "notnil"}}},
						}, {
							Name: "F91", Key: "F91", IsExported: true,
							Type: Type{
								Kind: TypeKindPtr,
								Elem: &Type{Kind: TypeKindString},
							},
							Tag:     tagutil.Tag{"is": []string{`notnil`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "notnil"}}},
						}, {
							Name: "F92", Key: "F92", IsExported: true,
							Type: Type{
								Kind: TypeKindSlice,
								Elem: &Type{
									Kind: TypeKindMap,
									Key: &Type{
										Kind: TypeKindPtr,
										Elem: &Type{
											Kind: TypeKindMap,
											Key: &Type{
												Kind: TypeKindString,
											},
											Elem: &Type{
												Kind: TypeKindString,
											},
										},
									},
									Elem: &Type{
										Kind: TypeKindSlice,
										Elem: &Type{
											Kind: TypeKindInt,
										},
									},
								},
							},
							Tag: tagutil.Tag{"is": []string{
								`[][[email]phone:us]len::10`,
								`[]rng:-54:256`,
							}},
							RuleTag: &TagNode{
								Elem: &TagNode{
									Key: &TagNode{
										Key: &TagNode{Rules: []*Rule{{Name: "email"}}},
										Elem: &TagNode{Rules: []*Rule{
											{Name: "phone", Options: []*RuleOption{
												{Value: "us", Type: OptionTypeString},
												//{Value: "ca", Type: OptionTypeString},
											}},
										}},
									},
									Elem: &TagNode{
										Rules: []*Rule{
											{Name: "len", Options: []*RuleOption{
												{Value: "", Type: OptionTypeUnknown},
												{Value: "10", Type: OptionTypeInt},
											}},
										},
										Elem: &TagNode{Rules: []*Rule{
											{Name: "rng", Options: []*RuleOption{
												{Value: "-54", Type: OptionTypeInt},
												{Value: "256", Type: OptionTypeInt},
											}},
										}},
									},
								},
							},
						}, {
							Name: "F93", Key: "F93", IsExported: true,
							Type: Type{Kind: TypeKindString},
							Tag:  tagutil.Tag{"is": []string{`runecount:28`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Options: []*RuleOption{
								{Value: "28", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F94", Key: "F94", IsExported: true,
							Type: Type{Kind: TypeKindString},
							Tag:  tagutil.Tag{"is": []string{`runecount:4:28`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
								{Value: "28", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F95", Key: "F95", IsExported: true,
							Type: Type{Kind: TypeKindSlice, Elem: &Type{
								Kind: TypeKindUint8, IsByte: true,
							}},
							Tag: tagutil.Tag{"is": []string{`runecount::28`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Options: []*RuleOption{
								{Value: "", Type: OptionTypeUnknown},
								{Value: "28", Type: OptionTypeInt},
							}}}},
						}, {
							Name: "F96", Key: "F96", IsExported: true,
							Type: Type{
								Kind:       TypeKindString,
								Name:       "someKind",
								PkgPath:    "github.com/frk/isvalid/internal/testdata",
								PkgName:    "testdata",
								PkgLocal:   "testdata",
								IsImported: true,
							},
							Tag: tagutil.Tag{"is": []string{`runecount:4:`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Options: []*RuleOption{
								{Value: "4", Type: OptionTypeInt},
								{Value: "", Type: OptionTypeUnknown},
							}}}},
						}, {
							Name: "F97", Key: "F97", IsExported: true,
							Type: Type{
								Kind: TypeKindSlice,
								Elem: &Type{Kind: TypeKindString},
							},
							Tag: tagutil.Tag{"is": []string{`required`, `[]email`}},
							RuleTag: &TagNode{
								Rules: []*Rule{{Name: "required"}},
								Elem: &TagNode{
									Rules: []*Rule{{Name: "email"}},
								},
							},
						}},
					},
				},
				IsExported: true,
			}},
		},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}

	anConf := Config{FieldKeyJoin: true, FieldKeySeparator: "."}
	anConf.customTypeMap = map[string]RuleType{
		// legit
		"utf8": RuleTypeFunc{
			FuncName:     "ValidString",
			PkgPath:      "unicode/utf8",
			FieldArgType: Type{Kind: TypeKindString},
		},

		// for testing success
		"timecheck": RuleTypeFunc{
			FuncName: "TimeCheck", PkgPath: "path/to/rule",
			FieldArgType: Type{Kind: TypeKindStruct, Name: "Time", PkgPath: "time", PkgName: "time"},
		},
		"ifacecheck": RuleTypeFunc{
			FuncName: "IfaceCheck", PkgPath: "path/to/rule",
			FieldArgType: Type{Kind: TypeKindInterface, IsEmptyInterface: true},
		},

		// for testing errors
		"rulefunc1": RuleTypeFunc{
			FuncName: "RuleFunc1", PkgPath: "path/to/rule",
			FieldArgType:   Type{Kind: TypeKindString},
			OptionArgTypes: []Type{{Kind: TypeKindInt}},
		},
		"rulefunc2": RuleTypeFunc{
			FuncName: "RuleFunc2", PkgPath: "path/to/rule",
			FieldArgType:   Type{Kind: TypeKindString},
			OptionArgTypes: []Type{{Kind: TypeKindInt}, {Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindBool}}},
			IsVariadic:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := getMatch(tt.name, t)
			got, err := anConf.Analyze(testast, match, &Info{})
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("Error: %v", e)
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Error(e)
			}

			if tt.printerr && err != nil {
				fmt.Println(err)
			}
		})
	}
}

func TestContainsRules(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "ContainsRulesTest1Validator"},
		{name: "ContainsRulesTest2Validator"},
		{name: "ContainsRulesTest3Validator"},
		{name: "ContainsRulesTest4Validator"},
		{name: "ContainsRulesTest5Validator"},
	}

	anConf := Config{FieldKeyJoin: true, FieldKeySeparator: "."}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match := getMatch(tt.name, t)
			vs, err := anConf.Analyze(testast, match, &Info{})
			if err != nil {
				t.Errorf("Error: %v", err)
			} else if len(vs.Fields) != 2 {
				t.Errorf("Bad number of fields: want=2 got=%d", len(vs.Fields))
			} else {
				for _, f := range vs.Fields {
					if f.Name == "yes" && !f.ContainsRules() {
						t.Errorf("%q does not contain rules.", f.Name)
					}
					if f.Name == "no" && f.ContainsRules() {
						t.Errorf("%q does contain rules.", f.Name)
					}
				}
			}
		})
	}
}

func getMatch(name string, t *testing.T) *search.Match {
	for _, file := range testpkg.Files {
		for _, match := range file.Matches {
			if match.Named.Obj().Name() == name {
				return match
			}
		}
	}

	t.Fatal(name, " not found")
	return nil
}
