package ie

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"IE", "IRL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+353871234567",
			"353831234567",
			"353851234567",
			"353861234567",
			"353871234567",
			"353881234567",
			"353891234567",
			"0871234567",
			"0851234567",
		},
		Fail: []string{
			"999",
			"+353341234567",
			"+33589484858",
			"353841234567",
			"353811234567",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"A65 TF12",
			"D02 AF30",
		},
		Fail: []string{
			"123",
			"A6W U9U9",
			"75690HG",
			"AW5  TF12",
			"AW5 TF12",
			"756  90HG",
			"A65T F12",
			"O62 O1O2",
		},
	}})
}
