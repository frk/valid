package types

type Ident interface {
	GetPkg() Pkg
	GetName() string
}

func (f *Func) GetPkg() Pkg     { return f.Type.Pkg }
func (f *Func) GetName() string { return f.Name }

func (c Const) GetPkg() Pkg     { return c.Pkg }
func (c Const) GetName() string { return c.Name }
