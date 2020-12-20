package nl

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"NL", "NLD"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0670123456",
			"0612345678",
			"31612345678",
			"31670123456",
			"+31612345678",
			"+31670123456",
			"+31(0)612345678",
			"0031612345678",
			"0031(0)612345678",
		},
		Fail: []string{
			"12345",
			"+3112345",
			"3112345",
			"06701234567",
			"012345678",
			"+3104701234567",
			"3104701234567",
			"0212345678",
			"021234567",
			"+3121234567",
			"3121234567",
			"+310212345678",
			"310212345678",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"1012 SZ",
			"3432FE",
			"1118 BH",
			"3950IO",
			"3997 GH",
		},
		Fail: []string{
			//
		},
	}, {
		Name: "VAT", Func: isvalid.VAT,
		Pass: []string{
			"NL999999999B01",
		},
		Fail: []string{
			//
		},
	}})
}
