package testdata

func RuleFunc(s string) bool {
	return false
}

// contains bad yaml
//
// valid:rule.yaml
//		name: 1234
func RuleFunc2(s string) bool {
	return false
}

// has invalid signature
func RuleFunc3(s string) (bool, int) {
	return false, 0
}

// has invalid signature #2
func RuleFunc4() bool {
	return false
}

func RuleFunc5(s string, opt uint) bool {
	return false
}

func RuleFunc6(s string, opt1 uint, opt2 bool) bool {
	return false
}

func RuleFunc7(s string, opt1 uint, opt2 bool, opt3 string) bool {
	return false
}

////////////////////////////////////////////////////////////////////////////////
// valid rules

type CheckHelper interface {
	Check(v string) bool
}

func RuleFunc8(s string, h CheckHelper) bool {
	return h.Check(s)
}

////////////////////////////////////////////////////////////////////////////////
// valid pre-processors

func PreProc1(s string) string {
	return ""
}

func PreProc2(s string, opt bool) string {
	return ""
}

func PreProc3(s string, opt byte, opts ...uint) string {
	return ""
}
