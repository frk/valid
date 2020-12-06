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

type AnalysisTestBAD_RuleOptionFieldUnknownValidator struct {
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

type AnalysisTestBAD_RuleOptionNumRequiredValidator struct {
	F string `is:"required:foobar"`
}

type AnalysisTestBAD_RuleOptionNumNotNilValidator struct {
	F []string `is:"notnil:foobar"`
}

type AnalysisTestBAD_TypeNilNotNilValidator struct {
	F [4]string `is:"notnil"`
}

type AnalysisTestBAD_RuleOptionNumEmailValidator struct {
	F string `is:"email:foo"`
}

type AnalysisTestBAD_TypeKindStringEmailValidator struct {
	F int `is:"email"`
}

type AnalysisTestBAD_RuleOptionNumURLValidator struct {
	F string `is:"url:foo"`
}

type AnalysisTestBAD_TypeKindStringURLValidator struct {
	F int `is:"url"`
}

type AnalysisTestBAD_RuleOptionNumPANValidator struct {
	F string `is:"pan:foo"`
}

type AnalysisTestBAD_TypeKindStringPANValidator struct {
	F bool `is:"pan"`
}

type AnalysisTestBAD_RuleOptionNumCVVValidator struct {
	F string `is:"cvv:foo"`
}

type AnalysisTestBAD_TypeKindStringCVVValidator struct {
	F bool `is:"cvv"`
}

type AnalysisTestBAD_RuleOptionNumSSNValidator struct {
	F string `is:"ssn:foo"`
}

type AnalysisTestBAD_TypeKindStringSSNValidator struct {
	F bool `is:"ssn"`
}

type AnalysisTestBAD_RuleOptionNumEINValidator struct {
	F string `is:"ein:foo"`
}

type AnalysisTestBAD_TypeKindStringEINValidator struct {
	F bool `is:"ein"`
}

type AnalysisTestBAD_RuleOptionNumNumericValidator struct {
	F string `is:"numeric:foo"`
}

type AnalysisTestBAD_TypeKindStringNumericValidator struct {
	F uint64 `is:"numeric"`
}

type AnalysisTestBAD_RuleOptionNumHexValidator struct {
	F string `is:"hex:foo"`
}

type AnalysisTestBAD_TypeKindStringHexValidator struct {
	F uint64 `is:"hex"`
}

type AnalysisTestBAD_RuleOptionNumHexcolorValidator struct {
	F string `is:"hexcolor:foo"`
}

type AnalysisTestBAD_TypeKindStringHexcolorValidator struct {
	F uint64 `is:"hexcolor"`
}

type AnalysisTestBAD_RuleOptionNumAlnumValidator struct {
	F string `is:"alnum:foo"`
}

type AnalysisTestBAD_TypeKindStringAlnumValidator struct {
	F uint64 `is:"alnum"`
}

type AnalysisTestBAD_RuleOptionNumCIDRValidator struct {
	F string `is:"cidr:foo"`
}

type AnalysisTestBAD_TypeKindStringCIDRValidator struct {
	F uint64 `is:"cidr"`
}

type AnalysisTestBAD_TypeKindStringPhoneValidator struct {
	F uint `is:"phone"`
}

type AnalysisTestBAD_RuleOptionTypePhoneValidator struct {
	F string `is:"phone:321"`
}

type AnalysisTestBAD_RuleOptionType2PhoneValidator struct {
	F string `is:"phone:true"`
}

type AnalysisTestBAD_RuleOptionType3PhoneValidator struct {
	F string `is:"phone:0.2"`
}

type AnalysisTestBAD_RuleOptionValueCountryCodePhoneValidator struct {
	F string `is:"phone:foo"`
}

type AnalysisTestBAD_RuleOptionValueCountryCode2PhoneValidator struct {
	F string `is:"phone:ab"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindPhoneValidator struct {
	F string `is:"phone:&x"`
	x int
}

type AnalysisTestBAD_TypeKindStringZipValidator struct {
	F uint `is:"zip"`
}

type AnalysisTestBAD_RuleOptionTypeZipValidator struct {
	F string `is:"zip:321"`
}

type AnalysisTestBAD_RuleOptionType2ZipValidator struct {
	F string `is:"zip:true"`
}

type AnalysisTestBAD_RuleOptionType3ZipValidator struct {
	F string `is:"zip:0.2"`
}

type AnalysisTestBAD_RuleOptionValueCountryCodeZipValidator struct {
	F string `is:"zip:foo"`
}

type AnalysisTestBAD_RuleOptionValueCountryCode2ZipValidator struct {
	F string `is:"zip:ab"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindZipValidator struct {
	F string `is:"zip:&x"`
	x uint
}

type AnalysisTestBAD_TypeKindStringUUIDValidator struct {
	F uint16 `is:"uuid"`
}

type AnalysisTestBAD_RuleOptionTypeUUIDValidator struct {
	F string `is:"uuid:-4"`
}

type AnalysisTestBAD_RuleOptionType2UUIDValidator struct {
	F string `is:"uuid:true"`
}

type AnalysisTestBAD_RuleOptionType3UUIDValidator struct {
	F string `is:"uuid:0.2"`
}

type AnalysisTestBAD_RuleOptionValueUUIDVerUUIDValidator struct {
	F string `is:"uuid:foo"`
}

type AnalysisTestBAD_RuleOptionValueUUIDVer2UUIDValidator struct {
	F string `is:"uuid:v8"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindUUIDValidator struct {
	F string `is:"uuid:&z"`
	z []byte
}

type AnalysisTestBAD_RuleOptionNumUUIDValidator struct {
	F string `is:"uuid:1:2:3:4:5:6"`
}

type AnalysisTestBAD_TypeKindStringIPValidator struct {
	F uint8 `is:"ip"`
}

type AnalysisTestBAD_RuleOptionTypeIPValidator struct {
	F string `is:"ip:-4"`
}

type AnalysisTestBAD_RuleOptionType2IPValidator struct {
	F string `is:"ip:true"`
}

type AnalysisTestBAD_RuleOptionType3IPValidator struct {
	F string `is:"ip:0.2"`
}

type AnalysisTestBAD_RuleOptionValueIPVerIPValidator struct {
	F string `is:"ip:v7"`
}

type AnalysisTestBAD_RuleOptionValueIPVer2IPValidator struct {
	F string `is:"ip:foo"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindIPValidator struct {
	F string `is:"ip:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumIPValidator struct {
	F string `is:"ip:v4:v6:v8"`
}

type AnalysisTestBAD_TypeKindStringMACValidator struct {
	F uint32 `is:"mac"`
}

type AnalysisTestBAD_RuleOptionTypeMACValidator struct {
	F string `is:"mac:-6"`
}

type AnalysisTestBAD_RuleOptionType2MACValidator struct {
	F string `is:"mac:true"`
}

type AnalysisTestBAD_RuleOptionType3MACValidator struct {
	F string `is:"mac:0.2"`
}

type AnalysisTestBAD_RuleOptionValueMACVerMACValidator struct {
	F string `is:"mac:v8"`
}

type AnalysisTestBAD_RuleOptionValueMACVer2MACValidator struct {
	F string `is:"mac:vv8"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindMACValidator struct {
	F string `is:"mac:&x:&y"`
	x string
	y float32
}

type AnalysisTestBAD_RuleOptionValueConflictMACValidator struct {
	F string `is:"mac:v6:6"`
}

type AnalysisTestBAD_RuleOptionNumMACValidator struct {
	F string `is:"mac:6:8:10"`
}

type AnalysisTestBAD_TypeKindStringRegexpValidator struct {
	F uint64 `is:"re:abc"`
}

type AnalysisTestBAD_RuleOptionValueRegexpRegexpValidator struct {
	F string `is:"re:^($"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindRegexpValidator struct {
	F string `is:"re:&x"`
	x uint32
}

type AnalysisTestBAD_RuleOptionNumRegexpValidator struct {
	F string `is:"re"`
}

type AnalysisTestBAD_RuleOptionNum2RegexpValidator struct {
	F string `is:"re:foo:bar"`
}

type AnalysisTestBAD_TypeKindStringPrefixValidator struct {
	F uint64 `is:"prefix:foo"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindPrefixValidator struct {
	F string `is:"prefix:&x:&y"`
	x string
	y uint32
}

type AnalysisTestBAD_RuleOptionNumPrefixValidator struct {
	F string `is:"prefix"`
}

type AnalysisTestBAD_TypeKindStringSuffixValidator struct {
	F uint64 `is:"suffix:foo"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindSuffixValidator struct {
	F string `is:"suffix:&x:&y"`
	x string
	y uint32
}

type AnalysisTestBAD_RuleOptionNumSuffixValidator struct {
	F string `is:"suffix"`
}

type AnalysisTestBAD_TypeKindStringContainsValidator struct {
	F uint64 `is:"contains:foo"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindContainsValidator struct {
	F string `is:"contains:&x:&y"`
	x string
	y uint32
}

type AnalysisTestBAD_RuleOptionNumContainsValidator struct {
	F string `is:"contains"`
}

type AnalysisTestBAD_RuleOptionNumEQValidator struct {
	F int `is:"eq"`
}

type AnalysisTestBAD_RuleOptionTypeStringEQValidator struct {
	F int `is:"eq:foo"`
}

type AnalysisTestBAD_RuleOptionTypeNintEQValidator struct {
	F uint `is:"eq:-123"`
}

type AnalysisTestBAD_RuleOptionTypeUintEQValidator struct {
	F []byte `is:"eq:123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatEQValidator struct {
	F int32 `is:"eq:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindEQValidator struct {
	F int32 `is:"eq:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumNEValidator struct {
	F int `is:"ne"`
}

type AnalysisTestBAD_RuleOptionTypeStringNEValidator struct {
	F int `is:"ne:foo"`
}

type AnalysisTestBAD_RuleOptionTypeNintNEValidator struct {
	F uint `is:"ne:-123"`
}

type AnalysisTestBAD_RuleOptionTypeUintNEValidator struct {
	F []byte `is:"ne:123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatNEValidator struct {
	F int32 `is:"ne:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindNEValidator struct {
	F int32 `is:"ne:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumGTValidator struct {
	F uint64 `is:"gt"`
}

type AnalysisTestBAD_RuleOptionNum2GTValidator struct {
	F uint64 `is:"gt:1:2"`
}

type AnalysisTestBAD_TypeNumericGTValidator struct {
	F string `is:"gt:123"`
}

type AnalysisTestBAD_RuleOptionTypeNintGTValidator struct {
	F uint `is:"gt:-123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatGTValidator struct {
	F int32 `is:"gt:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindGTValidator struct {
	F int32 `is:"gt:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumLTValidator struct {
	F uint64 `is:"lt"`
}

type AnalysisTestBAD_RuleOptionNum2LTValidator struct {
	F uint64 `is:"lt:1:2"`
}

type AnalysisTestBAD_TypeNumericLTValidator struct {
	F string `is:"lt:123"`
}

type AnalysisTestBAD_RuleOptionTypeNintLTValidator struct {
	F uint `is:"lt:-123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatLTValidator struct {
	F int32 `is:"lt:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindLTValidator struct {
	F int32 `is:"lt:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumGTEValidator struct {
	F uint64 `is:"gte"`
}

type AnalysisTestBAD_RuleOptionNum2GTEValidator struct {
	F uint64 `is:"gte:1:2"`
}

type AnalysisTestBAD_TypeNumericGTEValidator struct {
	F string `is:"gte:123"`
}

type AnalysisTestBAD_RuleOptionTypeNintGTEValidator struct {
	F uint `is:"gte:-123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatGTEValidator struct {
	F int32 `is:"gte:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindGTEValidator struct {
	F int32 `is:"gte:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumLTEValidator struct {
	F uint64 `is:"lte"`
}

type AnalysisTestBAD_RuleOptionNum2LTEValidator struct {
	F uint64 `is:"lte:1:2"`
}

type AnalysisTestBAD_TypeNumericLTEValidator struct {
	F string `is:"lte:123"`
}

type AnalysisTestBAD_RuleOptionTypeNintLTEValidator struct {
	F uint `is:"lte:-123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatLTEValidator struct {
	F int32 `is:"lte:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindLTEValidator struct {
	F int32 `is:"lte:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumMinValidator struct {
	F uint64 `is:"min"`
}

type AnalysisTestBAD_RuleOptionNum2MinValidator struct {
	F uint64 `is:"min:1:2"`
}

type AnalysisTestBAD_TypeNumericMinValidator struct {
	F string `is:"min:123"`
}

type AnalysisTestBAD_RuleOptionTypeNintMinValidator struct {
	F uint `is:"min:-123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatMinValidator struct {
	F int32 `is:"min:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindMinValidator struct {
	F int32 `is:"min:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumMaxValidator struct {
	F uint64 `is:"max"`
}

type AnalysisTestBAD_RuleOptionNum2MaxValidator struct {
	F uint64 `is:"max:1:2"`
}

type AnalysisTestBAD_TypeNumericMaxValidator struct {
	F string `is:"max:123"`
}

type AnalysisTestBAD_RuleOptionTypeNintMaxValidator struct {
	F uint `is:"max:-123"`
}

type AnalysisTestBAD_RuleOptionTypeFloatMaxValidator struct {
	F int32 `is:"max:1.23"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindMaxValidator struct {
	F int32 `is:"max:&x"`
	x string
}

type AnalysisTestBAD_RuleOptionNumRngValidator struct {
	F uint64 `is:"rng"`
}

type AnalysisTestBAD_RuleOptionNum2RngValidator struct {
	F uint64 `is:"rng:123"`
}

type AnalysisTestBAD_RuleOptionNum3RngValidator struct {
	F uint64 `is:"rng:1:2:3"`
}

type AnalysisTestBAD_TypeNumericRngValidator struct {
	F string `is:"rng:1:23"`
}

type AnalysisTestBAD_RuleOptionTypeStringRngValidator struct {
	F uint64 `is:"rng:foo:bar"`
}

type AnalysisTestBAD_RuleOptionTypeString2RngValidator struct {
	F uint64 `is:"rng:123:bar"`
}

type AnalysisTestBAD_RuleOptionTypeNintRngValidator struct {
	F uint `is:"rng:-123:0"`
}

type AnalysisTestBAD_RuleOptionTypeFloatRngValidator struct {
	F int32 `is:"rng::1.23"`
}

type AnalysisTestBAD_RuleOptionValueBoundsRngValidator struct {
	F float32 `is:"rng:2:1.23"`
}

type AnalysisTestBAD_RuleOptionValueBounds2RngValidator struct {
	F float32 `is:"rng::"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindRngValidator struct {
	F int32 `is:"rng:&x:&y"`
	x int32
	y string
}

type AnalysisTestBAD_RuleOptionNumLenValidator struct {
	F []string `is:"len"`
}

type AnalysisTestBAD_RuleOptionNum2LenValidator struct {
	F []string `is:"len:1:2:3"`
}

type AnalysisTestBAD_TypeLengthLenValidator struct {
	F int32 `is:"len:1:23"`
}

type AnalysisTestBAD_RuleOptionTypeLenValidator struct {
	F []string `is:"len:foo:bar"`
}

type AnalysisTestBAD_RuleOptionType2LenValidator struct {
	F map[string]struct{} `is:"len:-123:0"`
}

type AnalysisTestBAD_RuleOptionType3LenValidator struct {
	F []string `is:"len::1.23"`
}

type AnalysisTestBAD_RuleOptionValueBoundsLenValidator struct {
	F []string `is:"len:20:10"`
}

type AnalysisTestBAD_RuleOptionValueBounds2LenValidator struct {
	F []string `is:"len::"`
}

type AnalysisTestBAD_RuleOptionTypeReferenceKindLenValidator struct {
	F []int8 `is:"len:&x:&y"`
	x int32
	y string
}

type AnalysisTestBAD_RuleOptionNumRuneCountValidator struct {
	F string `is:"runecount"`
}

type AnalysisTestBAD_RuleOptionNum2RuneCountValidator struct {
	F []byte `is:"runecount:1:2:3"`
}

type AnalysisTestBAD_TypeRunelessRuneCountValidator struct {
	F map[string]rune `is:"runecount:1:23"`
}

type AnalysisTestBAD_TypeRuneless2RuneCountValidator struct {
	F []rune `is:"runecount:1:23"` // should use len
}

type AnalysisTestBAD_TypeRuneless3RuneCountValidator struct {
	F [12]byte `is:"runecount:1:23"`
}

type AnalysisTestBAD_RuleOptionTypeRuneCountValidator struct {
	F string `is:"runecount:foo:bar"`
}

type AnalysisTestBAD_RuleOptionType2RuneCountValidator struct {
	F []byte `is:"runecount:-123:0"`
}

type AnalysisTestBAD_RuleOptionType3RuneCountValidator struct {
	F string `is:"runecount::1.23"`
}

type AnalysisTestBAD_RuleOptionValueBoundsRuneCountValidator struct {
	F []byte `is:"runecount:20:10"`
}

type AnalysisTestBAD_RuleOptionValueBounds2RuneCountValidator struct {
	F string `is:"runecount::"`
}

type AnalysisTestBAD_RuleOptionTypeFieldKindRuneCountValidator struct {
	F []byte `is:"runecount:&x:&y"`
	x int32
	y string
}

type AnalysisTestBAD_RuleFuncRuleOptionCountValidator struct {
	F string `is:"rulefunc1"`
}

type AnalysisTestBAD_RuleFuncRuleOptionCount2Validator struct {
	F string `is:"rulefunc1:a:b:c"`
}

type AnalysisTestBAD_RuleFuncFieldOptionTypeValidator struct {
	F int `is:"rulefunc1:123"`
}

type AnalysisTestBAD_RuleFuncRuleOptionTypeValidator struct {
	F string `is:"rulefunc1:foo"`
}

type AnalysisTestBAD_RuleFuncRuleOptionType2Validator struct {
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

type AnalysisTestBAD_RuleOptionCountKeyValidator struct {
	F []map[string]string `is:"[][email:foo]"`
}

type AnalysisTestBAD_RuleOptionCountElemValidator struct {
	F []string `is:"[]email:foo"`
}

type AnalysisTestBAD_RuleOptionCountSubfieldValidator struct {
	F map[string]struct {
		F string `is:"email:foo"`
	}
}
