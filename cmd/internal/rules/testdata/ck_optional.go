package testdata

type Test_ERR_OPTIONAL_CONFLICT_1_Validator struct {
	F *string `is:"optional,required"`
}

type Test_ERR_OPTIONAL_CONFLICT_2_Validator struct {
	F *string `is:"required,optional"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_optional_Validator struct {
	F1 *string `is:"optional"`
	F2 *string `is:"omitnil"`
	F3 *string `is:"notnil,omitnil"`
}
