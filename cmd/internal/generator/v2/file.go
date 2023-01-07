package generate

import (
	"strconv"
	"strings"

	"github.com/frk/valid/cmd/internal/rules/v2"
	"github.com/frk/valid/cmd/internal/types"
)

type fileInfo struct {
	pkg     types.Pkg
	imports []*pkgInfo
	// A map of regular expressions that need
	// to be registered in the file's init func.
	reMap map[string]int
	// The generators used to generate the validator methods
	vms []*generator
}

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

// addRegExp adds a new regular expression string to the reMap.
func (f *fileInfo) addRegExp(a *rules.Arg) {
	if f.reMap == nil {
		f.reMap = make(map[string]int)
	}

	if _, ok := f.reMap[a.Value]; !ok {
		f.reMap[a.Value] = len(f.reMap)
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
