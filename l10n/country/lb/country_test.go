package lb

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"LB", "LBN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+96171234568",
			"+9613123456",
			"3456123",
			"3123456",
			"81978468",
			"77675798",
		},
		Fail: []string{
			"+961712345688888",
			"00912220000",
			"7767579888",
			"+0921110000",
			"+3123456888",
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
