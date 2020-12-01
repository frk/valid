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
				FuncName: "MyRule",
				PkgPath:  "github.com/frk/isvalid/internal/testdata/mypkg",
				ArgTypes: []Type{{Kind: TypeKindString}},
				typ:      &types.Func{},
			},
		}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyRule2",
		want: Config{customTypeMap: map[string]RuleType{
			"myrule": RuleTypeFunc{
				FuncName:   "MyRule2",
				PkgPath:    "github.com/frk/isvalid/internal/testdata/mypkg",
				ArgTypes:   []Type{{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindString}}},
				IsVariadic: true,
				typ:        &types.Func{},
			},
		}},
	}, {
		rulename: "myrule",
		pkgpath:  "github.com/frk/isvalid/internal/testdata/mypkg", funcname: "MyRule3",
		want: Config{customTypeMap: map[string]RuleType{
			"myrule": RuleTypeFunc{
				FuncName: "MyRule3",
				PkgPath:  "github.com/frk/isvalid/internal/testdata/mypkg",
				ArgTypes: []Type{
					{Kind: TypeKindInt64},
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
		name: "AnalysisTestBAD_RuleArgFieldUnknownValidator",
		err: &anError{Code: errRuleArgFieldUnknown, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleUnknownValidator",
		err:  &anError{Code: errRuleUnknown, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleUnknown2Validator",
		err:  &anError{Code: errRuleUnknown, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumRequiredValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foobar", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumNotNilValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foobar", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNilNotNilValidator",
		err:  &anError{Code: errRuleFieldNonNilable, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumEmailValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringEmailValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumURLValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringURLValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumURIValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringURIValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumPANValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringPANValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumCVVValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringCVVValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumSSNValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringSSNValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumEINValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringEINValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumNumericValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringNumericValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumHexValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringHexValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumHexcolorValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringHexcolorValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumAlnumValidator",
		err: &anError{Code: errRuleArgValueLanguageTag, a: &analysis{}, f: &StructField{},
			r:  &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringAlnumValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumCIDRValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringCIDRValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeKindStringPhoneValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypePhoneValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "321", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2PhoneValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3PhoneValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueCountryCodePhoneValidator",
		err: &anError{Code: errRuleArgValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueCountryCode2PhoneValidator",
		err: &anError{Code: errRuleArgValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "ab", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindPhoneValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringZipValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeZipValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "321", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2ZipValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3ZipValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueCountryCodeZipValidator",
		err: &anError{Code: errRuleArgValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueCountryCode2ZipValidator",
		err: &anError{Code: errRuleArgValueCountryCode, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "ab", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindZipValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringUUIDValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeUUIDValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-4", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2UUIDValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3UUIDValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueUUIDVerUUIDValidator",
		err: &anError{Code: errRuleArgValueUUIDVer, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueUUIDVer2UUIDValidator",
		err: &anError{Code: errRuleArgValueUUIDVer, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "v8", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindUUIDValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "z", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueConflictUUIDValidator",
		err: &anError{Code: errRuleArgValueConflict, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "4", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumUUIDValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{
				{Value: "1", Type: ArgTypeInt},
				{Value: "2", Type: ArgTypeInt},
				{Value: "3", Type: ArgTypeInt},
				{Value: "4", Type: ArgTypeInt},
				{Value: "5", Type: ArgTypeInt},
				{Value: "6", Type: ArgTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringIPValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeIPValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-4", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2IPValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3IPValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueIPVerIPValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "v7", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueIPVer2IPValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindIPValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumIPValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{
				{Value: "v4", Type: ArgTypeString},
				{Value: "v6", Type: ArgTypeString},
				{Value: "v8", Type: ArgTypeString},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringMACValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeMACValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-6", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2MACValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3MACValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueMACVerMACValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "v8", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueMACVer2MACValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "vv8", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindMACValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgValueConflictMACValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumMACValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{
				{Value: "6", Type: ArgTypeInt},
				{Value: "8", Type: ArgTypeInt},
				{Value: "10", Type: ArgTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringISOValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeISOValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2ISOValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3ISOValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindISOValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumISOValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2ISOValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "6", Type: ArgTypeInt}, {Value: "8", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringRFCValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeRFCValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2RFCValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "true", Type: ArgTypeBool},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3RFCValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "0.2", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindRFCValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumRFCValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2RFCValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "6", Type: ArgTypeInt}, {Value: "8", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringRegexpValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgValueRegexpRegexpValidator",
		err: &anError{Code: errRuleArgValueRegexp, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "^($", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindRegexpValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumRegexpValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2RegexpValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "foo", Type: ArgTypeString}, {Value: "bar", Type: ArgTypeString}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeKindStringPrefixValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindPrefixValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "y", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumPrefixValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeKindStringSuffixValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindSuffixValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "y", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumSuffixValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_TypeKindStringContainsValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindContainsValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "y", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumContainsValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNumEQValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeStringEQValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintEQValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeUintEQValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatEQValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindEQValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumNEValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeStringNEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintNEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeUintNEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatNEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindNEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumGTValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2GTValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "1", Type: ArgTypeInt}, {Value: "2", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericGTValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintGTValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatGTValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindGTValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumLTValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2LTValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "1", Type: ArgTypeInt}, {Value: "2", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericLTValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintLTValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatLTValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindLTValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumGTEValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2GTEValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "1", Type: ArgTypeInt}, {Value: "2", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericGTEValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintGTEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatGTEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindGTEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumLTEValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2LTEValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "1", Type: ArgTypeInt}, {Value: "2", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericLTEValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintLTEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatLTEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindLTEValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumMinValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2MinValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "1", Type: ArgTypeInt}, {Value: "2", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericMinValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintMinValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatMinValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindMinValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumMaxValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2MaxValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "1", Type: ArgTypeInt}, {Value: "2", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericMaxValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintMaxValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatMaxValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindMaxValidator",
		err: &anError{Code: errRuleBasicArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "x", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumRngValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2RngValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "123", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNum3RngValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{
				{Value: "1", Type: ArgTypeInt},
				{Value: "2", Type: ArgTypeInt},
				{Value: "3", Type: ArgTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeNumericRngValidator",
		err:  &anError{Code: errRuleFieldNonNumeric, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeStringRngValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeString2RngValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "bar", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeNintRngValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFloatRngValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueBoundsRngValidator",
		err: &anError{Code: errRuleArgValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "2", Type: ArgTypeInt}, {Value: "1.23", Type: ArgTypeFloat}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueBounds2RngValidator",
		err: &anError{Code: errRuleArgValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "", Type: ArgTypeUnknown}, {Value: "", Type: ArgTypeUnknown}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindRngValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "y", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumLenValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2LenValidator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{
				{Value: "1", Type: ArgTypeInt},
				{Value: "2", Type: ArgTypeInt},
				{Value: "3", Type: ArgTypeInt},
			}},
		},
	}, {
		name: "AnalysisTestBAD_TypeLengthLenValidator",
		err:  &anError{Code: errRuleFieldLengthless, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeLenValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType2LenValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgType3LenValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueBoundsLenValidator",
		err: &anError{Code: errRuleArgValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "20", Type: ArgTypeInt}, {Value: "10", Type: ArgTypeInt}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgValueBounds2LenValidator",
		err: &anError{Code: errRuleArgValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "", Type: ArgTypeUnknown}, {Value: "", Type: ArgTypeUnknown}}},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeReferenceKindLenValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "y", Type: ArgTypeField},
		},
	}, {
		name: "AnalysisTestBAD_RuleArgNumRuneCountValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgNum2RuneCountValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
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
		name: "AnalysisTestBAD_RuleArgTypeRuneCountValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "foo", Type: ArgTypeString}},
	}, {
		name: "AnalysisTestBAD_RuleArgType2RuneCountValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "-123", Type: ArgTypeInt}},
	}, {
		name: "AnalysisTestBAD_RuleArgType3RuneCountValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "1.23", Type: ArgTypeFloat}},
	}, {
		name: "AnalysisTestBAD_RuleArgValueBoundsRuneCountValidator",
		err:  &anError{Code: errRuleArgValueBounds, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgValueBounds2RuneCountValidator",
		err: &anError{Code: errRuleArgValueBounds, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{{Value: "", Type: ArgTypeUnknown}, {Value: "", Type: ArgTypeUnknown}}}},
	}, {
		name: "AnalysisTestBAD_RuleArgTypeFieldKindRuneCountValidator",
		err: &anError{Code: errRuleBasicArgTypeUint, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{Value: "y", Type: ArgTypeField}},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleArgCountValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleArgCount2Validator",
		err: &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{},
			r: &Rule{Args: []*RuleArg{
				{Value: "a", Type: ArgTypeString},
				{Value: "b", Type: ArgTypeString},
				{Value: "c", Type: ArgTypeString},
			}},
		},
	}, {
		name: "AnalysisTestBAD_RuleFuncFieldArgTypeValidator",
		err:  &anError{Code: errRuleFuncFieldType, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleArgTypeValidator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{"foo", ArgTypeString},
		},
	}, {
		name: "AnalysisTestBAD_RuleFuncRuleArgType2Validator",
		err: &anError{Code: errRuleFuncArgType, a: &analysis{}, f: &StructField{}, r: &Rule{},
			ra: &RuleArg{"abc", ArgTypeString},
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
		name: "AnalysisTestBAD_RuleArgCountKeyValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgCountElemValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
	}, {
		name: "AnalysisTestBAD_RuleArgCountSubfieldValidator",
		err:  &anError{Code: errRuleArgCount, a: &analysis{}, f: &StructField{}, r: &Rule{}},
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
				Name:    "UserInput",
				Tag:     tagutil.Tag{"isvalid": []string{"omitkey"}},
				RuleTag: &TagNode{},
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
							// NOTE the #key is dropped for now ...
							// }, {
							// 	Name: "F3", Key: "F3", IsExported: true,
							// 	Tag:     tagutil.Tag{"is": []string{"required:#key"}},
							// 	Type:    Type{Kind: TypeKindString},
							// 	RuleTag: &TagNode{Rules: []*Rule{{Name: "required"}, {Name: "#key"}}},
							// }, {
							// 	Name: "F4", Key: "F4", IsExported: true,
							// 	Tag:     tagutil.Tag{"is": []string{"required:@create:#key"}},
							// 	Type:    Type{Kind: TypeKindString},
							// 	RuleTag: &TagNode{Rules: []*Rule{{Name: "required", Context: "create"}, {Name: "#key"}}},
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
							Name: "F7", Key: "F7", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"uri"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uri"}}},
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
							RuleTag: &TagNode{Rules: []*Rule{{Name: "alnum", Args: []*RuleArg{
								{Value: "en", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F16", Key: "F16", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"cidr"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "cidr"}}},
						}, {
							Name: "F17", Key: "F17", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"phone"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "phone"}}},
						}, {
							Name: "F18", Key: "F18", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"phone:us:ca"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "phone", Args: []*RuleArg{
								{Value: "us", Type: ArgTypeString},
								{Value: "ca", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F19", Key: "F19", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"phone:&CountryCode"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "phone", Args: []*RuleArg{
								{Value: "CountryCode", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F20", Key: "F20", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"zip"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "zip"}}},
						}, {
							Name: "F21", Key: "F21", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"zip:deu:fin"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "zip", Args: []*RuleArg{
								{Value: "deu", Type: ArgTypeString},
								{Value: "fin", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F22", Key: "F22", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"zip:&CountryCode"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "zip", Args: []*RuleArg{
								{Value: "CountryCode", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F23", Key: "F23", IsExported: true,
							Tag:     tagutil.Tag{"is": []string{"uuid"}},
							Type:    Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid"}}},
						}, {
							Name: "F24", Key: "F24", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid:3"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Args: []*RuleArg{
								{Value: "3", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F25", Key: "F25", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid:v4"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F26", Key: "F26", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"uuid:&SomeVersion"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "uuid", Args: []*RuleArg{
								{Value: "SomeVersion", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F27", Key: "F27", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Args: []*RuleArg{
								{Value: "0", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F28", Key: "F28", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip:4"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F29", Key: "F29", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip:6"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Args: []*RuleArg{
								{Value: "6", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F30", Key: "F30", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"ip:&SomeVersion"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ip", Args: []*RuleArg{
								{Value: "SomeVersion", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F31", Key: "F31", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Args: []*RuleArg{
								{Value: "0", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F32", Key: "F32", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac:6"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Args: []*RuleArg{
								{Value: "6", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F33", Key: "F33", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac:8"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Args: []*RuleArg{
								{Value: "8", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F34", Key: "F34", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"mac:&SomeVersion"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "mac", Args: []*RuleArg{
								{Value: "SomeVersion", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F35", Key: "F35", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"iso:1234"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "iso", Args: []*RuleArg{
								{Value: "1234", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F36", Key: "F36", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{"rfc:1234"}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rfc", Args: []*RuleArg{
								{Value: "1234", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F37", Key: "F37", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`re:"^[a-z]+\[[0-9]+\]$"`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "re", Args: []*RuleArg{
								{Value: `^[a-z]+\[[0-9]+\]$`, Type: ArgTypeString},
							}}}},
						}, {
							Name: "F38", Key: "F38", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`re:"\w+"`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "re", Args: []*RuleArg{
								{Value: `\w+`, Type: ArgTypeString},
							}}}},
						}, {
							Name: "F39", Key: "F39", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`contains:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "contains", Args: []*RuleArg{
								{Value: "foo bar", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F40", Key: "F40", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`contains:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "contains", Args: []*RuleArg{
								{Value: "SomeValue", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F41", Key: "F41", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`prefix:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "prefix", Args: []*RuleArg{
								{Value: "foo bar", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F42", Key: "F42", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`prefix:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "prefix", Args: []*RuleArg{
								{Value: "SomeValue", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F43", Key: "F43", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`suffix:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "suffix", Args: []*RuleArg{
								{Value: "foo bar", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F44", Key: "F44", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`suffix:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "suffix", Args: []*RuleArg{
								{Value: "SomeValue", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F45", Key: "F45", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Args: []*RuleArg{
								{Value: "foo bar", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F46", Key: "F46", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:-123`}},
							Type: Type{Kind: TypeKindInt},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Args: []*RuleArg{
								{Value: "-123", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F47", Key: "F47", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:123.987`}},
							Type: Type{Kind: TypeKindFloat64},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Args: []*RuleArg{
								{Value: "123.987", Type: ArgTypeFloat},
							}}}},
						}, {
							Name: "F48", Key: "F48", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`eq:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "eq", Args: []*RuleArg{
								{Value: "SomeValue", Type: ArgTypeField},
							}}}},
						}, {
							Name: "F49", Key: "F49", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:foo bar`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Args: []*RuleArg{
								{Value: "foo bar", Type: ArgTypeString},
							}}}},
						}, {
							Name: "F50", Key: "F50", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:-123`}},
							Type: Type{Kind: TypeKindInt},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Args: []*RuleArg{
								{Value: "-123", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F51", Key: "F51", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:123.987`}},
							Type: Type{Kind: TypeKindFloat64},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Args: []*RuleArg{
								{Value: "123.987", Type: ArgTypeFloat},
							}}}},
						}, {
							Name: "F52", Key: "F52", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`ne:&SomeValue`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "ne", Args: []*RuleArg{
								{Value: "SomeValue", Type: ArgTypeField}}}}},
						}, {
							Name: "F53", Key: "F53", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gt:24`, `lt:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gt", Args: []*RuleArg{
									{Value: "24", Type: ArgTypeInt}}},
								{Name: "lt", Args: []*RuleArg{
									{Value: "128", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F54", Key: "F54", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gt:-128`, `lt:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gt", Args: []*RuleArg{
									{Value: "-128", Type: ArgTypeInt}}},
								{Name: "lt", Args: []*RuleArg{
									{Value: "-24", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F55", Key: "F55", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gt:0.24`, `lt:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gt", Args: []*RuleArg{
									{Value: "0.24", Type: ArgTypeFloat}}},
								{Name: "lt", Args: []*RuleArg{
									{Value: "1.28", Type: ArgTypeFloat}}},
							}},
						}, {
							Name: "F56", Key: "F56", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gte:24`, `lte:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gte", Args: []*RuleArg{
									{Value: "24", Type: ArgTypeInt}}},
								{Name: "lte", Args: []*RuleArg{
									{Value: "128", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F57", Key: "F57", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gte:-128`, `lte:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gte", Args: []*RuleArg{
									{Value: "-128", Type: ArgTypeInt}}},
								{Name: "lte", Args: []*RuleArg{
									{Value: "-24", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F58", Key: "F58", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`gte:0.24`, `lte:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "gte", Args: []*RuleArg{
									{Value: "0.24", Type: ArgTypeFloat}}},
								{Name: "lte", Args: []*RuleArg{
									{Value: "1.28", Type: ArgTypeFloat}}},
							}},
						}, {
							Name: "F59", Key: "F59", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`min:24`, `max:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "min", Args: []*RuleArg{
									{Value: "24", Type: ArgTypeInt}}},
								{Name: "max", Args: []*RuleArg{
									{Value: "128", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F60", Key: "F60", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`min:-128`, `max:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "min", Args: []*RuleArg{
									{Value: "-128", Type: ArgTypeInt}}},
								{Name: "max", Args: []*RuleArg{
									{Value: "-24", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F61", Key: "F61", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`min:0.24`, `max:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{
								{Name: "min", Args: []*RuleArg{
									{Value: "0.24", Type: ArgTypeFloat}}},
								{Name: "max", Args: []*RuleArg{
									{Value: "1.28", Type: ArgTypeFloat}}},
							}},
						}, {
							Name: "F62", Key: "F62", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`rng:24:128`}},
							Type: Type{Kind: TypeKindUint8},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rng", Args: []*RuleArg{
								{Value: "24", Type: ArgTypeInt},
								{Value: "128", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F63", Key: "F63", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`rng:-128:-24`}},
							Type: Type{Kind: TypeKindInt16},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rng", Args: []*RuleArg{
								{Value: "-128", Type: ArgTypeInt},
								{Value: "-24", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F64", Key: "F64", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`rng:0.24:1.28`}},
							Type: Type{Kind: TypeKindFloat32},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "rng", Args: []*RuleArg{
								{Value: "0.24", Type: ArgTypeFloat},
								{Value: "1.28", Type: ArgTypeFloat}}},
							}},
						}, {
							Name: "F65", Key: "F65", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F66", Key: "F66", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F67", Key: "F67", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:28`}},
							Type: Type{Kind: TypeKindMap, Key: &Type{Kind: TypeKindString}, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F68", Key: "F68", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F69", Key: "F69", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F70", Key: "F70", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:28`}},
							Type: Type{Kind: TypeKindMap, Key: &Type{Kind: TypeKindString}, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F71", Key: "F71", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len::28`}},
							Type: Type{Kind: TypeKindString},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "", Type: ArgTypeUnknown},
								{Value: "28", Type: ArgTypeInt}}},
							}},
						}, {
							Name: "F72", Key: "F72", IsExported: true,
							Tag:  tagutil.Tag{"is": []string{`len:4:`}},
							Type: Type{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindInt}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
								{Value: "", Type: ArgTypeUnknown}}},
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
							Tag: tagutil.Tag{"is": []string{`[phone:us:ca]zip:ca:us`}},
							RuleTag: &TagNode{
								Key: &TagNode{Rules: []*Rule{
									{Name: "phone", Args: []*RuleArg{
										{Value: "us", Type: ArgTypeString},
										{Value: "ca", Type: ArgTypeString},
									}},
								}},
								Elem: &TagNode{Rules: []*Rule{
									{Name: "zip", Args: []*RuleArg{
										{Value: "ca", Type: ArgTypeString},
										{Value: "us", Type: ArgTypeString},
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
										RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
											{Value: "2", Type: ArgTypeInt},
											{Value: "32", Type: ArgTypeInt},
										}}}},
									}, {
										Name: "F2", Key: "F87.F2", IsExported: true,
										Type: Type{Kind: TypeKindString},
										Tag:  tagutil.Tag{"is": []string{`len:2:32`}},
										RuleTag: &TagNode{Rules: []*Rule{{Name: "len", Args: []*RuleArg{
											{Value: "2", Type: ArgTypeInt},
											{Value: "32", Type: ArgTypeInt},
										}}}},
									}, {
										Name: "F3", Key: "F87.F3", IsExported: true,
										Type:    Type{Kind: TypeKindString},
										Tag:     tagutil.Tag{"is": []string{`phone`}},
										RuleTag: &TagNode{Rules: []*Rule{{Name: "phone"}}},
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
								`[][[email]phone:us:ca]len::10`,
								`[]rng:-54:256`,
							}},
							RuleTag: &TagNode{
								Elem: &TagNode{
									Key: &TagNode{
										Key: &TagNode{Rules: []*Rule{{Name: "email"}}},
										Elem: &TagNode{Rules: []*Rule{
											{Name: "phone", Args: []*RuleArg{
												{Value: "us", Type: ArgTypeString},
												{Value: "ca", Type: ArgTypeString},
											}},
										}},
									},
									Elem: &TagNode{
										Rules: []*Rule{
											{Name: "len", Args: []*RuleArg{
												{Value: "", Type: ArgTypeUnknown},
												{Value: "10", Type: ArgTypeInt},
											}},
										},
										Elem: &TagNode{Rules: []*Rule{
											{Name: "rng", Args: []*RuleArg{
												{Value: "-54", Type: ArgTypeInt},
												{Value: "256", Type: ArgTypeInt},
											}},
										}},
									},
								},
							},
						}, {
							Name: "F93", Key: "F93", IsExported: true,
							Type: Type{Kind: TypeKindString},
							Tag:  tagutil.Tag{"is": []string{`runecount:28`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Args: []*RuleArg{
								{Value: "28", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F94", Key: "F94", IsExported: true,
							Type: Type{Kind: TypeKindString},
							Tag:  tagutil.Tag{"is": []string{`runecount:4:28`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
								{Value: "28", Type: ArgTypeInt},
							}}}},
						}, {
							Name: "F95", Key: "F95", IsExported: true,
							Type: Type{Kind: TypeKindSlice, Elem: &Type{
								Kind: TypeKindUint8, IsByte: true,
							}},
							Tag: tagutil.Tag{"is": []string{`runecount::28`}},
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Args: []*RuleArg{
								{Value: "", Type: ArgTypeUnknown},
								{Value: "28", Type: ArgTypeInt},
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
							RuleTag: &TagNode{Rules: []*Rule{{Name: "runecount", Args: []*RuleArg{
								{Value: "4", Type: ArgTypeInt},
								{Value: "", Type: ArgTypeUnknown},
							}}}},
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
			FuncName: "ValidString",
			PkgPath:  "unicode/utf8",
			ArgTypes: []Type{{Kind: TypeKindString}},
		},

		// for testing success
		"timecheck": RuleTypeFunc{
			FuncName: "TimeCheck", PkgPath: "path/to/rule",
			ArgTypes: []Type{{Kind: TypeKindStruct, Name: "Time", PkgPath: "time", PkgName: "time"}},
		},
		"ifacecheck": RuleTypeFunc{
			FuncName: "IfaceCheck", PkgPath: "path/to/rule",
			ArgTypes: []Type{{Kind: TypeKindInterface, IsEmptyInterface: true}},
		},

		// for testing errors
		"rulefunc1": RuleTypeFunc{
			FuncName: "RuleFunc1", PkgPath: "path/to/rule",
			ArgTypes: []Type{{Kind: TypeKindString}, {Kind: TypeKindInt}},
		},
		"rulefunc2": RuleTypeFunc{
			FuncName: "RuleFunc2", PkgPath: "path/to/rule",
			ArgTypes: []Type{{Kind: TypeKindString}, {Kind: TypeKindInt},
				{Kind: TypeKindSlice, Elem: &Type{Kind: TypeKindBool}}},
			IsVariadic: true,
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
