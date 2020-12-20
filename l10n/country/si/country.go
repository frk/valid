package si

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// VAT: (https://web.archive.org/web/20150215110937/http://www.durs.gov.si/si/storitve/vpis_v_davcni_register_in_davcna_stevilka/vpis_v_davcni_register_in_davcna_stevilka_pojasnila/davcna_stevilka_splosno)
	//
	// Davčna številka je sestavljena iz osmih številk; sedem je naključno izbranih, osma pa je izračunana po modulu 11:
	//
	// prvih sedem mest je osnovna številka, ki je naključno izbrana iz nabora številk od 1.000.000 do 9.999.999,
	// na osmem mestu je kontrolna številka, izračunana po modulu 11.
	//
	// Postopek izračuna kontrolne številke:
	//
	// posamezno številko osnovne številke pomnožimo s konstantnimi ponderji 8, 7, 6, 5, 4, 3 in 2,
	// zmnožke seštejemo,
	// seštevek delimo z 11,
	// ostanek deljenja odštejemo od 11 in razlika je kontrolna številka.
	//
	// Če je ostanek deljenja 1 in je razlika zato 10, je kontrolna številka 0. Če je ostanek deljenja 0 in je razlika zato 11, se ta osnovna številka izključi iz nabora možnih davčnih številk.
	//
	// Primer:
	//
	// osnovna številka: 1 5 0 1 2 5 5
	// konstantni ponderji: 8,7,6,5,4,3,2
	// seštevek zmnožkov: 8 + 35 + 0 + 5 + 8 + 15 + 10 = 81
	// seštevek delimo z 11: 7
	// ostanek je: 4
	// razlika do 11: 7 (kontrolna št.)
	// davčna številka: 15012557
	rxvat := regexp.MustCompile(`^SI[1-9][0-9]{7}$`)
	weights := []int{8, 7, 6, 5, 4, 3, 2}

	country.Add(country.Country{
		A2: "SI", A3: "SVN", Num: "705",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+386[ ]?|0)(?:(?:[0-9]{1}[ ]?[0-9]{3}(?:[ ]?[0-9]{2}){2})|(?:[0-9]{2}(?:[ ]?[0-9]{3}){2}))$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			v = v[2:]

			sum := 0
			for i := 0; i < len(v)-1; i++ {
				n, _ := strconv.Atoi(string(v[i]))
				sum += n * weights[i]
			}

			res := (11 - sum%11)
			if res == 10 {
				res = 0
			}

			chk, _ := strconv.Atoi(string(v[len(v)-1]))
			return chk == res
		}),
	})
}
