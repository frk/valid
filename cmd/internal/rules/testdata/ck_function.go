package testdata

import (
	"github.com/frk/valid/cmd/internal/rules/testdata/mypkg"
)

type Test_ERR_FUNCTION_INTYPE_1_Validator struct {
	F int `is:"contains:foo"`
}

type Test_ERR_FUNCTION_ARGTYPE_1_Validator struct {
	F string `is:"uuid:v6"`
}

type Test_ERR_FUNCTION_ARGVALUE_1_Validator struct {
	F string `is:"alpha:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_2_Validator struct {
	F string `is:"alnum:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_3_Validator struct {
	F string `is:"ccy:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_4_Validator struct {
	F string `is:"decimal:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_5_Validator struct {
	F string `is:"hash:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_6_Validator struct {
	F string `is:"ip:5"`
}

type Test_ERR_FUNCTION_ARGVALUE_7_Validator struct {
	F string `is:"isbn:12"`
}

type Test_ERR_FUNCTION_ARGVALUE_8_Validator struct {
	F string `is:"iso639:3"`
}

type Test_ERR_FUNCTION_ARGVALUE_9_Validator struct {
	F string `is:"iso31661a:1"`
}

type Test_ERR_FUNCTION_ARGVALUE_10_Validator struct {
	F string `is:"mac:7"`
}

type Test_ERR_FUNCTION_ARGVALUE_11_Validator struct {
	F string `is:"re:[0-9)?"`
}

type Test_ERR_FUNCTION_ARGVALUE_12_Validator struct {
	F string `is:"uuid:6"`
}

type Test_ERR_FUNCTION_ARGVALUE_13_Validator struct {
	F string `is:"phone:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_14_Validator struct {
	F string `is:"vat:foo"`
}

type Test_ERR_FUNCTION_ARGVALUE_15_Validator struct {
	F string `is:"zip:foo"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_function_Validator struct {
	Fn1 string `is:"contains:foo"`
	Fn2 string `is:"contains:foo:bar:baz"`

	UUID1 string `is:"uuid"`
	UUID2 string `is:"uuid:3"`
	UUID3 string `is:"uuid:v3"`
	UUID4 string `is:"uuid:4"`
	UUID5 string `is:"uuid:v4"`
	UUID6 string `is:"uuid:5"`
	UUID7 string `is:"uuid:v5"`

	R8 string `is:"r8:&helper"`

	helper *mypkg.CheckHelper
}
