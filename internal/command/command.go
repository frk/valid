package command

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/frk/isvalid/internal/analysis"
	"github.com/frk/isvalid/internal/generator"
	"github.com/frk/isvalid/internal/search"
)

type Command struct {
	Config
}

func New(cfg Config) (*Command, error) {
	// update the working directory to its absolute path
	abs, err := filepath.Abs(cfg.WorkingDirectory.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path of working directory: %q -- %v",
			cfg.WorkingDirectory.Value, err)
	}
	cfg.WorkingDirectory.Value = abs

	// change to working dir
	if err := os.Chdir(cfg.WorkingDirectory.Value); err != nil {
		return nil, fmt.Errorf("failed to move to working directory: %q -- %v",
			cfg.WorkingDirectory.Value, err)
	}

	// check the config for errors
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &Command{cfg}, nil
}

func (cmd *Command) Run() error {
	var aConf analysis.Config
	aConf.FieldKeyTag = cmd.FieldKeyTag.Value
	aConf.FieldKeyBase = cmd.FieldKeyBase.Value
	aConf.FieldKeySeparator = cmd.FieldKeySeparator.Value

	// 1. search for validator types
	pkgs, err := search.Search(cmd.WorkingDirectory.Value, cmd.Recursive.Value, cmd.FileFilterFunc(), &aConf.AST)
	if err != nil {
		return err
	}

	// 2. find & analyze custom rule functions
	for _, rc := range cmd.CustomRules {
		f, err := search.FindFunc(rc.funcPkg, rc.funcName, aConf.AST)
		if err != nil {
			return err
		}
		if err := aConf.AddRuleFunc(rc.Name, f); err != nil {
			return err
		}
	}

	result := make([][]*outFile, len(pkgs))
	for i, pkg := range pkgs {
		outFiles := make([]*outFile, len(pkg.Files))

		for j, file := range pkg.Files {
			out := new(outFile)
			out.path = cmd.outFilePath(file.Path)
			out.targInfos = make([]*generator.TargetAnalysis, len(file.Matches))

			for k, match := range file.Matches {
				// 3. analyze targets
				aInfo := new(analysis.Info)
				vs, err := aConf.Analyze(pkg.Fset, match.Named, match.Pos, aInfo)
				if err != nil {
					return err
				}

				out.targInfos[k] = &generator.TargetAnalysis{ValidatorStruct: vs, Info: aInfo}
			}

			// 4. generate code
			if err := generator.Generate(&out.buf, pkg.Name, out.targInfos); err != nil {
				return err
			}

			outFiles[j] = out
		}
		result[i] = outFiles
	}

	// 5. write to file(s)
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
	name = fmt.Sprintf(cmd.OutputFileNameFormat.Value, name)
	if !strings.HasSuffix(name, ".go") {
		name = name + ".go"
	}

	return filepath.Join(dir, name)
}

type outFile struct {
	// absolute path of the output file
	path string
	// the type checked targets
	targInfos []*generator.TargetAnalysis
	// the generated code
	buf bytes.Buffer
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
	bs, err := format.Source(out.buf.Bytes())
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(bs)
	if _, err := io.Copy(f, buf); err != nil {
		return err
	}

	return f.Sync()
}
