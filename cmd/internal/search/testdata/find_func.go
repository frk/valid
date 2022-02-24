package testdata

func IsFoo(v string) bool {
	return false
}

// Unexported functions can be OK depending on where the rules are used.
//
// valid:rule.yaml
//	name: bar
//	args:
//	  - { default: "123" }
//	  - options:
//	    - { value: "1", alias: v1 }
//	    - { value: "2", alias: v2 }
//	    - { value: "3", alias: v3 }
//	error: { text: "bar is not valid" }
func isBar(v string, a1, a2 string) bool {
	return false
}
