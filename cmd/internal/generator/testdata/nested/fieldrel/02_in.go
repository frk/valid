package testdata

type T02Validator struct {
	*T02Params `is:"noguard,omitkey"`
}

type T02Params struct {
	F1 []int `is:"unique_ints"`
	F2 []int `is:"unique_ints:.F1"`
	F3 []*T02ElemParams
}

type T02ElemParams struct {
	X1 []int `is:"[]gt:.X2"`
	X2 int
}
