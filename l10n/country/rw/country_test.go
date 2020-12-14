package rw

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"RW", "RWA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+250728590432",
			"+250733875610",
			"250738590234",
			"0753346543",
			"0780459022",
		},
		Fail: []string{
			"999",
			"+254728590432",
			"+25089032",
			"123456789",
			"+250800723845",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}})
}
