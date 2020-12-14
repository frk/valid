package sl

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"SL", "SLE"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+94766661206",
			"94713114340",
			"0786642116",
			"078 7642116",
			"078-7642116",
		},
		Fail: []string{
			"9912349956789",
			"12345",
			"1678123456",
			"0731234567",
			"0749994567",
			"0797878674",
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
