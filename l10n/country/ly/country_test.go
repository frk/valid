package ly

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"LY", "LBY"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"912220000",
			"0923330000",
			"218945550000",
			"+218958880000",
			"212220000",
			"0212220000",
			"+218212220000",
		},
		Fail: []string{
			"9122220000",
			"00912220000",
			"09211110000",
			"+0921110000",
			"+2180921110000",
			"021222200000",
			"213333444444",
			"",
			"+212234",
			"+21",
			"02122333",
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
