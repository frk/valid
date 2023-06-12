package rules

// Kind represents the specific kind of a rule.
type Kind uint

func (k Kind) String() string {
	if int(k) < len(_kindstring) {
		return _kindstring[k]
	}
	return "<invalid>"
}

const (
	_ Kind = iota

	REQUIRED   // required, notnil
	COMPARABLE // =, !=
	ORDERED    // >, >=, <, <=, min, max
	LENGTH     // len, runecount, ...
	RANGE      // rng
	ENUM       // enum
	FUNCTION   // <custom/builtin/included func rules>
	METHOD     // isvalid (note: "isvalid" is added implicitly but can be specified explicitly)

	// "modifiers"
	OPTIONAL // omitnil [is the default rule for pointers] (ptr only), optional (ptr & base)
	NOGUARD  // nonilguard
	REMOVE   // remove implicit rule, e.g. "-isvalid"

	// "preprocessors"
	PREPROC // <custom/builtin/included func rules>
)

var _kindstring = [...]string{
	REQUIRED:   "REQUIRED",
	COMPARABLE: "COMPARABLE",
	ORDERED:    "ORDERED",
	LENGTH:     "LENGTH",
	RANGE:      "RANGE",
	ENUM:       "ENUM",
	FUNCTION:   "FUNCTION",
	METHOD:     "METHOD",
	OPTIONAL:   "OPTIONAL",
	NOGUARD:    "NOGUARD",
	REMOVE:     "REMOVE",
	PREPROC:    "PREPROC",
}
