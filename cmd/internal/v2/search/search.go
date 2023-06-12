// package search is used to find targets for the generator.
package search

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/frk/valid/cmd/internal/v2/config"

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

// AST is used to hold the packages that were loaded during a call to Search.
type AST struct {
	pkgs map[string]*packages.Package
	// used for resolving source code location
	fset *token.FileSet
}

func (a *AST) FileAndLine(obj interface{ Pos() token.Pos }) string {
	if p := a.fset.Position(obj.Pos()); p.IsValid() {
		return p.Filename + ":" + strconv.Itoa(p.Line)
	}
	return "[unknown-source-location]"
}

// add adds the given packages to the AST instance. If the given packages
// contain other imported packages then those will be added as well, and
// the imports of those packages will be added too, and so on.
func (a *AST) add(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		if _, ok := a.pkgs[pkg.PkgPath]; ok {
			// skip if already present
			continue

		}

		a.pkgs[pkg.PkgPath] = pkg

		if len(pkg.Imports) > 0 {
			imports := make([]*packages.Package, 0, len(pkg.Imports))
			for _, pkg := range pkg.Imports {
				imports = append(imports, pkg)
			}
			a.add(imports...)
		}
	}
}

// Match holds information on a matched validator struct type.
type Match struct {
	// The go/types.Named representation of the matched type.
	Named *types.Named `cmp:"+"`
	// The file set with which the matched type is associated.
	Fset *token.FileSet `cmp:"+"`
	// The source position of the matched type.
	Pos token.Pos `cmp:"+"`
}

// File represents a Go file that contains one or more matching validator struct types.
type File struct {
	Path    string
	Package *Package `cmp:"+"`
	Matches []*Match
}

// Package represents a Go package that contains one or more matching validator struct types.
type Package struct {
	Name  string
	Path  string
	Fset  *token.FileSet `cmp:"+"`
	Type  *types.Package `cmp:"+"`
	Info  *types.Info    `cmp:"+"`
	Files []*File
}

// Pkg describes a Go package. This is used by generator.Generate
// to identify the package for which the code should be generated.
type Pkg struct {
	Path string
	Name string
}

// Pkg returns the Pkg description of p.
func (p Package) Pkg() Pkg {
	return Pkg{
		Path: p.Path,
		Name: p.Name,
	}
}

// Search scans one or more Go packages looking for validator struct types
// whose names match the rxValidator regexp. The result will be a list of
// Packages, where each Package will contain a list of Files that belong to
// that Package, and each of these Files will contain a list of Matches each
// representing a validator struct type declared in that File.
//
// Scanned files and packages that do not contain any matching validator struct
// type declarations will be omitted from the result.
//
// Search will scan the Go package that is located in the specified directory and,
// optionally, if recursive is true, it will also scan the packages located in the
// hierarchy of the specified directory.
//
// If the *AST argument is not nil it will be populated with the list of
// packages that were loaded from the specified directory.
func Search(dir string, recursive bool, rxValidator *regexp.Regexp, filter func(filePath string) bool, a *AST) (out []*Package, err error) {
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

	if rxValidator == nil {
		c := config.DefaultConfig()
		rxValidator = regexp.MustCompile(c.ValidatorNamePattern.Value)
	}

	if a == nil {
		a = &AST{}
	}
	a.pkgs = make(map[string]*packages.Package)
	a.fset = token.NewFileSet()

	ldCfg := new(packages.Config)
	ldCfg.Mode = loadMode
	ldCfg.Dir = dir
	ldCfg.Fset = a.fset
	pkgs, err := packages.Load(ldCfg, pattern)
	if err != nil {
		return nil, err
	}

	// aggregate matches from all files in all packages
	for _, pkg := range pkgs {
		p := new(Package)
		p.Name = pkg.Name
		p.Path = pkg.PkgPath
		p.Fset = pkg.Fset
		p.Type = pkg.Types
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
					if !ok || !rxValidator.MatchString(typeSpec.Name.Name) || hasIgnoreDirective(typeSpec.Doc) {
						continue
					}
					if _, ok := typeSpec.Type.(*ast.StructType); !ok {
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

					match := new(Match)
					match.Named = named
					match.Fset = pkg.Fset
					match.Pos = typeName.Pos()
					f.Matches = append(f.Matches, match)
				}
			}

			// add file only if it contains any matches
			if len(f.Matches) > 0 {
				p.Files = append(p.Files, f)
			}
		}

		// add package only if it contains any matches
		if len(p.Files) > 0 {
			out = append(out, p)
		}
	}

	a.add(pkgs...)
	return out, nil
}

// hasIgnoreDirective reports whether or not the given documentation contains
// the "valid:ignore" directive indicating that the match should be ignored.
func hasIgnoreDirective(doc *ast.CommentGroup) bool {
	if doc != nil {
		for _, com := range doc.List {
			if strings.Contains(com.Text, "valid:ignore") {
				return true
			}
		}
	}
	return false
}

// findpkg searches for a *packages.Package by the given pkgpath and returns
// it if a match is found. The fname argument is optional, it would be the name
// of the function for which the package is being looked up, i.e. the "reason"
// for which findpkg was invoked, it is used for error info only.
func findpkg(pkgpath, fname string, a *AST) (*packages.Package, error) {
	pkgCache.Lock()
	defer pkgCache.Unlock()

	pkg, ok := pkgCache.m[pkgpath]
	if !ok {
		// no need to re-attempt load if it failed before
		if err, ok := pkgCache.err[pkgpath]; ok {
			return nil, err
		}

		// It is probable that the target package will already be loaded
		// in the AST instance supplied to the Search function, therefore
		// look there next and only if it's not there attempt to load it.
		if pkg, ok = a.pkgs[pkgpath]; !ok {
			ldCfg := &packages.Config{
				Mode: loadMode,
				Fset: a.fset,
			}
			pkgs, err := packages.Load(ldCfg, pkgpath)
			if err != nil {
				e := &Error{C: ERR_PKG_LOADFAIL, pkg: pkgpath, name: fname, err: err}
				pkgCache.err[pkgpath] = e
				return nil, e
			} else if len(pkgs) > 0 && len(pkgs[0].Errors) > 0 {
				e := &Error{C: ERR_PKG_ERROR, pkg: pkgpath, name: fname, err: pkgs[0].Errors[0]}
				pkgCache.err[pkgpath] = e
				return nil, e
			} else if len(pkgs) == 0 {
				e := &Error{C: ERR_PKG_NOTFOUND, pkg: pkgpath, name: fname}
				pkgCache.err[pkgpath] = e
				return nil, e
			}

			pkg = pkgs[0]
		}

		pkgCache.m[pkgpath] = pkg
	}
	return pkg, nil
}

// stores packages loaded by FindFunc.
var pkgCache = struct {
	sync.RWMutex
	m   map[string]*packages.Package
	err map[string]error
}{
	m:   make(map[string]*packages.Package),
	err: make(map[string]error),
}
