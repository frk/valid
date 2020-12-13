package ec

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"EC", ""}, testutil.List{{
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
	}})
}
