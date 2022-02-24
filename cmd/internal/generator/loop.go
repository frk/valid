package generator

import (
	"fmt"

	"github.com/frk/valid/cmd/internal/gotype"
	"github.com/frk/valid/cmd/internal/rules"

	GO "github.com/frk/ast/golang"
)

var _ = fmt.Println

func (b *bb) arrayForStmt(n *rules.Node) GO.ForRangeClause {
	E := n.Elem
	rc := b.arrayRangeClause(E)
	loop := &GO.ForStmt{Clause: rc}

	if E.Type.Kind != gotype.K_PTR && len(E.PreRules) > 0 {
		b.idx = GO.IndexExpr{X: rc.X, Index: rc.Key}
	} else {
		b.idx = nil
	}

	b.add(loop)
	b.use(&loop.Body)
	return rc
}

func (b *bb) arrayRangeClause(E *rules.Node) GO.ForRangeClause {
	rc := GO.ForRangeClause{Define: true, X: b.val}
	k, v := GO.Ident{}, GO.Ident{}

	if E.HasRules() {
		v = GO.Ident{"e"}
	}
	if E.Type.Kind != gotype.K_PTR && len(E.PreRules) > 0 {
		k = GO.Ident{"i"}
	} else if v.Name != "" {
		k = GO.Ident{"_"}
	}

	if x, ok := rc.X.(GO.Ident); ok {
		if x == k {
			k.Name += "2"
		}
		if x == v {
			v.Name += "2"
		}
	}

	if k.Name != "" {
		rc.Key = k
	}
	if v.Name != "" {
		rc.Value = v
	}
	return rc
}

func (b *bb) mapForStmt(n *rules.Node) GO.ForRangeClause {
	K, E := n.Key, n.Elem
	rc := b.mapRangeClause(K, E)
	loop := &GO.ForStmt{Clause: rc}

	if E.Type.Kind != gotype.K_PTR && len(E.PreRules) > 0 {
		b.idx = GO.IndexExpr{X: rc.X, Index: rc.Key}
	} else {
		b.idx = nil
	}

	b.add(loop)
	b.use(&loop.Body)
	return rc
}

func (b *bb) mapRangeClause(K, E *rules.Node) GO.ForRangeClause {
	rc := GO.ForRangeClause{Define: true, X: b.val}
	k, v := GO.Ident{}, GO.Ident{}

	if E.HasRules() {
		v = GO.Ident{"e"}
	}
	if K.HasRules() || (E.Type.Kind != gotype.K_PTR && len(E.PreRules) > 0) {
		k = GO.Ident{"k"}
	} else if v.Name != "" {
		k = GO.Ident{"_"}
	}

	if x, ok := rc.X.(GO.Ident); ok {
		if x == k {
			k.Name += "2"
		}
		if x == v {
			v.Name += "2"
		}
	}

	if k.Name != "" {
		rc.Key = k
	}
	if v.Name != "" {
		rc.Value = v
	}
	return rc
}
