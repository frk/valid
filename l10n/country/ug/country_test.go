package ug

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"UG", "UGA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+256728590432",
			"+256733875610",
			"256728590234",
			"0773346543",
			"0700459022",
		},
		Fail: []string{
			"999",
			"+254728590432",
			"+25489032",
			"123456789",
			"+254800723845",
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
