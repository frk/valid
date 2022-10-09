package testdata

type mytype string

func (mytype) IsValid() bool {
	// ...
	return true
}

type Validator struct {
	F1 mytype
	F2 mytype `is:"isvalid"`
	F3 mytype `is:"-isvalid"`
}
