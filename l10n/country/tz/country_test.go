package tz

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"TZ", "TZA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+255728590432",
			"+255733875610",
			"255628590234",
			"0673346543",
			"0600459022",
		},
		Fail: []string{
			"999",
			"+254728590432",
			"+25589032",
			"123456789",
			"+255800723845",
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
