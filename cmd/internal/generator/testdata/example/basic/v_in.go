package testdata

type UserCreateParams struct {
	FName string `is:"len:1:300" pre:"trim"`
	LName string `is:"len:1:300,required" pre:"trim"`
	Email string `is:"email,required" pre:"lower,trim"`
	Passw string `is:"strongpass,required" pre:"trim"`
	Age   int    `is:"min:3,max:150"`
}

type UserCreateParamsValidator struct {
	UserCreateParams
}
