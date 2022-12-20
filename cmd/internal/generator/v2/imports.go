package generate

import (
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/types"
)

// The pkgInfo type holds info that will be used to generate an import spec.
type pkgInfo struct {
	// The package path.
	path string
	// The package name.
	name string
	// If set, it indicates that the name should be used
	// to specify the local name of the package.
	local bool
	// The number of package's with the same name. This value
	// is used by those packages to modify their name in order
	// to not cause an import conflict.
	num int
}

// genImports produces an import declaration for all
// the packages that need to be imported in the file.
func (f *fileInfo) genImports() {
	if len(f.imports) == 0 {
		return
	}

	// Split the imports into 3 groups separated by a new line.
	// - the 1st group will contain imports from the standard library
	// - the 3rd group will contain imports from github.com/frk/valid...
	// - and the 2nd group will contain the rest of the imports.
	var p1, p2, p3 []*pkgInfo
	for _, p := range f.imports {
		if strings.HasPrefix(string(p.path), `github.com/frk/valid`) {
			p3 = append(p3, p)
		} else if i := strings.IndexByte(string(p.path), '.'); i >= 0 {
			p2 = append(p2, p)
		} else {
			p1 = append(p1, p)
		}
	}

	f.wr.nl()
	f.wr.ln("import (")
	emptyline := false
	for i, p := range [][]*pkgInfo{p1, p2, p3} {
		if emptyline {
			f.wr.nl()
		}

		if len(p) == 0 {
			emptyline = false
			continue
		}

		f.genImportSpecs(p)
		emptyline = true
	}
	f.wr.ln(")")
}

// genImportSpecs gens the ImportSpecs for the provided list of packages.
func (f *fileInfo) genImportSpecs(pp []*pkgInfo) {
	for _, p := range pp {
		if p.local {
			f.wr.ln("\t$0 \"$1\"", p.name, p.path)
		} else {
			f.wr.ln("\t\"$0\"", p.path)
		}
	}
}

// addImport adds a new pkgInfo to the import
// set if it is not already a member of that set.
func (f *fileInfo) addImport(pkg types.Pkg) *pkgInfo {
	if pkg.Name == "" {
		pkg.Name = pkg.Path
		if i := strings.LastIndexByte(pkg.Name, '/'); i > -1 {
			pkg.Name = pkg.Name[i+1:]
		}
	}

	var sameName *pkgInfo
	for _, p := range f.imports {
		// already added, exit
		if p.path == pkg.Path {
			return p
		}

		// retain package that has the same name
		if p.name == pkg.Name {
			sameName = p
		}
	}

	p := &pkgInfo{path: pkg.Path, name: pkg.Name}
	if sameName != nil {
		sameName.num += 1
		p.name = pkg.Name + strconv.Itoa(sameName.num)
		p.local = true
	}

	f.imports = append(f.imports, p)
	return p
}
