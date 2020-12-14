package ke

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"KE", "KEN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+254728590432",
			"+254733875610",
			"254728590234",
			"0733346543",
			"0700459022",
			"0110934567",
			"+254110456794",
			"254198452389",
		},
		Fail: []string{
			"999",
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
