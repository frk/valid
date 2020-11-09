package parser

import (
	"testing"
)

func TestParseFunc(t *testing.T) {
	tests := []struct {
		fpkg      string
		fname     string
		wantError bool
	}{
		{"strings", "Contains", false},
		{"strings", "Abracadabra", true},
		{"sgnirts", "Contains", true},
	}

	for i, tt := range tests {
		f, err := ParseFunc(tt.fpkg, tt.fname, nil)
		if err != nil && !tt.wantError {
			t.Errorf("#%d: ParseFunc(%q, %q) want err=<nil>; got err=%v", i, tt.fpkg, tt.fname, err)
		} else if err == nil && tt.wantError {
			t.Errorf("#%d: ParseFunc(%q, %q) want err=<non-nil>; got err=<nil>", i, tt.fpkg, tt.fname)
		}

		if err == nil && !tt.wantError {
			if p := f.Pkg(); p.Path() != tt.fpkg || f.Name() != tt.fname {
				t.Errorf("#%d: want=%s.%s; got err=%v", i, tt.fpkg, tt.fname, f)
			}
		}
	}
}
