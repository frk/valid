package fr

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"FR", "FRA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0612457898",
			"+33612457898",
			"33612457898",
			"0712457898",
			"+33712457898",
			"33712457898",
		},
		Fail: []string{
			"061245789",
			"06124578980",
			"0112457898",
			"0212457898",
			"0312457898",
			"0412457898",
			"0512457898",
			"0812457898",
			"0912457898",
			"+34612457898",
			"+336124578980",
			"+3361245789",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"75008",
			"44 522",
			"98025",
			"38 499",
			"39940",
		},
		Fail: []string{
			//
		},
	}})
}
