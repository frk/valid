package be

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"BE", "BEL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0470123456",
			"+32470123456",
			"32470123456",
			"021234567",
			"+3221234567",
			"3221234567",
		},
		Fail: []string{
			"12345",
			"+3212345",
			"3212345",
			"04701234567",
			"+3204701234567",
			"3204701234567",
			"0212345678",
			"+320212345678",
			"320212345678",
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
