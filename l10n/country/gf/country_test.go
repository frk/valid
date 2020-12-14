package gf

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"GF", "GUF"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0612457898",
			"+594612457898",
			"594612457898",
			"0712457898",
			"+594712457898",
			"594712457898",
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
			"+54612457898",
			"+5946124578980",
			"+59461245789",
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
