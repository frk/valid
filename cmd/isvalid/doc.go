// isvalid is a tool for generating struct field validation.
package main

// TODO(mkopriva): turn the stuff below into proper documentation

// IsValider is here for the purposes of documentation only.
type IsValider interface {
	//
	IsValid() bool
}

// BeforeValidator is here for the purposes of documentation only.
type BeforeValidator interface {
	BeforeValidate() error
}

// AfterValidator is here for the purposes of documentation only.
type AfterValidator interface {
	AfterValidate() error
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
