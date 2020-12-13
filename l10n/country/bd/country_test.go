package bd

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"BD", "BGD"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+8801794626846",
			"01399098893",
			"8801671163269",
			"01717112029",
			"8801898765432",
			"+8801312345678",
			"01494676946",
		},
		Fail: []string{
			"",
			"0174626346",
			"017943563469",
			"18001234567",
			"0131234567",
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
