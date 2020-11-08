package isvalid

func Email(v string) bool {
	return false
}

func URL(v string) bool {
	return false
}

func URI(v string) bool {
	return false
}

func PAN(v string) bool {
	return false
}

func CVV(v string) bool {
	return false
}

func SSN(v string) bool {
	return false
}

func EIN(v string) bool {
	return false
}

func Numeric(v string) bool {
	return false
}

func Hex(v string) bool {
	return false
}

func HexColor(v string) bool {
	return false
}

func Alphanum(v string) bool {
	return false
}

func CIDR(v string) bool {
	return false
}

func Phone(v string, cc ...string) bool {
	return false
}

func Zip(v string, cc ...string) bool {
	return false
}

func UUID(v string, ver ...int) bool {
	return false
}

func IP(v string, ver ...int) bool {
	return false
}

func MAC(v string, ver ...int) bool {
	return false
}

func RFC(v string, num int) bool {
	return false
}

func ISO(v string, num int) bool {
	return false
}

func Match(v string, expr string) bool {
	return regexpCache.m[expr].MatchString(v)
}
