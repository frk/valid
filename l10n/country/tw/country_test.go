package tw

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"TW", "TWN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0987123456",
			"+886987123456",
			"886987123456",
			"+886-987123456",
			"886-987123456",
		},
		Fail: []string{
			"12345",
			"",
			"Vml2YW11cyBmZXJtZtesting123",
			"0-987123456",
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
