package generate

type context uint

const (
	ctx_top_level context = iota
	ctx_import_stmt
	ctx_if_stmt
	ctx_if_block
	ctx_elif_stmt
	ctx_elif_block
	ctx_else_stmt
	ctx_else_block
	ctx_for_stmt
	ctx_for_block
)

func (g *generator) fffffff() {
	// ???
}
