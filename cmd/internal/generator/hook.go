package generator

import (
	"github.com/frk/valid/cmd/internal/xtypes"

	GO "github.com/frk/ast/golang"
)

func hookAST(h *xtypes.MethodInfo, b bb) {
	call := GO.CallExpr{}
	call.Fun = GO.SelectorExpr{X: b.g.recv, Sel: GO.Ident{h.Name}}

	ifs := new(GO.IfStmt)
	ifs.Init = GO.AssignStmt{Token: GO.AssignDefine, Lhs: ERR, Rhs: call}
	ifs.Cond = GO.BinaryExpr{Op: GO.BinaryNeq, X: ERR, Y: NIL}
	ifs.Body = GO.BlockStmt{[]GO.StmtNode{GO.ReturnStmt{ERR}}}

	b.add(ifs)
}
