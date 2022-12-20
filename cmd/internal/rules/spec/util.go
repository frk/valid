package spec

// AddCustom is intended to be used by tests in other packages
// that don't normally have write access to the _custom map.
func AddCustom(key string, s *Spec) {
	_custom[key] = s
}
