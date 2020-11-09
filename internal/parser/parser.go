package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

var (
	// Matches names of types that are valid targets for the generator.
	rxTargetName = regexp.MustCompile(`(?i:validator)$`)
)

type Target struct {
	Named *types.Named
	Pos   token.Pos
}

type File struct {
	Path    string
	Package *Package
	Targets []*Target
}

type Package struct {
	Name  string
	Path  string
	Fset  *token.FileSet
	Info  *types.Info
	Files []*File
}

type AST struct {
	pkgs []*packages.Package
}

func (a *AST) add(pkgs ...*packages.Package) {
	if a == nil {
		return
	}

loop:
	for _, pkg := range pkgs {
		for i := range a.pkgs {
			if pkg.PkgPath == a.pkgs[i].PkgPath {
				// skip if already in slice
				continue loop
			}
		}
		a.pkgs = append(a.pkgs, pkg)
	}
}

// Parse parses Go packages at the given dir / pattern. Only packages that contain
// files with type declarations that match the standard "isvalid" targets will be
// included in the returned slice.
func Parse(dir string, recursive bool, filter func(filePath string) bool, a *AST) (out []*Package, err error) {
	// resolve absolute dir path
	if dir, err = filepath.Abs(dir); err != nil {
		return nil, err
	}

	// if no filter was provided, pass all files
	if filter == nil {
		filter = func(string) bool { return true }
	}

	// initialize the pattern to use with packages.Load
	pattern := "."
	if recursive {
		pattern = "./..."
	}

	loadConfig := new(packages.Config)
	loadConfig.Mode = loadMode
	loadConfig.Dir = dir
	loadConfig.Fset = token.NewFileSet()
	pkgs, err := packages.Load(loadConfig, pattern)
	if err != nil {
		return nil, err
	}

	// aggregate targets from all files in all packages
	for _, pkg := range pkgs {
		p := new(Package)
		p.Name = pkg.Name
		p.Path = pkg.PkgPath
		p.Fset = pkg.Fset
		p.Info = pkg.TypesInfo

		for i, syn := range pkg.Syntax {
			// ignore file?
			if filePath := pkg.CompiledGoFiles[i]; !filter(filePath) {
				continue
			}

			f := new(File)
			f.Path = pkg.CompiledGoFiles[i]
			f.Package = p
			for _, dec := range syn.Decls {
				gd, ok := dec.(*ast.GenDecl)
				if !ok || gd.Tok != token.TYPE || hasIgnoreDirective(gd.Doc) {
					continue
				}

				for _, spec := range gd.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || !rxTargetName.MatchString(typeSpec.Name.Name) || hasIgnoreDirective(typeSpec.Doc) {
						continue
					}

					obj, ok := p.Info.Defs[typeSpec.Name]
					if !ok {
						continue
					}
					typeName, ok := obj.(*types.TypeName)
					if !ok {
						continue
					}
					named, ok := typeName.Type().(*types.Named)
					if !ok {
						continue
					}

					target := new(Target)
					target.Named = named
					target.Pos = typeName.Pos()
					f.Targets = append(f.Targets, target)
				}
			}

			// add file only if it declares targets
			if len(f.Targets) > 0 {
				p.Files = append(p.Files, f)
			}
		}

		// add package only if it has files with targets
		if len(p.Files) > 0 {
			out = append(out, p)
		}
	}

	a.add(pkgs...)
	return out, nil
}

// ParseFunc parses the given package looking for the specified function and,
// if successful returns the go/types.Func representation of that function.
// The provided fpkg should be the import path of a single package, if, instead,
// a pattern is provided then the result is undefined.
func ParseFunc(fpkg, fname string, a *AST) (*types.Func, error) {
	pkgCache.Lock()
	defer pkgCache.Unlock()

	pkg, ok := pkgCache.m[fpkg]
	if !ok {
		cfg := &packages.Config{Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo}
		pkgs, err := packages.Load(cfg, fpkg)
		if err != nil {
			return nil, fmt.Errorf("failed to load package %q for function %q: %v", fpkg, fname, err)
		}

		pkg = pkgs[0]
		pkgCache.m[fpkg] = pkg
	}

	for _, syn := range pkg.Syntax {
		for _, dec := range syn.Decls {
			fd, ok := dec.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || !fd.Name.IsExported() {
				continue
			}

			if fd.Name.Name == fname {
				obj, ok := pkg.TypesInfo.Defs[fd.Name]
				if !ok {
					continue
				}

				if f, ok := obj.(*types.Func); ok {
					a.add(pkg)
					return f, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("could not find function %q in package %q", fname, fpkg)
}

// FindConstantsByType looks for and returns all declared constants of the
// type identified by pkgpath and name in the given AST.
func FindConstantsByType(pkgpath, name string, a AST) (consts []*types.Const) {
	done := map[string]struct{}{}
	for _, pkg := range a.pkgs {
		cc := findConstantsByType(pkgpath, name, pkg, done)
		consts = append(consts, cc...)
	}
	return consts
}

// findConstantsByType recursively looks for and returns all declared constants of
// the type identified by pkgpath and name in the given pkg and all its imported pkgs.
func findConstantsByType(pkgpath, name string, pkg *packages.Package, done map[string]struct{}) (consts []*types.Const) {
	// already done, exit
	if _, ok := done[pkg.PkgPath]; ok {
		return nil
	}

	// does not import & is not the target package, exit
	if _, ok := pkg.Imports[pkgpath]; !ok && pkgpath != pkg.PkgPath {
		return nil
	}

	for _, syn := range pkg.Syntax {
		for _, dec := range syn.Decls {
			gd, ok := dec.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}

			for _, spec := range gd.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for _, id := range valueSpec.Names {
					obj, ok := pkg.TypesInfo.Defs[id]
					if !ok {
						continue
					}

					if c, ok := obj.(*types.Const); ok {
						named, ok := c.Type().(*types.Named)
						if !ok {
							continue
						}

						tn := named.Obj()
						if tn.Name() != name || tn.Pkg().Path() != pkgpath {
							continue
						}

						consts = append(consts, c)
					}
				}
			}
		}

	}
	done[pkg.PkgPath] = struct{}{}

	for _, p := range pkg.Imports {
		cc := findConstantsByType(pkgpath, name, p, done)
		consts = append(consts, cc...)
	}
	return consts
}

// hasIgnoreDirective reports whether or not the given documentation contains the "isvalid:ignore" directive.
func hasIgnoreDirective(doc *ast.CommentGroup) bool {
	if doc != nil {
		for _, com := range doc.List {
			if strings.Contains(com.Text, "isvalid:ignore") {
				return true
			}
		}
	}
	return false
}

var pkgCache = struct {
	sync.RWMutex
	m map[string]*packages.Package
}{m: make(map[string]*packages.Package)}
