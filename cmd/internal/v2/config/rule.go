package config

type RuleConfig struct {
	// The function's name qualified with
	// the path of the function's package.
	Func ObjectIdent `yaml:"func"`
	// The configuration for the function's rule. This is
	// optional if the func's doc already has a valid config.
	Rule *Rule `yaml:"rule"`
}

type Rule struct {
	// The name of the rule.
	Name string `yaml:"name"`
	// The configuration for the rule's arguments.
	//
	// NOTE: If args is NOT empty, then it MUST contain a number
	// of elements that is compatible with the parameters of the
	// rule's function.
	Args []RuleArg `yaml:"args"`
	// ArgMin can be used to enforce the number of arguments
	// that the rule must accept.
	//
	// - ArgMin is optional, if omitted its value will
	//   be inferred from the function's signature.
	// - When ArgMin is provided it will be used to enforce
	//   the valid number of arguments for variadic functions.
	// - When ArgMin is provided it MUST be compatible
	//   with the function's signature, if it isn't then
	//   the tool will exit with an error.
	ArgMin *uint `yaml:"arg_min"`
	// ArgMax can be used to override the upper limit of the number
	// of arguments that the rule should be allowed to accept.
	// A negative ArgMax can be used to indicate that there's no
	// upper limit to the number of arguments.
	//
	// - ArgMax is optional, if omitted its value will
	//   be inferred from the function's signature.
	// - When ArgMax is provided it will be used to enforce
	//   the valid number of arguments for variadic functions.
	// - When ArgMax is provided it MUST be compatible
	//   with the function's signature, if it isn't then
	//   the tool will exit with an error.
	ArgMax *int `yaml:"arg_max"`
	// The configuration for the error that should be generated for the rule.
	Error RuleErrMesg `yaml:"error"`
	// The join operator that should be used to join
	// multiple instances of the rule into a single one.
	//
	// The value MUST be one of: "AND", "OR", or "NOT" (case insensitive).
	JoinOp JoinOp `yaml:"join_op"`
}

type RuleArg struct {
	// The rule argument's default value. If nil, then the
	// rule argument's value MUST be provided in the struct tag.
	//
	// If not nil, the value must be a scalar value.
	Default *Scalar `yaml:"default"`
	// If options is empty, then ANY value can be provided for
	// the argument in the rule's struct tag.
	//
	// If options is NOT empty, then it is considered to represent, together
	// with the default value, the *complete* set of valid values that can be
	// provided as the argument in the rule's struct tag.
	Options []RuleArgOption `yaml:"options"`
}

type RuleArgOption struct {
	// Value specifies the value that the generator should supply
	// as the rule's argument in the generated code.
	Value Scalar `yaml:"value"`
	// Alias is an alternative identifier of the argument's value that
	// can be used within the rule's struct tag. This field is optional.
	Alias string `yaml:"alias"`
}

type RuleErrMesg struct {
	// The text of the error message.
	Text string `yaml:"text,omitempty"`
	// If true the generated error message
	// will include the rule's arguments.
	WithArgs bool `yaml:"with_args,omitempty"`
	// The separator used to join the rule's
	// arguments for the error message.
	ArgSep string `yaml:"arg_sep,omitempty"`
	// The text to be appended after the list of arguments.
	ArgSuffix string `yaml:"arg_suffix,omitempty"`
}
