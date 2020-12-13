package om

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"OM", "OMN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+96891212121",
			"0096899999999",
			"93112211",
			"99099009",
		},
		Fail: []string{
			"+96890212121",
			"0096890999999",
			"0090999999",
			"+9689021212",
			"",
			"+212234",
			"+21",
			"02122333",
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
