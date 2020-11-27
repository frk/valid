package testdata

type AlnumValidator struct {
	F1 string   `is:"alnum"`
	F2 **string `is:"alnum:en"`
	F3 **string `is:"required,alnum:sk:cs"`

	F4 string   `is:"alpha"`
	F5 **string `is:"alpha:en"`
	F6 **string `is:"required,alpha:sk:cs"`
}
