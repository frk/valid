package testdata

type ContextOptionValidator struct {
	F1      string  `is:"required:@new"`
	F2      *string `is:"len::21:@update"`
	context string
}
