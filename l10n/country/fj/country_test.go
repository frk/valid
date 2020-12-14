package fj

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"FJ", "FJI"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+6799898679",
			"6793788679",
			"+679 989 8679",
			"679 989 8679",
			"679 3456799",
			"679908 8909",
		},
		Fail: []string{
			"12345",
			"",
			"04555792",
			"902w99900030900000000099",
			"8uiuiuhhyy&GUU88d",
			"010-38238383",
			"19676338855",
			"679 9 89 8679",
			"6793 45679",
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
