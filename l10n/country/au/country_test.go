package au

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"AU", "AUS"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"61404111222",
			"+61411222333",
			"0417123456",
		},
		Fail: []string{
			"082123",
			"08212312345",
			"21821231234",
			"+21821231234",
			"+0821231234",
			"04123456789",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"4000",
			"2620",
			"3000",
			"2017",
			"0800",
		},
		Fail: []string{
			//
		},
	}})
}
