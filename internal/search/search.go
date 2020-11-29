// package search is used to find targets for the generator.
package search

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

// AST is used to hold the packages that were loaded during a call to Search.
type AST struct {
	pkgs map[string]*packages.Package
}

// add adds the given packages to the AST instance. If the given packages
// contain other imported packages then those will be added as well, and
// the imports of those packages will be added too, and so on.
func (a *AST) add(pkgs ...*packages.Package) {
	if a == nil {
		return
	}
	if a.pkgs == nil {
		a.pkgs = make(map[string]*packages.Package)
	}

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
	Named *types.Named
	// The file set with which the matched type is associated.
	Fset *token.FileSet
	// The source position of the matched type.
	Pos token.Pos
}

// File represents a Go file that contains one or more matching validator struct types.
type File struct {
	Path    string
	Package *Package
	Matches []*Match
}

// Package represents a Go package that contains one or more matching validator struct types.
type Package struct {
	Name  string
	Path  string
	Fset  *token.FileSet
	Info  *types.Info
	Files []*File
}

// Search scans one or more Go packages looking for named struct types that have
// their name suffixed with "Validator", e.g. "type InputValidator struct { ...".
// The result will be a list of Packages, where each Package will contain a list
// of Files that belong to that Package, and each of these Files will contain a
// list of Matches each representing a validator struct type declared in that File.
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
func Search(dir string, recursive bool, filter func(filePath string) bool, a *AST) (out []*Package, err error) {
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

	// aggregate matches from all files in all packages
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
// the "isvalid:ignore" directive indicating that the match should be ignored.
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

// FindConstantsByType scans the given AST looking for all declared constants
// of the type identified by pkgpath and name. On success the result will be
// a slice of go/types.Const instances that represent those constants.
//
// FindConstantsByType is exepcted to be invoked *after* Search and the AST argument is
// expected to be the same as the one given to Search for caching the packages it loads.
func FindConstantsByType(pkgpath, name string, a AST) (consts []*types.Const) {
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
func findpkg(pkgpath, fname string, a AST) (*packages.Package, error) {
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
			cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax |
				packages.NeedTypes | packages.NeedTypesInfo}
			pkgs, err := packages.Load(cfg, pkgpath)
			if err != nil || len(pkgs) == 0 {
				pe := pkgLoadError{pkgpath, fname, err}
				pkgCache.err[pkgpath] = pe
				return nil, pe
			} else if len(pkgs[0].Errors) > 0 {
				pe := pkgLoadError{pkgpath, fname, pkgs[0].Errors[0]}
				pkgCache.err[pkgpath] = pe
				return nil, pe
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
func FindFunc(pkgpath, name string, a AST) (*types.Func, error) {
	pkg, err := findpkg(pkgpath, name, a)
	if err != nil {
		return nil, err
	}

	for _, syn := range pkg.Syntax {
		for _, dec := range syn.Decls {
			fd, ok := dec.(*ast.FuncDecl)
			if !ok || fd.Recv != nil || !fd.Name.IsExported() {
				continue
			}

			if fd.Name.Name == name {
				obj, ok := pkg.TypesInfo.Defs[fd.Name]
				if !ok {
					continue
				}

				if f, ok := obj.(*types.Func); ok {
					return f, nil
				}
			}
		}
	}

	return nil, findFuncError{pkgpath, name}
}

// LoadBuiltinFuncs
//
// TODO(mkopriva): make this not blow up if package can't be found
//                 on the system, because of the following:
//
// 	It is possible that the user of the cmd/isvalid tool does not
// 	have github.com/frk/isvalid source on the user's machine, which
// 	is ok because the source would be downloaded automatically as
// 	soon as the user attempts to run the generated code, or maybe
// 	the user does not intend to use the builtin rules, or perhaps
// 	the user has supplied a set of custom rules that override
// 	the builtin ones anyway.
//
// 	In case the error is genuine the code should keep working without
// 	issues, it's just that the reporting of user errors will be poorer.
func LoadBuiltinFuncs(a AST, callback func([]byte, *types.Func) error) error {
	pkg, err := findpkg("github.com/frk/isvalid", "", a)
	if err != nil {
		return err
	}

	for i, syn := range pkg.Syntax {
		// all the builtin funcs are in the isvalid.go file,
		// if this is not it; next
		if !strings.HasSuffix(pkg.GoFiles[i], "isvalid.go") {
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

			if f, ok := obj.(*types.Func); ok {
				confjson := getrulejson(fd.Doc)
				if len(confjson) == 0 {
					continue
				}

				if err := callback(confjson, f); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// getrulejson returns the json bytes as parsed from the "isvalid:rule"
// directive in the given documentation, if no "isvalid:rule" directive is
// found, nil will be returned instead.
func getrulejson(doc *ast.CommentGroup) (out []byte) {
	const directive = "isvalid:rule"

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

		// the rest of the doc text after a directive is expected to be json
		if hasdirective {
			text = strings.TrimLeft(text, "/") // remove leading //
			text = strings.TrimSpace(text)
			out = append(out, text...)
		}
	}
	return out
}

type pkgLoadError struct {
	pkgpath string
	fname   string
	err     error `cmp:"-"`
}

func (e pkgLoadError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("failed to load package %q for function %q: %v", e.pkgpath, e.fname, e.err)
	}
	return fmt.Sprintf("failed to load package %q for function %q.", e.pkgpath, e.fname)
}

type findFuncError struct {
	pkgpath string
	fname   string
}

func (e findFuncError) Error() string {
	return fmt.Sprintf("could not find function %q in package %q", e.fname, e.pkgpath)
}
