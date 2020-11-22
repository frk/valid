package analysis

import (
	"github.com/frk/isvalid/internal/testdata/mypkg"
)

type AnalysisTestBAD_ErrorHandlerFieldConflictValidator struct {
	F string `is:"required"`
	mypkg.MyErrorConstructor
	mypkg.MyErrorAggregator
}

type AnalysisTestBAD_ErrorHandlerFieldConflict2Validator struct {
	F string `is:"required"`
	mypkg.MyErrorAggregator
	mypkg.MyErrorConstructor
}

type AnalysisTestBAD_ContextOptionFieldConflictValidator struct {
	F       string `is:"required"`
	Context string
	context string
}

type AnalysisTestBAD_ContextOptionFieldTypeValidator struct {
	F       string `is:"required"`
	Context int
}

type AnalysisTestBAD_ContextOptionFieldRequiredValidator struct {
	F string `is:"required:@ctx"`
}

type AnalysisTestBAD_RuleArgFieldUnknownValidator struct {
	F int64 `is:"gt:&x"`
}

type AnalysisTestBAD_ValidatorNoFieldValidator struct {
	// ...
}

type AnalysisTestBAD_ValidatorNoField2Validator struct {
	F string `is:"-"`
}

type AnalysisTestBAD_ValidatorNoField3Validator struct {
	_ struct {
		F string `is:"required"`
	}
}

type AnalysisTestBAD_RuleUnknownValidator struct {
	F string `is:"abracadabra:foobar"`
}

type AnalysisTestBAD_RuleUnknown2Validator struct {
	F string `is:"notempty"`
}

type AnalysisTestBAD_RuleArgNumRequiredValidator struct {
	F string `is:"required:foobar"`
}

type AnalysisTestBAD_RuleArgNumNotNilValidator struct {
	F []string `is:"notnil:foobar"`
}

type AnalysisTestBAD_TypeNilNotNilValidator struct {
	F [4]string `is:"notnil"`
}

type AnalysisTestBAD_RuleArgNumEmailValidator struct {
	F string `is:"email:foo"`
}

type AnalysisTestBAD_TypeKindStringEmailValidator struct {
	F int `is:"email"`
}

type AnalysisTestBAD_RuleArgNumURLValidator struct {
	F string `is:"url:foo"`
}

type AnalysisTestBAD_TypeKindStringURLValidator struct {
	F int `is:"url"`
}

type AnalysisTestBAD_RuleArgNumURIValidator struct {
	F string `is:"uri:foo"`
}

type AnalysisTestBAD_TypeKindStringURIValidator struct {
	F int `is:"uri"`
}

type AnalysisTestBAD_RuleArgNumPANValidator struct {
	F string `is:"pan:foo"`
}

type AnalysisTestBAD_TypeKindStringPANValidator struct {
	F bool `is:"pan"`
}

type AnalysisTestBAD_RuleArgNumCVVValidator struct {
	F string `is:"cvv:foo"`
}

type AnalysisTestBAD_TypeKindStringCVVValidator struct {
	F bool `is:"cvv"`
}

type AnalysisTestBAD_RuleArgNumSSNValidator struct {
	F string `is:"ssn:foo"`
}

type AnalysisTestBAD_TypeKindStringSSNValidator struct {
	F bool `is:"ssn"`
}

type AnalysisTestBAD_RuleArgNumEINValidator struct {
	F string `is:"ein:foo"`
}

type AnalysisTestBAD_TypeKindStringEINValidator struct {
	F bool `is:"ein"`
}

type AnalysisTestBAD_RuleArgNumNumericValidator struct {
	F string `is:"numeric:foo"`
}

type AnalysisTestBAD_TypeKindStringNumericValidator struct {
	F uint64 `is:"numeric"`
}

type AnalysisTestBAD_RuleArgNumHexValidator struct {
	F string `is:"hex:foo"`
}

type AnalysisTestBAD_TypeKindStringHexValidator struct {
	F uint64 `is:"hex"`
}

type AnalysisTestBAD_RuleArgNumHexcolorValidator struct {
	F string `is:"hexcolor:foo"`
}

type AnalysisTestBAD_TypeKindStringHexcolorValidator struct {
	F uint64 `is:"hexcolor"`
}

type AnalysisTestBAD_RuleArgNumAlphanumValidator struct {
	F string `is:"alphanum:foo"`
}

type AnalysisTestBAD_TypeKindStringAlphanumValidator struct {
	F uint64 `is:"alphanum"`
}

type AnalysisTestBAD_RuleArgNumCIDRValidator struct {
	F string `is:"cidr:foo"`
}

type AnalysisTestBAD_TypeKindStringCIDRValidator struct {
	F uint64 `is:"cidr"`
}

type AnalysisTestBAD_TypeKindStringPhoneValidator struct {
	F uint `is:"phone"`
}

type AnalysisTestBAD_RuleArgTypePhoneValidator struct {
	F string `is:"phone:321"`
}

type AnalysisTestBAD_RuleArgType2PhoneValidator struct {
	F string `is:"phone:true"`
}

type AnalysisTestBAD_RuleArgType3PhoneValidator struct {
	F string `is:"phone:0.2"`
}

type AnalysisTestBAD_RuleArgValueCountryCodePhoneValidator struct {
	F string `is:"phone:foo"`
}

type AnalysisTestBAD_RuleArgValueCountryCode2PhoneValidator struct {
	F string `is:"phone:ab"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindPhoneValidator struct {
	F string `is:"phone:&x"`
	x int
}

type AnalysisTestBAD_TypeKindStringZipValidator struct {
	F uint `is:"zip"`
}

type AnalysisTestBAD_RuleArgTypeZipValidator struct {
	F string `is:"zip:321"`
}

type AnalysisTestBAD_RuleArgType2ZipValidator struct {
	F string `is:"zip:true"`
}

type AnalysisTestBAD_RuleArgType3ZipValidator struct {
	F string `is:"zip:0.2"`
}

type AnalysisTestBAD_RuleArgValueCountryCodeZipValidator struct {
	F string `is:"zip:foo"`
}

type AnalysisTestBAD_RuleArgValueCountryCode2ZipValidator struct {
	F string `is:"zip:ab"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindZipValidator struct {
	F string `is:"zip:&x"`
	x uint
}

type AnalysisTestBAD_TypeKindStringUUIDValidator struct {
	F uint16 `is:"uuid"`
}

type AnalysisTestBAD_RuleArgTypeUUIDValidator struct {
	F string `is:"uuid:-4"`
}

type AnalysisTestBAD_RuleArgType2UUIDValidator struct {
	F string `is:"uuid:true"`
}

type AnalysisTestBAD_RuleArgType3UUIDValidator struct {
	F string `is:"uuid:0.2"`
}

type AnalysisTestBAD_RuleArgValueUUIDVerUUIDValidator struct {
	F string `is:"uuid:foo"`
}

type AnalysisTestBAD_RuleArgValueUUIDVer2UUIDValidator struct {
	F string `is:"uuid:v8"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindUUIDValidator struct {
	F string `is:"uuid:&x:&y:&z"`
	x int32
	y uint16
	z []byte
}

type AnalysisTestBAD_RuleArgValueConflictUUIDValidator struct {
	F string `is:"uuid:v4:5:4"`
}

type AnalysisTestBAD_RuleArgNumUUIDValidator struct {
	F string `is:"uuid:1:2:3:4:5:6"`
}

type AnalysisTestBAD_TypeKindStringIPValidator struct {
	F uint8 `is:"ip"`
}

type AnalysisTestBAD_RuleArgTypeIPValidator struct {
	F string `is:"ip:-4"`
}

type AnalysisTestBAD_RuleArgType2IPValidator struct {
	F string `is:"ip:true"`
}

type AnalysisTestBAD_RuleArgType3IPValidator struct {
	F string `is:"ip:0.2"`
}

type AnalysisTestBAD_RuleArgValueIPVerIPValidator struct {
	F string `is:"ip:v7"`
}

type AnalysisTestBAD_RuleArgValueIPVer2IPValidator struct {
	F string `is:"ip:foo"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindIPValidator struct {
	F string `is:"ip:&x:&y"`
	x string
	y float32
}

type AnalysisTestBAD_RuleArgValueConflictIPValidator struct {
	F string `is:"ip:v4:4"`
}

type AnalysisTestBAD_RuleArgNumIPValidator struct {
	F string `is:"ip:v4:v6:v8"`
}

type AnalysisTestBAD_TypeKindStringMACValidator struct {
	F uint32 `is:"mac"`
}

type AnalysisTestBAD_RuleArgTypeMACValidator struct {
	F string `is:"mac:-6"`
}

type AnalysisTestBAD_RuleArgType2MACValidator struct {
	F string `is:"mac:true"`
}

type AnalysisTestBAD_RuleArgType3MACValidator struct {
	F string `is:"mac:0.2"`
}

type AnalysisTestBAD_RuleArgValueMACVerMACValidator struct {
	F string `is:"mac:v10"`
}

type AnalysisTestBAD_RuleArgValueMACVer2MACValidator struct {
	F string `is:"mac:vv8"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindMACValidator struct {
	F string `is:"mac:&x:&y"`
	x string
	y float32
}

type AnalysisTestBAD_RuleArgValueConflictMACValidator struct {
	F string `is:"mac:v6:6"`
}

type AnalysisTestBAD_RuleArgNumMACValidator struct {
	F string `is:"mac:6:8:10"`
}

type AnalysisTestBAD_TypeKindStringISOValidator struct {
	F complex128 `is:"iso:1234"`
}

type AnalysisTestBAD_RuleArgTypeISOValidator struct {
	F string `is:"iso:foo"`
}

type AnalysisTestBAD_RuleArgType2ISOValidator struct {
	F string `is:"iso:true"`
}

type AnalysisTestBAD_RuleArgType3ISOValidator struct {
	F string `is:"iso:0.2"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindISOValidator struct {
	F string `is:"iso:&x"`
	x bool
}

type AnalysisTestBAD_RuleArgNumISOValidator struct {
	F string `is:"iso"`
}

type AnalysisTestBAD_RuleArgNum2ISOValidator struct {
	F string `is:"iso:6:8"`
}

type AnalysisTestBAD_TypeKindStringRFCValidator struct {
	F uint64 `is:"rfc:1234"`
}

type AnalysisTestBAD_RuleArgTypeRFCValidator struct {
	F string `is:"rfc:foo"`
}

type AnalysisTestBAD_RuleArgType2RFCValidator struct {
	F string `is:"rfc:true"`
}

type AnalysisTestBAD_RuleArgType3RFCValidator struct {
	F string `is:"rfc:0.2"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindRFCValidator struct {
	F string `is:"rfc:&x"`
	x bool
}

type AnalysisTestBAD_RuleArgNumRFCValidator struct {
	F string `is:"rfc"`
}

type AnalysisTestBAD_RuleArgNum2RFCValidator struct {
	F string `is:"rfc:6:8"`
}

type AnalysisTestBAD_TypeKindStringRegexpValidator struct {
	F uint64 `is:"re:abc"`
}

type AnalysisTestBAD_RuleArgValueRegexpRegexpValidator struct {
	F string `is:"re:^($"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindRegexpValidator struct {
	F string `is:"re:&x"`
	x uint32
}

type AnalysisTestBAD_RuleArgNumRegexpValidator struct {
	F string `is:"re"`
}

type AnalysisTestBAD_RuleArgNum2RegexpValidator struct {
	F string `is:"re:foo:bar"`
}

type AnalysisTestBAD_TypeKindStringPrefixValidator struct {
	F uint64 `is:"prefix:foo"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindPrefixValidator struct {
	F string `is:"prefix:&x:&y"`
	x string
	y uint32
}

type AnalysisTestBAD_RuleArgNumPrefixValidator struct {
	F string `is:"prefix"`
}

type AnalysisTestBAD_TypeKindStringSuffixValidator struct {
	F uint64 `is:"suffix:foo"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindSuffixValidator struct {
	F string `is:"suffix:&x:&y"`
	x string
	y uint32
}

type AnalysisTestBAD_RuleArgNumSuffixValidator struct {
	F string `is:"suffix"`
}

type AnalysisTestBAD_TypeKindStringContainsValidator struct {
	F uint64 `is:"contains:foo"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindContainsValidator struct {
	F string `is:"contains:&x:&y"`
	x string
	y uint32
}

type AnalysisTestBAD_RuleArgNumContainsValidator struct {
	F string `is:"contains"`
}

type AnalysisTestBAD_RuleArgNumEQValidator struct {
	F int `is:"eq"`
}

type AnalysisTestBAD_RuleArgTypeStringEQValidator struct {
	F int `is:"eq:foo"`
}

type AnalysisTestBAD_RuleArgTypeNintEQValidator struct {
	F uint `is:"eq:-123"`
}

type AnalysisTestBAD_RuleArgTypeUintEQValidator struct {
	F []byte `is:"eq:123"`
}

type AnalysisTestBAD_RuleArgTypeFloatEQValidator struct {
	F int32 `is:"eq:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindEQValidator struct {
	F int32 `is:"eq:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumNEValidator struct {
	F int `is:"ne"`
}

type AnalysisTestBAD_RuleArgTypeStringNEValidator struct {
	F int `is:"ne:foo"`
}

type AnalysisTestBAD_RuleArgTypeNintNEValidator struct {
	F uint `is:"ne:-123"`
}

type AnalysisTestBAD_RuleArgTypeUintNEValidator struct {
	F []byte `is:"ne:123"`
}

type AnalysisTestBAD_RuleArgTypeFloatNEValidator struct {
	F int32 `is:"ne:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindNEValidator struct {
	F int32 `is:"ne:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumGTValidator struct {
	F uint64 `is:"gt"`
}

type AnalysisTestBAD_RuleArgNum2GTValidator struct {
	F uint64 `is:"gt:1:2"`
}

type AnalysisTestBAD_TypeNumericGTValidator struct {
	F string `is:"gt:123"`
}

type AnalysisTestBAD_RuleArgTypeNintGTValidator struct {
	F uint `is:"gt:-123"`
}

type AnalysisTestBAD_RuleArgTypeFloatGTValidator struct {
	F int32 `is:"gt:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindGTValidator struct {
	F int32 `is:"gt:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumLTValidator struct {
	F uint64 `is:"lt"`
}

type AnalysisTestBAD_RuleArgNum2LTValidator struct {
	F uint64 `is:"lt:1:2"`
}

type AnalysisTestBAD_TypeNumericLTValidator struct {
	F string `is:"lt:123"`
}

type AnalysisTestBAD_RuleArgTypeNintLTValidator struct {
	F uint `is:"lt:-123"`
}

type AnalysisTestBAD_RuleArgTypeFloatLTValidator struct {
	F int32 `is:"lt:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindLTValidator struct {
	F int32 `is:"lt:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumGTEValidator struct {
	F uint64 `is:"gte"`
}

type AnalysisTestBAD_RuleArgNum2GTEValidator struct {
	F uint64 `is:"gte:1:2"`
}

type AnalysisTestBAD_TypeNumericGTEValidator struct {
	F string `is:"gte:123"`
}

type AnalysisTestBAD_RuleArgTypeNintGTEValidator struct {
	F uint `is:"gte:-123"`
}

type AnalysisTestBAD_RuleArgTypeFloatGTEValidator struct {
	F int32 `is:"gte:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindGTEValidator struct {
	F int32 `is:"gte:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumLTEValidator struct {
	F uint64 `is:"lte"`
}

type AnalysisTestBAD_RuleArgNum2LTEValidator struct {
	F uint64 `is:"lte:1:2"`
}

type AnalysisTestBAD_TypeNumericLTEValidator struct {
	F string `is:"lte:123"`
}

type AnalysisTestBAD_RuleArgTypeNintLTEValidator struct {
	F uint `is:"lte:-123"`
}

type AnalysisTestBAD_RuleArgTypeFloatLTEValidator struct {
	F int32 `is:"lte:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindLTEValidator struct {
	F int32 `is:"lte:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumMinValidator struct {
	F uint64 `is:"min"`
}

type AnalysisTestBAD_RuleArgNum2MinValidator struct {
	F uint64 `is:"min:1:2"`
}

type AnalysisTestBAD_TypeNumericMinValidator struct {
	F string `is:"min:123"`
}

type AnalysisTestBAD_RuleArgTypeNintMinValidator struct {
	F uint `is:"min:-123"`
}

type AnalysisTestBAD_RuleArgTypeFloatMinValidator struct {
	F int32 `is:"min:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindMinValidator struct {
	F int32 `is:"min:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumMaxValidator struct {
	F uint64 `is:"max"`
}

type AnalysisTestBAD_RuleArgNum2MaxValidator struct {
	F uint64 `is:"max:1:2"`
}

type AnalysisTestBAD_TypeNumericMaxValidator struct {
	F string `is:"max:123"`
}

type AnalysisTestBAD_RuleArgTypeNintMaxValidator struct {
	F uint `is:"max:-123"`
}

type AnalysisTestBAD_RuleArgTypeFloatMaxValidator struct {
	F int32 `is:"max:1.23"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindMaxValidator struct {
	F int32 `is:"max:&x"`
	x string
}

type AnalysisTestBAD_RuleArgNumRngValidator struct {
	F uint64 `is:"rng"`
}

type AnalysisTestBAD_RuleArgNum2RngValidator struct {
	F uint64 `is:"rng:123"`
}

type AnalysisTestBAD_RuleArgNum3RngValidator struct {
	F uint64 `is:"rng:1:2:3"`
}

type AnalysisTestBAD_TypeNumericRngValidator struct {
	F string `is:"rng:1:23"`
}

type AnalysisTestBAD_RuleArgTypeStringRngValidator struct {
	F uint64 `is:"rng:foo:bar"`
}

type AnalysisTestBAD_RuleArgTypeString2RngValidator struct {
	F uint64 `is:"rng:123:bar"`
}

type AnalysisTestBAD_RuleArgTypeNintRngValidator struct {
	F uint `is:"rng:-123:0"`
}

type AnalysisTestBAD_RuleArgTypeFloatRngValidator struct {
	F int32 `is:"rng::1.23"`
}

type AnalysisTestBAD_RuleArgValueBoundsRngValidator struct {
	F float32 `is:"rng:2:1.23"`
}

type AnalysisTestBAD_RuleArgValueBounds2RngValidator struct {
	F float32 `is:"rng::"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindRngValidator struct {
	F int32 `is:"rng:&x:&y"`
	x int32
	y string
}

type AnalysisTestBAD_RuleArgNumLenValidator struct {
	F []string `is:"len"`
}

type AnalysisTestBAD_RuleArgNum2LenValidator struct {
	F []string `is:"len:1:2:3"`
}

type AnalysisTestBAD_TypeLengthLenValidator struct {
	F int32 `is:"len:1:23"`
}

type AnalysisTestBAD_RuleArgTypeLenValidator struct {
	F []string `is:"len:foo:bar"`
}

type AnalysisTestBAD_RuleArgType2LenValidator struct {
	F map[string]struct{} `is:"len:-123:0"`
}

type AnalysisTestBAD_RuleArgType3LenValidator struct {
	F []string `is:"len::1.23"`
}

type AnalysisTestBAD_RuleArgValueBoundsLenValidator struct {
	F []string `is:"len:20:10"`
}

type AnalysisTestBAD_RuleArgValueBounds2LenValidator struct {
	F []string `is:"len::"`
}

type AnalysisTestBAD_RuleArgTypeReferenceKindLenValidator struct {
	F []int8 `is:"len:&x:&y"`
	x int32
	y string
}

type AnalysisTestBAD_RuleFuncRuleArgCountValidator struct {
	F string `is:"rulefunc1"`
}

type AnalysisTestBAD_RuleFuncRuleArgCount2Validator struct {
	F string `is:"rulefunc1:a:b:c"`
}

type AnalysisTestBAD_RuleFuncFieldArgTypeValidator struct {
	F int `is:"rulefunc1:123"`
}

type AnalysisTestBAD_RuleFuncRuleArgTypeValidator struct {
	F string `is:"rulefunc1:foo"`
}

type AnalysisTestBAD_RuleFuncRuleArgType2Validator struct {
	F string `is:"rulefunc2:123:true:false:abc"`
}

type AnalysisTestBAD_RuleEnumTypeUnnamedValidator struct {
	F string `is:"enum"`
}

type AnalysisTestBAD_RuleEnumTypeUnnamed2Validator struct {
	F uint32 `is:"enum"`
}

type AnalysisTestBAD_RuleEnumTypeValidator struct {
	F mypkg.MyNonBasicType `is:"enum"`
}

type AnalysisTestBAD_RuleEnumType2Validator struct {
	F mypkg.MyNonBasicType2 `is:"enum"`
}

type AnalysisTestBAD_RuleEnumTypeNoConstValidator struct {
	F mypkg.MyNoConst `is:"enum"`
}

type AnalysisTestBAD_RuleKeyValidator struct {
	F string `is:"[email]"`
}

type AnalysisTestBAD_RuleKey2Validator struct {
	F []string `is:"[email]"`
}

type AnalysisTestBAD_RuleKey3Validator struct {
	F map[string][]string `is:"[][email]"`
}

type AnalysisTestBAD_RuleElemValidator struct {
	F string `is:"[]email"`
}

type AnalysisTestBAD_RuleElem2Validator struct {
	F []string `is:"[][]email"`
}

type AnalysisTestBAD_RuleElem3Validator struct {
	F []map[string]string `is:"[][]email,[]email"`
}

type AnalysisTestBAD_RuleArgCountKeyValidator struct {
	F []map[string]string `is:"[][email:foo]"`
}

type AnalysisTestBAD_RuleArgCountElemValidator struct {
	F []string `is:"[]email:foo"`
}

type AnalysisTestBAD_RuleArgCountSubfieldValidator struct {
	F map[string]struct {
		F string `is:"email:foo"`
	}
}
