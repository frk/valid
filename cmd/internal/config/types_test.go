package config

import (
	"fmt"
	"testing"

	"github.com/frk/compare"
	"gopkg.in/yaml.v3"
)

func Test_String_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		data string
		dst  String
		want String
		err  error
		show bool
	}{{
		name: "dst should match want",
		data: `"foo"`,
		dst:  String{},
		want: String{Value: "foo", IsSet: true},
	}, {
		name: "override if IsSet is false",
		data: `"foo"`,
		dst:  String{Value: "bar", IsSet: false},
		want: String{Value: "foo", IsSet: true},
	}, {
		name: "empty string is a valid value",
		data: `""`,
		dst:  String{Value: "foo", IsSet: false},
		want: String{Value: "", IsSet: true},
	}, {
		name: "don't touch if IsSet is true",
		data: `"foo"`,
		dst:  String{Value: "bar", IsSet: true},
		want: String{Value: "bar", IsSet: true},
	}, {
		name: "don't touch if nil",
		data: `!!nil`,
		dst:  String{Value: "foo", IsSet: false},
		want: String{Value: "foo", IsSet: false},
	}, {
		name: "blow up if non-nil & non-stringy",
		data: `[ foo ]`,
		dst:  String{},
		want: String{},
		err: &Error{C: ERR_YAML_ERROR, tt: &String{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}, {
		name: "blow up if non-nil & non-stringy #2",
		data: `{ foo: bar }`,
		dst:  String{},
		want: String{},
		err: &Error{C: ERR_YAML_ERROR, tt: &String{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := yaml.Unmarshal([]byte(tt.data), &tt.dst)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if e := compare.Compare(tt.dst, tt.want); e != nil {
				t.Error(e)
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}

func Test_Bool_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		data string
		dst  Bool
		want Bool
		err  error
		show bool
	}{{
		name: "dst should match want",
		data: `true`,
		dst:  Bool{},
		want: Bool{Value: true, IsSet: true},
	}, {
		name: "dst should match want",
		data: `false`,
		dst:  Bool{},
		want: Bool{Value: false, IsSet: true},
	}, {
		name: "accept yes/no",
		data: `YES`,
		dst:  Bool{},
		want: Bool{Value: true, IsSet: true},
	}, {
		name: "override if IsSet is false",
		data: `false`,
		dst:  Bool{Value: true, IsSet: false},
		want: Bool{Value: false, IsSet: true},
	}, {
		name: "don't touch if nil",
		data: `!!nil`,
		dst:  Bool{Value: true, IsSet: false},
		want: Bool{Value: true, IsSet: false},
	}, {
		name: "don't touch if IsSet is true",
		data: `true`,
		dst:  Bool{Value: false, IsSet: true},
		want: Bool{Value: false, IsSet: true},
	}, {
		name: "blow up if non-booly",
		data: `"foo"`,
		dst:  Bool{Value: false, IsSet: false},
		want: Bool{Value: false, IsSet: false},
		err: &Error{C: ERR_YAML_ERROR, tt: &Bool{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}, {
		name: "blow up if non-booly #2",
		data: `[ foo ]`,
		dst:  Bool{},
		want: Bool{},
		err: &Error{C: ERR_YAML_ERROR, tt: &Bool{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}, {
		name: "blow up if non-booly #3",
		data: `{ foo: bar }`,
		dst:  Bool{},
		want: Bool{},
		err: &Error{C: ERR_YAML_ERROR, tt: &Bool{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := yaml.Unmarshal([]byte(tt.data), &tt.dst)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if e := compare.Compare(tt.dst, tt.want); e != nil {
				t.Error(e)
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}

func Test_StringSlice_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		data string
		dst  StringSlice
		want StringSlice
		err  error
		show bool
	}{{
		name: "dst should match want",
		data: `[ foo ]`,
		dst:  StringSlice{},
		want: StringSlice{Value: []string{"foo"}, IsSet: true},
	}, {
		name: "accept scalar",
		data: `foo`,
		dst:  StringSlice{},
		want: StringSlice{Value: []string{"foo"}, IsSet: true},
	}, {
		name: "override if IsSet is false",
		data: `[ foo ]`,
		dst:  StringSlice{Value: []string{"bar"}, IsSet: false},
		want: StringSlice{Value: []string{"foo"}, IsSet: true},
	}, {
		name: "override if IsSet is false (scalar)",
		data: `foo`,
		dst:  StringSlice{Value: []string{"bar"}, IsSet: false},
		want: StringSlice{Value: []string{"foo"}, IsSet: true},
	}, {
		name: "empty string is a valid value",
		data: `[ foo, "" ]`,
		dst:  StringSlice{},
		want: StringSlice{Value: []string{"foo", ""}, IsSet: true},
	}, {
		name: "don't touch if IsSet is true",
		data: `[ foo ]`,
		dst:  StringSlice{Value: []string{"bar"}, IsSet: true},
		want: StringSlice{Value: []string{"bar"}, IsSet: true},
	}, {
		name: "don't touch if IsSet is true (scalar)",
		data: `[ foo ]`,
		dst:  StringSlice{Value: []string{"bar"}, IsSet: true},
		want: StringSlice{Value: []string{"bar"}, IsSet: true},
	}, {
		name: "don't touch if nil",
		data: `!!nil`,
		dst:  StringSlice{Value: []string{"foo"}, IsSet: false},
		want: StringSlice{Value: []string{"foo"}, IsSet: false},
	}, {
		name: "blow up if non-stringy & non-sequency",
		data: `{ foo: bar }`,
		dst:  StringSlice{},
		want: StringSlice{},
		err:  &Error{C: ERR_YAML_TYPE, tt: &StringSlice{}, node: &yaml.Node{}},
	}, {
		name: "blow up if sequence of non-stringy values",
		data: `[ { foo: bar } ]`,
		dst:  StringSlice{},
		want: StringSlice{},
		err: &Error{C: ERR_YAML_ERROR, tt: &StringSlice{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := yaml.Unmarshal([]byte(tt.data), &tt.dst)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if e := compare.Compare(tt.dst, tt.want); e != nil {
				t.Error(e)
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}

func Test_ObjectIdent_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		data string
		dst  ObjectIdent
		want ObjectIdent
		err  error
		show bool
	}{{
		name: "dst should match want",
		data: `foo.Bar`,
		want: ObjectIdent{Pkg: "foo", Name: "Bar", IsSet: true},
	}, {
		name: "dst should match want #2",
		data: `foo.bar.baz`,
		want: ObjectIdent{Pkg: "foo.bar", Name: "baz", IsSet: true},
	}, {
		name: "dst should match want #3",
		data: `github.com/me/mod/pkg.FooBarBaz`,
		want: ObjectIdent{Pkg: "github.com/me/mod/pkg", Name: "FooBarBaz", IsSet: true},
	}, {
		name: "should blow up if bad format",
		data: `Foo`,
		err:  &Error{C: ERR_OBJECT_IDENT, tt: &ObjectIdent{}, node: &yaml.Node{}, val: `Foo`},
	}, {
		name: "should blow up if bad format #2",
		data: `.Foo`,
		err:  &Error{C: ERR_OBJECT_IDENT, tt: &ObjectIdent{}, node: &yaml.Node{}, val: `.Foo`},
	}, {
		name: "should blow up if bad format #3",
		data: `Foo.`,
		err:  &Error{C: ERR_OBJECT_IDENT, tt: &ObjectIdent{}, node: &yaml.Node{}, val: `Foo.`},
	}, {
		name: "should blow up if non-stringy",
		data: `[ foo.Bar ]`,
		err: &Error{C: ERR_YAML_ERROR, tt: &ObjectIdent{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}, {
		name: "should blow up if non-stringy #2",
		data: `{ foo: bar }`,
		err: &Error{C: ERR_YAML_ERROR, tt: &ObjectIdent{},
			node: &yaml.Node{}, err: &yaml.TypeError{}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := yaml.Unmarshal([]byte(tt.data), &tt.dst)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if e := compare.Compare(tt.dst, tt.want); e != nil {
				t.Error(e)
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}

func Test_Scalar_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		data string
		want Scalar
		err  error
		show bool
	}{{
		name: "decode !!nil",
		data: `!!nil`,
		want: Scalar{Type: NIL},
	}, {
		name: "decode bool",
		data: `false`,
		want: Scalar{Type: BOOL, Value: "false"},
	}, {
		name: "decode string",
		data: `foo`,
		want: Scalar{Type: STRING, Value: "foo"},
	}, {
		name: "decode float",
		data: `3.14`,
		want: Scalar{Type: FLOAT, Value: "3.14"},
	}, {
		name: "decode int",
		data: `42`,
		want: Scalar{Type: INT, Value: "42"},
	}, {
		name: "blow up if not scalar",
		data: `[ 42 ]`,
		want: Scalar{},
		err:  &Error{C: ERR_YAML_TYPE, tt: &Scalar{}, node: &yaml.Node{}},
	}, {
		name: "blow up if not scalar #2",
		data: `{ foo: bar }`,
		want: Scalar{},
		err:  &Error{C: ERR_YAML_TYPE, tt: &Scalar{}, node: &yaml.Node{}},
	}}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := Scalar{}
			err := yaml.Unmarshal([]byte(tt.data), &dst)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if e := compare.Compare(dst, tt.want); e != nil {
				t.Error(e)
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}

func Test_JoinOp_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		data string
		want JoinOp
		err  error
		show bool
	}{
		{name: "skip !!nil", data: `!!nil`, want: 0},
		{name: "decode lower-case", data: `and`, want: JOIN_AND},
		{name: "decode lower-case", data: `or`, want: JOIN_OR},
		{name: "decode lower-case", data: `not`, want: JOIN_NOT},
		{name: "decode upper-case", data: `AND`, want: JOIN_AND},
		{name: "decode upper-case", data: `OR`, want: JOIN_OR},
		{name: "decode upper-case", data: `NOT`, want: JOIN_NOT},
		{
			name: "blow up if not valid join op",
			data: `foo`,
			err:  &Error{C: ERR_JOIN_OP, tt: _ptr[JoinOp](0), val: "foo", node: &yaml.Node{}},
		}, {
			name: "blow up if not valid join op",
			data: `123`,
			err:  &Error{C: ERR_JOIN_OP, tt: _ptr[JoinOp](0), val: "123", node: &yaml.Node{}},
		}, {
			name: "blow up if non-stringy",
			data: `[ and ]`,
			err:  &Error{C: ERR_YAML_TYPE, tt: _ptr[JoinOp](0), node: &yaml.Node{}},
		}, {
			name: "blow up if non-stringy",
			data: `{ and: or }`,
			err:  &Error{C: ERR_YAML_TYPE, tt: _ptr[JoinOp](0), node: &yaml.Node{}},
		},
	}

	compare := compare.Config{ObserveFieldTag: "cmp"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst JoinOp
			err := yaml.Unmarshal([]byte(tt.data), &dst)
			if e := compare.Compare(err, tt.err); e != nil {
				t.Errorf("%#v\n%v\n", err, e)
			}
			if e := compare.Compare(dst, tt.want); e != nil {
				t.Error(e)
			}
			if tt.show && tt.err != nil {
				fmt.Println(err)
			}
		})
	}
}
