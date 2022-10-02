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
	F string `is:"uuid:6"`
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
