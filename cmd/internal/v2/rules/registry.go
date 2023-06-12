package rules

type Registry interface {
	Lookup(name string) (r Rule, ok bool)
}

type registry struct {
	builtin  map[string]Rule
	included map[string]Rule
	custom   map[string]Rule
}

func (rr registry) Lookup(name string) (r Rule, ok bool) {
	if r, ok := rr.custom[name]; ok {
		return r, ok
	}
	if r, ok := rr.included[name]; ok {
		return r, ok
	}
	if r, ok := rr.builtin[name]; ok {
		return r, ok
	}
	return Rule{}, false
}

////////////////////////////////////////////////////////////////////////////////
