package gl

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"GL", "GRL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"123456",
			"12 34 56",
			"299 123456",
			"299123456",
			"299 12 34 56",
			"+299 12 34 56",
		},
		Fail: []string{
			"",
			"+4501020304",
			"ASDFGJKLmZXJtZtesting123",
			"12345678",
			"12 34 56 78",
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
