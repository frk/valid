package bg

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"BG", "BGR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+359897123456",
			"+359898888888",
			"0897123123",
		},
		Fail: []string{
			"",
			"0898123",
			"+359212555666",
			"18001234567",
			"12125559999",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"1000",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: isvalid.VAT,
		Pass: []string{
			"BG123456789",
			"BG1234567890",
		},
		Fail: []string{
			"BG12345678",
			"BG12345678901",
		},
	}})
}
