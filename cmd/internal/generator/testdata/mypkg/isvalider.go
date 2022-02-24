package mypkg

type MyString string

func (MyString) IsValid() bool {
	// ...
	return false
}

type MyInt int

func (*MyInt) IsValid() bool {
	// ...
	return false
}
