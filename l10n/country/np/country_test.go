package np

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"NP", "NPL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+9779817385479",
			"+9779717385478",
			"+9779862002615",
			"+9779853660020",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+97796123456789",
			"+9771234567",
			"+977981234",
			"4736338855",
			"66338855",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"10811",
			"32600",
			"56806",
			"977",
		},
		Fail: []string{
			"11977",
			"asds",
			"13 32",
			"-977",
			"97765",
		},
	}})
}
