package se

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SE", "SWE"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+46701234567",
			"46701234567",
			"0721234567",
			"073-1234567",
			"0761-234567",
			"079-123 45 67",
		},
		Fail: []string{
			"12345",
			"+4670123456",
			"+46301234567",
			"+0731234567",
			"0731234 56",
			"+7312345678",
			"",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"12994",
			"284 39",
			"39556",
			"489 39",
			"499 49",
		},
		Fail: []string{
			//
		},
	}})
}
