package testdata

type Test_ERR_METHOD_TYPE_1_Validator struct {
	F int `is:"isvalid"`
}

////////////////////////////////////////////////////////////////////////////////
// valid test cases
////////////////////////////////////////////////////////////////////////////////

type Typ struct{}

func (*Typ) IsValid() bool { return true }

type Test_method_Validator struct {
	F1 Typ  `is:"isvalid"`
	F2 *Typ `is:"isvalid"`
}
