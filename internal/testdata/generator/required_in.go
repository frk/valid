package testdata

type RequiredValidator struct {
	F1 *string                `is:"required"`
	F2 string                 `is:"required"`
	F3 int                    `is:"required"`
	F4 *float32               `is:"required"`
	F5 *interface{}           `is:"required"`
	F6 **bool                 `is:"required"`
	F7 map[string]interface{} `is:"required"`
	F8 *[]int32               `is:"required"`

	G1 struct {
		F1 *string `is:"required"`
		G2 *struct {
			F1 *string `is:"required"`
		} `is:"required"`
		F2 *string `is:"required"`
		G3 *struct {
			F1 string
			F2 int
		} `is:"required"`
	}

	// should still work even though...
	FX *****interface{} `is:"required"`
	GX *****struct {
		F1 string  `is:"required"`
		F2 int     `is:"required"`
		F3 float64 `is:"required"`
	} `is:"required"`
}
