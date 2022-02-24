package mypkg

// valid:rule.yaml
//	name: myrule
func MyRule(v string) bool {
	// ...
	return false
}

// valid:rule.yaml
//	name: myrule2
func MyRule2(v ...string) bool {
	// ...
	return false
}

// valid:rule.yaml
//	name: myrule3
func MyRule3(v int64, i int, f float64, s string, b bool) bool {
	// ...
	return false
}
