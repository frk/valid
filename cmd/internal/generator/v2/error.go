package generate

import (
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
	"github.com/frk/valid/cmd/internal/types/global"
)

func (g *generator) gen_error_expr(r *rules.Rule) {
	switch {
	case g.info.Validator.ErrorHandlerField != nil:
		g.gen_handler_error(r)
	case global.ErrorAggregator != nil:
		g.gen_global_error(r, true)
	case global.ErrorConstructor != nil:
		g.gen_global_error(r, false)
	default:
		g.gen_default_error(r)
	}
}

func (g *generator) gen_handler_error(r *rules.Rule) {
	// args := make(GO.ExprList, 3)
	// args[0] = GO.StringLit(b.key)
	// args[1] = b.rootv()
	// args[2] = GO.StringLit(r.Name)

	// x := GO.ExprNode(nil)
	// h := b.g.info.Validator.ErrorHandlerField

	// x = GO.SelectorExpr{X: b.g.recv, Sel: GO.Ident{h.Name}}
	// x = GO.SelectorExpr{X: x, Sel: GO.Ident{"Error"}}
	// x = GO.CallExpr{Fun: x, Args: GO.ArgsList{List: args}}
	// if h.IsAggregator {
	// 	body.Add(GO.ExprStmt{x})
	// } else {
	// 	body.Add(GO.ReturnStmt{Result: x})
	// }
}

func (g *generator) gen_global_error(r *rules.Rule, is_agg bool) {
	var f any
	if is_agg {
		// TODO this needs to be a method call on
		// the generated instance of the aggregator
		f = global.ErrorAggregator
	} else {
		f = global.ErrorConstructor
	}

	// TODO include enums in args if rule is "enum"

	o := g.info.RuleObjMap[r]
	g.P(`${0}(${1:fk}, ${1}, ${2}`, f, o, r.Name)
	if len(r.Args) > 0 {
		g.P(`, ${0:any}`, r.Args)
	}
	g.P(`)`)
}

func (g *generator) gen_default_error(r *rules.Rule) {
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
