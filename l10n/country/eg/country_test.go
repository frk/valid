package eg

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"EG", "EGY"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+201004513789",
			"+201111453489",
			"+201221204610",
			"+201144621154",
			"+201200124304",
			"+201011201564",
			"+201124679001",
			"+201064790156",
			"+201274652177",
			"+201280134679",
			"+201090124576",
			"+201583728900",
			"201599495596",
			"201090124576",
			"01090124576",
			"01538920744",
			"1593075993",
			"1090124576",
		},
		Fail: []string{
			"+221004513789",
			"+201404513789",
			"12345",
			"",
			"+9639626626262",
			"+963332210972",
			"0114152198",
			"962796477263",
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
