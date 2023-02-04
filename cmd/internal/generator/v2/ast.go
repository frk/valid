package generate

type node interface {
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
	list []stmt
}

type if_stmt struct {
	init stmt
	cond expr
	body *block_stmt
	els  stmt // block-or-if
}

type for_stmt struct {
	init stmt
	cond expr
	post stmt
	body *block_stmt
}

type range_stmt struct {
	key  expr       // may be nil
	val  expr       // may be nil
	x    expr       // value to range over
	tt   token_type // = or :=
	body *block_stmt
}

type return_stmt struct {
	res []expr
}

type assign_stmt struct {
	lhs []expr
	rhs []expr
	tt  token_type // = or :=
}

type incdec_stmt struct {
	x  expr
	tt token_type
}

type expr_stmt struct {
	x expr
}

// line_comment produces a single //-style comment.
type line_comment struct {
	text string
}

// line_comment_group produces a group of //-style comments.
type line_comment_group struct {
	list []string
}

type lit_expr struct {
	tt    token_type
	value string
}

type ident_expr struct {
	name string
}

type param_expr struct {
	arg any
}

type unary_expr struct {
	root *unary_expr
	op   token_type // operator
	x    expr       // operand
}

type binary_expr struct {
	left  expr
	op    token_type
	right expr
}

type index_expr struct {
	x     expr
	index expr
}

type call_expr struct {
	fun      expr
	args     []expr
	ellipsis bool
}

type selector_expr struct {
	x   expr
	sel expr
}

type paren_expr struct {
	x expr
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
