package testdata

type ContextOption4Validator struct {
	Context string
	G0      *struct {
		F1 []string             `is:"required:@foo,[]email:@bar,len::128:@baz"`
		F2 *map[string]string   `is:"required:@foo,[email:@bar,len:8:128:@baz]phone:@bar,ssn:@baz"`
		F3 []*map[string]string `is:"[]notnil:@foo,len:5:10:@bar,[email:@bar,len:8:128:@baz]phone:@bar,ssn:@baz"`
	} `is:"required:@foo,required:@bar"`
}
