package ee

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"EE", "EST"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+372 512 34 567",
			"372 512 34 567",
			"+37251234567",
			"51234567",
			"81234567",
			"+372842345678",
		},
		Fail: []string{
			"12345",
			"",
			"NotANumber",
			"+333 51234567",
			"61234567",
			"+51234567",
			"+372 539 57 4",
			"+372 900 1234",
			"12345678",
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
