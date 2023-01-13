package generate

import (
	"fmt"

	"github.com/frk/valid/cmd/internal/rules/specs"
	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

var _ = fmt.Println

func (g *generator) gen_required_pointer_code(o *types.Obj) {
	list := p_rules(o)
	base := p_end(o)

	switch {
	case len(base.IsRules) == 0:
		g.L(`if ${0[||]} {`, list)
		g.P(`	return ${0:e}`, list[0])
		g.L(`}`)

	case len(base.IsRules) > 0:
		g.L(`if ${0[||]} {`, list)
		g.L(`	return ${0:e}`, list[0])
		for _, r := range base.IsRules {
			g.L(`} else if ${0} {`, r)
			g.L(`	return ${0:e}`, r)
		}
		g.L(`}`)
	}
}

func (g *generator) gen_optional_pointer_code(o *types.Obj) {
	list := p_rules(o)
	base := p_end(o)

	switch {
	case can_join_opt(base):
		r := base.IsRules[0]
		g.L(`if ${0[&&]} && ${1:p} {`, list, r)
		g.L(`	return ${0:e}`, r)
		g.L(`}`)

	default:
		g.L(`if ${0[&&]} {`, list)
		g.P(`	${0:g}`, base)
		g.L(`}`)
	}
}

func (g *generator) gen_pre_rules_code(o *types.Obj) {
	g.L(`$0 = ${1[@]}`, o, o.PreRules)

	o.PreRules = nil
	g.genObjCode(o)
}

func (g *generator) gen_is_rules_code(o *types.Obj) {
	r := o.IsRules[0]
	list := o.IsRules[1:]

	// DO NOT move nor remove this or else the can_join_opt(o) call
	// may not produce the expected result.
	//
	// TODO: make the writer less easily breakable
	o.IsRules = list

	switch {
	case len(list) == 0:
		g.L(`if ${0} {`, r)
		g.L(`	return ${0:e}`, r)

	case r.Is(rules.OPTIONAL) && can_join_opt(o):
		r2 := list[0]
		g.L(`if ${0} && ${1:p} {`, r, r2)
		g.L(`	return ${0:e}`, r2)

	case len(list) > 0:
		g.L(`if ${0} {`, r)
		g.L(`	return ${0:e}`, r)
		for _, r := range list {
			g.L(`} else if ${0} {`, r)
			g.L(`	return ${0:e}`, r)
		}
	}

	o.IsRules = nil
	if o.Type.HasRules() {
		g.L(`} else {`)
		g.P(`	${0:g}`, o)
	}
	g.L(`}`)
}

func (g *generator) gen_rule_expr(r *rules.Rule, can_group bool) {
	switch {
	case can_group && g.is_multi_expr(r):
		g.P(`(${0})`, r)

	case r.Spec.Kind == rules.PREPROC:
		g.P(`${0}`, specs.GetFunc(r.Spec))

	default:
		o := g.info.RuleObjMap[r]
		g.genIsRuleExpr(o, r)
	}
}

func (g *generator) gen_rule_list_expr(rr []*rules.Rule, op exprOp) {
	switch {
	case len(rr) == 0:
		// nothing to do

	case op.has(func_call):
		// the last func is the leftmost
		n := len(rr) - 1
		r := rr[n]
		o := g.info.RuleObjMap[r]

		switch {
		case len(rr) == 1 && len(r.Args) == 0:
			g.P(`${0}(${1})`, r, o)
		case len(rr) == 1 && len(r.Args) > 0:
			g.P(`${0}(${1}, ${2})`, r, o, r.Args)
		case len(rr) > 1 && len(r.Args) == 0:
			g.P(`${0}(${1[@]})`, r, rr[:n])
		case len(rr) > 1 && len(r.Args) > 0:
			g.P(`${0}(${1[@]}, ${2})`, r, rr[:n], r.Args)
		}

	default:
		if op.has(unary_not) {
			g.P(`!${0} `, rr[0])
		} else {
			g.P(`${0} `, rr[0])
		}
		for _, r := range rr[1:] {
			switch {
			case op.has(unary_not, bool_and):
				g.P(`&& !${0}`, r)
			case op.has(unary_not, bool_or):
				g.P(`|| !${0}`, r)
			case op.has(bool_and):
				g.P(`&& ${0}`, r)
			case op.has(bool_or):
				g.P(`|| ${0}`, r)
			}
		}
	}
}
