package uy

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"UY", "URY"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+59899123456",
			"099123456",
			"+59894654321",
			"091111111",
		},
		Fail: []string{
			"54321",
			"montevideo",
			"",
			"+598099123456",
			"090883338",
			"099 999 999",
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
