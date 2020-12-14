package do

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"DO", "DOM"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+18096622563",
			"+18295614488",
			"+18495259567",
			"8492283478",
			"8092324576",
			"8292387713",
		},
		Fail: []string{
			"+18091",
			"+1849777777",
			"-18296643245",
			"+18086643245",
			"+18396643245",
			"8196643245",
			"+38492283478",
			"6492283478",
			"8192283478",
		},
	}, {
		Name: "Zip", Func: isvalid.Zip,
		Pass: []string{
			"12345",
		},
		Fail: []string{
			"A1234",
			"123",
			"123456",
		},
	}})
}
