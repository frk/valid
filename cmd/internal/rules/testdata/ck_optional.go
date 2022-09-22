package testdata

type Test_ERR_OPTIONAL_CONFLICT_1_Validator struct {
	F *string `is:"optional,required"`
}

type Test_ERR_OPTIONAL_CONFLICT_2_Validator struct {
	F *string `is:"required,optional"`
}

type Test_ERR_OPTIONAL_CONFLICT_3_Validator struct {
	F *string `is:"notnil,optional"`
}

type Test_ERR_OPTIONAL_CONFLICT_4_Validator struct {
	F *string `is:"omitnil,notnil"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_optional_Validator struct {
	F1 *string `is:"optional"`
	F2 *string `is:"omitnil"`
	F3 *string `is:"required"`
	F4 *string `is:"notnil"`
	F5 *string `is:"omitnil,required"`
}
