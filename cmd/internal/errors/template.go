package errors

import (
	"fmt"
	"strings"
	"text/template"
)

func ParseTemplate(text string) { t = template.Must(t.Parse(text)) }

func String(name string, data interface{}) string {
	sb := new(strings.Builder)
	if err := t.ExecuteTemplate(sb, name, data); err != nil {
		panic(err)
	}
	return sb.String()
}

var t = template.New("t").Funcs(template.FuncMap{
	// red color (terminal)
	"r":  func(v ...string) string { return getcolor("\033[0;31m", v) },
	"rb": func(v ...string) string { return getcolor("\033[1;31m", v) },
	"ri": func(v ...string) string { return getcolor("\033[3;31m", v) },
	"ru": func(v ...string) string { return getcolor("\033[4;31m", v) },
	// yellow color (terminal)
	"y":  func(v ...string) string { return getcolor("\033[0;33m", v) },
	"yb": func(v ...string) string { return getcolor("\033[1;33m", v) },
	"yi": func(v ...string) string { return getcolor("\033[3;33m", v) },
	"yu": func(v ...string) string { return getcolor("\033[4;33m", v) },
	// white color (terminal)
	"w":  func(v ...string) string { return getcolor("\033[0;37m", v) },
	"wb": func(v ...string) string { return getcolor("\033[1;37m", v) },
	"wi": func(v ...string) string { return getcolor("\033[3;37m", v) },
	"wu": func(v ...string) string { return getcolor("\033[4;37m", v) },
	// cyan color (terminal)
	"c":  func(v ...string) string { return getcolor("\033[0;36m", v) },
	"cb": func(v ...string) string { return getcolor("\033[1;36m", v) },
	"ci": func(v ...string) string { return getcolor("\033[3;36m", v) },
	"cu": func(v ...string) string { return getcolor("\033[4;36m", v) },

	/////////////////////////////////////////////////////////////////////////
	// High Intensity
	/////////////////////////////////////////////////////////////////////////

	// red color HI (terminal)
	"R":  func(v ...string) string { return getcolor("\033[0;91m", v) },
	"Rb": func(v ...string) string { return getcolor("\033[1;91m", v) },
	"Ri": func(v ...string) string { return getcolor("\033[3;91m", v) },
	"Ru": func(v ...string) string { return getcolor("\033[4;91m", v) },
	// green color HI (terminal)
	"G":  func(v ...string) string { return getcolor("\033[0;92m", v) },
	"Gb": func(v ...string) string { return getcolor("\033[1;92m", v) },
	"Gi": func(v ...string) string { return getcolor("\033[3;92m", v) },
	"Gu": func(v ...string) string { return getcolor("\033[4;92m", v) },
	// yellow color HI (terminal)
	"Y":  func(v ...string) string { return getcolor("\033[0;93m", v) },
	"Yb": func(v ...string) string { return getcolor("\033[1;93m", v) },
	"Yi": func(v ...string) string { return getcolor("\033[3;93m", v) },
	"Yu": func(v ...string) string { return getcolor("\033[4;93m", v) },
	// blue color HI (terminal)
	"B":  func(v ...string) string { return getcolor("\033[0;94m", v) },
	"Bb": func(v ...string) string { return getcolor("\033[1;94m", v) },
	"Bi": func(v ...string) string { return getcolor("\033[3;94m", v) },
	"Bu": func(v ...string) string { return getcolor("\033[4;94m", v) },
	// cyan color HI (terminal)
	"C":  func(v ...string) string { return getcolor("\033[0;96m", v) },
	"Cb": func(v ...string) string { return getcolor("\033[1;96m", v) },
	"Ci": func(v ...string) string { return getcolor("\033[3;96m", v) },
	"Cu": func(v ...string) string { return getcolor("\033[4;96m", v) },
	// white color HI (terminal)
	"W":  func(v ...string) string { return getcolor("\033[0;97m", v) },
	"Wb": func(v ...string) string { return getcolor("\033[1;97m", v) },
	"Wi": func(v ...string) string { return getcolor("\033[3;97m", v) },
	"Wu": func(v ...string) string { return getcolor("\033[4;97m", v) },

	// no color (terminal)
	"off": func() string { return "\033[0m" },

	"raw":   func(s string) string { return "`" + s + "`" },
	"quote": func(s string) string { return `"` + s + `"` },
	"Up":    strings.ToUpper,

	// error note
	"ERROR":  func() string { return getcolor("\033[0;91m", []string{"ERROR:"}) },
	"ERRCFG": func() string { return getcolor("\033[0;91m", []string{"ERROR(config):"}) },
	// new line
	"NL": func() string { return "\n" },
	// new line & tab
	"NT": func() string { return "\n\t" },
})

func getcolor(c string, v []string) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v\033[0m", c, stringsStringer(v))
	}
	return c
}

type stringsStringer []string

func (s stringsStringer) String() string {
	return strings.Join([]string(s), "")
}
