package testdata

type V1Validator struct {
	*V1Params `is:"noguard,omitkey"`
}

type V1Params struct {
	F1 []int `is:"unique_ints"`
	F2 []int `is:"unique_ints:.F1"`
	F3 []*V1ElemParams
}

type V1ElemParams struct {
	X1 []int `is:"[]gt:.X2"`
	X2 int
}
