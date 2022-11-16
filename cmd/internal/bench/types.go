package bench

type FrkValidator struct {
	F01  string              `is:"required"`
	F02  *string             `is:"required"`
	F03  *float32            `is:"eq:3.14"`
	F04  *int                `is:"gt:123"`
	F05  int8                `is:"lte:42"`
	F06  float64             `is:"rng:3.14:42"`
	F07  float64             `is:"min:3.14,max:42"`
	F08  string              `is:"len:8:256"`
	F09  []uint8             `is:"len:12"`
	F10  map[string]struct{} `is:"len::12"`
	F11  []byte              `is:"runecount:4:"`
	F12  string              `is:"runecount::15"`
	F13  *string             `is:"prefix:hello,required"`
	F14  *string             `is:"contains:hello"`
	F15  string              `is:"alnum:en"`
	F16  string              `is:"cidr"`
	F17  *string             `is:"fqdn"`
	F18  *string             `is:"ip"`
	F19  string              `is:"email"`
	F20  string              `is:"phone"`
	Sub1 *FrkSub1Type
	Sub2 *FrkSub2Type

	Arg1 string
	Arg2 int32
	Arg3 float64
}

type FrkSub1Type struct {
	F01 []string           `is:"len::10,[]email"`
	F02 []string           `is:"required,[]phone"`
	F03 *map[string]string `is:"required,[email]"`
	F04 map[string]string  `is:"[email]phone"`
	F05 map[string]int     `is:"len:5:,[phone]rng:21:99"`
	F06 *string            `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
	F07 string             `is:"suffix:goodbye"`
	F08 *string            `is:"eq:&Arg1"`
	F09 int32              `is:"lt:&Arg2"`
	F10 float64            `is:"gt:&Arg3"`
}

type FrkSub2Type struct {
	F01 []map[string][]string `is:"len:3,[][uuid]len:1:,[]zip"`
	F02 map[[2]string]any     `is:"[[]ip]notnil"`
	F03 *struct {
		f1 string `is:"eq:&Sub2.Arg1"`
		f2 int32  `is:"lt:&Arg2,required"`
	} `is:"omitnil"`
	F04 string    `is:"prefix:foo,contains:bar,suffix:baz:quux,len:9:64"`
	F05 *string   `is:"email,optional"`
	F06 []*string `is:"[]email,omitnil"`
	F07 int       `is:"max:128,optional"`
	F08 *string   `is:"eq:.Arg1"`
	F09 int32     `is:"lt:.Arg2"`
	F10 float64   `is:"gt:.Arg3"`
	Sub *FrkSub1Type

	Arg1 string
	Arg2 int32
	Arg3 float64
}

////////////////////////////////////////////////////////////////////////////////

type GoPlaygroundParams struct {
	F01  string              `validate:"required"`
	F02  *string             `validate:"required"`
	F03  *float64            `validate:"eq=3.14"` // doesn't seem to be able to handle float32 properly
	F04  *int                `validate:"gt=123"`
	F05  int8                `validate:"lte=42"`
	F06  float64             `validate:"min=3.14,max=42"` // no range/between i think
	F07  float64             `validate:"min=3.14,max=42"`
	F08  string              `validate:"min=8,max=256"` // no len-range i think
	F09  []uint8             `validate:"len=12"`
	F10  map[string]struct{} `validate:"max=12"`
	F11  []byte              `validate:"min=4"`  // no rune count
	F12  string              `validate:"max=15"` // no rune count
	F13  *string             `validate:"startswith=hello,required"`
	F14  *string             `validate:"contains=hello"`
	F15  string              `validate:"alphanum"`
	F16  string              `validate:"cidr"`
	F17  *string             `validate:"fqdn"`
	F18  *string             `validate:"ip_addr"`
	F19  string              `validate:"email"`
	F20  string              `validate:"e164"`
	Sub1 *GoPlaygroundSub1Type
	Sub2 *GoPlaygroundSub2Type

	Arg1 string
	Arg2 int32
	Arg3 float64
}

type GoPlaygroundSub1Type struct {
	F01 []string           `validate:"max=10,dive,email"`
	F02 []string           `validate:"required,dive,e164"`
	F03 *map[string]string `validate:"required,dive,keys,email,endkeys"`
	F04 map[string]string  `validate:"dive,keys,email,endkeys,e164"`
	F05 map[string]int     `validate:"min=5,dive,keys,e164,endkeys,min=21,max=99"`
	F06 *string            `validate:"alphanum,lowercase"` // no regexp support
	F07 string             `validate:"endswith=goodbye"`
	F08 *string            `validate:"eqfield=Arg1"`
	F09 int32              `validate:"ltfield=Arg2"`
	F10 float64            `validate:"gtfield=Arg3"`

	Arg1 string
	Arg2 int32
	Arg3 float64

	// doesn't support absolute (from root) field reference ...
}

type GoPlaygroundSub2Type struct {
	F01 []map[string][]string `validate:"len=3,dive,dive,keys,uuid,endkeys,min=1,dive,e164"` // no zip :(
	F02 map[[2]string]any     `validate:"dive,keys,dive,ip_addr,endkeys,required"`
	F03 *struct {
		f1 string `validate:"eqcsfield=Sub2.Arg1"`
		f2 int32  `validate:"ltcsfield=Arg2,required"`
	} `validate:"omitempty"`
	F04 string    `validate:"startswith=foo,contains=bar,endswith=baz,min=9,max=64"`
	F05 *string   `validate:"email,omitempty"`
	F06 []*string `validate:"dive,email,omitempty"`
	F07 int       `validate:"max=128,omitempty"`
	F08 *string   `validate:"eqcsfield=Sub.Arg1"`
	F09 int32     `validate:"ltcsfield=Sub.Arg2"`
	F10 float64   `validate:"gtcsfield=Sub.Arg3"`
	Sub *GoPlaygroundSub1Type
}