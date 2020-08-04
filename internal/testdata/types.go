package testdata

type UserInput struct {
	F1 string `is:"len:8:512"`
	F2 int    `is:"min:8,max:512"`
	F3 string `is:"re:\"\\w+\""`
	F4 string `is:"ne:\"foo\"bar\",ne:\"baz\""`
}
