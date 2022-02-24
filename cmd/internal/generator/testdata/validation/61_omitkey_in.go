package testdata

type T61FieldSet1 struct {
	F1 *string `is:"email"`
}

type T61FieldSet2 struct {
	F2 *string `is:"required"`
	T61FieldSet3
}

type T61FieldSet3 struct {
	F3           *string `is:"len:42"`
	T61FieldSet4 `is:"omitkey"`
}

type T61FieldSet4 struct {
	F4 string `is:"required"`
}

type T61Validator struct {
	F T61FieldSet1 `is:"omitkey"`
	T61FieldSet2
}
