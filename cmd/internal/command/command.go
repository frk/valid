package command

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frk/valid/cmd/internal/config"
	"github.com/frk/valid/cmd/internal/generator"
	"github.com/frk/valid/cmd/internal/global"
	"github.com/frk/valid/cmd/internal/rules"
	"github.com/frk/valid/cmd/internal/search"
)

type Command struct {
	Cfg config.Config
}

func New(cfg config.Config) (*Command, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = printUsage

	// unmarshal cli flags into the config.
	if err := parseFlags(&cfg, fs, os.Args[1:]); err != nil {
		return nil, err
	}

	// merge with config file and then validate
	if err := cfg.MergeAndCheck(); err != nil {
		return nil, fmt.Errorf("failed to load config from file: %v", err)
	}

	// change to working dir
	if err := os.Chdir(cfg.WorkDir.Value); err != nil {
		return nil, fmt.Errorf("failed to move to working directory: %q -- %v",
			cfg.WorkDir.Value, err)
	}

	return &Command{cfg}, nil
}

func parseFlags(c *config.Config, fs *flag.FlagSet, osArgs []string) error {
	fs.Var(&c.File, "c", "")
	fs.Var(&c.WorkDir, "wd", "")
	fs.Var(&c.Recursive, "r", "")
	fs.Var(&c.FileList, "f", "")
	fs.Var(&c.FilePatternList, "rx", "")
	fs.Var(&c.OutNameFormat, "o", "")

	fs.Var(&c.ErrorHandling.FieldKey.Tag, "fk.tag", "")
	fs.Var(&c.ErrorHandling.FieldKey.Join, "fk.join", "")
	fs.Var(&c.ErrorHandling.FieldKey.Separator, "fk.sep", "")

	fs.Var(&c.ErrorHandling.Constructor, "error.constructor", "")
	fs.Var(&c.ErrorHandling.Aggregator, "error.aggregator", "")

	if err := fs.Parse(osArgs); err != nil {
		return err
	}
	return nil
}

func (cmd *Command) Run() error {
	// 1. search for validator structs
	var AST search.AST
	pkgs, err := search.Search(
		cmd.Cfg.WorkDir.Value,
		cmd.Cfg.Recursive.Value,
		cmd.Cfg.ValidatorRegexp(),
		cmd.Cfg.FileFilterFunc(),
		&AST,
	)
	if err != nil {
		return err
	}

	// 2. initialize globals, if any were specified in the config
	if err := global.Init(cmd.Cfg, &AST); err != nil {
		return err
	}

	// 3. initialize rule types
	if err := rules.InitSpecs(cmd.Cfg, &AST); err != nil {
		return err
	}

	result := make([][]*outFile, len(pkgs))
	for i, pkg := range pkgs {
		outFiles := make([]*outFile, len(pkg.Files))

		for j, file := range pkg.Files {
			out := new(outFile)
			out.path = cmd.outFilePath(file.Path)

			infos := make([]*rules.Info, len(file.Matches))
			for k, match := range file.Matches {
				// 4. rule-check the matched validator structs
				info := new(rules.Info)
				fkCfg := cmd.Cfg.ErrorHandling.FieldKey
				checker := rules.NewChecker(&AST, pkg.Pkg(), &fkCfg, info)
				if err := checker.Check(match); err != nil {
					return err
				}
				infos[k] = info
			}

			// 5. generate code
			code, err := generator.Generate(pkg.Pkg(), infos)
			if err != nil {
				return err
			}

			out.code = code
			outFiles[j] = out
		}
		result[i] = outFiles
	}

	// 6. write to file(s)
	for _, outFiles := range result {
		for _, out := range outFiles {
			if err := cmd.writeOutFile(out); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cmd *Command) outFilePath(inFilePath string) string {
	dir := filepath.Dir(inFilePath)

	name := strings.TrimSuffix(filepath.Base(inFilePath), ".go")
	outf := cmd.Cfg.OutNameFormat.Value
	if i := strings.IndexByte(outf, '%'); i > -1 {
		name = outf[:i] + name + outf[i+1:]
	}
	if !strings.HasSuffix(name, ".go") {
		name = name + ".go"
	}

	return filepath.Join(dir, name)
}

type outFile struct {
	// absolute path of the output file
	path string
	// the generated code
	code []byte
}

func (cmd *Command) writeOutFile(out *outFile) (err error) {
	f, err := os.Create(out.path)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
		if err != nil {
			os.Remove(out.path)
		}
	}()

	// make it look pretty
	bs, err := format.Source(out.code)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(bs)
	if _, err := io.Copy(f, buf); err != nil {
		return err
	}

	return f.Sync()
}
