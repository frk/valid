package check

import (
	"github.com/frk/valid/cmd/internal/types"
)

func (c *checker) makeFPath(ff types.FieldChain) (path string) {
	sep := "."
	for _, f := range ff {
		path += f.Name + sep
	}
	if len(path) > 0 {
		path = path[:len(path)-1] // drop the last "."
	}
	return path
}
