package testdata

type T04Validator struct {
	F1 string        `is:"ne:foo"`
	F2 *int          `is:"ne:123::321"`
	F3 **interface{} `is:"required,ne:foo:123:false:3.14"`
}
