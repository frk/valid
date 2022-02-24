package testdata

type Test_ERR_NOTNIL_TYPE_1_Validator struct {
	F string `is:"notnil"`
}

type Test_ERR_NOTNIL_TYPE_2_Validator struct {
	F bool `is:"notnil"`
}

type Test_ERR_NOTNIL_TYPE_3_Validator struct {
	F float64 `is:"notnil"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Test_required_Validator struct {
	F1 string            `is:"required"`
	F2 *string           `is:"notnil"`
	F3 []byte            `is:"notnil"`
	F4 map[string]string `is:"notnil"`
	F5 interface{}       `is:"notnil"`
	F6 int               `is:"required"`
	F7 bool              `is:"required"`
}
