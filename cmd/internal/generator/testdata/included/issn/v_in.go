package testdata

type Validator struct {
	F1 string  `is:"issn:true:true"`
	F2 *string `is:"issn:false:false"`
}
