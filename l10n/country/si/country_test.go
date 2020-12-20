package si

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SI", "SVN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			//
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: isvalid.VAT,
		Pass: []string{
			"SI99662981",
			"SI19136234",
		},
		Fail: []string{
			"SI99662982",
			"SI19136235",
		},
	}})
}
