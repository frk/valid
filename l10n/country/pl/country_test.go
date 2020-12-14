package pl

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"PL", "POL"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+48512689767",
			"+48 56 376 87 47",
			"56 566 78 46",
			"657562855",
			"+48657562855",
			"+48 887472765",
			"+48 56 6572724",
			"+48 67 621 5461",
			"48 67 621 5461",
		},
		Fail: []string{
			"+48  67 621 5461",
			"+55657562855",
			"3454535",
			"teststring",
			"",
			"1800-88-8687",
			"+6019-5830837",
			"357562855",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"47-260",
			"12-930",
			"78-399",
			"39-490",
			"38-483",
		},
		Fail: []string{
			"360",
			"90312",
			"399",
			"935",
			"38842",
		},
	}})
}
