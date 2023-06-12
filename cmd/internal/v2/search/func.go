package search

import (
	"go/ast"
	"go/types"
	"strings"

	"gopkg.in/yaml.v3"
)

// FindFunc scans the package identified by pkgpath looking for a function
// with the given name and, if successful, returns the go/types.Func
// representation of that function.
//
// FindFunc is exepcted to be invoked *after* Search and the AST argument is expected
// to be the same as the one given to Search for caching the packages it loads.
//
// The pkgpath parameter should be the import path of a single package,
// if it's a pattern or something else then the result is undefined.
func FindFunc(pkgpath, name string, a *AST) (fn *types.Func, cd ConfigDecoder, err error) {
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
					rawCfg := extractRuleYAML(fd.Doc)
					return fn, configDecoder{rawCfg}, nil
				}
			}
		}
	}

	return nil, nil, &Error{C: ERR_FUNC_NOTFOUND, pkg: pkgpath, name: name}
}

// FindIncludedFuncs
//
// TODO(mkopriva): make this not blow up if package can't be found
// on the system, because of the following:
//
// It is possible that the user of the cmd/validgen tool does not
// have github.com/frk/valid source on the user's machine, which
// is ok because the source would be downloaded automatically as
// soon as the user attempts to run the generated code, or maybe
// the user does not intend to use the included rules, or perhaps
// the user has supplied a set of custom rules that override
// the included ones anyway.
//
// In case the error is genuine the code should keep working without
// issues, it's just that the reporting of user errors will be poorer.
func FindIncludedFuncs(a *AST, callback func(fn *types.Func, cd ConfigDecoder) error) error {
	pkg, err := findpkg("github.com/frk/valid", "", a)
	if err != nil {
		return err
	}

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
				rawCfg := extractRuleYAML(fd.Doc)
				if len(rawCfg) == 0 {
					continue
				}

				cd := configDecoder{rawCfg}
				if err := callback(fn, cd); err != nil {
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

////////////////////////////////////////////////////////////////////////////////
// helpers

type ConfigDecoder interface {
	DecodeConfig(v any) error
}

type configDecoder struct {
	rawYAML []byte
}

func (d configDecoder) DecodeConfig(v any) error {
	return yaml.Unmarshal(d.rawYAML, v)
}
