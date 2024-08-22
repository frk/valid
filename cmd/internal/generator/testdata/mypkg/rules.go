package mypkg

func HasUniqueInts(v []int, vv ...[]int) bool {
	// ...
	return false
}

func MyRule(v string) bool {
	// ...
	return false
}

func MyRule2(v ...string) bool {
	// ...
	return false
}

func MyRule3(v int64, i int, f float64, s string, b bool) bool {
	// ...
	return false
}

func RuleWithErr1(v string) (ok bool, err error) {
	// ...
	return false, nil
}

func RuleWithErr2(v string, x, y int) (ok bool, err error) {
	// ...
	return false, nil
}

func PreWithOpt(v string, opts any) string {
	// ...
	return v
}

func PreWithOpt2(v string, opts string) string {
	// ...
	return v
}
