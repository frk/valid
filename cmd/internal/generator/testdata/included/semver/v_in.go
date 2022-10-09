package testdata

type Validator struct {
	F1 string  `is:"semver"`
	F2 *string `is:"semver"`
}
