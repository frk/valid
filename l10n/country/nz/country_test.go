package nz

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"NZ", "NZL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+6427987035",
			"642240512347",
			"0293981646",
			"029968425",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"+642956696123566",
			"+02119620856",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"7843",
			"3581",
			"0449",
			"0984",
			"4144",
		},
		Fail: []string{
			//
		},
	}})
}
