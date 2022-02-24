package testdata

type T60FieldSet1 struct {
	F1 *string `is:"email"`
}

type T60FieldSet2 struct {
	F2 *string `is:"required"`
	T60FieldSet3
}

type T60FieldSet3 struct {
	F3 *string `is:"len:42"`
	T60FieldSet4
}

type T60FieldSet4 struct {
	F4 string `is:"required"`
}

type T60Validator struct {
	T60FieldSet1
	T60FieldSet2
}
