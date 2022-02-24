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

	"github.com/frk/valid/cmd/internal/config"

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

// FindConstantsByType scans the given AST looking for all declared constants
// of the type identified by pkgpath and name. On success the result will be
// a slice of go/types.Const instances that represent those constants.
//
// FindConstantsByType is exepcted to be invoked *after* Search and the AST argument is
// expected to be the same as the one given to Search for caching the packages it loads.
func FindConstantsByType(pkgpath, name string, a *AST) (consts []*types.Const) {
	for _, pkg := range a.pkgs {
		if pkg.PkgPath != pkgpath {
			if _, ok := pkg.Imports[pkgpath]; !ok {
				// If pkg is not the target package, and it also
				// does not import the target package, go to next
				continue
			}
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
	}
	return consts
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

// FindFunc scans the package identified by pkgpath looking for a function
// with the given name and, if successful, returns the go/types.Func
// representation of that function.
//
// FindFunc is exepcted to be invoked *after* Search and the AST argument is expected
// to be the same as the one given to Search for caching the packages it loads.
//
// The pkgpath parameter should be the import path of a single package,
// if it's a pattern or something else then the result is undefined.
func FindFunc(pkgpath, name string, a *AST) (fn *types.Func, rawCfg []byte, err error) {
	pkg, err := findpkg(pkgpath, name, a)
	if err != nil {
		return nil, nil, err
	}

	for _, syn := range pkg.Syntax {
		for _, dec := range syn.Decls {
			fd, ok := dec.(*ast.FuncDecl)
			if !ok || fd.Recv != nil {
				continue
			}

			if fd.Name.Name == name {
				obj, ok := pkg.TypesInfo.Defs[fd.Name]
				if !ok {
					continue
				}

				if fn, ok := obj.(*types.Func); ok {
					rawCfg = extractRuleYAML(fd.Doc)
					return fn, rawCfg, nil
				}
			}
		}
	}

	return nil, nil, &Error{C: ERR_FUNC_NOTFOUND, pkg: pkgpath, name: name}
}

// FindObject returns a top-level declared object that matches
// the given pkgpath and name. The returned object will either
// be a top-level declared type or a top-level declared function.
func FindObject(pkgpath, name string, a *AST) (obj types.Object, err error) {
	pkg, err := findpkg(pkgpath, name, a)
	if err != nil {
		return nil, err
	}

	for _, syn := range pkg.Syntax {
		for _, d := range syn.Decls {
			switch d := d.(type) {
			case *ast.GenDecl:
				if d.Tok != token.TYPE {
					continue
				}
				for _, s := range d.Specs {
					s, ok := s.(*ast.TypeSpec)
					if !ok || s.Name.Name != name {
						continue
					}
					def, ok := pkg.TypesInfo.Defs[s.Name]
					if !ok {
						continue
					}
					if tn, ok := def.(*types.TypeName); ok {
						return tn, nil
					}
				}
			case *ast.FuncDecl:
				if d.Recv != nil || d.Name.Name != name {
					continue
				}
				def, ok := pkg.TypesInfo.Defs[d.Name]
				if !ok {
					continue
				}
				if f, ok := def.(*types.Func); ok {
					return f, nil
				}
			}
		}
	}

	return nil, &Error{C: ERR_OBJECT_NOTFOUND, pkg: pkgpath, name: name}
}

// FindBuiltinFuncs
//
// TODO(mkopriva): make this not blow up if package can't be found
//                 on the system, because of the following:
//
// 	It is possible that the user of the cmd/validgen tool does not
// 	have github.com/frk/valid source on the user's machine, which
// 	is ok because the source would be downloaded automatically as
// 	soon as the user attempts to run the generated code, or maybe
// 	the user does not intend to use the builtin rules, or perhaps
// 	the user has supplied a set of custom rules that override
// 	the builtin ones anyway.
//
// 	In case the error is genuine the code should keep working without
// 	issues, it's just that the reporting of user errors will be poorer.
func FindBuiltinFuncs(a *AST, callback func(fn *types.Func, rawCfg []byte) error) error {
	pkg, err := findpkg("github.com/frk/valid", "", a)
	if err != nil {
		return err
	}

	for i, syn := range pkg.Syntax {
		// all the builtin funcs are in the valid.go file,
		// if this is not it; next
		if !strings.HasSuffix(pkg.GoFiles[i], "valid.go") {
			continue
		}

		for _, dec := range syn.Decls {
			fd, ok := dec.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || !fd.Name.IsExported() {
				continue
			}

			obj, ok := pkg.TypesInfo.Defs[fd.Name]
			if !ok {
				continue
			}

			if fn, ok := obj.(*types.Func); ok {
				rawCfg := extractRuleYAML(fd.Doc)
				if len(rawCfg) == 0 {
					continue
				}

				if err := callback(fn, rawCfg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// extractRuleYAML returns the yaml bytes as parsed from the "valid:rule.yaml"
// directive in the given documentation, if no "valid:rule.yaml" directive is
// found, nil will be returned instead.
func extractRuleYAML(doc *ast.CommentGroup) (out []byte) {
	const directive = "valid:rule.yaml"

	if doc == nil {
		return nil
	}

	hasdirective := false
	for _, com := range doc.List {
		text := com.Text

		// look for directive if not yet found
		if !hasdirective {
			if i := strings.Index(text, directive); i > -1 {
				hasdirective = true
				text = text[i+len(directive):]
			}
		}

		// the rest of the doc text after a directive is expected to be yaml
		if hasdirective {
			text = strings.TrimPrefix(text, "//\t") + "\n"
			out = append(out, text...)
		}
	}
	return out
}
