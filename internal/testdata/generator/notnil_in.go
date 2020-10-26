package testdata

type NotnilValidator struct {
	F1 *string                `is:"notnil"`
	F2 **bool                 `is:"notnil"`
	F3 *[]float32             `is:"notnil"`
	F4 *interface{}           `is:"notnil"`
	F5 map[string]interface{} `is:"notnil"`

	G1 *struct {
		F1 *string `is:"notnil"`
		G2 *struct {
			F1 *string `is:"notnil"`
		} `is:"notnil"`
		F2 *string `is:"notnil"`
	}

	// should still work even though...
	FX *****interface{} `is:"notnil"`
	GX *****struct {
		F1 *string    `is:"notnil"`
		F2 **int      `is:"notnil"`
		F3 ***float64 `is:"notnil"`
	} `is:"notnil"`
}
