package sg

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SG", "SGP"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"87654321",
			"98765432",
			"+6587654321",
			"+6598765432",
			"+6565241234",
		},
		Fail: []string{
			"987654321",
			"876543219",
			"8765432",
			"9876543",
			"12345678",
			"+98765432",
			"+9876543212",
			"+15673628910",
			"19876543210",
			"8005552222",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"308215",
			"546080",
		},
		Fail: []string{
			//
		},
	}})
}
