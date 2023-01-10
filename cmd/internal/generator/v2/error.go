package generate

import (
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

func (g *generator) genErrorExpr(r *rules.Rule) {
	errSpec := r.ErrSpec()
	text := errSpec.Text
	if len(text) == 0 {
		text = "is not valid"
	}

	o := g.info.RuleObjMap[r]
	f := g.info.ObjFieldMap[o]
	text = g.info.FKeyMap[f] + " " + text

	//var refs GO.ExprList
	if errSpec.WithArgs {
		var args []string
		for _, a := range r.Args {
			// A rule argument of unknown kind for
			// a numeric type can be treated as 0.
			if a.Type == rules.ARG_UNKNOWN && o.Type.Kind.IsNumeric() {
				a = &rules.Arg{rules.ARG_INT, "0"}
			}
			// skip empty
			if a.Value == "" {
				continue
			}

			switch a.Type {
			case rules.ARG_FIELD_ABS, rules.ARG_FIELD_REL:
				// x := b.g.recv
				// for _, f := range b.g.info.KeyMap[a.Value].Selector {
				// 	x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
				// }
				args = append(args, "%v")
				// refs = append(refs, x)
			case rules.ARG_STRING:
				args = append(args, strconv.Quote(a.Value))
			default:
				args = append(args, a.Value)
			}
		}
		if len(args) > 0 {
			text += ": " + strings.Join(args, errSpec.ArgSep)
			if len(errSpec.ArgSuffix) > 0 {
				text += " " + errSpec.ArgSuffix
			}
		}
	}

	text = strconv.Quote(text)

	//if len(refs) > 0 {
	//	pkg := b.g.addImport(xtypes.Pkg{Path: "fmt"})
	//	body.Add(GO.ReturnStmt{GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, "Errorf"},
	//		Args: GO.ArgsList{List: append(GO.ExprList{textExpr}, refs...)}}})
	//} else {

	pkg := g.file.addImport(types.Pkg{Path: "errors"})
	g.P("$0.New($1)", pkg.name, text)
}

func (g *generator) ErrExpr(o *types.Obj, r *rules.Rule) (F func()) {
	errSpec := r.ErrSpec()
	text := errSpec.Text
	if len(text) == 0 {
		text = "is not valid"
	}
	f := g.info.ObjFieldMap[o]
	text = g.info.FKeyMap[f] + " " + text

	//var refs GO.ExprList
	if errSpec.WithArgs {
		var args []string
		for _, a := range r.Args {
			// A rule argument of unknown kind for
			// a numeric type can be treated as 0.
			if a.Type == rules.ARG_UNKNOWN && o.Type.Kind.IsNumeric() {
				a = &rules.Arg{rules.ARG_INT, "0"}
			}
			// skip empty
			if a.Value == "" {
				continue
			}

			switch a.Type {
			case rules.ARG_FIELD_ABS, rules.ARG_FIELD_REL:
				// x := b.g.recv
				// for _, f := range b.g.info.KeyMap[a.Value].Selector {
				// 	x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
				// }
				args = append(args, "%v")
				// refs = append(refs, x)
			case rules.ARG_STRING:
				args = append(args, strconv.Quote(a.Value))
			default:
				args = append(args, a.Value)
			}
		}
		if len(args) > 0 {
			text += ": " + strings.Join(args, errSpec.ArgSep)
			if len(errSpec.ArgSuffix) > 0 {
				text += " " + errSpec.ArgSuffix
			}
		}
	}

	text = strconv.Quote(text)

	//if len(refs) > 0 {
	//	pkg := b.g.addImport(xtypes.Pkg{Path: "fmt"})
	//	body.Add(GO.ReturnStmt{GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, "Errorf"},
	//		Args: GO.ArgsList{List: append(GO.ExprList{textExpr}, refs...)}}})
	//} else {

	pkg := g.file.addImport(types.Pkg{Path: "errors"})

	return func() {
		g.P("$0.New($1)", pkg.name, text)
	}
}
