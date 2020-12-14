package fo

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"FO", "FRO"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"123456",
			"12 34 56",
			"298 123456",
			"298123456",
			"298 12 34 56",
			"+298 12 34 56",
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
