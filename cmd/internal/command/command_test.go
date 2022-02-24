package command

import (
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/frk/valid/cmd/internal/config"

	"github.com/frk/compare"
)

func Test_parseFlags(t *testing.T) {
	tests := []struct {
		args []string
		want config.Config
		err  error
	}{{
		args: []string{},
		want: config.Config{},
	}, {
		args: []string{
			`-c`, `/path/to/my/config`,
			`-wd`, `/path/to/my/dir`,
			`-r=false`,
			`-f=/path/to/my/file1.go`,
			`-f`, `/path/to/my/file_2.go`,
			`-rx=^\/path\/to\/my\/\w+_foo.go$`,
			`-rx`, `^\/path\/to\/my\/\w+_bar.go$`,
			`-o`, `%_out.go`,
			`-fk.tag`, `json`,
			`-fk.join`,
			`-fk.sep=.`,
			`-error.constructor`, `example.com/me/mymod/mypkg.NewError`,
			`-error.aggregator`, `example.com/me/mymod/mypkg.MyErrorAggregator`,
		},
		want: config.Config{
			File:      config.String{Value: "/path/to/my/config", IsSet: true},
			WorkDir:   config.String{Value: "/path/to/my/dir", IsSet: true},
			Recursive: config.Bool{Value: false, IsSet: true},
			FileList: config.StringSlice{
				Value: []string{
					"/path/to/my/file1.go",
					"/path/to/my/file_2.go",
				},
				IsSet: true,
			},
			FilePatternList: config.StringSlice{
				Value: []string{
					"^\\/path\\/to\\/my\\/\\w+_foo.go$",
					"^\\/path\\/to\\/my\\/\\w+_bar.go$",
				},
				IsSet: true,
			},
			OutNameFormat: config.String{Value: "%_out.go", IsSet: true},
			ErrorHandling: config.ErrorHandlingConfig{
				FieldKey: config.FieldKeyConfig{
					Tag:       config.String{Value: "json", IsSet: true},
					Join:      config.Bool{Value: true, IsSet: true},
					Separator: config.String{Value: ".", IsSet: true},
				},
				Constructor: config.ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "NewError",
					IsSet: true,
				},
				Aggregator: config.ObjectIdent{
					Pkg:   "example.com/me/mymod/mypkg",
					Name:  "MyErrorAggregator",
					IsSet: true,
				},
			},
		},
	}, {
		args: []string{`-r=foo`},
		err:  fmt.Errorf(`invalid boolean value "foo" for -r: strconv.ParseBool: parsing "foo": invalid syntax`),
	}, {
		args: []string{`-fk.join=123`},
		err:  fmt.Errorf(`invalid boolean value "123" for -fk.join: strconv.ParseBool: parsing "123": invalid syntax`),
	}}

	for _, tt := range tests {
		var got config.Config

		fs := flag.NewFlagSet("test", 0)
		fs.SetOutput(io.Discard)

		err := parseFlags(&got, fs, tt.args)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Errorf("%v\n%#v\n", e, err)
		}
		if e := compare.Compare(got, tt.want); e != nil {
			t.Error(e)
		}
	}
}
