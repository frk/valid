package mypkg

// incompatible in/out
func MyPreproc1(v float64) string {
	// ...
	return ""
}

// incompatible in/out
func MyPreproc2(v string, opt uint) float64 {
	// ...
	return 0
}

////////////////////////////////////////////////////////////////////////////////

func MyPreproc3(v string) string {
	return ""
}

func MyPreproc4(v string, opt uint) string {
	return ""
}
