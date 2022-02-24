package testdata

type Test_ERR_PREPROC_INTYPE_1_Validator struct {
	F string `pre:"p1"`
}

type Test_ERR_PREPROC_OUTTYPE_1_Validator struct {
	F string `pre:"p2"`
}

type Test_ERR_PREPROC_ARGTYPE_1_Validator struct {
	F string `pre:"p4:foo"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_preproc_Validator struct {
	F1 string `pre:"trim"`
}
