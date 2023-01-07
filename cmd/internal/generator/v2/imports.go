package generate

import (
	"strings"
)

// genImports produces an import declaration for all
// the packages that need to be imported in the file.
func (g *generator) genImports() {
	if len(g.file.imports) == 0 {
		return
	}

	// Split the imports into 3 groups separated by a new line.
	// - the 1st group will contain imports from the standard library
	// - the 3rd group will contain imports from github.com/frk/valid...
	// - and the 2nd group will contain the rest of the imports.
	var p1, p2, p3 []*pkgInfo
	for _, p := range g.file.imports {
		if strings.HasPrefix(string(p.path), `github.com/frk/valid`) {
			p3 = append(p3, p)
		} else if i := strings.IndexByte(string(p.path), '.'); i >= 0 {
			p2 = append(p2, p)
		} else {
			p1 = append(p1, p)
		}
	}

	g.L(``)
	g.L("import (")
	emptyline := false
	for _, p := range [][]*pkgInfo{p1, p2, p3} {
		if len(p) == 0 {
			continue
		}
		if emptyline {
			g.L(``)
		}
		g.genImportSpecs(p)
		emptyline = true
	}
	g.L(")")
	g.L(``)
}

// genImportSpecs gens the ImportSpecs for the provided list of packages.
func (g *generator) genImportSpecs(pp []*pkgInfo) {
	for _, p := range pp {
		if p.local {
			g.L("\t$0 \"$1\"", p.name, p.path)
		} else {
			g.L("\t\"$0\"", p.path)
		}
	}
}
