package testdata

import (
	"github.com/frk/valid/cmd/internal/generator/testdata/mypkg"
)

type myenum uint

const (
	myenum0 myenum = iota
	myenum1
	myenum2
	_
	myenum4
	_
	myenum6
)

type gibberish string

const (
	gibfoo  gibberish = "foo"
	gibbar  gibberish = "bar"
	gibbaz  gibberish = "baz"
	gibquux gibberish = "quux"
)

type T14Validator struct {
	F1 myenum       `is:"enum"`
	F2 mypkg.MyKind `is:"enum"`
	F3 **gibberish  `is:"len:3,enum,required"`

	F4 *mypkg.MyEnumIsValid `is:"enum,-isvalid"`

	ea errorAggregator
}
