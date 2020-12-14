package pa

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"PA", "PAN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+5076784565",
			"+5074321557",
			"5073331112",
			"+50723431212",
		},
		Fail: []string{
			"+50755555",
			"+207123456",
			"2001236542",
			"+507987643254",
			"+507jjjghtf",
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
