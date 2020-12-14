package hk

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"HK", "HKG"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"91234567",
			"9123-4567",
			"61234567",
			"51234567",
			"+85291234567",
			"+852-91234567",
			"+852-9123-4567",
			"+852 9123 4567",
			"9123 4567",
			"852-91234567",
		},
		Fail: []string{
			"999",
			"+852-912345678",
			"123456789",
			"+852-1234-56789",
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
