package fi

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"FI", "FIN"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+358505557171",
			"0455571",
			"0505557171",
			"358505557171",
			"04412345",
			"0457 123 45 67",
			"+358457 123 45 67",
			"+358 50 555 7171",
		},
		Fail: []string{
			"12345",
			"",
			"045557",
			"045555717112312332423423421",
			"Vml2YW11cyBmZXJtZtesting123",
			"010-38238383",
			"+3-585-0555-7171",
			"+9676338855",
			"19676338855",
			"6676338855",
			"+99676338855",
			"044123",
			"019123456789012345678901",
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
