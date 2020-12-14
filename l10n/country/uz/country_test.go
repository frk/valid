package uz

import (
	"testing"

	"github.com/frk/isvalid"
	"github.com/frk/isvalid/internal/testutil"
)

func Test(t *testing.T) {
	testutil.Run(t, []string{"UZ", "UZB"}, testutil.List{{
		Name: "Phone", Func: isvalid.Phone,
		Pass: []string{
			"+998664835244",
			"998664835244",
			"664835244",
			"+998957124555",
			"998957124555",
			"957124555",
		},
		Fail: []string{
			"+998644835244",
			"998644835244",
			"644835244",
			"+99664835244",
			"ASDFGJKLmZXJtZtesting123",
			"123456789",
			"870123456",
			"",
			"+998",
			"998",
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
