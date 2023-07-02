package source

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"gopkg.in/yaml.v3"
)

// Obj represents an object with a position in the source.
type Obj interface {
	Pos() token.Pos
}

// Match holds information on a matched target struct type.
type Match struct {
	// The go/types.Named representation of the matched type.
	Named *types.Named `cmp:"+"`
	// The file set with which the matched type is associated.
	Fset *token.FileSet `cmp:"+"`
	// The source position of the matched type.
	Pos token.Pos `cmp:"+"`
}

// File represents a Go file that contains one or more matching target struct types.
type File struct {
	Path    string
	Package *Package `cmp:"+"`
	Matches []*Match
}

// Package represents a Go package that contains one or more matching target struct types.
type Package struct {
	Name  string
	Path  string
	Fset  *token.FileSet `cmp:"+"`
	Type  *types.Package `cmp:"+"`
	Info  *types.Info    `cmp:"+"`
	Files []*File
}

// Pkg returns the Pkg description of p.
func (p Package) Pkg() Pkg {
	return Pkg{
		Path: p.Path,
		Name: p.Name,
	}
}

// Pkg describes a Go package. This is used by generator.Generate
// to identify the package for which the code should be generated.
type Pkg struct {
	Path string
	Name string
}

// Func contains information about a Go function's type.
type Func struct {
	*types.Func
	// The raw config associated with this specific function, if any.
	config []byte
}

// DecodeConfig will decode the Func's config, if any, into v.
func (f Func) DecodeConfig(v any) error {
	if len(f.config) > 0 {
		return yaml.Unmarshal(f.config, v)
	}
	return nil
}

type Source struct {
	// Dir is the directory in which to search for the source information.
	// If Dir is empty, search will run in the current directory.
	Dir string
	// Recursive can be set to indicate that source information can also
	// be searched for in the packages located within Dir's hierarchy.
	Recursive bool
	// TargetNamePattern is a regular expression string that
	// will be used to identify the target struct types.
	TargetNamePattern string
	// IgnoreDirective can be set to a string that will be checked
	// for in the target type's documentation to determine whether
	// to include the target in the search results or not.
	IgnoreDirective string
	// Filter can be set to filter the files that should be searched. If the
	// func returns false for a filePath then that file will not be searched.
	Filter func(filePath string) bool

	// used for resolving source code location
	fset *token.FileSet
	// the set of loaded packages
	pkgs map[string]*packages.Package
	// the regexp created from TargetNamePattern
	rxTarget *regexp.Regexp
	// the packages' load mode
	mode packages.LoadMode

	// stores packages loaded by find_func.
	cache struct {
		sync.RWMutex
		m   map[string]*packages.Package
		err map[string]error
	}
}

func (s *Source) Load() (out []*Package, err error) {
	const load_mode = packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo

	// resolve absolute dir path
	if s.Dir, err = filepath.Abs(s.Dir); err != nil {
		return nil, err
	}

	// if no filter was provided, pass all files
	if s.Filter == nil {
		s.Filter = func(string) bool { return true }
	}

	// initialize the pattern to use with packages.Load
	pattern := "."
	if s.Recursive {
		pattern = "./..."
	}

	s.mode = load_mode
	s.fset = token.NewFileSet()
	s.pkgs = make(map[string]*packages.Package)
	s.rxTarget = regexp.MustCompile(s.TargetNamePattern)
	s.cache.m = make(map[string]*packages.Package)
	s.cache.err = make(map[string]error)

	cfg := new(packages.Config)
	cfg.Mode = s.mode
	cfg.Dir = s.Dir
	cfg.Fset = s.fset
	pkgs, err := packages.Load(cfg, pattern)
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
			if filePath := pkg.CompiledGoFiles[i]; !s.Filter(filePath) {
				continue
			}

			f := new(File)
			f.Path = pkg.CompiledGoFiles[i]
			f.Package = p
			for _, dec := range syn.Decls {
				gd, ok := dec.(*ast.GenDecl)
				if !ok || gd.Tok != token.TYPE || s.has_ignore_directive(gd.Doc) {
					continue
				}

				for _, spec := range gd.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || !s.rxTarget.MatchString(typeSpec.Name.Name) || s.has_ignore_directive(typeSpec.Doc) {
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

	s.add_pkgs(pkgs...)
	return out, nil
}

// FileAndLine returns the "filename:line" source position of obj.
func (s *Source) FileAndLine(obj Obj) string {
	if p := s.fset.Position(obj.Pos()); p.IsValid() {
		return p.Filename + ":" + strconv.Itoa(p.Line)
	}
	return "[unknown-source-location]"
}

// FindObject returns a top-level declared object that matches the given
// pkgpath and name. The returned object will either be a top-level declared
// type or a top-level declared function.
func (s *Source) FindObject(pkgpath, name string) (obj types.Object, err error) {
	pkg, err := s.find_pkg(pkgpath, name)
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

// FindConsts scans the source looking for all declared constants of the type
// identified by pkgpath and name. On success the result will be a slice of
// go/types.Const instances that represent those constants.
//
// FindConsts is exepcted to be invoked *after* Load.
func (s *Source) FindConsts(pkgpath, name string) (consts []*types.Const) {
	for _, pkg := range s.pkgs {
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

// FindFunc scans the package identified by pkgpath looking for a function
// with the given name and, if successful, returns the Func representation
// of that function.
//
// The pkgpath parameter should be the import path of a single package,
// if it's a pattern or something else then the result is undefined.
//
// FindFunc is exepcted to be invoked *after* Load.
func (s *Source) FindFunc(pkgpath, name string) (*Func, error) {
	pkg, err := s.find_pkg(pkgpath, name)
	if err != nil {
		return nil, err
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
					config := s.extract_rule_config(fd.Doc)
					f := &Func{Func: fn, config: config}
					return f, nil
				}
			}
		}
	}

	return nil, &Error{C: ERR_FUNC_NOTFOUND, pkg: pkgpath, name: name}
}

// GetIncludedRuleFuncs returns a slice of Funcs that represent "rule functions"
// declared in the github.com/frk/valid/valid.go file.
//
// NOTE(mkopriva): this method assumes that the github.com/frk/valid package
// is available on the host on which the cmd/validgen tool is being executed.
func (s *Source) GetIncludedRuleFuncs() ([]*Func, error) {
	pkg, err := s.find_pkg("github.com/frk/valid", "")
	if err != nil {
		return nil, err
	}

	funcs := make([]*Func, 0)
	for i, syn := range pkg.Syntax {
		// all the included funcs are in the valid.go file, if this is not it; next
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
				config := s.extract_rule_config(fd.Doc)
				if len(config) == 0 {
					continue
				}

				f := &Func{Func: fn, config: config}
				funcs = append(funcs, f)
			}
		}
	}
	return funcs, nil
}

// has_ignore_directive reports whether or not the given documentation contains
// the Source's IgnoreDirective indicating that the match should be ignored.
func (s *Source) has_ignore_directive(doc *ast.CommentGroup) bool {
	if s.IgnoreDirective == "" {
		return false
	}

	if doc != nil {
		for _, com := range doc.List {
			if strings.Contains(com.Text, s.IgnoreDirective) {
				return true
			}
		}
	}
	return false
}

// add_pkgs adds the given packages to the AST instance. If the given packages
// contain other imported packages then those will be added as well, and the
// imports of those packages will be added too, and so on.
func (s *Source) add_pkgs(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		if _, ok := s.pkgs[pkg.PkgPath]; ok {
			// skip if already present
			continue

		}

		s.pkgs[pkg.PkgPath] = pkg

		if len(pkg.Imports) > 0 {
			imports := make([]*packages.Package, 0, len(pkg.Imports))
			for _, pkg := range pkg.Imports {
				imports = append(imports, pkg)
			}
			s.add_pkgs(imports...)
		}
	}
}

// find_pkg searches for a *packages.Package by the given pkgpath
// and returns it if a match is found.
//
// The fname argument is optional, it would be the name of the function for
// which the package is being looked up, i.e. the "reason" for which find_pkg
// was invoked, it is used for error info only.
func (s *Source) find_pkg(pkgpath, fname string) (*packages.Package, error) {
	s.cache.Lock()
	defer s.cache.Unlock()

	pkg, ok := s.cache.m[pkgpath]
	if !ok {
		// no need to re-attempt load if it failed before
		if err, ok := s.cache.err[pkgpath]; ok {
			return nil, err
		}

		// It is probable that the target package will already be loaded
		// in the AST instance supplied to the Search function, therefore
		// look there next and only if it's not there attempt to load it.
		if pkg, ok = s.pkgs[pkgpath]; !ok {
			cfg := &packages.Config{
				Mode: s.mode,
				Fset: s.fset,
			}
			pkgs, err := packages.Load(cfg, pkgpath)
			if err != nil {
				e := &Error{C: ERR_PKG_LOADFAIL, pkg: pkgpath, name: fname, err: err}
				s.cache.err[pkgpath] = e
				return nil, e
			} else if len(pkgs) > 0 && len(pkgs[0].Errors) > 0 {
				e := &Error{C: ERR_PKG_ERROR, pkg: pkgpath, name: fname, err: pkgs[0].Errors[0]}
				s.cache.err[pkgpath] = e
				return nil, e
			} else if len(pkgs) == 0 {
				e := &Error{C: ERR_PKG_NOTFOUND, pkg: pkgpath, name: fname}
				s.cache.err[pkgpath] = e
				return nil, e
			}

			pkg = pkgs[0]
		}

		s.cache.m[pkgpath] = pkg
	}
	return pkg, nil
}

// extract_rule_config returns the config bytes as parsed from the "valid:rule.yaml"
// directive in the given documentation, if no "valid:rule.yaml" directive is found,
// nil will be returned instead.
func (s *Source) extract_rule_config(doc *ast.CommentGroup) (out []byte) {
	const directive = "valid:rule.yaml"

	if doc == nil {
		return nil
	}

	has_directive := false
	for _, com := range doc.List {
		text := com.Text

		// look for directive if not yet found
		if !has_directive {
			if i := strings.Index(text, directive); i > -1 {
				has_directive = true
				text = text[i+len(directive):]
			}
		}

		if strings.TrimSpace(text) == "//" {
			continue
		}

		// the rest of the doc text after a directive is expected to be yaml
		if has_directive {
			text = strings.TrimPrefix(text, "//\t") + "\n"
			out = append(out, text...)
		}
	}
	return out
}
