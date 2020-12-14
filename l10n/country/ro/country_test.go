package ro

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"RO", "ROU"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+40740123456",
			"+40 740123456",
			"+40740 123 456",
			"+40740.123.456",
			"+40740-123-456",
			"40740123456",
			"40 740123456",
			"40740 123 456",
			"40740.123.456",
			"40740-123-456",
			"0740123456",
			"0740/123456",
			"0740 123 456",
			"0740.123.456",
			"0740-123-456",
		},
		Fail: []string{
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"123456",
			"740123456",
			"+40640123456",
			"+40210123456",
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
