package testdata

type Validator struct {
	F1 string  `is:"latlong"`
	F2 *string `is:"latlong:dms"`
}
