package testdata

type Test_ERR_FIELD_UNKNOWN_1_Validator struct {
	F int `is:"gt:&num"`
}

type Test_ERR_FIELD_UNKNOWN_2_Validator struct {
	F string `pre:"p4:&num"`
}

type Test_ERR_FIELD_UNKNOWN_3_Validator struct {
	num int
	X   struct {
		F int `is:"gt:.num"`
	}
}

type Test_ERR_FIELD_UNKNOWN_4_Validator struct {
	num int
	X   struct {
		F string `pre:"p4:.num"`
	}
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

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_checker_Validator struct {
	F1 int
	F2 int
	S1 struct {
		F1 int
		F2 int

		T1 int `is:"gt:.F1"`
		T2 int `is:"gt:&F2"`
		T3 int `is:"gt:.S2.F1"`

		S2 struct {
			F1 int
		}
	}

	T1 int `is:"gt:&F1"`
	T2 int `is:"gt:&S1.F2"`
	T3 int `is:"gt:.S1.F1"`
	T4 int `is:"gt:.S1.S2.F1"`
}
