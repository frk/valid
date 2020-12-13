package az

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AZ", "AZE"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+994707007070",
			"0707007070",
			"+994502111111",
			"0505436743",
			"0554328772",
			"0993301022",
			"+994776007139",
		},
		Fail: []string{
			"wronumber",
			"",
			"994707007070",
			"++9945005050",
			"556007070",
			"1234566",
			"+994778008080a",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"AZ0100",
			"AZ0121",
			"AZ3500",
		},
		Fail: []string{
			"",
			" AZ0100",
			"AZ100",
			"AZ34340",
			"EN2020",
			"AY3030",
		},
	}})
}
