package testdata

import (
	"time"

	"github.com/frk/isvalid/internal/testdata/mypkg"
)

type UserInput struct {
	CountryCode string
	SomeVersion int
	SomeValue   string

	F0 string `is:"-"`
	F1 string `is:"required"`
	F2 string `is:"required:@create"`

	// NOTE the #key is dropped for now ...
	// F3  string `is:"required:#key"`
	// F4  string `is:"required:@create:#key"`

	F5  string `is:"email"`
	F6  string `is:"url"`
	F7  string `is:"uri"`
	F8  string `is:"pan"`
	F9  string `is:"cvv"`
	F10 string `is:"ssn"`
	F11 string `is:"ein"`
	F12 string `is:"numeric"`
	F13 string `is:"hex"`
	F14 string `is:"hexcolor"`
	F15 string `is:"alnum"`
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

	F41 string `is:"prefix:foo bar"`
	F42 string `is:"prefix:&SomeValue"`

	F43 string `is:"suffix:foo bar"`
	F44 string `is:"suffix:&SomeValue"`

	F45 string  `is:"eq:foo bar"`
	F46 int     `is:"eq:-123"`
	F47 float64 `is:"eq:123.987"`
	F48 string  `is:"eq:&SomeValue"`

	F49 string  `is:"ne:foo bar"`
	F50 int     `is:"ne:-123"`
	F51 float64 `is:"ne:123.987"`
	F52 string  `is:"ne:&SomeValue"`

	F53 uint8   `is:"gt:24,lt:128"`
	F54 int16   `is:"gt:-128,lt:-24"`
	F55 float32 `is:"gt:0.24,lt:1.28"`

	F56 uint8   `is:"gte:24,lte:128"`
	F57 int16   `is:"gte:-128,lte:-24"`
	F58 float32 `is:"gte:0.24,lte:1.28"`

	F59 uint8   `is:"min:24,max:128"`
	F60 int16   `is:"min:-128,max:-24"`
	F61 float32 `is:"min:0.24,max:1.28"`

	F62 uint8   `is:"rng:24:128"`
	F63 int16   `is:"rng:-128:-24"`
	F64 float32 `is:"rng:0.24:1.28"`

	F65 string         `is:"len:28"`
	F66 []int          `is:"len:28"`
	F67 map[string]int `is:"len:28"`

	F68 string         `is:"len:4:28"`
	F69 []int          `is:"len:4:28"`
	F70 map[string]int `is:"len:4:28"`

	F71 string `is:"len::28"`
	F72 []int  `is:"len:4:"`

	G1 struct {
		F1 string `is:"required"`
	}

	// custom rule...
	F73 string    `is:"utf8"`
	F74 time.Time `is:"timecheck,ifacecheck"`

	// custom isValider
	F75 mypkg.MyString  `is:"required"`
	F76 mypkg.MyInt     `is:"required"`
	F77 *mypkg.MyString `is:"-isvalid"`
	F78 **mypkg.MyInt   `is:"isvalid:@create"`
	F79 *interface {
		IsValid() bool
	}

	// enums
	F80 someKind      `is:"enum"`
	F81 mypkg.MyKind  `is:"enum"`
	F82 **someKind    `is:"enum"`
	F83 *mypkg.MyKind `is:"enum:@create"`

	// elements & keys
	F84 []string          `is:"[]email"`
	F85 map[string]string `is:"[email]"`
	F86 map[string]string `is:"[phone:us:ca]zip:ca:us"`
	F87 map[string]*struct {
		F1 string `is:"len:2:32"`
		F2 string `is:"len:2:32"`
		F3 string `is:"phone"`
	} `is:"[email]"`

	// notnil
	F88 []string          `is:"notnil"`
	F89 map[string]string `is:"notnil"`
	F90 interface{}       `is:"notnil"`
	F91 *string           `is:"notnil"`

	// more nested elements & keys
	F92 []map[*map[string]string][]int `is:"[][[email]phone:us:ca]len::10,[]rng:-54:256"`
}

// local enum
type someKind string

const (
	SomeFoo someKind = "foo"
	SomeBar someKind = "bar"
	SomeBaz someKind = "baz"
)
