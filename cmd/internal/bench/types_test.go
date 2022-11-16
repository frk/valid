package bench

import (
	"testing"

	goplayground "github.com/go-playground/validator/v10"
)

func Test(t *testing.T) {
	v1 := goodFrk()
	if err := v1.Validate(); err != nil {
		t.Errorf("frk/valid failed: %v", err)
	}

	v2 := goodGoPlayground()
	vv := goplayground.New()
	if err := vv.Struct(v2); err != nil {
		t.Errorf("go-playground/validator failed: %v", err)
	}
}

func BenchmarkFrk(b *testing.B) {
	value := goodFrk()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := value.Validate(); err != nil {
			panic(err) // shouldn't fail
		}
	}
}

func BenchmarkGoPlayground(b *testing.B) {
	vv := goplayground.New()
	value := goodGoPlayground()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := vv.Struct(value); err != nil {
			panic(err) // shouldn't fail
		}
	}
}

func _ptr[T any](v T) *T { return &v }

func goodFrk() *FrkValidator {
	return &FrkValidator{
		F01: "foobar",                                       // `is:"required"`
		F02: _ptr("foobar"),                                 // `is:"required"`
		F03: _ptr(float32(3.14)),                            // `is:"eq:3.14"`
		F04: _ptr(int(128)),                                 // `is:"gt:123"`
		F05: 42,                                             // `is:"lte:42"`
		F06: 5.12,                                           // `is:"rng:3.14:42"`
		F07: 32,                                             // `is:"min:3.14,max:42"`
		F08: "1234567890",                                   // `is:"len:8:256"`
		F09: []uint8("0123456789AB"),                        // `is:"len:12"`
		F10: map[string]struct{}{"a": {}, "b": {}, "c": {}}, // `is:"len::12"`
		F11: []byte("世界世界"),                                 // `is:"runecount:4:"`
		F12: "Hello, 世界",                                    // `is:"runecount::15"`
		F13: _ptr("hello world"),                            // `is:"prefix:hello,required"`
		F14: _ptr("hi & hello, world!"),                     // `is:"contains:hello"`
		F15: "abc123",                                       // `is:"alnum:en"`
		F16: "135.104.0.1/24",                               // `is:"cidr"`
		F17: _ptr("foo--bar.com"),                           // `is:"fqdn"`
		F18: _ptr("127.0.0.1"),                              // `is:"ip"`
		F19: "ben@example.com",                              // `is:"email"`
		F20: "8005552222",                                   // `is:"phone"`
		Sub1: &FrkSub1Type{
			F01: []string{"foo@bbb.com", "bar@bbb.com", "baz@bbb.com"},           // `is:"len::10,[]email"`
			F02: []string{"8005552222", "19876543210"},                           // `is:"required,[]phone"`
			F03: &map[string]string{"foo@bar.com": "foo", "bar@quux.com": "bar"}, // `is:"required,[email]"`
			F04: map[string]string{
				"foo@bar.com":  "8005552222",
				"bar@quux.com": "19876543210",
			}, // `is:"[email]phone"`
			F05: map[string]int{
				"15673628910": 33,
				"8005552222":  43,
				"19876543210": 50,
				"15673628912": 55,
				"8005553333":  87,
			}, // `is:"len:5:,[phone]rng:21:99"`
			F06: _ptr("abc[123]"),    // `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
			F07: "world goodbye",     // `is:"suffix:goodbye"`
			F08: _ptr("foo and bar"), // `is:"eq:&Arg1"`
			F09: 511,                 // `is:"lt:&Arg2"`
			F10: 42.2,                // `is:"gt:&Arg3"`
		},
		Sub2: &FrkSub2Type{
			F01: []map[string][]string{
				{"713ae7e3-cb32-45f9-adcb-7c4fa86b90c1": {"20505", "22313"}},
				{"625e63f3-58f5-40b7-83a1-a72ad31acffb": {"08109"}},
				{
					"57b73598-8764-4ad0-a76a-679bb6640eb1": {"23220"},
					"9c858901-8a57-4791-81fe-4c455b099bc9": {"33311"},
				},
			}, // `is:"len:3,[][uuid]len:1:,[]zip"`
			F02: map[[2]string]any{
				{"127.0.0.1", "0.0.0.0"}:       "foo",
				{"255.255.255.255", "1.2.3.4"}: 123,
			}, // `is:"[[]ip]notnil"`
			F03: &struct {
				f1 string `is:"eq:&Sub2.Arg1"`
				f2 int32  `is:"lt:&Arg2,required"`
			}{"bar and baz", 128},
			F04: "foo bar quux",          // `is:"prefix:foo,contains:bar,suffix:baz:quux,len:9:64"`
			F05: _ptr("ben@example.com"), // `is:"email,optional"`
			F06: []*string{
				_ptr("ben@example.com"),
				_ptr("mark@example.com"),
			}, // `is:"[]email,omitnil"`
			F07: 111,                 // `is:"max:128,optional"`
			F08: _ptr("bar and baz"), // `is:"eq:.Arg1"`
			F09: 111,                 // `is:"lt:.Arg2"`
			F10: 43.2,                // `is:"gt:.Arg3"`
			Sub: &FrkSub1Type{
				F01: []string{"foo@bbb.com", "bar@bbb.com", "baz@bbb.com"},           // `is:"len::10,[]email"`
				F02: []string{"8005552222", "19876543210"},                           // `is:"required,[]phone"`
				F03: &map[string]string{"foo@bar.com": "foo", "bar@quux.com": "bar"}, // `is:"required,[email]"`
				F04: map[string]string{
					"foo@bar.com":  "8005552222",
					"bar@quux.com": "19876543210",
				}, // `is:"[email]phone"`
				F05: map[string]int{
					"15673628910": 33,
					"8005552222":  43,
					"19876543210": 50,
					"15673628912": 55,
					"8005553333":  87,
				}, // `is:"len:5:,[phone]rng:21:99"`
				F06: _ptr("abc[123]"),    // `is:"re:\"^[a-z]+\\[[0-9]+\\]$\""`
				F07: "world goodbye",     // `is:"suffix:goodbye"`
				F08: _ptr("foo and bar"), // `is:"eq:&Arg1"`
				F09: 511,                 // `is:"lt:&Arg2"`
				F10: 42.2,                // `is:"gt:&Arg3"`
			},

			Arg1: "bar and baz",
			Arg2: 128,
			Arg3: 43.1,
		},

		Arg1: "foo and bar",
		Arg2: 512,
		Arg3: 42.1,
	}
}

func goodGoPlayground() *GoPlaygroundParams {
	return &GoPlaygroundParams{
		F01: "foobar",                                       // `validate:"required"`
		F02: _ptr("foobar"),                                 // `validate:"required"`
		F03: _ptr(float64(3.14)),                            // `validate:"eq=3.14"`
		F04: _ptr(int(128)),                                 // `validate:"gt=123"`
		F05: 42,                                             // `validate:"lte=42"`
		F06: 5.12,                                           // `validate:"min=3.14,max=42"` // no range/between i think
		F07: 32,                                             // `validate:"min=3.14,max=42"`
		F08: "1234567890",                                   // `validate:"min=8,max=256"` // no len-range i think
		F09: []uint8("0123456789AB"),                        // `validate:"len=12"`
		F10: map[string]struct{}{"a": {}, "b": {}, "c": {}}, // `validate:"max=12"`
		F11: []byte("世界世界"),                                 // `validate:"min=4"`  // no rune count
		F12: "Hello, 世界",                                    // `validate:"max=15"` // no rune count
		F13: _ptr("hello world"),                            // `validate:"startswith=hello,required"`
		F14: _ptr("hi & hello, world!"),                     // `validate:"contains=hello"`
		F15: "abc123",                                       // `validate:"alphanum"`
		F16: "135.104.0.1/24",                               // `validate:"cidr"`
		F17: _ptr("foo--bar.com"),                           // `validate:"fqdn"`
		F18: _ptr("127.0.0.1"),                              // `validate:"ip_addr"`
		F19: "ben@example.com",                              // `validate:"email"`
		F20: "+8005552222",                                  // `validate:"e164"`
		Sub1: &GoPlaygroundSub1Type{
			F01: []string{"foo@bbb.com", "bar@bbb.com", "baz@bbb.com"},           // `validate:"max=10,dive,email"`
			F02: []string{"+1123456789", "+4124418875"},                          // `validate:"required,dive,e164"`
			F03: &map[string]string{"foo@bar.com": "foo", "bar@quux.com": "bar"}, // `validate:"required,dive,keys,email,endkeys"`
			F04: map[string]string{
				"foo@bar.com":  "+1123456789",
				"bar@quux.com": "+4124418875",
			}, // `validate:"dive,keys,email,endkeys,e164"`
			F05: map[string]int{
				"+15673628910": 33,
				"+8005552222":  43,
				"+19876543210": 50,
				"+15673627910": 55,
				"+1123456789":  87,
			}, // `validate:"min=5,dive,keys,e164,endkeys,min=21,max=99"`
			F06: _ptr("abc123"),      // `validate:"alphanum,lowercase"` // no regexp support
			F07: "world goodbye",     // `validate:"endswith=goodbye"`
			F08: _ptr("foo and bar"), // `validate:"eqfield=Arg1"`
			F09: 511,                 // `validate:"ltfield=Arg2"`
			F10: 42.2,                // `validate:"gtfield=Arg3"`

			Arg1: "foo and bar",
			Arg2: 512,
			Arg3: 42.1,
		},
		Sub2: &GoPlaygroundSub2Type{
			F01: []map[string][]string{
				{"713ae7e3-cb32-45f9-adcb-7c4fa86b90c1": {"+15673628910", "+8005552222"}},
				{"625e63f3-58f5-40b7-83a1-a72ad31acffb": {"+19876543210"}},
				{
					"57b73598-8764-4ad0-a76a-679bb6640eb1": {"+15673628910"},
					"9c858901-8a57-4791-81fe-4c455b099bc9": {"+1123456789"},
				},
			}, // `validate:"len=3,dive,dive,keys,uuid,endkeys,min=1,dive,e164"` // no zip :(
			F02: map[[2]string]any{
				{"127.0.0.1", "0.0.0.0"}:       "foo",
				{"255.255.255.255", "1.2.3.4"}: 123,
			}, // `validate:"dive,keys,dive,ip_addr,endkeys,required"`
			F03: &struct {
				f1 string `validate:"eqcsfield=Sub2.Arg1"`
				f2 int32  `validate:"ltcsfield=Arg2,required"`
			}{"bar and baz", 128},
			F04: "foo bar quux baz",      // `validate:"startswith=foo,contains=bar,endswith=baz,min=9,max=64"`
			F05: _ptr("ben@example.com"), // `validate:"email,omitempty"`
			F06: []*string{
				_ptr("ben@example.com"),
				_ptr("mark@example.com"),
			}, // `validate:"dive,email,omitempty"`
			F07: 111,                 // `validate:"max=128,omitempty"`
			F08: _ptr("bar and baz"), // `validate:"eqcsfield=Arg1"`
			F09: 111,                 // `validate:"ltcsfield=Arg2"`
			F10: 43.2,                // `validate:"gtcsfield=Arg3"`
			Sub: &GoPlaygroundSub1Type{
				F01: []string{"foo@bbb.com", "bar@bbb.com", "baz@bbb.com"},           // `validate:"max=10,dive,email"`
				F02: []string{"+1123456789", "+4124418875"},                          // `validate:"required,dive,e164"`
				F03: &map[string]string{"foo@bar.com": "foo", "bar@quux.com": "bar"}, // `validate:"required,dive,keys,email,endkeys"`
				F04: map[string]string{
					"foo@bar.com":  "+1123456789",
					"bar@quux.com": "+4124418875",
				}, // `validate:"dive,keys,email,endkeys,e164"`
				F05: map[string]int{
					"+15673628910": 33,
					"+8005552222":  43,
					"+19876543210": 50,
					"+15673628917": 55,
					"+1123456789":  87,
				}, // `validate:"min=5,dive,keys,e164,endkeys,min=21,max=99"`
				F06: _ptr("abc123"),      // `validate:"alphanum,lowercase"` // no regexp support
				F07: "world goodbye",     // `validate:"endswith=goodbye"`
				F08: _ptr("bar and baz"), // `validate:"eqfield=Arg1"`
				F09: 111,                 // `validate:"ltfield=Arg2"`
				F10: 43.2,                // `validate:"gtfield=Arg3"`

				Arg1: "bar and baz",
				Arg2: 128,
				Arg3: 43.1,
			},
		},
	}
}
