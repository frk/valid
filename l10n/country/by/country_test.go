package by

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"BY", "BLR"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+375241234567",
			"+375251234567",
			"+375291234567",
			"+375331234567",
			"+375441234567",
			"375331234567",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+0821231234",
			"12345",
			"",
			"ASDFGJKLmZXJtZtesting123",
			"010-38238383",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"225320",
			"211120",
			"247710",
			"231960",
		},
		Fail: []string{
			//
		},
	}})
}
