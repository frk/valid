package my

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"MY", "MYS"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+60128228789",
			"+60195830837",
			"+6019-5830837",
			"+6019-5830837",
			"+6010-4357675",
			"+60172012370",
			"0128737867",
			"0172012370",
			"01468987837",
			"01112347345",
			"016-2838768",
			"016 2838768",
		},
		Fail: []string{
			"12345",
			"601238788657",
			"088387675",
			"16-2838768",
			"032551433",
			"6088-387888",
			"088-261987",
			"1800-88-8687",
			"088-320000",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"56000",
			"12000",
			"79502",
		},
		Fail: []string{
			//
		},
	}})
}
