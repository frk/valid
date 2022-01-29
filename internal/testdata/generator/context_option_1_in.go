package testdata

type ContextOption1Validator struct {
	F1      string  `is:"required:@new"`
	F2      *string `is:"len::21:@update"`
	F3      *int    `is:"lt:8:@new"`
	context string
}
