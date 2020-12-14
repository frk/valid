package mo

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"MO", "MAC"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"61234567",
			"+85361234567",
			"+853-61234567",
			"+853-6123-4567",
			"+853 6123 4567",
			"6123 4567",
			"853-61234567",
		},
		Fail: []string{
			"999",
			"12345678",
			"612345678",
			"+853-12345678",
			"+853-22345678",
			"+853-82345678",
			"+853-612345678",
			"+853-1234-5678",
			"+853 1234 5678",
			"+853-6123-45678",
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
