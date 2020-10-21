package testdata

type UserInput struct {
	CountryCode string
	SomeVersion string
	SomeValue   string

	f0  string `is:"-"`
	f1  string `is:"required"`
	f2  string `is:"required:@create"`
	f3  string `is:"required:#key"`
	f4  string `is:"required:@create:#key"`
	f5  string `is:"email"`
	f6  string `is:"url"`
	f7  string `is:"uri"`
	f8  string `is:"pan"`
	f9  string `is:"cvv"`
	F10 string `is:"ssn"`
	F11 string `is:"ein"`
	F12 string `is:"numeric"`
	F13 string `is:"hex"`
	F14 string `is:"hexcolor"`
	F15 string `is:"alphanum"`
	F16 string `is:"cidr"`

	F17 string `is:"phone"`
	F18 string `is:"phone:us:ca"`
	F19 string `is:"phone:&CountryCode"`

	F20 string `is:"zip"`
	F21 string `is:"zip:deu:fin"`
	F22 string `is:"zip:&CountryCode"`

	F23 string `is:"uuid"`
	F24 string `is:"uuid:3"`
	F25 string `is:"uuid:v4"`
	F26 string `is:"uuid:&SomeVersion"`

	F27 string `is:"ip"`
	F28 string `is:"ip:4"`
	F29 string `is:"ip:v6"`
	F30 string `is:"ip:&SomeVersion"`

	F31 string `is:"mac"`
	F32 string `is:"mac:6"`
	F33 string `is:"mac:v8"`
	F34 string `is:"mac:&SomeVersion"`

	F35 string `is:"iso:1234"`
	F36 string `is:"rfc:1234"`

	F37 string `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
	F38 string `is:"re:\"\\w+\""`

	F39 string `is:"contains:foo bar"`
	F40 string `is:"contains:&SomeValue"`

	f41 string `is:"prefix:foo bar"`
	f42 string `is:"prefix:&SomeValue"`

	f43 string `is:"suffix:foo bar"`
	f44 string `is:"suffix:&SomeValue"`

	f45 string  `is:"eq:foo bar"`
	f46 int     `is:"eq:-123"`
	f47 float64 `is:"eq:123.987"`
	f48 string  `is:"eq:&SomeValue"`

	f49 string  `is:"ne:foo bar"`
	f50 int     `is:"ne:-123"`
	f51 float64 `is:"ne:123.987"`
	f52 string  `is:"ne:&SomeValue"`

	f53 uint8   `is:"gt:24,lt:128"`
	f54 int16   `is:"gt:-128,lt:-24"`
	f55 float32 `is:"gt:0.24,lt:1.28"`

	f56 uint8   `is:"gte:24,lte:128"`
	f57 int16   `is:"gte:-128,lte:-24"`
	f58 float32 `is:"gte:0.24,lte:1.28"`

	f59 uint8   `is:"min:24,max:128"`
	f60 int16   `is:"min:-128,max:-24"`
	f61 float32 `is:"min:0.24,max:1.28"`

	f62 uint8   `is:"rng:24:128"`
	f63 int16   `is:"rng:-128:-24"`
	f64 float32 `is:"rng:0.24:1.28"`

	f65 string         `is:"len:28"`
	f66 []int          `is:"len:28"`
	f67 map[string]int `is:"len:28"`

	f68 string         `is:"len:4:28"`
	f69 []int          `is:"len:4:28"`
	f70 map[string]int `is:"len:4:28"`

	f71 string `is:"len::28"`
	f72 []int  `is:"len:4:"`
}
