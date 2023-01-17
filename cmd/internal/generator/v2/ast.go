package generate

type node interface {
	build(t *token) node
	writeTo(w *writer)
}

type decl interface {
	node
	decl()
}

type expr interface {
	node
	expr()
}

type stmt interface {
	node
	stmt()
}

////////////////////////////////////////////////////////////////////////////////

type method_decl struct {
	recv    *ident_expr
	typ     expr // T or *T
	name    string
	params  field_list
	results field_list
	body    *block_stmt
}

type block_stmt struct {
	outer stmt

	list []stmt
}

type if_stmt struct {
	outer stmt // block-or-if
	init  stmt
	cond  expr
	body  *block_stmt
	els   stmt // block-or-if
}

type for_stmt struct {
	outer *block_stmt
	init  stmt
	cond  expr
	post  stmt
	body  *block_stmt
}

type range_stmt struct {
	outer *block_stmt
	key   expr       // may be nil
	val   expr       // may be nil
	x     expr       // value to range over
	tt    token_type // = or :=
	body  *block_stmt
}

type return_stmt struct {
	outer *block_stmt
	res   []expr
}

type assign_stmt struct {
	outer stmt
	lhs   []expr
	rhs   []expr
	tt    token_type // = or :=
}

type incdec_stmt struct {
	outer stmt
	x     expr
	tt    token_type
}

type expr_stmt struct {
	x expr
}

type line_comment struct {
	text string
}

type lit_expr struct {
	outer stmt
	tt    token_type
	value string
}

type ident_expr struct {
	outer stmt
	name  string
}

type param_expr struct {
	outer stmt
	arg   any
}

type unary_expr struct {
	outer stmt
	root  *unary_expr
	op    token_type // operator
	x     expr       // operand
}

type binary_expr struct {
	outer stmt
	left  expr
	op    token_type
	right expr
}

type index_expr struct {
	outer stmt
	x     expr
	index expr
}

type call_expr struct {
	outer    stmt
	fun      expr
	args     []expr
	ellipsis bool
}

type selector_expr struct {
	outer stmt
	x     expr
	sel   expr
}

type paren_expr struct {
	outer expr
	x     expr
}

type expr_list struct {
	items []expr
}

type field_list []*field_item

type field_item struct {
	names    []*ident_expr
	typ      expr
	variadic bool
}

////////////////////////////////////////////////////////////////////////////////

func (*method_decl) decl() {}

func (*block_stmt) stmt()   {}
func (*if_stmt) stmt()      {}
func (*for_stmt) stmt()     {}
func (*range_stmt) stmt()   {}
func (*return_stmt) stmt()  {}
func (*assign_stmt) stmt()  {}
func (*incdec_stmt) stmt()  {}
func (*expr_stmt) stmt()    {}
func (*line_comment) stmt() {}

func (*lit_expr) expr()      {}
func (*ident_expr) expr()    {}
func (*unary_expr) expr()    {}
func (*binary_expr) expr()   {}
func (*index_expr) expr()    {}
func (*call_expr) expr()     {}
func (*selector_expr) expr() {}
func (*paren_expr) expr()    {}
func (*expr_list) expr()     {}
func (*param_expr) expr()    {}
