package generate

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_writer(t *testing.T) {
	tests := []struct {
		skip bool
		text string
		args []any
		vars map[string]any
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
	}, {
		text: "hello $key",
		vars: map[string]any{"key": "world"},
		want: "hello world",
	}, {
		text: "hi $k1, $k2, and $0",
		args: []any{"baz"},
		vars: map[string]any{"k1": "foo", "k2": "bar"},
		want: "hi foo, bar, and baz",
	}, {
		text: "hello ${key}world",
		vars: map[string]any{"key": "blue"},
		want: "hello blueworld",
	}, {
		text: "hello ${k1}world${k2}sky",
		vars: map[string]any{"k1": "blue", "k2": "red"},
		want: "hello blueworldredsky",
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

		w := writer{vars: tt.vars}
		w.p(tt.text, tt.args...)
		if !reflect.DeepEqual(tt.err, w.err) {
			t.Errorf("got error: %v\nwant error: %v", w.err, tt.err)
		}
		if got := w.buf.String(); got != tt.want {
			t.Errorf("got text: %v\nwant text: %v", got, tt.want)
		}
	}

}
