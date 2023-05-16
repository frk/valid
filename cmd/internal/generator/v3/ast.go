package generate

type node interface {
	build(g *gg, args []any) error
}

type decl interface {
	node
	_decl()
}

type expr interface {
	node
	_expr()
}

type stmt interface {
	node
	_stmt()
}

////////////////////////////////////////////////////////////////////////////////

type file struct {
	doc []string
	pkg *pkgclause
}

func (n *file) build(g *gg, args []any) error {
	switch g.t.t {
	case t_comment:
		n.doc = append(n.doc, g.t.v)
	case t_package:
		n.pkg = new(pkg_clause)
		g.build(n.pkg, args)
		// if err := g.next(); err != nil {
		// 	return err
		// }
		// if g.t.t != t_ident {
		// 	// ...
		// }
		// ... = g.t.v
	}
	return nil
}

type pkg_clause struct {
	name string
}

func (n *pkg_clause) build(g *gg, args []any) error {
	// ...
}
