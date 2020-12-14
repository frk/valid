package cr

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"CR", "CRI"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+50688888888",
			"+50665408090",
			"+50640895069",
			"25789563",
			"85789563",
		},
		Fail: []string{
			"+5081",
			"+5067777777",
			"+50188888888",
			"+50e987643254",
			"+506e4t4",
			"-50688888888",
			"50688888888",
			"12345678",
			"98765432",
			"01234567",
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
