package testdata

type EqualsValidator struct {
	F1 string        `is:"eq:foo"`
	F2 *int          `is:"eq:123::321"`
	F3 **interface{} `is:"required,eq:foo:123:false:3.14"`
	F4 **interface{} `is:"eq:foo:123:false:3.14"`
}
