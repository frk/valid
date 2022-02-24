package testdata

type Test_ERR_RULE_ELEM_1_is_Validator struct {
	F string `is:"[]email"`
}

type Test_ERR_RULE_ELEM_2_is_Validator struct {
	F []string `is:"[][]email"`
}

type Test_ERR_RULE_ELEM_3_is_Validator struct {
	F []map[string]string `is:"[][]email,[]email"`
}

type Test_ERR_RULE_KEY_1_is_Validator struct {
	F string `is:"[email]"`
}

type Test_ERR_RULE_KEY_2_is_Validator struct {
	F []string `is:"[email]"`
}

type Test_ERR_RULE_KEY_3_is_Validator struct {
	F map[string][]string `is:"[][email]"`
}

type Test_ERR_RULE_ELEM_1_pre_Validator struct {
	F string `pre:"[]trim"`
}

type Test_ERR_RULE_ELEM_2_pre_Validator struct {
	F []string `pre:"[][]trim"`
}

type Test_ERR_RULE_ELEM_3_pre_Validator struct {
	F []map[string]string `pre:"[][]trim,[]trim"`
}

type Test_ERR_RULE_KEY_1_pre_Validator struct {
	F string `pre:"[trim]"`
}

type Test_ERR_RULE_KEY_2_pre_Validator struct {
	F []string `pre:"[trim]"`
}

type Test_ERR_RULE_KEY_3_pre_Validator struct {
	F map[string][]string `pre:"[][trim]"`
}

type Test_ERR_RULE_UNDEFINED_1_is_Validator struct {
	F string `is:"email,trim"`
}

type Test_ERR_RULE_UNDEFINED_2_is_Validator struct {
	F *string `is:"email,trim"`
}

type Test_ERR_RULE_UNDEFINED_1_pre_Validator struct {
	F string `pre:"trim,email"`
}

type Test_ERR_RULE_UNDEFINED_2_pre_Validator struct {
	F *string `pre:"trim,email"`
}
