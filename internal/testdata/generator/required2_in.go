package testdata

type Required2Validator struct {
	G1 struct {
		F1 *string `is:"required"`
		G3 *struct {
			F1 string
			F2 int
		} `is:"required"`
	}
}
