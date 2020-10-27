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

// An implementation of the ErrorConstructor interface can be used to make the
// generated validation code return custom, application specific errors.
//
// The cmd/isvalid tool checks the target validator types for a field whose
// type implements the ErrorConstructor interface and, if a validator type
// does have such a field, the tool will generate code that invokes the Error
// method when an invalid value is encountered, passing it the relevant
// information, and yielding the error returned from it.
type ErrorConstructor interface {
	// The intended implementation of the Error method should construct a
	// new custom error value based on the given parameters and return it.
	//
	// The key parameter will hold the key of the field whose value failed
	// validation. The val parameter holds the value that failed validation.
	// The rule parameter holds the name of the validation rule which the
	// failed value did not pass. The args parameter holds the rule's
	// arguments specified by the `is` tag.
	Error(key string, val interface{}, rule string, args ...interface{}) error
}

// An implementation of the ErrorAggregator interface can be used to make the
// generated validation code return custom, application specific errors.
//
// The cmd/isvalid tool checks the target validator types for a field whose type
// implements the ErrorAggregator interface and, if a validator type does have
// such a field, the tool will generate code that invokes the Error method each
// time an invalid value is encountered, passing it the relevant information,
// at the end it will invoke the Out method to return the final result.
type ErrorAggregator interface {
	// The intended implementation of the Error method should construct a
	// new custom error value based on the given parameters and retain it
	// for until the generated validation code is done.
	//
	// The key parameter will hold the key of the field whose value failed
	// validation. The val parameter holds the value that failed validation.
	// The rule parameter holds the name of the validation rule which the
	// failed value did not pass. The args parameter holds the rule's
	// arguments specified by the `is` tag.
	Error(key string, val interface{}, rule string, args ...interface{})
	// The Out method will be invoked by the generated validation code at
	// the end to yield the error value it returns.
	Out() error
}
