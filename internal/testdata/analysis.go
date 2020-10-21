package testdata

type AnalysisTestOK_Validator struct {
	UserInput *UserInput
}

type AnalysisTestBAD_EmptyValidator struct {
	// ...
}

type AnalysisTestBAD_Empty2Validator struct {
	F string `is:"-"`
}

type AnalysisTestBAD_Empty3Validator struct {
	_ struct {
		F string `is:"required"`
	}
}

type AnalysisTestBAD_RuleArgNumRequiredValidator struct {
	F string `is:"required:@create:#group-key:#foobar"`
}

type AnalysisTestBAD_RuleArgKindRequiredValidator struct {
	F string `is:"required:@create:foobar"`
}

type AnalysisTestBAD_RuleArgKindConflictRequiredValidator struct {
	F string `is:"required:@create:@update"`
}

type AnalysisTestBAD_RuleArgKindConflict2RequiredValidator struct {
	F string `is:"required:#key1:#key2"`
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

type AnalysisTestBAD_RuleArgKindPhoneValidator struct {
	F string `is:"phone:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2PhoneValidator struct {
	F string `is:"phone:#key"`
}

type AnalysisTestBAD_TypeKindStringPhoneValidator struct {
	F []byte `is:"phone"`
}

type AnalysisTestBAD_RuleArgValueCountryCodePhoneValidator struct {
	F string `is:"phone:foo"`
}

type AnalysisTestBAD_RuleArgKindZipValidator struct {
	F string `is:"zip:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2ZipValidator struct {
	F string `is:"zip:#key"`
}

type AnalysisTestBAD_TypeKindStringZipValidator struct {
	F []byte `is:"zip"`
}

type AnalysisTestBAD_RuleArgValueCountryCodeZipValidator struct {
	F string `is:"zip:foo"`
}

type AnalysisTestBAD_RuleArgKindUUIDValidator struct {
	F string `is:"uuid:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2UUIDValidator struct {
	F string `is:"uuid:#key"`
}

type AnalysisTestBAD_TypeKindStringUUIDValidator struct {
	F []byte `is:"uuid"`
}

type AnalysisTestBAD_RuleArgValueUUIDVerUUIDValidator struct {
	F string `is:"uuid:foo"`
}

type AnalysisTestBAD_RuleArgKindIPValidator struct {
	F string `is:"ip:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2IPValidator struct {
	F string `is:"ip:#key"`
}

type AnalysisTestBAD_TypeKindStringIPValidator struct {
	F []byte `is:"ip"`
}

type AnalysisTestBAD_RuleArgValueIPVerIPValidator struct {
	F string `is:"ip:foo"`
}

type AnalysisTestBAD_RuleArgNumIPValidator struct {
	F string `is:"ip:v4:v6"`
}

type AnalysisTestBAD_RuleArgKindMACValidator struct {
	F string `is:"mac:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2MACValidator struct {
	F string `is:"mac:#key"`
}

type AnalysisTestBAD_TypeKindStringMACValidator struct {
	F []byte `is:"mac"`
}

type AnalysisTestBAD_RuleArgValueMACVerMACValidator struct {
	F string `is:"mac:foo"`
}

type AnalysisTestBAD_RuleArgNumMACValidator struct {
	F string `is:"mac:6:8"`
}

type AnalysisTestBAD_RuleArgKindISOValidator struct {
	F string `is:"iso:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2ISOValidator struct {
	F string `is:"iso:#key"`
}

type AnalysisTestBAD_TypeKindStringISOValidator struct {
	F []byte `is:"iso:1234"`
}

type AnalysisTestBAD_RuleArgValueISOStdISOValidator struct {
	F string `is:"iso:foo"`
}

type AnalysisTestBAD_RuleArgNumISOValidator struct {
	F string `is:"iso"`
}

type AnalysisTestBAD_RuleArgNum2ISOValidator struct {
	F string `is:"iso:6:8"`
}

type AnalysisTestBAD_RuleArgKindRFCValidator struct {
	F string `is:"rfc:@ctx"`
}

type AnalysisTestBAD_RuleArgKind2RFCValidator struct {
	F string `is:"rfc:#key"`
}

type AnalysisTestBAD_TypeKindStringRFCValidator struct {
	F []byte `is:"rfc:1234"`
}

type AnalysisTestBAD_RuleArgValueRFCStdRFCValidator struct {
	F string `is:"rfc:foo"`
}

type AnalysisTestBAD_RuleArgNumRFCValidator struct {
	F string `is:"rfc"`
}

type AnalysisTestBAD_RuleArgNum2RFCValidator struct {
	F string `is:"rfc:6:8"`
}
