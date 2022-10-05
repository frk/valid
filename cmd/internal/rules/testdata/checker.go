package testdata

type Test_ERR_FIELD_UNKNOWN_1_Validator struct {
	F int `is:"gt:&num"`
}

type Test_ERR_FIELD_UNKNOWN_2_Validator struct {
	F string `pre:"p4:&num"`
}

type Test_ERR_RULE_ARGMIN_1_Validator struct {
	F int `is:"gt"`
}

type Test_ERR_RULE_ARGMIN_2_Validator struct {
	F string `pre:"p4"`
}

type Test_ERR_RULE_ARGMAX_1_Validator struct {
	F int `is:"gt:4:5"`
}

type Test_ERR_RULE_ARGMAX_2_Validator struct {
	F string `pre:"p4:1:2:3"`
}

type Test_ERR_PREPROC_INNVALID_1_Validator struct {
	F string `pre:"p0"`
}
