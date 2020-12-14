package gh

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"GH", "GHA"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"0202345671",
			"0502345671",
			"0242345671",
			"0542345671",
			"0272345671",
			"0572345671",
			"0262345671",
			"0562345671",
			"0232345671",
			"0282345671",
			"+233202345671",
			"+233502345671",
			"+233242345671",
			"+233542345671",
			"+233272345671",
			"+233572345671",
			"+233262345671",
			"+233562345671",
			"+233232345671",
			"+233282345671",
		},
		Fail: []string{
			"082123",
			"232345671",
			"0292345671",
			"+233292345671",
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
