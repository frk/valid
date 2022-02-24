package testdata

type Test_ERR_ORDERED_TYPE_1_Validator struct {
	F []string `is:"min:8"`
}

type Test_ERR_ORDERED_TYPE_2_Validator struct {
	F []int `is:"gt:8"`
}

type Test_ERR_ORDERED_ARGTYPE_1_Validator struct {
	F int `is:"gte:0.8"`
}

type Test_ERR_ORDERED_ARGTYPE_2_Validator struct {
	F float64 `is:"gte:foo"`
}

type Test_ERR_ORDERED_ARGTYPE_3_Validator struct {
	F string `is:"lte:&S.F"`
	S struct{ F int }
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_ordered_Validator struct {
	F1 int     `is:"gt:8"`
	F2 float64 `is:"lt:0.029"`
	F3 uint16  `is:"gte:128"`
	F4 int16   `is:"lte:-128"`
	F5 int16   `is:"min:8,max:128"`
	F6 string  `is:"lt:foo,min:bar"`
	F7 int     `is:"lt:&S.F1,gt:&S.F2"`
	S  struct {
		F1 int
		F2 float64
	}
}
