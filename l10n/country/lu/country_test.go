package lu

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"LU", "LUX"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"601123456",
			"+352601123456",
		},
		Fail: []string{
			"NaN",
			"791234",
			"+352791234",
			"26791234",
			"+35226791234",
			"+112039812",
			"+352703123456",
			"1234",
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
