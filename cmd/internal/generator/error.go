package generator

import (
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/global"
	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

func (b *bb) err(n *rules.Node, r *rules.Rule, body *GO.BlockStmt) {
	if body == nil {
		// default to current block statement
		body = b.cur
	}

	switch {
	case b.g.info.Validator.ErrorHandlerField != nil:
		b.errHandler(n, r, body)
	case global.ErrorAggregator != nil:
		b.errGlobalHandler(n, r, body, true)
	case global.ErrorConstructor != nil:
		b.errGlobalHandler(n, r, body, false)
	default:
		b.errDefault(n, r, body)
	}
}

func (b *bb) errGlobalHandler(n *rules.Node, r *rules.Rule, body *GO.BlockStmt, isAgg bool) {
	args := make(GO.ExprList, 3)
	args[0] = GO.StringLit(b.key)
	args[1] = b.rootv()
	args[2] = GO.StringLit(r.Name)

	if r.Spec.Kind == rules.ENUM {
		enums := b.g.enumap[r]
		args = append(args, enums...)
	} else {
		for _, a := range r.Args {
			switch a.Type {
			case rules.ARG_FIELD_ABS, rules.ARG_FIELD_REL:
				x := b.g.recv
				for _, f := range b.g.info.KeyMap[a.Value].Selector {
					x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
				}
				args = append(args, x)
			case rules.ARG_STRING:
				args = append(args, GO.StringLit(a.Value))
			case rules.ARG_UNKNOWN:
				args = append(args, GO.StringLit(""))
			default:
				args = append(args, GO.ValueLit(a.Value))
			}
		}
	}

	if isAgg {
		x := GO.ExprNode(nil)
		x = GO.SelectorExpr{X: GO.Ident{"ea"}, Sel: GO.Ident{"Error"}}
		x = GO.CallExpr{Fun: x, Args: GO.ArgsList{List: args}}
		body.Add(GO.ExprStmt{x})
	} else {
		ctor := global.ErrorConstructor
		pkg := b.g.addImport(ctor.Type.Pkg)

		x := GO.ExprNode(nil)
		x = GO.QualifiedIdent{pkg.name, ctor.Name}
		x = GO.CallExpr{Fun: x, Args: GO.ArgsList{List: args}}
		body.Add(GO.ReturnStmt{Result: x})
	}
}

func (b *bb) errHandler(n *rules.Node, r *rules.Rule, body *GO.BlockStmt) {
	args := make(GO.ExprList, 3)
	args[0] = GO.StringLit(b.key)
	args[1] = b.rootv()
	args[2] = GO.StringLit(r.Name)

	if r.Spec.Kind == rules.ENUM {
		enums := b.g.enumap[r]
		args = append(args, enums...)
	} else {
		for _, a := range r.Args {
			switch a.Type {
			case rules.ARG_FIELD_ABS, rules.ARG_FIELD_REL:
				x := b.g.recv
				for _, f := range b.g.info.KeyMap[a.Value].Selector {
					x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
				}
				args = append(args, x)
			case rules.ARG_STRING:
				args = append(args, GO.StringLit(a.Value))
			case rules.ARG_UNKNOWN:
				args = append(args, GO.StringLit(""))
			default:
				args = append(args, GO.ValueLit(a.Value))
			}
		}
	}

	x := GO.ExprNode(nil)
	h := b.g.info.Validator.ErrorHandlerField

	x = GO.SelectorExpr{X: b.g.recv, Sel: GO.Ident{h.Name}}
	x = GO.SelectorExpr{X: x, Sel: GO.Ident{"Error"}}
	x = GO.CallExpr{Fun: x, Args: GO.ArgsList{List: args}}
	if h.IsAggregator {
		body.Add(GO.ExprStmt{x})
	} else {
		body.Add(GO.ReturnStmt{Result: x})
	}
}

func (b *bb) errDefault(n *rules.Node, r *rules.Rule, body *GO.BlockStmt) {
	cfg := r.Spec.Err
	if len(r.Spec.ErrOpts) > 0 && len(r.Args) > 0 {
		var key string
		for _, a := range r.Args {
			key += ":"
			if len(a.Value) > 0 {
				key += "x"
			}
		}

		key = key[1:]
		if c, ok := r.Spec.ErrOpts[key]; ok {
			cfg = c
		}
	}

	text := cfg.Text
	if len(text) == 0 {
		text = "is not valid"
	}
	text = b.key + " " + text
	//////////////////////////////

	var refs GO.ExprList
	if cfg.WithArgs {
		var args []string
		for _, arg := range r.Args {
			// A rule argument of unknown kind for
			// a numeric type can be treated as 0.
			if arg.Type == rules.ARG_UNKNOWN && n.Type.Kind.IsNumeric() {
				arg = &rules.Arg{rules.ARG_INT, "0"}
			}

			// skip empty
			if arg.Value == "" {
				continue
			}

			switch arg.Type {
			case rules.ARG_FIELD_ABS, rules.ARG_FIELD_REL:
				x := b.g.recv
				for _, f := range b.g.info.KeyMap[arg.Value].Selector {
					x = GO.SelectorExpr{X: x, Sel: GO.Ident{f.Name}}
				}
				args = append(args, "%v")
				refs = append(refs, x)
			case rules.ARG_STRING:
				args = append(args, strconv.Quote(arg.Value))
			default:
				args = append(args, arg.Value)
			}
		}
		if len(args) > 0 {
			text += ": " + strings.Join(args, cfg.ArgSep)
			if len(cfg.ArgSuffix) > 0 {
				text += " " + cfg.ArgSuffix
			}
		}
	}
	textExpr := GO.ValueLit(strconv.Quote(text))

	if len(refs) > 0 {
		pkg := b.g.addImport(gotype.Pkg{Path: "fmt"})
		body.Add(GO.ReturnStmt{GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, "Errorf"},
			Args: GO.ArgsList{List: append(GO.ExprList{textExpr}, refs...)}}})
	} else {
		pkg := b.g.addImport(gotype.Pkg{Path: "errors"})
		body.Add(GO.ReturnStmt{GO.CallExpr{Fun: GO.QualifiedIdent{pkg.name, "New"},
			Args: GO.ArgsList{List: GO.ExprList{textExpr}}}})
	}
}
