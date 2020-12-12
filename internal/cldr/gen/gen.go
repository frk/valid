package main

import (
	"bytes"
	"fmt"
	"go/format"
	"golang.org/x/text/unicode/cldr"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	GO "github.com/frk/ast/golang"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer f.Close()

	var d cldr.Decoder
	d.SetDirFilter("main", "supplemental")

	cldr, err := d.DecodeZip(f)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	setnumerics(cldr)

	locs := getLocales(cldr)
	inheritMissingInfo(locs, nil)

	file := buildTableFile(locs)
	if err := writeTableFile(file); err != nil {
		fmt.Println("ERROR:", err)
		return
	}
}

var numerics map[string]string

type locale struct {
	lang    string
	digits  string
	decimal rune
	group   rune
}

// collect numberingSystems type=numeric data.
func setnumerics(cldr *cldr.CLDR) {
	numerics = make(map[string]string)

	s := cldr.Supplemental()
	for _, ns := range s.NumberingSystems.NumberingSystem {
		if ns.Type == "numeric" {
			numerics[ns.Id] = ns.Digits
		}
	}
}

func getLocales(cldr *cldr.CLDR) (locs []locale) {
	for _, lang := range cldr.Locales() {
		if lang == "root" {
			continue
		}

		loc := locale{lang: lang}
		ldml := cldr.RawLDML(loc.lang)
		getLocaleNumberSymbols(ldml, &loc)
		locs = append(locs, loc)
	}
	return locs
}

func getLocaleNumberSymbols(ldml *cldr.LDML, loc *locale) {
	nums := ldml.Numbers
	if nums == nil {
		return
	}

	var nsid string
	if len(nums.DefaultNumberingSystem) > 0 {
		nsid = strings.TrimSpace(nums.DefaultNumberingSystem[0].Data())
		loc.digits = numerics[nsid]
	} else if len(nums.Symbols) > 0 {
		nsid = strings.TrimSpace(nums.Symbols[0].NumberSystem)
		loc.digits = numerics[nsid]
	}

	for _, sym := range nums.Symbols {
		if sym.NumberSystem == nsid {
			if len(sym.Decimal) > 0 {
				sep := sym.Decimal[0].Data()
				r, _ := utf8.DecodeRune([]byte(sep))
				if r != utf8.RuneError {
					loc.decimal = r
				}
			}
			if len(sym.Group) > 0 {
				sep := sym.Group[0].Data()
				r, _ := utf8.DecodeRune([]byte(sep))
				if r != utf8.RuneError {
					loc.group = r
				}
			}
			break
		}
	}
}

func inheritMissingInfo(locs []locale, loc *locale) {
	if loc == nil {
		for i, loc := range locs {
			if loc.group > 0 && loc.decimal > 0 && len(loc.digits) > 0 {
				continue
			}

			inheritMissingInfo(locs, &loc)
			locs[i] = loc
		}
		return
	}

	lang := loc.lang
	for strings.ContainsRune(lang, '_') {
		lang = lang[:strings.LastIndexByte(lang, '_')]

		for _, parent := range locs {
			if parent.lang == lang {
				if loc.group == 0 {
					loc.group = parent.group
				}
				if loc.decimal == 0 {
					loc.decimal = parent.decimal
				}
				if loc.digits == "" {
					loc.digits = parent.digits
				}

				if loc.group > 0 && loc.decimal > 0 && len(loc.digits) > 0 {
					return // done
				} else {
					break // try next parent
				}

			}
		}
	}
}

func buildTableFile(locs []locale) *GO.File {
	locales := buildLocaleInfoSlice(locs)

	file := new(GO.File)
	file.PkgName = "cldr"
	file.Decls = append(file.Decls, locales)
	return file
}

func buildLocaleInfoSlice(locs []locale) (decl GO.VarDecl) {
	slice := GO.SliceLit{Type: GO.SliceType{GO.Ident{"LocaleInfo"}}}
	elems := GO.ExprList{}
	for _, loc := range locs {
		f1 := GO.FieldElement{Field: "Lang", Value: GO.StringLit(loc.lang)}

		f2 := GO.FieldElement{Field: "SepDecimal"}
		if loc.decimal > 0 {
			f2.Value = GO.RuneLit(loc.decimal)
		} else {
			f2.Value = GO.IntLit(0)
		}

		f3 := GO.FieldElement{Field: "SepGroup"}
		if loc.group > 0 {
			f3.Value = GO.RuneLit(loc.group)
		} else {
			f3.Value = GO.IntLit(0)
		}

		f4 := GO.FieldElement{Field: "DigitZero"}
		if r, _ := utf8.DecodeRune([]byte(loc.digits)); r != utf8.RuneError {
			f4.Value = GO.RuneLit(r)
		} else {
			f4.Value = GO.IntLit(0)
		}

		f5 := GO.FieldElement{Field: "DigitNine"}
		if r, _ := utf8.DecodeLastRune([]byte(loc.digits)); r != utf8.RuneError {
			f5.Value = GO.RuneLit(r)
		} else {
			f5.Value = GO.IntLit(0)
		}

		elem := GO.StructLit{Elems: []GO.FieldElement{f1, f2, f3, f4, f5}, Compact: true}
		elems = append(elems, elem)
	}
	slice.Elems = elems
	decl.Spec = GO.ValueSpec{Names: GO.Ident{"localeslice"}, Values: slice}
	return decl
}

func writeTableFile(file *GO.File) (err error) {
	buf := &bytes.Buffer{}
	if err := GO.Write(file, buf); err != nil {
		return err
	}

	path := "../tables.go"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		if err != nil {
			os.Remove(path)
		}
	}()

	// make it look pretty
	bs, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	buf = bytes.NewBuffer(bs)
	if _, err := io.Copy(f, buf); err != nil {
		return err
	}

	return f.Sync()
}
