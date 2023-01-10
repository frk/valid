package generate

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func Test_generatorP(t *testing.T) {
	tests := []struct {
		skip bool
		text string
		args []any
		// varz map[string]string
		want string
		err  error
	}{{
		text: "hello world",
		want: "hello world",
	}, {
		text: "hello world$",
		want: "hello world$",
	}, {
		text: "hello $$world",
		want: "hello $world",
	}, {
		text: "hello $0",
		args: []any{"world"},
		want: "hello world",
	}, {
		text: "hello $0$1",
		args: []any{"world", "!!"},
		want: "hello world!!",
	}, {
		text: "hi $0, $1, and $2",
		args: []any{"foo", "bar", "baz"},
		want: "hi foo, bar, and baz",
		//}, {
		//	text: "hello $key",
		//	varz: map[string]string{"key": "world"},
		//	want: "hello world",
		//}, {
		//	text: "hi $k1, $k2, and $0",
		//	args: []any{"baz"},
		//	varz: map[string]string{"k1": "foo", "k2": "bar"},
		//	want: "hi foo, bar, and baz",
		//}, {
		//	text: "hello ${key}world",
		//	varz: map[string]string{"key": "blue"},
		//	want: "hello blueworld",
		//}, {
		//	text: "hello ${k1}world${k2}sky",
		//	varz: map[string]string{"k1": "blue", "k2": "red"},
		//	want: "hello blueworldredsky",
	}, {
		text: "hello $26world$27sky",
		args: append(make([]any, 26), "blue", "red"),
		want: "hello blueworldredsky",
	}, {
		text: "hello $1world$2sky",
		args: []any{"blue", "red"},
		err:  fmt.Errorf("index $2 out of range"),
	}, {
		skip: true,
		// for this error to occur the isDigit func would
		// have to be broken, so just skip the test case
		text: "hello $0world$1asky",
		args: []any{"blue", "red"},
		err:  fmt.Errorf("failed to parse $1a: strconv.Atoi: parsing \"1a\": invalid syntax"),
	}}

	for _, tt := range tests {
		if tt.skip {
			continue
		}

		g := generator{ /*varz: tt.varz*/ }
		g.P(tt.text, tt.args...)
		if !reflect.DeepEqual(tt.err, g.werr) {
			t.Errorf("got error: %v\nwant error: %v", g.werr, tt.err)
		}
		if g.werr == nil {
			if got := g.buf.String(); got != tt.want {
				t.Errorf("got text: %v\nwant text: %v", got, tt.want)
			}
		}
	}

}

func Test_generatorRL(t *testing.T) {
	tests := []struct {
		skip bool
		buf  string
		text string
		args []any
		want string
		err  error
	}{{
		buf:  "hello world",
		text: "goodbye!",
		want: "goodbye!\n",
	}, {
		buf:  "hello\nworld",
		text: "goodbye!",
		want: "hello\ngoodbye!\n",
	}, {
		buf:  "hello\nworld\n",
		text: "goodbye!",
		want: "hello\ngoodbye!\n",
	}}

	for _, tt := range tests {
		if tt.skip {
			continue
		}

		buf := bytes.NewBuffer([]byte(tt.buf))

		g := generator{buf: *buf}
		g.RL(tt.text, tt.args...)
		if !reflect.DeepEqual(tt.err, g.werr) {
			t.Errorf("got error: %v\nwant error: %v", g.werr, tt.err)
		}
		if g.werr == nil {
			if got := g.buf.String(); got != tt.want {
				t.Errorf("got text: %v\nwant text: %v", got, tt.want)
			}
		}
	}

}
